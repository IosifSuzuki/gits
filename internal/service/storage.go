package service

import (
	"encoding/base64"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"gits/internal/utils"
)

type Storage interface {
	AccountByCredential(credential *app.Credential) (*app.Account, error)
}

type storage struct {
	container.Container
	storageProvider provider.Storage
}

func NewStorage(container container.Container, storageProvider provider.Storage) Storage {
	return &storage{
		Container:       container,
		storageProvider: storageProvider,
	}
}

func (s *storage) AccountByCredential(credential *app.Credential) (*app.Account, error) {
	conn := s.storageProvider.GetConnection()

	var account stor.Account
	if err := conn.Where("username = ?", credential.Username).First(&account).Error; err != nil {
		return nil, err
	}
	decodedHashPassword, err := base64.StdEncoding.DecodeString(account.HashPassword)
	if err != nil {
		return nil, err
	}
	match, err := utils.CompareHashAndPassword(string(decodedHashPassword), credential.Password)
	if err != nil {
		return nil, err
	} else if !match {
		return nil, errs.NotMatchCredentialsError
	}
	var appAccount = app.Account{
		Id:           int(account.ID),
		Username:     account.Username,
		HashPassword: account.HashPassword,
		Role:         app.ParseString(string(account.Role)),
	}
	return &appAccount, nil
}
