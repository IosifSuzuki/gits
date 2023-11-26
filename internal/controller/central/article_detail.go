package central

import (
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"gits/internal/model/html"
	"go.uber.org/zap"
	"html/template"
)

func (m *mainController) Article(articleIdentifier *dto.ArticleIdentifier) (*html.Article, error) {
	log := m.GetLogger()

	storArticle, err := m.storageDAO.GetArticleRepository().Article(articleIdentifier.ID)
	if err != nil {
		log.Error("fail to retrieve article by id", zap.Error(err))

		return nil, err
	}

	storAccount, err := m.storageDAO.GetAccountRepository().AccountByIdentifier(storArticle.PublisherId)
	if err != nil {
		log.Error("fail to retrieve publisher from article", zap.Error(err))

		return nil, err
	} else if storAccount == nil {
		err = errs.NilError
		log.Error("publisher contains nil value", zap.Error(err))

		return nil, err
	}

	categories := make([]dto.Category, 0, len(storArticle.Categories))
	for _, storCategory := range storArticle.Categories {
		category := dto.NewCategory(storCategory)
		categories = append(categories, *category)
	}

	account := dto.NewAccount(*storAccount)
	htmlArticle := html.Article{
		Author:      account,
		Title:       storArticle.Title,
		Location:    storArticle.Location,
		Date:        &storArticle.UpdatedAt,
		ReadingTime: storArticle.ReadingTime,
		Content:     template.HTML(*storArticle.Content),
		Categories:  categories,
	}
	return &htmlArticle, nil
}
