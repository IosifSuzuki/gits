package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gits/internal/container"
	"gits/internal/model/dto"
	"gits/internal/model/errs"
	"gits/internal/provider"
	"go.uber.org/zap"
	"time"
)

type AccountSession interface {
	CreateAccountSession(ctx context.Context, account *dto.Account) (*dto.AccountSession, error)
	RetrieveAccountSession(ctx context.Context, accountSession *dto.AccountSession) (*dto.Account, error)
}

type accountSession struct {
	container     container.Container
	cacheProvider provider.Cache
}

func NewSession(container container.Container, cache provider.Cache) AccountSession {
	return &accountSession{
		container:     container,
		cacheProvider: cache,
	}
}

func (s *accountSession) CreateAccountSession(ctx context.Context, account *dto.Account) (*dto.AccountSession, error) {
	uuidKey := uuid.New()
	conf := s.container.GetConfig()
	ttl := conf.Cache.SessionTTL

	log := s.container.GetLogger()

	if err := s.SaveJSONData(ctx, account, uuidKey.String(), ttl); err != nil {
		log.Error("fail to save json to service", zap.Error(err))
		return nil, err
	}
	return &dto.AccountSession{
		SessionId: uuidKey.String(),
	}, nil
}

func (s *accountSession) RetrieveAccountSession(ctx context.Context, accountSession *dto.AccountSession) (*dto.Account, error) {
	var (
		account dto.Account
		uuidKey = accountSession.SessionId
	)
	log := s.container.GetLogger()

	if err := s.GetJSONData(ctx, uuidKey, &account); err != nil {
		log.Error("fail to get json from service", zap.Error(err))
		return nil, err
	}
	return &account, nil
}

func (s *accountSession) SaveJSONData(ctx context.Context, value interface{}, key string, ttl time.Duration) error {
	log := s.container.GetLogger()
	conn := s.cacheProvider.GetConnection()

	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Error("marshal json for accountSession has failed", zap.Error(err))
		return err
	}
	resultCmd, err := conn.Set(ctx, key, jsonData, ttl).Result()
	if err != nil {
		log.Error("save value to accountSession has failed", zap.Error(err))
		return err
	}
	log.Debug("save value to accountSession has completed", zap.String("message", resultCmd))
	return nil
}

func (s *accountSession) GetJSONData(ctx context.Context, key string, value interface{}) error {
	log := s.container.GetLogger()
	conn := s.cacheProvider.GetConnection()

	jsonText, err := conn.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Error("accountSession doesn't contains json in cache", zap.String("key-path", key))
		return errs.NilError
	} else if err != nil {
		log.Error("accountSession returns error when made get operation", zap.Error(err))
		return err
	}
	if err := json.Unmarshal([]byte(jsonText), value); err != nil {
		log.Error("cannot unmarshal data to expert struct", zap.Error(err))
		return err
	}
	return nil
}
