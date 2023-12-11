package central

import (
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"gits/internal/model/html"
	"go.uber.org/zap"
)

func (m *mainController) NewCategory(account *dto.Account) (*html.NewCategory, error) {
	log := m.GetLogger()

	if account == nil {
		err := errs.NilError
		log.Error("account has nil value", zap.Error(err))
		return nil, err
	}

	return &html.NewCategory{
		PublisherName: account.FullName(),
	}, nil
}
