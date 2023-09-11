package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/constant"
	"gits/internal/service"
	"go.uber.org/zap"
)

type Observer interface {
	Observer() gin.HandlerFunc
}

type observer struct {
	container.Container
	AccountSession service.AccountSession
	IpService      service.Ip
	storageService service.Storage
}

func NewObserver(container container.Container, accountSession service.AccountSession, Ip service.Ip, storageService service.Storage) Observer {
	return &observer{
		container,
		accountSession,
		Ip,
		storageService,
	}
}

func (o *observer) Observer() gin.HandlerFunc {
	log := o.GetLogger()
	sessionCtx := context.Background()

	return func(ctx *gin.Context) {
		ctx.Next()

		var (
			account            *app.Account
			retrieveAccountErr error
			geoLocation        *app.GeoLocation
		)
		accountValue, ok := ctx.Get(constant.AccountAppKey)
		if ok {
			account = accountValue.(*app.Account)
		}
		if account == nil {
			account, retrieveAccountErr = retrieveAccount(ctx, sessionCtx, o.AccountSession)
		}
		if retrieveAccountErr != nil {
			log.Error("retrieve account by session id has failed", zap.Error(retrieveAccountErr))
		}
		var accountId *int
		if account != nil {
			accountId = &account.Id
		}
		ip := ctx.ClientIP()
		exists, err := o.storageService.ExistsIp(ip)
		if err != nil {
			log.Error("exist ip operation has failed", zap.Error(err))
		}
		if exists {
			geoLocation = &app.GeoLocation{
				Ip: &ip,
			}
		} else {
			geoLocation, err = o.IpService.GeoLocation(ip)
		}
		if err != nil {
			log.Error("retrieve geo location has failed", zap.Error(err))
			return
		}
		ua := ctx.GetHeader(constant.UserAgentHeaderKey)
		userAgentModel := useragent.Parse(ua)

		ipAppModel := app.Ip{
			ID:       nil,
			Ip:       &ip,
			Hostname: geoLocation.Hostname,
			City:     geoLocation.City,
			Region:   geoLocation.Region,
			Country:  geoLocation.Country,
			Loc:      geoLocation.Loc,
			Org:      geoLocation.Org,
			Postal:   geoLocation.Postal,
			Timezone: geoLocation.Timezone,
		}
		observable := app.Observable{
			AccountId: accountId,
			Ip:        &ipAppModel,
			Browser:   &userAgentModel.Name,
			OS:        &userAgentModel.OS,
			OSVersion: &userAgentModel.OSVersion,
			Path:      &ctx.Request.URL.Path,
			Device:    &userAgentModel.Device,
		}
		ok, err = o.storageService.CreateNewObservable(&observable)
		if err != nil {
			log.Error("create new observable has failed", zap.Error(err))
			return
		} else if !ok {
			log.Error("unsuccessfully create new observable")
			return
		}
	}
}
