package app

import (
	"gits/internal/container"
	"gits/internal/controller/central"
	"gits/internal/middleware"
	"gits/internal/model"
	"gits/internal/provider"
	"gits/internal/router"
	"gits/internal/service"
	"gits/internal/service/storage"
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

	storageDAO := storage.NewDAO(box, storageProvider)
	cacheProvider := provider.NewCache(box)
	sessionService := service.NewSession(box, cacheProvider)
	decompressorFile := service.NewDecompressorFile(box)
	md := service.NewMD(box)
	ip := service.NewIp(box)
	attachmentStorage, err := service.NewAttachmentStorage(box)
	if err != nil {
		return err
	}

	mainController := central.NewMainController(
		box, sessionService, decompressorFile, attachmentStorage, md, storageDAO,
	)

	authMiddleware := middleware.NewAuth(box, sessionService)
	observable := middleware.NewObserver(box, sessionService, ip, storageDAO)
	errorHandlerMiddleware := middleware.NewErrorHandler(box)

	r := router.NewRouter(box, mainController, authMiddleware, errorHandlerMiddleware, observable)

	return r.SetupHandlers()
}
