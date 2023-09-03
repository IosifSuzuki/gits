package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/constant"
	"gits/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type Auth interface {
	Authorization() gin.HandlerFunc
}

type auth struct {
	container.Container
	service.AccountSession
}

func NewAuth(container container.Container, session service.AccountSession) Auth {
	return &auth{
		Container:      container,
		AccountSession: session,
	}
}

func (a *auth) Authorization() gin.HandlerFunc {
	log := a.GetLogger()
	sessionCtx := context.Background()
	return func(ctx *gin.Context) {
		sessionId, err := ctx.Cookie(constant.CookieSessionKey)
		if err != nil {
			log.Error("cookie doesn't have session", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		accountSession := app.AccountSession{
			SessionId: sessionId,
		}
		account, err := a.AccountSession.RetrieveAccountSession(sessionCtx, &accountSession)
		if err != nil {
			log.Error("session doesn't exist", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(constant.AccountAppKey, account)
		ctx.Next()
	}
}
