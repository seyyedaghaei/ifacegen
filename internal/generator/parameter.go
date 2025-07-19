package generator

import (
	"go/ast"
	"go/types"
	"strings"
)

func (d *packageData) formatSignatureWithParamNames(name string, sig *types.Signature, fnDecl *ast.FuncDecl) string {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString("(")

	params := sig.Params()
	paramNames := extractParamNames(fnDecl)

	for i := 0; i < params.Len(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}

		pname := "_"
		if i < len(paramNames) && paramNames[i] != "" {
			pname = paramNames[i]
		}

		b.WriteString(pname)
		b.WriteString(" ")
		b.WriteString(d.typeString(params.At(i).Type()))
	}

	b.WriteString(")")

	results := sig.Results()
	switch results.Len() {
	case 0:
	case 1:
		b.WriteString(" ")
		b.WriteString(d.typeString(results.At(0).Type()))
	default:
		b.WriteString(" (")
		for i := 0; i < results.Len(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(d.typeString(results.At(i).Type()))
		}
		b.WriteString(")")
	}

	return b.String()
}

func extractParamNames(fnDecl *ast.FuncDecl) []string {
	if fnDecl == nil || fnDecl.Type == nil || fnDecl.Type.Params == nil {
		return nil
	}

	var names []string
	for _, field := range fnDecl.Type.Params.List {
		if len(field.Names) == 0 {
			names = append(names, "")
		} else {
			for _, nameIdent := range field.Names {
				names = append(names, nameIdent.Name)
			}
		}
	}
	return names
}
