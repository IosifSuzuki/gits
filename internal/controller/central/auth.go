package central

import (
	"context"
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"go.uber.org/zap"
)

func (m *mainController) Auth(authentication *dto.Authentication) (*dto.AccountSession, error) {
	log := m.GetLogger()
	ctx := context.Background()

	storAccount, err := m.storageDAO.
		GetAccountRepository().
		AccountByCredential(&authentication.Username, &authentication.Password)
	if err != nil {
		log.Error("retrieve account by credential has failed", zap.Error(err))
		return nil, err
	} else if storAccount == nil {
		err = errs.NilError
		log.Error("get nil instead object", zap.Error(err))
		return nil, err
	}

	account := dto.NewAccount(*storAccount)

	accountSession, err := m.AccountSession.CreateAccountSession(ctx, account)
	if err != nil {
		log.Error("create session has failed", zap.Error(err))
		return nil, err
	}

	return accountSession, nil
}
