package storage

import (
	"encoding/base64"
	"gits/internal/container"
	"gits/internal/model/errs"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"gits/internal/utils"
	"go.uber.org/zap"
)

type AccountRepository interface {
	AccountByCredential(username *string, password *string) (*stor.Account, error)
	AccountByIdentifier(id int) (*stor.Account, error)
}

type accountRepository struct {
	container.Container
	storageProvider provider.Storage
}

func NewAccount(container container.Container, storageProvider provider.Storage) AccountRepository {
	return &accountRepository{
		Container:       container,
		storageProvider: storageProvider,
	}
}

func (a *accountRepository) AccountByCredential(username *string, password *string) (*stor.Account, error) {
	conn := a.storageProvider.GetConnection()
	log := a.GetLogger()

	if username == nil || password == nil {
		err := errs.NilError
		log.Error("nil error", zap.Error(err))
		return nil, err
	}

	var account stor.Account
	if err := conn.Where("username = ?", username).First(&account).Error; err != nil {
		log.Error("retrieve account has failed")
		return nil, err
	}

	decodedHashPassword, err := base64.StdEncoding.DecodeString(account.HashPassword)
	if err != nil {
		log.Error("decoding hash password has failed", zap.Error(err))
		return nil, err
	}

	match, err := utils.CompareHashAndPassword(string(decodedHashPassword), *password)
	if err != nil {
		log.Error("error during comparing password", zap.Error(err))
		return nil, err
	} else if !match {
		log.Error("couldn't found entity")
		return nil, errs.NotMatchCredentialsError
	}

	return &account, nil
}

func (a *accountRepository) AccountByIdentifier(id int) (*stor.Account, error) {
	conn := a.storageProvider.GetConnection()
	log := a.GetLogger()

	var account stor.Account
	if err := conn.First(&account, id).Error; err != nil {
		log.Error("failed to search for account by id", zap.Error(err))
		return nil, err
	}

	return &account, nil
}
