package service

import (
	"github.com/alecthomas/chroma"
	chromaHtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"gits/internal/container"
	"gits/internal/utils"
	"go.uber.org/zap"
	"io"
	"strings"
)

type MD interface {
	RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error)
	RenderMdToPreviewHTML(md []byte, words uint) ([]byte, error)
}

type md struct {
	container.Container
	htmlCodeFormatter *chromaHtml.Formatter
}

func NewMD(container container.Container) MD {
	htmlFormatter := chromaHtml.New(
		chromaHtml.WithLineNumbers(true),
		chromaHtml.TabWidth(2),
	)
	return &md{
		Container:         container,
		htmlCodeFormatter: htmlFormatter,
	}
}

func (m *md) RenderMdToPreviewHTML(md []byte, words uint) ([]byte, error) {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md)
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: m.PreparePreviewHook(words),
	}
	renderer := html.NewRenderer(opts)
	return markdown.Render(doc, renderer), nil
}

func (m *md) RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error) {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md)
	doc = m.modifyAct(doc, attachmentIdentifiers)
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: m.RenderContentHook,
	}
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

func (m *md) RenderContentHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	log := m.GetLogger()

	if codeBlock, ok := node.(*ast.CodeBlock); ok {
		defaultLang := ""
		lang := string(codeBlock.Info)
		sourceCode := string(codeBlock.Literal)

		if err := m.htmlHighlight(w, sourceCode, lang, defaultLang); err != nil {
			log.Error("render source code has failed", zap.Error(err))
			return ast.GoToNext, false
		}

		return ast.GoToNext, true
	}

	return ast.GoToNext, false
}

func (m *md) PreparePreviewHook(words uint) html.RenderNodeFunc {
	var currentWords uint = 0

	return func(_ io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		switch node := node.(type) {
		case *ast.Heading:
			return ast.SkipChildren, true
		case *ast.CodeBlock:
			return ast.SkipChildren, true
		case *ast.Image:
			return ast.SkipChildren, true
		case *ast.Link:
			return ast.SkipChildren, true
		case *ast.List:
			return ast.SkipChildren, true
		case *ast.Table:
			return ast.SkipChildren, true
		case *ast.Paragraph:
			if !utils.MarkdownHasTextNode(node) {
				return ast.SkipChildren, true
			} else if !entering {
				return ast.GoToNext, false
			} else if currentWords > words {
				return ast.Terminate, true
			}

			text := utils.MarkdownTextContent(node)
			words := strings.Split(text, " ")
			currentWords += uint(len(words))

			return ast.GoToNext, false
		}

		return ast.GoToNext, false
	}
}

func (m *md) htmlHighlight(w io.Writer, sourceCode, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(sourceCode)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, sourceCode)
	if err != nil {
		return err
	}

	styleName := "xcode"
	highlightStyle := styles.Get(styleName)
	return m.htmlCodeFormatter.Format(w, highlightStyle, it)
}
