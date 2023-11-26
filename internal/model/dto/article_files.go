package dto

import (
	"io"
)

type ArticleFiles struct {
	MDFile      io.Reader
	Attachments map[string]io.Reader
}

func NewArticleFiles(mdFile io.Reader, attachments map[string]io.Reader) *ArticleFiles {
	return &ArticleFiles{
		MDFile:      mdFile,
		Attachments: attachments,
	}
}
