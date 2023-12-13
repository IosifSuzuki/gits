package central

import (
	"fmt"
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"gits/internal/service"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) Articles(page *dto.Page) (*html.Articles, error) {
	const BatchSize uint = 10

	log := m.GetLogger()

	pagination, err := m.prepareArticlesPagination(page, BatchSize)
	if err != nil {
		log.Error("fail to prepare pagination for articles", zap.Error(err))
		return nil, err
	}

	storArticles, err := m.storageDAO.GetArticleRepository().Articles(int(page.Page), int(BatchSize))
	if err != nil {
		log.Error("retrieve articles from storage by pagination has failed", zap.Error(err))
		return nil, err
	}

	previewArticles := make([]html.PreviewArticle, 0, len(storArticles))
	for _, storArticle := range storArticles {
		htmlContent, err := m.convertMdToHtmlPreview(&storArticle, 70)
		if err != nil {
			log.Error("convert mark down to html has failed", zap.Error(err))

			return nil, err
		}

		contentHTML := template.HTML(*htmlContent)
		date := storArticle.UpdatedAt
		previewArticles = append(previewArticles, html.PreviewArticle{
			Id:      int(storArticle.ID),
			Title:   storArticle.Title,
			Date:    &date,
			Content: &contentHTML,
		})
	}

	return &html.Articles{
		Articles:   previewArticles,
		Pagination: pagination,
	}, nil
}

func (m *mainController) prepareArticlesPagination(page *dto.Page, batchSize uint) (*html.Pagination, error) {
	log := m.GetLogger()

	countArticles, err := m.storageDAO.GetArticleRepository().LenArticles()
	if err != nil {
		log.Error("obtain count of articles has failed", zap.Error(err))
		return nil, err
	}

	paginationBuilder := service.NewPagination(page.Page, batchSize, countArticles)

	paginationItemPreviousBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("/articles/%d", page.Page-1)
		title := template.HTML("&laquo")
		return html.NewPaginationItem(path, title, false)
	}
	paginationItemBuilderFunc := func(idx uint) *html.PaginationItem {
		path := fmt.Sprintf("/articles/%d", idx)
		title := template.HTML(fmt.Sprintf("%d", idx))
		isActive := idx == page.Page

		return html.NewPaginationItem(path, title, isActive)
	}
	paginationItemNextBuilderFunc := func() *html.PaginationItem {
		path := fmt.Sprintf("/articles/%d", page.Page+1)
		title := template.HTML("&raquo")
		return html.NewPaginationItem(path, title, false)
	}

	pagination := paginationBuilder.Build(
		paginationItemPreviousBuilderFunc,
		paginationItemBuilderFunc,
		paginationItemNextBuilderFunc,
	)

	return pagination, nil
}
