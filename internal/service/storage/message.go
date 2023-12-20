package storage

import (
	"gits/internal/container"
	stor "gits/internal/model/storage"
	"gits/internal/provider"
	"go.uber.org/zap"
)

type messageRepository struct {
	container       container.Container
	storageProvider provider.Storage
}

type MessageRepository interface {
	Create(message *stor.Message) (*uint, error)
}

func NewMessageRepository(container container.Container, storageProvider provider.Storage) MessageRepository {
	return &messageRepository{
		container:       container,
		storageProvider: storageProvider,
	}
}

func (m *messageRepository) Create(message *stor.Message) (*uint, error) {
	log := m.container.GetLogger()
	conn := m.storageProvider.GetConnection()

	if err := conn.Save(message).Error; err != nil {
		log.Error("save message has failed", zap.Error(err))
		return nil, err
	}

	return &message.ID, nil
}
