package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"gits/internal/container"
	"gits/internal/model/constant"
	"gits/internal/model/dto"
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
		account, err := retrieveAccount(ctx, sessionCtx, a.AccountSession)
		if err != nil {
			log.Error("session doesn't exist", zap.Error(err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(constant.AccountAppKey, account)
		ctx.Next()
	}
}

func retrieveAccount(ctx *gin.Context, sessionCtx context.Context, accountSession service.AccountSession) (*dto.Account, error) {
	sessionId, err := ctx.Cookie(constant.CookieSessionKey)
	if err != nil {
		return nil, err
	}

	accountSessionModel := dto.AccountSession{
		SessionId: sessionId,
	}

	return accountSession.RetrieveAccountSession(sessionCtx, &accountSessionModel)
}
