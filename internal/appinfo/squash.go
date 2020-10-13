package appinfo

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

type appInfo interface {
	PackageInfo
	HasUsedC(pb *build.Package) bool
	IsUsed(ast.Node) bool
	IsMethod(ast.Node) bool
	IsInit(ast.Node) bool
	GetEntryPointDecl() *ast.FuncDecl
	GetMethods(*ast.Ident) []*ast.Ident
	GetFuncDecl(*ast.Ident) *ast.FuncDecl
	GetFile(nd ast.Node) *ast.File
	GetPackage(nd ast.Node) *build.Package
	GetReferrings(nd ast.Node) []*ast.Ident
}

// SquashedApp represents an application combined into a single file
type SquashedApp struct {
	// pkgName is the name of package that has the entry point of the application
	pkgName string

	// xDecls are slices of *ast.GenDecl/*ast.FuncDecl of the combined application.
	// You can use it for printing out the application as a single file.
	// Note that each node in it is derived from various postions of files including token.NoPos.
	importDecls []ast.Decl
	otherDecls  []ast.Decl
	mainDecls   []ast.Decl

	// Fset keeps FileSet for the Syntax
	fset *token.FileSet
}

// newSquashedApp build SquashedApp from appInfo
func newSquashedApp(ai appInfo) (*SquashedApp, error) {
	// collect used items for the SquashedApp
	ingr := &ingredients{}
	fset := ai.FileSet()

	for _, bp := range ai.Packages() {
		if ai.HasUsedC(bp) {
			for _, f := range bp.CFiles {
				fname := filepath.Join(bp.ImportPath, f)
				fpath := filepath.Join(bp.Dir, f)
				decl, err := createImportDeclFromCFile(fset, fname, fpath)
				if err != nil {
					return nil, fmt.Errorf("creating func decls for C files: %w", err)
				}
				ingr.importCDecls = append(ingr.importCDecls, decl)
			}

		}
	}

	forEachFile(ai, func(bp *build.Package, f *ast.File) {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				switch decl.Tok {
				case token.IMPORT:
					specs, cdecls, cspecs := collectUsedImport(ai, decl)
					ingr.importSpecs = append(ingr.importSpecs, specs...)
					ingr.importCDecls = append(ingr.importCDecls, cdecls...)
					ingr.importCSpecs = append(ingr.importCSpecs, cspecs...)
				case token.CONST:
					if hasUsedValueSpec(ai, decl.Specs) {
						ingr.constDecls = append(ingr.constDecls, decl)
					}
				case token.TYPE:
					if hasUsedTypeSpec(ai, decl) {
						ingr.typeDecls = append(ingr.typeDecls, decl)
					}
				case token.VAR:
					if hasUsedValueSpec(ai, decl.Specs) {
						ingr.varDecls = append(ingr.varDecls, decl)
					}
				}
			case *ast.FuncDecl:
				id := decl.Name
				if !ai.IsUsed(id) {
					break
				}
				if ai.IsInit(id) {
					ingr.initDecls = append(ingr.initDecls, decl)
					break
				}
				if ai.GetEntryPointDecl() == decl {
					break
				}
				ingr.funcDecls = append(ingr.funcDecls, decl)
			}
		}
	})
	ingr.mainDecl = ai.GetEntryPointDecl()

	return ingr.newSquashedApp(ai), nil

}

// Fprint formats the source code and writes it to the given io.Writer
func (sa *SquashedApp) Fprint(w io.Writer) error {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "package %s\n\n", sa.pkgName)
	pcfg := printer.Config{
		Mode:     printer.TabIndent | printer.UseSpaces,
		Tabwidth: 8,
		Indent:   0,
	}

	for _, decl := range sa.importDecls {
		pcfg.Fprint(buf, sa.fset, decl)
		fmt.Fprintf(buf, "\n\n")
	}

	// The printer.SourcePos should be used but some Cgo comments got borken by them.
	// So don't print source pos during import Decls.
	pcfg.Mode |= printer.SourcePos

	if 0 < len(sa.otherDecls) {
		fmt.Fprintf(buf, "// =============================================================================\n")
		fmt.Fprintf(buf, "// Populated Libiraries\n")
		fmt.Fprintf(buf, "// =============================================================================\n\n")
	}
	fprintDecls(pcfg, buf, sa.fset, sa.otherDecls)

	if 0 < len(sa.otherDecls) {
		fmt.Fprintf(buf, "// =============================================================================\n")
		fmt.Fprintf(buf, "// Original Main Package\n")
		fmt.Fprintf(buf, "// =============================================================================\n\n")
	}
	fprintDecls(pcfg, buf, sa.fset, sa.mainDecls)

	b, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("formatting: %w", err)
	}
	w.Write(b)
	return nil
}

