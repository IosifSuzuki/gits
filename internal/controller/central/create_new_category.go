package central

import (
	"gits/internal/model/dto"
	stor "gits/internal/model/storage"
	"go.uber.org/zap"
)

func (m *mainController) CreateNewCategory(account *dto.Account, form *dto.FormCategory) error {
	log := m.GetLogger()

	var newCategory = stor.Category{
		PublisherId: account.ID,
		Title:       form.Title,
	}

	err := m.storageDAO.GetArticleRepository().CreateNewCategory(&newCategory)
	if err != nil {
		log.Error("create new category has failed", zap.Error(err))
		return err
	}

	return nil
}
