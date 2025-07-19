package generator

import (
	"go/ast"
	"strings"
)

func extractCommentFromDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	var b strings.Builder
	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "ifacegen:generate") {
			continue
		}
		b.WriteString(comment.Text)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}

func extractComment(fn *ast.FuncDecl) string {
	if fn == nil || fn.Doc == nil {
		return ""
	}
	var b strings.Builder
	for _, comment := range fn.Doc.List {
		b.WriteString(comment.Text)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String())
}
