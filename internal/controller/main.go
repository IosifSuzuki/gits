package controller

import (
	"context"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/form"
	"gits/internal/model/response"
	"gits/internal/service"
	"go.uber.org/zap"
)

type mainController struct {
	container.Container
	service.Storage
	service.Session
}

type MainController interface {
	Auth(authentication *form.Authentication) (*response.AuthSessionResponse, error)
}

func NewMainController(container container.Container, storage service.Storage, session service.Session) MainController {
	return &mainController{
		Container: container,
		Storage:   storage,
		Session:   session,
	}
}

func (m *mainController) Auth(authentication *form.Authentication) (*response.AuthSessionResponse, error) {
	log := m.GetLogger()
	ctx := context.Background()

	credential := &app.Credential{
		Username: authentication.Username,
		Password: authentication.Password,
	}
	account, err := m.Storage.AccountByCredential(credential)
	if err != nil {
		log.Error("retrieve account by credential has failed", zap.Error(err))
		return nil, err
	}
	accountSession, err := m.Session.CreateAccountSession(ctx, account)
	if err != nil {
		log.Error("create session has failed", zap.Error(err))
		return nil, err
	}

	return &response.AuthSessionResponse{
		SessionId: accountSession.SessionId,
	}, nil
}
