package service

import (
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"go.uber.org/zap"
)

type Storage interface {
	ExistsIp(ip string) (bool, error)
	CreateNewIp(ip *app.Ip) (*uint, error)
	RetrieveIp(ip string) (*app.Ip, error)
	CreateNewObservable(observable *app.Observable) (bool, error)
	RetrieveObservables() ([]app.Observable, error)
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

func (s *storage) ExistsIp(ip string) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	var exists bool
	if err := conn.Model(stor.Ip{}).Select("count(*) > 0").Where("ip = ?", ip).Find(&exists).Error; err != nil {
		log.Error("checking is exist ip has failed", zap.Error(err))
		return false, err
	}
	return exists, nil
}

func (s *storage) RetrieveIp(ip string) (*app.Ip, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	var storageIp stor.Ip
	if err := conn.Where("ip = ?", ip).First(&storageIp).Error; err != nil {
		log.Error("retrieve ip model from db has failed", zap.Error(err))
		return nil, err
	}
	return &app.Ip{
		ID:       &storageIp.ID,
		Ip:       storageIp.Ip,
		Hostname: storageIp.Hostname,
		City:     storageIp.City,
		Region:   storageIp.Region,
		Country:  storageIp.Country,
		Loc:      storageIp.Loc,
		Org:      storageIp.Org,
		Postal:   storageIp.Postal,
		Timezone: storageIp.Timezone,
	}, nil
}

func (s *storage) CreateNewIp(ip *app.Ip) (*uint, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	storageIp := stor.Ip{
		Ip:         ip.Ip,
		Hostname:   ip.Hostname,
		City:       ip.City,
		Region:     ip.Region,
		Country:    ip.Country,
		Loc:        ip.Loc,
		Org:        ip.Org,
		Postal:     ip.Postal,
		Timezone:   ip.Timezone,
		Observable: nil,
	}

	if err := conn.Save(&storageIp).Error; err != nil {
		log.Error("save ip to db has failed", zap.Error(err))
		return nil, err
	}
	return &storageIp.ID, nil
}

func (s *storage) CreateNewObservable(observable *app.Observable) (bool, error) {
	conn := s.storageProvider.GetConnection()
	log := s.GetLogger()

	ipAppModel := observable.Ip
	if ipAppModel == nil {
		err := errs.NilError
		log.Error("ip has nil value", zap.Error(err))
		return false, err
	}
	existIpStorageModel, err := s.ExistsIp(*ipAppModel.Ip)
	if err != nil {
		log.Error("exists ip operation has failed", zap.Error(err))
		return false, err
	}
	if !existIpStorageModel {
		ipID, err := s.CreateNewIp(ipAppModel)
		if err != nil {
			log.Error("create new ip has failed", zap.Error(err))
			return false, err
		}
		ipAppModel.ID = ipID
	} else {
		ipAppModel, err = s.RetrieveIp(*ipAppModel.Ip)
		if err != nil {
			log.Error("retrieve ip has failed", zap.Error(err))
			return false, err
		}
	}
	storageObservable := stor.Observable{
		AccountId: observable.AccountId,
		IpId:      *ipAppModel.ID,
		Browser:   observable.Browser,
		OS:        observable.OS,
		OSVersion: observable.OSVersion,
		Path:      observable.Path,
		Device:    observable.Device,
	}

	if err := conn.Save(&storageObservable).Error; err != nil {
		log.Error("observable model has failed to save in db", zap.Error(err))
		return false, err
	}
	return true, nil
}

func (s *storage) RetrieveObservables() ([]app.Observable, error) {
	log := s.GetLogger()
	conn := s.storageProvider.GetConnection()

	var observables []stor.Observable
	if err := conn.Order("updated_at desc").Preload("Ip").Preload("Account").Find(&observables).Error; err != nil {
		log.Error("retrieve observables has failed", zap.Error(err))
		return nil, err
	}
	var appObservables = make([]app.Observable, 0, len(observables))
	for _, observable := range observables {
		var (
			appAccount *app.Account
			appIp      *app.Ip
		)
		storageAccount := observable.Account
		if storageAccount != nil {
			appAccount = &app.Account{
				Id:           int(storageAccount.ID),
				Username:     storageAccount.Username,
				HashPassword: storageAccount.HashPassword,
				Role:         app.Role(storageAccount.Role),
			}
		}
		storageIp := observable.Ip
		if storageIp != nil {
			appIp = &app.Ip{
				ID:       &storageIp.ID,
				Ip:       storageIp.Ip,
				Hostname: storageIp.Hostname,
				City:     storageIp.City,
				Region:   storageIp.Region,
				Country:  storageIp.Country,
				Loc:      storageIp.Loc,
				Org:      storageIp.Org,
				Postal:   storageIp.Postal,
				Timezone: storageIp.Timezone,
			}
		}
		appObservables = append(appObservables, app.Observable{
			AccountId: observable.AccountId,
			Account:   appAccount,
			Ip:        appIp,
			Browser:   observable.Browser,
			OS:        observable.OS,
			OSVersion: observable.OSVersion,
			Path:      observable.Path,
			Device:    observable.Device,
			UpdatedAt: &observable.UpdatedAt,
		})
	}
	return appObservables, nil
}