func fprintDecls(pcfg printer.Config, w io.Writer, fset *token.FileSet, decls []ast.Decl) {
	for _, d := range decls {
		pcfg.Fprint(w, fset, d)
		fmt.Fprintf(w, "\n\n")
	}
}

// ingredients for SquashedApp
type ingredients struct {
	importCDecls []*ast.GenDecl    // In `import "C"` notation, the C source is associated to GenDecl, not to ImportSpec, so the whole GenDecl is needed
	importCSpecs []*ast.ImportSpec // Not likely but `import ( ... /* <C Source> */\n<white spaces>"C" ...` format is allowed and this field is for that case
	importSpecs  []*ast.ImportSpec // For normal ImportSpecs. they can be packed into a single import ( ... ) notation

	constDecls []*ast.GenDecl // for used const GenDecls

	typeDecls []*ast.GenDecl // for used type GenDecls

	varDecls []*ast.GenDecl // for used var GenDecls

	initDecls []*ast.FuncDecl // init should be above for readability
	funcDecls []*ast.FuncDecl // normal functions

	mainDecl *ast.FuncDecl // last function shoud be main()
}

func (ingr *ingredients) squashRest(ai appInfo, used map[string]bool, filter func(node ast.Node) bool) []ast.Decl {
	var res []ast.Decl
	for _, d := range ingr.constDecls {
		if !filter(d) {
			continue
		}
		renameGenDecl(ai, used, d)
		res = append(res, d)
	}

	for _, d := range ingr.typeDecls {
		if !filter(d) {
			continue
		}
		renameGenDecl(ai, used, d)
		res = append(res, d)
	}

	for _, d := range ingr.varDecls {
		if !filter(d) {
			continue
		}
		renameGenDecl(ai, used, d)
		res = append(res, d)
	}

	for _, d := range ingr.initDecls {
		if !filter(d) {
			continue
		}
		// init() doesn't need renaming
		res = append(res, d)
	}

	for _, d := range ingr.funcDecls {
		if !filter(d) {
			continue
		}
		// Methods doesn't need renaming because they are type scope
		if d.Recv == nil {
			renameFuncDecl(ai, used, d)
		}
		res = append(res, d)
	}
	return res
}

func (ingr *ingredients) squashImports(ai appInfo, used map[string]bool) []ast.Decl {
	var res []ast.Decl

	// `import "C"` and `import ( ... )`
	for _, d := range ingr.importCDecls {
		res = append(res, d)
	}

	idecl := &ast.GenDecl{
		TokPos: token.NoPos,
		Tok:    token.IMPORT,
		Lparen: token.NoPos,
		Rparen: token.NoPos,
	}
	for _, s := range ingr.importCSpecs {
		idecl.Specs = append(idecl.Specs, s)
	}

	ispecs := squashImportSpecs(ai, used, ingr.importSpecs)

	for _, s := range ispecs {
		idecl.Specs = append(idecl.Specs, s)
	}
	if 0 < len(idecl.Specs) {
		res = append(res, idecl)
	}
	return res
}

func (ingr *ingredients) newUsedNames(ai appInfo) map[string]bool {
	// used identity in the target file
	used := make(map[string]bool)
	// collect names in scopes of each function body
	decls := make([]*ast.FuncDecl, 0, len(ingr.initDecls)+len(ingr.funcDecls)+1)
	for _, ds := range [][]*ast.FuncDecl{ingr.initDecls, ingr.funcDecls} {
		for _, d := range ds {
			decls = append(decls, d)
		}
	}
	decls = append(decls, ingr.mainDecl)
	for _, d := range decls {
		ast.Inspect(d, func(node ast.Node) bool {
			id, ok := node.(*ast.Ident)
			if !ok {
				return true
			}
			if d.Name == id {
				// the function name has a object in Defs but it is package scope
				return false
			}
			obj := ai.TypesInfo().Defs[id]
			if obj == nil {
				return false
			}
			// obj != nil means this identity is function scope
			used[id.Name] = true
			return false
		})
	}

	// entry point should not be renamed, so register its name here
	used[ingr.mainDecl.Name.Name] = true
	return used
}

