package central

import (
	"gits/internal/model/dto"
	"gits/internal/model/html"
	"go.uber.org/zap"
)

func (m *mainController) NewArticle(account *dto.Account) (*html.NewArticle, error) {
	log := m.GetLogger()

	storCategories, err := m.storageDAO.GetArticleRepository().AvailableCategories()
	if err != nil {
		log.Error("retrieve available categories has failed", zap.Error(err))
		return nil, err
	}

	categories := make([]dto.Category, 0, len(storCategories))
	for _, storCategory := range storCategories {
		category := dto.NewCategory(storCategory)
		categories = append(categories, *category)
	}

	return &html.NewArticle{
		PublisherName:       account.Username,
		AvailableCategories: categories,
	}, err
}
