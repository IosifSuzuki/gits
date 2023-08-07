package container

import (
	"gits/internal/model"
	"gits/pkg/logger"
)

type Container interface {
	GetConfig() *model.Config
	GetLogger() logger.Logger
}

type container struct {
	config *model.Config
	log    logger.Logger
}

func NewContainer(config *model.Config, log logger.Logger) Container {
	return &container{
		config: config,
		log:    log,
	}
}

func (c *container) GetConfig() *model.Config {
	return c.config
}

func (c *container) GetLogger() logger.Logger {
	return c.log
}
