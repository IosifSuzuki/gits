package html

import (
	"html/template"
	"time"
)

type PreviewArticle struct {
	Id      int
	Title   *string
	Date    *time.Time
	Content *template.HTML
}
