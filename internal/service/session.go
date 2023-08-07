package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gits/internal/container"
	"gits/internal/model/app"
	"gits/internal/model/errs"
	"gits/internal/provider"
	"go.uber.org/zap"
	"time"
)

type Session interface {
	CreateAccountSession(ctx context.Context, account *app.Account) (*app.AccountSession, error)
	RetrieveAccountSession(ctx context.Context, accountSession *app.AccountSession) (*app.Account, error)
}

type session struct {
	container     container.Container
	cacheProvider provider.Cache
}

func NewSession(container container.Container, cache provider.Cache) Session {
	return &session{
		container:     container,
		cacheProvider: cache,
	}
}

func (s *session) CreateAccountSession(ctx context.Context, account *app.Account) (*app.AccountSession, error) {
	uuidKey := uuid.New()
	conf := s.container.GetConfig()
	ttl := conf.Cache.SessionTTL

	if err := s.SaveJSONData(ctx, account, uuidKey.String(), ttl); err != nil {
		return nil, err
	}
	return &app.AccountSession{
		SessionId: uuidKey.String(),
	}, nil
}

func (s *session) RetrieveAccountSession(ctx context.Context, accountSession *app.AccountSession) (*app.Account, error) {
	var (
		account app.Account
		uuidKey = accountSession.SessionId
	)
	if err := s.GetJSONData(ctx, uuidKey, &account); err != nil {
		return nil, err
	}
	return &account, nil
}

func (s *session) SaveJSONData(ctx context.Context, value interface{}, key string, ttl time.Duration) error {
	log := s.container.GetLogger()
	conn := s.cacheProvider.GetConnection()

	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Error("marshal json for session has failed", zap.Error(err))
		return err
	}
	resultCmd, err := conn.Set(ctx, key, jsonData, ttl).Result()
	if err != nil {
		log.Error("save value to session has failed", zap.Error(err))
		return err
	}
	log.Debug("save value to session has completed", zap.String("message", resultCmd))
	return nil
}

func (s *session) GetJSONData(ctx context.Context, key string, value interface{}) error {
	log := s.container.GetLogger()
	conn := s.cacheProvider.GetConnection()

	jsonText, err := conn.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Error("session doesn't contains json in cache", zap.String("key-path", key))
		return errs.NilError
	} else if err != nil {
		log.Error("session returns error when made get operation", zap.Error(err))
		return err
	}
	if err := json.Unmarshal([]byte(jsonText), value); err != nil {
		log.Error("cannot unmarshal data to expert struct", zap.Error(err))
		return err
	}
	return nil
}