// newSquashedApp populates used items to a single *ast.Node deduping and renameing if needed.
// Note that the oriiginal ast.Nodes are modified so they are no longer used for rebuilding the original source
func (ingr *ingredients) newSquashedApp(ai appInfo) *SquashedApp {
	mainPkg := ai.Root()

	res := &SquashedApp{
		pkgName: mainPkg.Name,
		fset:    ai.FileSet(),
	}

	// memo for used identity in the target file
	used := ingr.newUsedNames(ai)

	res.importDecls = ingr.squashImports(ai, used)

	pred := func(node ast.Node) bool {
		return ai.GetPackage(node) == mainPkg
	}
	mainDecls := ingr.squashRest(ai, used, pred)
	otherDecls := ingr.squashRest(ai, used, func(node ast.Node) bool { return !pred(node) })
	mainDecls = append(mainDecls, ingr.mainDecl)

	res.mainDecls = removeInvalidSelector(mainDecls)
	res.otherDecls = removeInvalidSelector(otherDecls)

	return res
}

// *ast.SelectorExpr `lib.Foo()` in the original source is renamed to invalid `.Foo()`
// by squashImportSpecs() if the `lib` is a third party, non-standard, library.
// This function replaces the invalid SelectorExpr to valid *ast.Ident like `Foo()`
func removeInvalidSelector(decls []ast.Decl) []ast.Decl {
	var res []ast.Decl
	fn := func(c *astutil.Cursor) bool {
		if sel, ok := c.Node().(*ast.SelectorExpr); ok {
			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "" {
				c.Replace(sel.Sel)
			}
			return false
		}
		return true
	}
	for _, d := range decls {
		res = append(res, astutil.Apply(d, fn, nil).(ast.Decl))
	}
	return res
}

// rename function name if needed.
func renameFuncDecl(ai appInfo, used map[string]bool, decl *ast.FuncDecl) {
	if !used[decl.Name.Name] {
		used[decl.Name.Name] = true
		return
	}
	bp := ai.GetPackage(decl)
	prefix := bp.Name
	renameIdents(ai, used, prefix, []*ast.Ident{decl.Name})
}

// rename type, const, var name if needed.
// It scans all specs in the decl and rename all ident when one of them need
// renaming.  It is done for readability.
func renameGenDecl(ai appInfo, used map[string]bool, decl *ast.GenDecl) {
	var ids []*ast.Ident
	var needRename bool
	for _, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.TypeSpec:
			id := spec.Name
			ids = append(ids, id)
			needRename = needRename || used[id.Name]
		case *ast.ValueSpec:
			for _, id := range spec.Names {
				ids = append(ids, id)
				needRename = needRename || used[id.Name]
			}
		}
	}

	if needRename {
		bp := ai.GetPackage(decl)
		prefix := bp.Name
		renameIdents(ai, used, prefix, ids)
	} else {
		for _, id := range ids {
			used[id.Name] = true
		}
	}
}

// Rename a set of identities specified by the given ids adding the same name prefix
func renameIdents(ai appInfo, used map[string]bool, prefix string, ids []*ast.Ident) {
	// Find the safe prefix for the identifiers.
	// Initially the prefx candidate is package name, e.g. "foo".
	// When the candidate is not safe, add 'x' to the prefix, e.g. "xfoo"
	tns := make([]string, len(ids))

	for needRename := true; needRename; {
		needRename = false
		for i, id := range ids {
			tns[i] = fmt.Sprintf("%s_%s", prefix, id.Name)
			if used[tns[i]] {
				needRename = true
				prefix = "x" + prefix
				break
			}
		}
	}

	for i, id := range ids {
		id.Name = tns[i]
		renameTo(ai, id, tns[i])
		used[tns[i]] = true
	}
}

func hasUsedValueSpec(ai appInfo, specs []ast.Spec) bool {
	for _, spec := range specs {
		spec := spec.(*ast.ValueSpec)
		for _, id := range spec.Names {
			if ai.IsUsed(id) {
				return true
			}
		}
	}
	return false
}

