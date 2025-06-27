package appinfo

import (
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
)

type PackageInfo interface {
	FileSet() *token.FileSet
	TypesInfo() *types.Info
	Root() *build.Package
	Packages() []*build.Package
	AllPackages() []*build.Package
	GetAstFiles(bp *build.Package) []*ast.File
	GetBuildPackage(path, dir string) *build.Package
}

type ApplicationInfo struct {
	PackageInfo
	entrypointName string

	// cache
	defs  map[*ast.Ident]ast.Node
	refs  map[ast.Node][]*ast.Ident
	used  map[ast.Node]bool
	pkgs  map[ast.Node]*build.Package
	decls map[*ast.Ident]ast.Node // XDecl or XSpec defining the key
}

func NewApplicationInfo(pi PackageInfo, entrypointName string) *ApplicationInfo {
	return &ApplicationInfo{
		PackageInfo:    pi,
		entrypointName: entrypointName,
	}
}

func (ai *ApplicationInfo) HasUsedC(bp *build.Package) bool {
	for _, f := range ai.GetAstFiles(bp) {
		var res bool
		ast.Inspect(f, func(node ast.Node) bool {
			switch node := node.(type) {
			case *ast.ImportSpec:
				if node.Path.Value == `"C"` && ai.IsUsed(node) {
					res = true
				}
				return false
			}
			return true
		})
		if res {
			return true
		}
	}
	return false
}

func (ai *ApplicationInfo) GetDefinition(id *ast.Ident) ast.Node {
	if ai.defs == nil {
		ai.defs, ai.refs = createDefsAndReferrings(ai.TypesInfo())
	}
	return ai.defs[id]
}

func (ai *ApplicationInfo) GetReferrings(nd ast.Node) []*ast.Ident {
	if ai.refs == nil {
		ai.defs, ai.refs = createDefsAndReferrings(ai.TypesInfo())
	}
	return ai.refs[nd]
}

// IsUsed() returned whether the given nd is used in call graph started from the entrypoint.
func (ai *ApplicationInfo) IsUsed(nd ast.Node) bool {
	if ai.used != nil {
		return ai.used[nd]
	}

	used := make(map[ast.Node]bool)

	var rec func(nd ast.Node)
	rec = func(nd ast.Node) {
		if nd == nil {
			return
		}

		if used[nd] {
			return
		}
		used[nd] = true

		switch nd := nd.(type) {
		case *ast.Ident:
			id := nd
			if bp := ai.GetPackage(nd); bp != nil && bp.Goroot {
				return
			}
			decl := ai.getDecl(id)
			if decl != nil {
				rec(decl)
			}
			def := ai.GetDefinition(id)
			if def != nil {
				rec(def)
			}
		case *ast.TypeSpec:
			ids := ai.GetMethods(nd.Name)
			for _, id := range ids {
				rec(id)
			}
			ast.Inspect(nd, func(nd ast.Node) bool {
				rec(nd)
				return true
			})
		default:
			ast.Inspect(nd, func(nd ast.Node) bool {
				rec(nd)
				return true
			})
		}
	}

	ep := ai.GetEntryPointDecl()
	if ep == nil {
		return false
	}
	rec(ep)

	// check initializers
	isVar := make(map[*ast.Ident]bool)
	var initDecls []*ast.FuncDecl
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			for _, decl := range f.Decls {
				switch decl := decl.(type) {
				case *ast.FuncDecl:
					if decl.Recv == nil && decl.Name.Name == "init" {
						initDecls = append(initDecls, decl)
					}
				case *ast.GenDecl:
					if decl.Tok != token.VAR {
						break
					}
					for _, spec := range decl.Specs {
						spec := spec.(*ast.ValueSpec)
						for _, name := range spec.Names {
							isVar[name] = true
						}
					}
				}
			}
		}
	}

	memo := make(map[ast.Node]bool)
	var rec2 func(node ast.Node) bool
	rec2 = func(node ast.Node) bool {
		if node == nil {
			return false
		}

		if res, ok := memo[node]; ok {
			return res
		}

		var res bool
		switch node := node.(type) {
		case *ast.Ident:
			id := node
			bp := ai.GetPackage(id)
			if bp != nil && bp.Goroot {
				break
			}
			if isVar[id] && used[id] {
				res = true
				break
			}
			decl := ai.getDecl(id)
			if decl != nil {
				res = rec2(decl)
				if res {
					break
				}
			}
			def := ai.GetDefinition(id)
			if def != nil {
				res = rec2(def)
				if res {
					break
				}
			}
		case *ast.FuncDecl:
			ast.Inspect(node.Body, func(node ast.Node) bool {
				if res {
					return false
				}
				res = rec2(node)
				return !res
			})
		}
		memo[node] = res
		return res
	}

	for ok := true; ok; {
		prev := len(used)
		for _, decl := range initDecls {
			if used[decl] {
				continue
			}
			if rec2(decl) {
				rec(decl)
			}
		}
		ok = prev != len(used)
	}

	ai.used = used
	return used[nd]
}

