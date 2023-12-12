package service

import (
	"gits/internal/model/html"
	"gits/internal/utils"
)

type Pagination struct {
	page       uint
	batchSize  uint
	countItems uint
}

func NewPagination(page uint, batchSize uint, countItems uint) Pagination {
	return Pagination{
		page:       page,
		batchSize:  batchSize,
		countItems: countItems,
	}
}

func (p *Pagination) Build(
	paginationPreviousItemFunc func() *html.PaginationItem,
	paginationItemFunc func(idx uint) *html.PaginationItem,
	paginationNextItemFunc func() *html.PaginationItem,
) *html.Pagination {
	const VisiblePagesPerSide uint = 2

	availablePages := p.countItems / p.batchSize
	if p.countItems%p.batchSize > 0 {
		availablePages += 1
	}

	paginationItems := make([]html.PaginationItem, 0, availablePages)
	beginPageIndex := uint(utils.Max(int(p.page)-int(VisiblePagesPerSide), 1))
	endPageIndex := utils.Min(p.page+VisiblePagesPerSide+1, availablePages)

	for i := beginPageIndex; i <= endPageIndex; i++ {
		if paginationItem := paginationItemFunc(i); paginationItem != nil {
			paginationItems = append(paginationItems, *paginationItem)
		}
	}

	nextPaginationItem := paginationNextItemFunc()
	nextPaginationItem.Active = p.page < availablePages

	previousPaginationItem := paginationPreviousItemFunc()
	previousPaginationItem.Active = p.page > 1

	return &html.Pagination{
		Page:         int(p.page),
		Batch:        int(p.batchSize),
		NextItem:     nextPaginationItem,
		Items:        paginationItems,
		PreviousItem: previousPaginationItem,
	}
}
