// pkginfo keeps package information of an application for gottani.
// It also implements types.Importer for parsing files and checking types.
package pkginfo // import "github.com/ktateish/gottani/internal/pkginfo"

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"
)

// fakeCbpkg is *build.Pacakge for `import "C"`
var fakeCbpkg *build.Package

func init() {
	fakeCbpkg = &build.Package{
		Name:       "C",
		ImportPath: "C",
		Dir:        filepath.Join(build.Default.GOROOT, "src", "C"),
		Goroot:     true,
	}
}

// key for maps
type pkgKey struct {
	importPath string
	dir        string
}

// PackageInfo represents information of packages used by a applicaion for gottani.
// It also implements types.Importer for parsing and type-checking.
type PackageInfo struct {
	// mapping (dir, importPath) => *build.Packages
	pkgs map[pkgKey]*build.Package

	// typesPkgs keeps mapping from *build.Package to *types.Package
	typesPkgs map[*build.Package]*types.Package

	// astFiles keeps mapping from *build.Package to []*ast.File
	astFiles map[*build.Package][]*ast.File

	fset  *token.FileSet // sotre all files of whole application
	tinfo *types.Info    // store type information of whole application

	rootPackage *build.Package

	// memo for Pacakges() and AllPackages()
	pkgSlice    []*build.Package
	allPkgSlice []*build.Package
}

// New creates PackageInfo with default setting and then Load the given dir
func New(dir string) (*PackageInfo, error) {
	fset := token.NewFileSet()
	tinfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	pi := NewPackageInfo(fset, tinfo)
	err := pi.Load(dir)
	return pi, err
}

// NewPackageInfo creates PackageInfo
func NewPackageInfo(fset *token.FileSet, tinfo *types.Info) *PackageInfo {
	return &PackageInfo{
		fset:  token.NewFileSet(),
		tinfo: tinfo,

		pkgs: make(map[pkgKey]*build.Package),

		typesPkgs: make(map[*build.Package]*types.Package),
		astFiles:  make(map[*build.Package][]*ast.File),
	}
}

// Load loads package information of the application given by the dir
func (ip *PackageInfo) Load(dir string) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("getting absolute path %q: %w", dir, err)
	}
	bp, err := build.ImportDir(abs, 0)
	if err != nil {
		return fmt.Errorf("importing %q: %w", abs, err)
	}
	ip.rootPackage = bp
	ip.pkgs[pkgKey{".", abs}] = bp
	ip.pkgs[pkgKey{bp.ImportPath, "."}] = bp

	tp, err := ip.typeCheck(bp)
	if err != nil {
		return fmt.Errorf("type checking: %w", err)
	}
	ip.typesPkgs[bp] = tp

	return nil
}

// Packages returns a slice of *build.Packages that are depended by the loaded package and non standard ones.
func (ip *PackageInfo) Packages() []*build.Package {
	if ip.pkgSlice != nil {
		return ip.pkgSlice
	}
	var res []*build.Package
	ip.WalkPackages(func(bp *build.Package) bool { return bp.Goroot }, nil, func(bp *build.Package, _ *types.Package, _ []*ast.File) {
		res = append(res, bp)
	})
	ip.pkgSlice = res
	return res
}

// Packages returns a slice of *build.Packages that are depended by the loaded package.
func (ip *PackageInfo) AllPackages() []*build.Package {
	if ip.allPkgSlice != nil {
		return ip.allPkgSlice
	}
	var res []*build.Package
	ip.WalkPackages(nil, nil, func(bp *build.Package, _ *types.Package, _ []*ast.File) {
		res = append(res, bp)
	})
	ip.allPkgSlice = res
	return res
}

// Root() returns the package that is in the directory given on Load().
func (ip *PackageInfo) Root() *build.Package {
	return ip.rootPackage
}

// FileSet() returns the *token.FileSet that holds all source code for whole the application except standard packages.
func (ip *PackageInfo) FileSet() *token.FileSet {
	return ip.fset
}

// TypesInfo() returns the *types.Info  that holds all type information for whole the application.
func (ip *PackageInfo) TypesInfo() *types.Info {
	return ip.tinfo
}

// GetBuildPackage() returns the *build.Package coresponds for the given importPath on the given dir.
func (ip *PackageInfo) GetBuildPackage(importPath, dir string) *build.Package {
	bp, err := ip.getBuildPackage(importPath, dir)
	if err != nil {
		panic(err)
	}
	return bp
}

// GetTypesPackage() returns the *type.Package for the package specified by the given bp.
func (ip *PackageInfo) GetTypesPackage(bp *build.Package) *types.Package {
	return ip.typesPkgs[bp]
}

// GetTypesPackage() returns the slice of *ast.File for the package specified by the given bp.
// It includs Go files and Cgo files.
func (ip *PackageInfo) GetAstFiles(bp *build.Package) []*ast.File {
	return ip.astFiles[bp]
}

// Walker is alias to a function type for PackageInfo.WalkPacakge().
type Walker func(bp *build.Package, tp *types.Package, asts []*ast.File)

