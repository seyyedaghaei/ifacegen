package generator

import (
	"go/ast"
	"path"
	"strings"
)

// shouldIncludeStruct checks if a struct should be included for interface generation
func (d *packageData) shouldIncludeStruct(name string, matches []string, doc *ast.CommentGroup) bool {
	if doc != nil {
		for _, c := range doc.List {
			if strings.Contains(c.Text, "ifacegen:skip") {
				return false
			}

			if strings.Contains(c.Text, "ifacegen:generate") {
				return true
			}
		}
	}
	for _, s := range matches {
		if matched, _ := path.Match(s, name); matched {
			return true
		}
	}
	return false
}
