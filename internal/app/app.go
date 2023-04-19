package app

import (
	"gits/internal/container"
	"gits/internal/controller"
	"gits/internal/model"
	"gits/internal/services"
)

func Run() error {
	conf, err := model.NewConfig()
	if err != nil {
		return err
	}
	box := container.NewContainer(conf)
	md := services.NewMD()
	publisher := services.NewPublisher(md)
	c := controller.NewController(box, publisher)
	return c.SetupHandlers()
}
