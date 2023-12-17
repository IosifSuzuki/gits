package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"gits/internal/container"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
	stor "gits/internal/model/storage"
	"gits/internal/service"
	"gits/internal/service/storage"
	"go.uber.org/zap"
)

type Observer interface {
	Observer() gin.HandlerFunc
}

type observer struct {
	container.Container
	AccountSession service.AccountSession
	IPService      service.Ip
	dao            storage.DAO
}

func NewObserver(container container.Container, accountSession service.AccountSession, Ip service.Ip, dao storage.DAO) Observer {
	return &observer{
		container,
		accountSession,
		Ip,
		dao,
	}
}

func (o *observer) Observer() gin.HandlerFunc {
	log := o.GetLogger()
	sessionCtx := context.Background()

	return func(ctx *gin.Context) {
		ctx.Next()

		var (
			account          *dto.Account
			obtainAccountErr error
		)
		accountValue, ok := ctx.Get(constant.AccountAppKey)
		if ok {
			account = accountValue.(*dto.Account)
		}

		if account == nil {
			account, obtainAccountErr = retrieveAccount(ctx, sessionCtx, o.AccountSession)
		}
		if obtainAccountErr != nil {
			log.Error("retrieve account by session id has failed", zap.Error(obtainAccountErr))
		}

		var accountId *int
		if account != nil {
			accountId = &account.ID
		}

		ip := ctx.ClientIP()
		storIP, err := o.ObtainStorageIp(ip)
		if err != nil {
			log.Error("obtain storage ip has failed", zap.Error(err))
			storIP, err = o.SaveStorageIp(ip)
		}
		if err != nil {
			log.Error("save storage ip has failed", zap.Error(err))
		}

		ua := ctx.GetHeader(constant.UserAgentHeaderKey)
		userAgentModel := useragent.Parse(ua)

		storObservable := stor.Observable{
			AccountId: accountId,
			IpId:      storIP.ID,
			Browser:   &userAgentModel.Name,
			OS:        &userAgentModel.OS,
			OSVersion: &userAgentModel.OSVersion,
			Path:      &ctx.Request.URL.Path,
			Device:    &userAgentModel.Device,
		}

		if err := o.dao.GetObservableRepository().CreateNewObservable(&storObservable); err != nil {
			log.Error("create new observable has failed", zap.Error(err))
		}
	}
}

func (o *observer) ObtainStorageIp(ip string) (*stor.Ip, error) {
	log := o.GetLogger()
	var (
		storIP *stor.Ip
		err    error
	)

	existsIP, err := o.dao.GetObservableRepository().ExistsIp(ip)
	if existsIP {
		storIP, err = o.GetStorageIp(ip)
	} else {
		storIP, err = o.GetStorageIp(ip)
	}

	if err != nil {
		log.Error("obtain storage ip has failed", zap.Error(err))
		return nil, err
	}

	return storIP, nil
}

func (o *observer) GetStorageIp(ip string) (*stor.Ip, error) {
	log := o.GetLogger()

	storIP, err := o.dao.GetObservableRepository().RetrieveIp(ip)
	if err != nil {
		log.Error("fail to get ip from db", zap.Error(err))
	}

	return storIP, err
}

func (o *observer) SaveStorageIp(ip string) (*stor.Ip, error) {
	log := o.GetLogger()

	geoLocation, err := o.IPService.GeoLocation(ip)
	if err != nil {
		log.Error("get geolocation has failed", zap.Error(err))
		return nil, err
	}

	storIP := stor.Ip{
		Ip:       geoLocation.Ip,
		Hostname: geoLocation.Hostname,
		City:     geoLocation.City,
		Region:   geoLocation.Region,
		Country:  geoLocation.Country,
		Loc:      geoLocation.Loc,
		Org:      geoLocation.Org,
		Postal:   geoLocation.Postal,
		Timezone: geoLocation.Timezone,
	}
	_, err = o.dao.GetObservableRepository().CreateNewIp(&storIP)
	if err != nil {
		log.Error("fail to create new ip model in db", zap.Error(err))
		return nil, err
	}

	return &storIP, nil
}
