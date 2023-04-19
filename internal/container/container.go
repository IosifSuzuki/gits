package container

import (
	"gits/internal/model"
)

type Container interface {
	GetConfig() *model.Config
}

type container struct {
	config *model.Config
}

func NewContainer(config *model.Config) Container {
	return &container{
		config: config,
	}
}

func (c *container) GetConfig() *model.Config {
	return c.config
}