func (ai *ApplicationInfo) IsMethod(nd ast.Node) bool {
	id, ok := nd.(*ast.Ident)
	if !ok {
		return false
	}
	var res bool
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			ast.Inspect(f, func(node ast.Node) bool {
				switch decl := node.(type) {
				case *ast.GenDecl:
					return false
				case *ast.FuncDecl:
					if decl.Name == id {
						res = decl.Recv != nil
					}
					return false
				}
				return true
			})
		}
	}
	return res
}

func (ai *ApplicationInfo) IsInit(nd ast.Node) bool {
	id, ok := nd.(*ast.Ident)
	if !ok {
		return false
	}
	if id.Name != "init" {
		return false
	}
	var res bool
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			ast.Inspect(f, func(node ast.Node) bool {
				switch decl := node.(type) {
				case *ast.GenDecl:
					return false
				case *ast.FuncDecl:
					if decl.Name == id {
						res = decl.Name.Name == "init" && decl.Recv == nil
					}
					return false
				}
				return true
			})
		}
	}
	return res
}

func (ai *ApplicationInfo) GetEntryPointDecl() *ast.FuncDecl {
	var res *ast.FuncDecl
	for _, f := range ai.GetAstFiles(ai.Root()) {
		ast.Inspect(f, func(node ast.Node) bool {
			switch decl := node.(type) {
			case *ast.GenDecl:
				return false
			case *ast.FuncDecl:
				if decl.Name.Name == ai.entrypointName && decl.Recv == nil {
					res = decl
				}
				return false
			}
			return true
		})
	}
	return res
}

// Returns methods of the given id if the id is the Name (*ast.Ident) of the *ast.FuncDecl
func (ai *ApplicationInfo) GetMethods(id *ast.Ident) []*ast.Ident {
	obj, ok := ai.TypesInfo().Defs[id]
	if !ok {
		return nil
	}
	if obj == nil {
		return nil
	}

	for _, p := range ai.Packages() {
		var res []*ast.Ident
		for _, f := range ai.GetAstFiles(p) {
			ast.Inspect(f, func(node ast.Node) bool {
				switch decl := node.(type) {
				case *ast.GenDecl:
					return false
				case *ast.FuncDecl:
					if decl.Recv == nil {
						return false
					}
					var isObjsMethod bool
					ast.Inspect(decl.Recv.List[0].Type, func(node ast.Node) bool {
						tid, ok := node.(*ast.Ident)
						if !ok {
							return true
						}
						robj, ok := ai.TypesInfo().Uses[tid]
						if !ok {
							return false
						}
						if robj == nil {
							return false
						}
						isObjsMethod = isObjsMethod || obj == robj
						return false
					})
					if !isObjsMethod {
						return false
					}
					res = append(res, decl.Name)
					return false
				}
				return true
			})
		}
		if 0 < len(res) {
			return res
		}
	}

	return nil
}

