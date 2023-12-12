package html

import "html/template"

type PaginationItem struct {
	Path   string
	Title  template.HTML
	Active bool
}

func NewPaginationItem(path string, title template.HTML, active bool) *PaginationItem {
	return &PaginationItem{
		Path:   path,
		Title:  title,
		Active: active,
	}
}
