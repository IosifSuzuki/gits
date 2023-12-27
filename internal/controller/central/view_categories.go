package central

import (
	"fmt"
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/service"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) ViewCategories(page *dto.Page) (*html.Categories, error) {
	const BatchSize uint = 20

	log := m.GetLogger()

	pagination, err := m.prepareCategoriesPagination(page, BatchSize)
	if err != nil {
		log.Error("fail to prepare pagination for action", zap.Error(err))
		return nil, err
	}

	storCategories, err := m.storageDAO.GetArticleRepository().PaginationCategories(int(page.Page), int(BatchSize))
	if err != nil {
		log.Error("pagination categories retrieve has failed", zap.Error(err))
		return nil, err
	}

	categories := make([]html.Category, 0, len(storCategories))
	for _, storCategory := range storCategories {
		categories = append(categories, html.Category{
			Id:    int(storCategory.ID),
			Title: storCategory.Title,
		})
	}

	return &html.Categories{
		Categories: categories,
		Pagination: pagination,
	}, nil
}

func (m *mainController) prepareCategoriesPagination(page *dto.Page, batchSize uint) (*html.Pagination, error) {
	const Path = "/admin/view/categories"
	log := m.GetLogger()

	countCategories, err := m.storageDAO.GetArticleRepository().LenCategories(page.Page, batchSize)
	if err != nil {
		log.Error("obtain count of categories has failed", zap.Error(err))
		return nil, err
	}

	paginationBuilder := service.NewPagination(page.Page, batchSize, countCategories)

	paginationItemBuilderFunc := func(idx uint) *html.PaginationItem {
		path := fmt.Sprintf("%s%d", Path, idx)
		title := template.HTML(fmt.Sprintf("%d", idx))
		isActive := idx == page.Page

		return html.NewPaginationItem(path, title, isActive)
	}
	paginationItemNextBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("%s%d", Path, page.Page+1)
		title := template.HTML("&raquo")
		return html.NewPaginationItem(path, title, false)
	}
	paginationItemPreviousBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("%s%d", Path, page.Page-1)
		title := template.HTML("&laquo")
		return html.NewPaginationItem(path, title, false)
	}

	pagination := paginationBuilder.Build(
		paginationItemPreviousBuilderFunc,
		paginationItemBuilderFunc,
		paginationItemNextBuilderFunc,
	)

	return pagination, nil
}
