package storage

import (
	"gits/internal/container"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"go.uber.org/zap"
)

type ObservableRepository interface {
	CreateNewIp(ip *stor.Ip) (*uint, error)
	RetrieveIp(ip string) (*stor.Ip, error)
	ExistsIp(ip string) (bool, error)
	CreateNewObservable(observable *stor.Observable) error
	RetrieveObservables() ([]stor.Observable, error)
}

type observableRepository struct {
	container       container.Container
	storageProvider provider.Storage
}

func NewObservableRepository(container container.Container, storageProvider provider.Storage) ObservableRepository {
	return &observableRepository{
		container:       container,
		storageProvider: storageProvider,
	}
}

func (o *observableRepository) CreateNewObservable(observable *stor.Observable) error {
	log := o.container.GetLogger()
	conn := o.storageProvider.GetConnection()

	if err := conn.Save(observable).Error; err != nil {
		log.Error("can't save observable", zap.Error(err))

		return err
	}

	return nil
}

func (o *observableRepository) CreateNewIp(storIP *stor.Ip) (*uint, error) {
	conn := o.storageProvider.GetConnection()
	log := o.container.GetLogger()

	if err := conn.Save(&storIP).Error; err != nil {
		log.Error("can't save ip to db", zap.Error(err))
		return nil, err
	}

	return &storIP.ID, nil
}

func (o *observableRepository) RetrieveIp(ip string) (*stor.Ip, error) {
	conn := o.storageProvider.GetConnection()
	log := o.container.GetLogger()

	var storageIP stor.Ip
	if err := conn.Where("ip = ?", ip).First(&storageIP).Error; err != nil {
		log.Error("retrieve ip model from db has failed", zap.Error(err))
		return nil, err
	}

	return &storageIP, nil
}

func (o *observableRepository) ExistsIp(ip string) (bool, error) {
	conn := o.storageProvider.GetConnection()
	log := o.container.GetLogger()

	var exists bool
	if err := conn.Model(stor.Ip{}).Select("count(*) > 0").Where("ip = ?", ip).Find(&exists).Error; err != nil {
		log.Error("checking is exist ip has failed", zap.Error(err))
		return false, err
	}

	return exists, nil
}

func (o *observableRepository) RetrieveObservables() ([]stor.Observable, error) {
	log := o.container.GetLogger()
	conn := o.storageProvider.GetConnection()

	var observables []stor.Observable
	if err := conn.Order("updated_at desc").Preload("Ip").Preload("Account").Find(&observables).Error; err != nil {
		log.Error("retrieve observables has failed", zap.Error(err))
		return nil, err
	}

	return observables, nil
}