func hasUsedTypeSpec(ai appInfo, decl *ast.GenDecl) bool {
	if ai.IsUsed(decl) {
		return true
	}
	for _, spec := range decl.Specs {
		spec := spec.(*ast.TypeSpec)
		if ai.IsUsed(spec) {
			return true
		}
	}
	return false
}

func collectUsedImportC(ai appInfo, ingr *ingredients, decl *ast.GenDecl) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		if spec.Path.Value != `"C"` {
			continue
		}
		if len(decl.Specs) == 1 {
			ingr.importCDecls = append(ingr.importCDecls, decl)
		} else {
			ingr.importCSpecs = append(ingr.importCSpecs, spec)
		}
	}
}

func collectUsedImport(ai appInfo, decl *ast.GenDecl) (specs []*ast.ImportSpec, cdecls []*ast.GenDecl, cspecs []*ast.ImportSpec) {
	for _, spec := range decl.Specs {
		spec, ok := spec.(*ast.ImportSpec)
		if !ok {
			continue
		}
		var isUsed bool
		ast.Inspect(spec, func(node ast.Node) bool {
			if isUsed {
				return false
			}
			if ai.IsUsed(node) {
				isUsed = true
				return false
			}
			return true
		})
		if !isUsed {
			continue
		}
		if spec.Path.Value == `"C"` {
			if len(decl.Specs) == 1 {
				cdecls = append(cdecls, decl)
			} else {
				cspecs = append(cspecs, spec)
			}
			continue
		}
		specs = append(specs, spec)
	}
	return
}

func createImportDeclFromCFile(fset *token.FileSet, name, path string) (*ast.GenDecl, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}
	format := `package main
/*
%s
*/
import "C"
`
	src := fmt.Sprintf(format, string(b))
	decl, err := parser.ParseFile(fset, name, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}
	return decl.Decls[0].(*ast.GenDecl), nil
}

func squashImportSpecs(ai appInfo, used map[string]bool, specs []*ast.ImportSpec) []*ast.ImportSpec {
	collected := make(map[string]bool)
	var res []*ast.ImportSpec
	for i, spec := range specs {
		path := strings.Trim(spec.Path.Value, `"`)
		f := ai.GetFile(spec)
		bp := ai.GetBuildPackage(path, f.Name.Name)
		if bp == nil {
			panic(fmt.Sprintf("unknown package: %s", path))
		}

		// non-standard packages will be embedded into the target source file
		if !bp.Goroot {
			// *ast.SelectorExpr using this empty name will be replaced by its .Sel.
			if spec.Name != nil {
				// ai.GetReferrings keeps referring identities for spec.Name if it isn't nil
				renameTo(ai, spec.Name, "")
			} else {
				renameTo(ai, spec, "")
			}
			continue
		}

		if collected[path] {
			continue
		}
		collected[path] = true

		sames := []*ast.ImportSpec{spec}
		candidate := bp.Name

		// find same imports
		for j := i; j < len(specs); j++ {
			jpath := strings.Trim(specs[j].Path.Value, `"`)
			if path != jpath {
				continue
			}
			sames = append(sames, specs[j])
		}

		name := candidate
		for used[name] {
			name = "x" + name
		}
		used[name] = true

		if spec.Name == nil {
			obj := ai.TypesInfo().Implicits[spec]
			if obj.Name() != name {
				spec.Name = &ast.Ident{
					NamePos: token.NoPos,
					Name:    name,
				}
			}
		} else if spec.Name.Name != name {
			spec.Name.Name = name
		}
		res = append(res, spec)

		for _, sp := range sames {
			renameTo(ai, sp, name)
		}
	}
	return res
}

// rename the name of identities referring the given node.
func renameTo(ai appInfo, node ast.Node, to string) {
	for _, uid := range ai.GetReferrings(node) {
		if uid.Name == to {
			continue
		}
		uid.Name = to
	}
}

// call f() for each *ast.File in all packages except standard ones.
func forEachFile(ai PackageInfo, fn func(bp *build.Package, f *ast.File)) {
	for _, bp := range ai.Packages() {
		for _, f := range ai.GetAstFiles(bp) {
			fn(bp, f)
		}
	}
}
