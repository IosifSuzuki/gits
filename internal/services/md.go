package services

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"strings"
)

type MD interface {
	RenderMdToHTML(md []byte) ([]byte, error)
}

type md struct {
}

func NewMD() MD {
	return &md{}
}

func (m *md) RenderMdToHTML(md []byte) ([]byte, error) {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md)
	doc = modifyAct(doc)
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer), nil
}

func modifyAct(doc ast.Node) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if imageNode, ok := node.(*ast.Image); ok && entering {
			destination := imageNode.Destination
			var destinationBuilder strings.Builder
			_, err := destinationBuilder.Write([]byte("/assets/articles/middleware/"))
			if err != nil {
				return ast.GoToNext
			}
			_, err = destinationBuilder.Write([]byte(destination))
			if err != nil {
				return ast.GoToNext
			}
			imageNode.Destination = []byte(destinationBuilder.String())
		}
		return ast.GoToNext
	})
	return doc
}
