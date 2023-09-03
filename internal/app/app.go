package app

import (
	"gits/internal/container"
	"gits/internal/controller"
	"gits/internal/middleware"
	"gits/internal/model"
	"gits/internal/provider"
	"gits/internal/router"
	"gits/internal/service"
	"gits/pkg/logger"
)

func Run() error {
	conf, err := model.NewConfig()
	if err != nil {
		return err
	}
	log := logger.NewLogger(logger.DebugLevel)
	box := container.NewContainer(conf, log)
	storageProvider, err := provider.NewStorage(box)
	if err != nil {
		return err
	}
	storageService := service.NewStorage(box, storageProvider)
	cacheProvider := provider.NewCache(box)
	sessionService := service.NewSession(box, cacheProvider)
	decompressorFile := service.NewDecompressorFile(box)
	md := service.NewMD(box)
	attachmentStorage, err := service.NewAttachmentStorage(box)
	if err != nil {
		return err
	}
	mainController := controller.NewMainController(box, storageService, sessionService, decompressorFile, attachmentStorage, md)
	authMiddleware := middleware.NewAuth(box, sessionService)
	errorHandlerMiddleware := middleware.NewErrorHandler(box)
	r := router.NewRouter(box, mainController, authMiddleware, errorHandlerMiddleware)
	return r.SetupHandlers()
}