func (ai *ApplicationInfo) GetFuncDecl(id *ast.Ident) *ast.FuncDecl {
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			var res *ast.FuncDecl
			ast.Inspect(f, func(node ast.Node) bool {
				decl, ok := node.(*ast.FuncDecl)
				if !ok {
					if _, ok := node.(*ast.GenDecl); ok {
						return false
					}
					return true
				}
				fid := decl.Name
				if id == fid {
					res = decl
					return false
				}
				return false
			})
			if res != nil {
				return res
			}
		}
	}
	return nil
}

func (ai *ApplicationInfo) GetFile(nd ast.Node) *ast.File {
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			var res *ast.File
			ast.Inspect(f, func(node ast.Node) bool {
				if nd == node {
					res = f
					return false
				}
				return true
			})
			if res != nil {
				return res
			}
		}
	}
	return nil
}

func (ai *ApplicationInfo) GetPackage(nd ast.Node) *build.Package {
	if ai.pkgs != nil {
		return ai.pkgs[nd]
	}
	pkgs := make(map[ast.Node]*build.Package)
	for _, p := range ai.AllPackages() {
		for _, f := range ai.GetAstFiles(p) {
			ast.Inspect(f, func(node ast.Node) bool {
				pkgs[node] = p
				return true
			})
		}
	}
	ai.pkgs = pkgs
	return pkgs[nd]
}

func (ai *ApplicationInfo) Squash() (*SquashedApp, error) {
	return newSquashedApp(ai)
}

// getDecl returns thedeclaration node (XxxSpec, FuncDecl, ...) for the given identity or nil
func (ai *ApplicationInfo) getDecl(id *ast.Ident) ast.Node {
	if ai.decls != nil {
		return ai.decls[id]
	}

	decls := make(map[*ast.Ident]ast.Node)
	for _, p := range ai.Packages() {
		for _, f := range ai.GetAstFiles(p) {
			for _, decl := range f.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					switch decl.Tok {
					case token.IMPORT:
						for _, spec := range decl.Specs {
							spec := spec.(*ast.ImportSpec)
							if spec.Name != nil {
								decls[spec.Name] = spec
							}
						}
					case token.CONST, token.VAR:
						for _, spec := range decl.Specs {
							spec := spec.(*ast.ValueSpec)
							for _, name := range spec.Names {
								decls[name] = spec
							}
						}
					case token.TYPE:
						for _, spec := range decl.Specs {
							spec := spec.(*ast.TypeSpec)
							decls[spec.Name] = spec
						}
					}
				case *ast.FuncDecl:
					decls[decl.Name] = decl
				}
			}
		}
	}

	ai.decls = decls
	return decls[id]
}

// It creates and returns 2 maps. The first one is defs; map from *ast.Ident/*ast.SelectorExpr to *ast.ImportSpec/*ast.Ident.
// The second one is refs: map from *ast.ImportSpec/*ast.Ident => Slice of *ast.Ident/*ast.SelectorExpr
func createDefsAndReferrings(tinfo *types.Info) (map[*ast.Ident]ast.Node, map[ast.Node][]*ast.Ident) {
	obj2def := make(map[types.Object]ast.Node)
	for id, obj := range tinfo.Defs {
		obj2def[obj] = id
	}
	for nd, obj := range tinfo.Implicits {
		obj2def[obj] = nd
	}

	defs := make(map[*ast.Ident]ast.Node)
	refs := make(map[ast.Node][]*ast.Ident)
	for id, obj := range tinfo.Uses {
		def, ok := obj2def[obj]
		if !ok || def == nil {
			continue
		}
		defs[id] = def
		refs[def] = append(refs[def], id)
	}
	for expr, sel := range tinfo.Selections {
		if sel == nil || sel.Obj() == nil {
			continue
		}
		def, ok := obj2def[sel.Obj()]
		if !ok || def == nil {
			continue
		}
		defs[expr.Sel] = def
		refs[def] = append(refs[def], expr.Sel)
	}
	return defs, refs
}