// WalkPackages travarse packages and do the given pre and/or post function for each package
func (ip *PackageInfo) WalkPackages(prune func(*build.Package) bool, pre Walker, post Walker) {
	visited := make(map[*build.Package]bool)
	var rec func(bp *build.Package)
	rec = func(bp *build.Package) {
		if visited[bp] {
			return
		}
		visited[bp] = true
		if prune != nil && prune(bp) {
			return
		}
		if pre != nil {
			pre(bp, ip.typesPkgs[bp], ip.astFiles[bp])
		}
		for _, ipath := range bp.Imports {
			next, err := ip.getBuildPackage(ipath, bp.ImportPath)
			if err != nil {
				msg := fmt.Sprintf("package %q not found on %q", ipath, bp.ImportPath)
				panic(msg)
			}
			rec(next)
		}
		if post != nil {
			post(bp, ip.typesPkgs[bp], ip.astFiles[bp])
		}
	}
	rec(ip.rootPackage)
}

// getBuildPackage finds *build.Packages matched to the given importPath on the given dir.
func (pi *PackageInfo) getBuildPackage(importPath, dir string) (*build.Package, error) {
	key := pkgKey{importPath, dir}
	if bp, ok := pi.pkgs[key]; ok {
		return bp, nil
	}

	var bp *build.Package
	if importPath == "C" {
		// Always returns fake package for importPath "C" on any directory because "C" package doesn't exist
		bp = fakeCbpkg
	} else {
		abs, err := getImportDirAbs(importPath, dir)
		absKey := pkgKey{".", abs}
		var ok bool
		bp, ok = pi.pkgs[absKey]
		if !ok {
			bp, err = build.ImportDir(abs, build.AllowBinary)
			if err != nil {
				return nil, err
			}
			bp.ImportPath = importPath
			pi.pkgs[absKey] = bp
		}
	}
	pi.pkgs[key] = bp
	return bp, nil
}

// methods implements types.ImporterFrom and helper functions

// Import imports package specified by the path and returns its type information
func (ip *PackageInfo) Import(path string) (*types.Package, error) {
	return ip.ImportFrom(path, ".", 0)
}

// Import imports package specified by the given path and dir, then returns its type information
// The mode must be set 0. It is reserved for future use.
func (ip *PackageInfo) ImportFrom(path, dir string, mode types.ImportMode) (*types.Package, error) {
	bp, err := ip.getBuildPackage(path, dir)
	if err != nil {
		return nil, fmt.Errorf("getting *build.Package for %q on %q: %w", path, dir, err)
	}
	if tp, ok := ip.typesPkgs[bp]; ok {
		return tp, nil
	}

	tp, err := ip.typeCheck(bp)
	if err != nil {
		return nil, err
	}

	ip.typesPkgs[bp] = tp

	return tp, nil
}

func (ip *PackageInfo) typeCheck(bp *build.Package) (*types.Package, error) {
	if bp.Goroot && bp.ImportPath == "unsafe" {
		return types.Unsafe, nil
	}

	files, err := parsePackage(ip.fset, bp)
	if err != nil {
		return nil, fmt.Errorf("parsing package: %w", err)
	}

	var hardErrors, softErrors []error
	tcfg := types.Config{
		IgnoreFuncBodies: bp.Goroot, // doesn't check function body if the package is standard (in GOROOT/src)
		FakeImportC:      true,
		Error: func(err error) {
			if terr, ok := err.(types.Error); ok && !terr.Soft {
				hardErrors = append(hardErrors, err)
			} else {
				softErrors = append(softErrors, err)
			}
		},
		Importer: ip,
	}

	tp, err := tcfg.Check(bp.ImportPath, ip.fset, files, ip.tinfo)
	if err != nil {
		if 0 < len(hardErrors) {
			return nil, fmt.Errorf("type checking: %w", hardErrors[0])
		} else {
			return nil, fmt.Errorf("type checking: %w", err)
		}
	}

	var imps []*types.Package
	for _, ipath := range bp.Imports {
		ibp, err := ip.getBuildPackage(ipath, bp.ImportPath)
		if err != nil {
			return nil, fmt.Errorf("getting *build.Package for %q on %q: %w", ipath, bp.ImportPath, err)
		}
		p, ok := ip.typesPkgs[ibp]
		if !ok || p == nil {
			continue
		}
		imps = append(imps, p)
	}
	tp.SetImports(imps)

	ip.astFiles[bp] = files
	ip.typesPkgs[bp] = tp

	return tp, nil
}

func parsePackage(fset *token.FileSet, bp *build.Package) ([]*ast.File, error) {
	files := make([]string, 0, len(bp.GoFiles)+len(bp.CgoFiles))
	for _, f := range bp.GoFiles {
		files = append(files, f)
	}
	for _, f := range bp.CgoFiles {
		files = append(files, f)
	}

	res := make([]*ast.File, 0, len(files))
	for _, f := range files {
		var fname string
		if bp.Name == "main" {
			fname = f
		} else {
			fname = filepath.Join(bp.ImportPath, f)
		}
		path := filepath.Join(bp.Dir, f)
		af, err := parseFile(fset, fname, path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		res = append(res, af)
	}
	return res, nil
}

func parseFile(fset *token.FileSet, fname, path string) (*ast.File, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	f, err := parser.ParseFile(fset, fname, b, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return f, nil
}

// getImportDirAbs finds abs path of the package pointed by the gvien importPath on the given dir.
func getImportDirAbs(importPath, dir string) (string, error) {
	bp, err := build.Import(importPath, dir, build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("finding %q on %q: %w", importPath, dir, err)
	}
	return bp.Dir, nil
}
