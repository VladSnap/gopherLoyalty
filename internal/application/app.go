package application

import (
	"context"
	"fmt"
	"os"

	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
	"go.uber.org/zap"
)

var resourceManager *services.ResourceManager

func Run() {
	resourceManager = services.NewResourceManager()
	defer func() {
		err := resourceManager.Cleanup()

		if err != nil {
			panic(fmt.Errorf("failed resourceManager clean: %w", err))
		}
	}()
	// Регаем функцию Sync Zap логов
	resourceManager.Register(log.Close)

	log.Zap.Info("run shorneter server", zap.Strings("Args", os.Args))

	confValidator := &config.OptionsValidator{}
	conf, err := config.LoadConfig(confValidator)
	if err != nil {
		panic(err)
	}

	apiServer, accrualWorker, err := CreateApiServer(conf, resourceManager)
	if err != nil {
		panic(err)
	}

	accrualWorker.StartPolling(context.Background())

	err = apiServer.Start()

	if err != nil {
		log.Zap.Errorf("failed stop server: %w", err)
	}
	log.Zap.Info("app.Run end")
}
