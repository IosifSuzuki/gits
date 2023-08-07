package provider

import (
	"gits/internal/container"
	"gits/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage interface {
	GetConnection() *gorm.DB
	CloseConnection() error
}

type storage struct {
	container container.Container
	db        *gorm.DB
}

func NewStorage(container container.Container) (Storage, error) {
	db, err := connectToDB(container.GetConfig())
	if err != nil {
		return nil, err
	}
	return &storage{
		container: container,
		db:        db,
	}, nil
}

func connectToDB(config *model.Config) (*gorm.DB, error) {
	db, err := gorm.Open(
		postgres.Open(config.Storage.Address()),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			IgnoreRelationshipsWhenMigrating:         true,
			Logger:                                   logger.Default.LogMode(logger.Info),
		})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (s *storage) GetConnection() *gorm.DB {
	return s.db
}

func (s *storage) CloseConnection() error {
	db, err := s.db.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	return err
}
