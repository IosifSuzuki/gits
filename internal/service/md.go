package service

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gits/internal/container"
)

type MD interface {
	RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error)
}

type md struct {
	container.Container
}

func NewMD(container container.Container) MD {
	return &md{
		Container: container,
	}
}

func (m *md) RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error) {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md)
	doc = m.modifyAct(doc, attachmentIdentifiers)
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer), nil
}

func (m *md) modifyAct(doc ast.Node, attachmentIdentifiers map[string]string) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if imageNode, ok := node.(*ast.Image); ok && entering {
			destination := imageNode.Destination
			destinationString := string(destination)
			attachmentIdentifier, ok := attachmentIdentifiers[destinationString]
			if !ok {
				return ast.GoToNext
			}
			imageNode.Destination = []byte(attachmentIdentifier)
		}
		return ast.GoToNext
	})
	return doc
}
