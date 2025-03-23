package application

import (
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/features/getBalance"
	"github.com/VladSnap/gopherLoyalty/internal/features/getOrders"
	"github.com/VladSnap/gopherLoyalty/internal/features/getWithdrawals"
	"github.com/VladSnap/gopherLoyalty/internal/features/loginUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/registrationUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/uploadOrder"
	"github.com/VladSnap/gopherLoyalty/internal/features/withdrawBalance"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/repositories"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
)

func CreateApiServer(config *config.AppConfig, resMng *services.ResourceManager) (ApiServer, error) {
	database, err := repositories.NewDatabaseLoyalty(config.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("failed create DatabaseShortener: %w", err)
	}
	resMng.Register(database.Close)
	err = database.InitDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed init Database: %w", err)
	}

	userRepo := repositories.NewUserImplRepository(database)
	passService := services.NewPasswordServiceImpl()

	regCmd := registrationUser.NewRegistrationUserCmdHandler(userRepo, passService)
	registerUseCase := registrationUser.NewRegistrationUserUseCase(regCmd)

	loginUseCase := loginUser.NewLoginUserUseCase()

	uploadOrderUseCase := uploadOrder.NewUploadOrderUseCase()

	getOrdersUseCase := getOrders.NewGetOrdersUseCase()

	getBalanceUseCase := getBalance.NewGetBalanceUseCase()

	withdrawBalanceUseCase := withdrawBalance.NewWithdrawBalanceUseCase()

	getWithdrawalsUseCase := getWithdrawals.NewGetWithdrawalsUseCase()

	server := NewApiServer(config,
		registerUseCase,
		loginUseCase,
		uploadOrderUseCase,
		getOrdersUseCase,
		getBalanceUseCase,
		withdrawBalanceUseCase,
		getWithdrawalsUseCase)

	return server, nil
}
