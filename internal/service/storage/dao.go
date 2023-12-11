package storage

import (
	"gits/internal/container"
	"gits/internal/provider"
)

type DAO interface {
	GetAccountRepository() AccountRepository
	GetArticleRepository() ArticleRepository
	GetObservableRepository() ObservableRepository
}

type dao struct {
	container       container.Container
	storageProvider provider.Storage
}

func NewDAO(container container.Container, storageProvider provider.Storage) DAO {
	return &dao{
		container:       container,
		storageProvider: storageProvider,
	}
}

func (d *dao) GetAccountRepository() AccountRepository {
	return NewAccount(d.container, d.storageProvider)
}

func (d *dao) GetArticleRepository() ArticleRepository {
	return NewArticleRepository(d.container, d.storageProvider)
}

func (d *dao) GetObservableRepository() ObservableRepository {
	return NewObservableRepository(d.container, d.storageProvider)
}
