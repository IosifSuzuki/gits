package central

import (
	"fmt"
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/service"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) ViewArticles(page *dto.Page) (*html.Articles, error) {
	const BatchSize uint = 20

	log := m.GetLogger()

	pagination, err := m.prepareViewArticlesPagination(page, BatchSize)
	if err != nil {
		log.Error("fail to prepare pagination for view articles", zap.Error(err))
		return nil, err
	}

	storArticles, err := m.storageDAO.GetArticleRepository().Articles(int(page.Page), int(BatchSize))
	if err != nil {
		log.Error("pagination categories retrieve has failed", zap.Error(err))
		return nil, err
	}

	previewArticles := make([]html.PreviewArticle, 0, len(storArticles))

	for _, storArticle := range storArticles {
		previewArticles = append(previewArticles, html.PreviewArticle{
			Id:      int(storArticle.ID),
			Title:   storArticle.Title,
			Date:    nil,
			Content: nil,
		})
	}
	return &html.Articles{
		Articles:   previewArticles,
		Pagination: pagination,
	}, nil
}

func (m *mainController) prepareViewArticlesPagination(page *dto.Page, batchSize uint) (*html.Pagination, error) {
	const Path = "/admin/view/articles/"
	log := m.GetLogger()

	countArticles, err := m.storageDAO.GetArticleRepository().LenArticles()
	if err != nil {
		log.Error("obtain count of categories has failed", zap.Error(err))
		return nil, err
	}

	paginationBuilder := service.NewPagination(page.Page, batchSize, countArticles)

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
