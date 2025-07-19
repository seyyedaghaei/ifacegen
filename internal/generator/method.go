package generator

import (
	"go/ast"
	"go/types"
	"strings"
)

func (d *packageData) collectMethods(files []*ast.File, info *types.Info, typ types.Type) []methodInfo {
	var methods []methodInfo

	ms := types.NewMethodSet(types.NewPointer(typ))
	for i := 0; i < ms.Len(); i++ {
		m := ms.At(i)
		obj := m.Obj()

		if !obj.Exported() {
			continue
		}

		fn, ok := obj.(*types.Func)
		if !ok {
			continue
		}

		if skipMethod(fn, files, info) {
			continue
		}

		sig, ok := obj.Type().(*types.Signature)
		if !ok {
			continue
		}

		fnDecl := findFuncDecl(fn, files, info)

		methods = append(methods, methodInfo{
			Comment:   extractComment(fnDecl),
			Signature: d.formatSignatureWithParamNames(obj.Name(), sig, fnDecl),
		})
	}

	return methods
}

func findFuncDecl(obj *types.Func, files []*ast.File, info *types.Info) *ast.FuncDecl {
	for _, f := range files {
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name == nil {
				continue
			}
			if info.Defs[fn.Name] == obj {
				return fn
			}
		}
	}
	return nil
}

func skipMethod(obj *types.Func, files []*ast.File, info *types.Info) bool {
	fnDecl := findFuncDecl(obj, files, info)
	if fnDecl == nil || fnDecl.Doc == nil {
		return false
	}
	for _, c := range fnDecl.Doc.List {
		if strings.Contains(c.Text, "ifacegen:skip") {
			return true
		}
	}
	return false
}
