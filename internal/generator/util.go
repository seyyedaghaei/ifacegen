package generator

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"
)

func (d *packageData) typeString(t types.Type) string {
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name()
	case *types.Pointer:
		return "*" + d.typeString(tt.Elem())
	case *types.Slice:
		return "[]" + d.typeString(tt.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", tt.Len(), d.typeString(tt.Elem()))
	case *types.Map:
		return "map[" + d.typeString(tt.Key()) + "]" + d.typeString(tt.Elem())
	case *types.Chan:
		dir := ""
		switch tt.Dir() {
		case types.SendRecv:
			dir = "chan "
		case types.SendOnly:
			dir = "chan<- "
		case types.RecvOnly:
			dir = "<-chan "
		}
		return dir + d.typeString(tt.Elem())
	case *types.Named:
		obj := tt.Obj()
		pkg := obj.Pkg()
		if pkg != nil && pkg.Name() != "" && pkg.Name() != "main" {
			pkgPath := pkg.Path()

			if _, ok := d.pkgImps[pkgPath]; !ok {
				name := pkg.Name()
				d.pkgCounter[name]++
				counter := d.pkgCounter[name]
				d.pkgImps[pkgPath] = name
				imp := fmt.Sprintf(`"%s"`, pkgPath)
				if counter > 1 {
					d.pkgImps[pkgPath] += strconv.Itoa(counter)
					imp = d.pkgImps[pkgPath] + " " + imp
				}
				d.Imports = append(d.Imports, imp)
			}
			return d.pkgImps[pkgPath] + "." + obj.Name()
		}
		return obj.Name()
	case *types.Interface:
		return "any"
	case *types.Signature:
		return "func"
	default:
		return tt.String()
	}
}

func export(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
