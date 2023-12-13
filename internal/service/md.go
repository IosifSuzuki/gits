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
	"go.uber.org/zap"
	"io"
)

type MD interface {
	RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error)
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

func (m *md) RenderMdToHTML(md []byte, attachmentIdentifiers map[string]string) ([]byte, error) {
	ext := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(ext)
	doc := p.Parse(md)
	doc = m.modifyAct(doc, attachmentIdentifiers)
	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: m.RenderHook,
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

func (m *md) RenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
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

func (m *md) htmlHighlight(w io.Writer, sourceCode, lang, defaultLang string) error {
	log := m.GetLogger()

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

	log.Debug("detect language in code", zap.String("lang", l.Config().Name))

	it, err := l.Tokenise(nil, sourceCode)
	if err != nil {
		return err
	}

	styleName := "xcode"
	highlightStyle := styles.Get(styleName)
	log.Debug("obtain highlight style", zap.String("style", highlightStyle.Name))
	return m.htmlCodeFormatter.Format(w, highlightStyle, it)
}
