package application

import (
	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/features/getBalance"
	"github.com/VladSnap/gopherLoyalty/internal/features/getOrders"
	"github.com/VladSnap/gopherLoyalty/internal/features/getWithdrawals"
	"github.com/VladSnap/gopherLoyalty/internal/features/loginUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/registrationUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/uploadOrder"
	"github.com/VladSnap/gopherLoyalty/internal/features/withdrawBalance"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
)

func CreateApiServer(config *config.AppConfig, resMng *services.ResourceManager) (ApiServer, error) {
	registerUseCase := registrationUser.NewRegistrationUserUseCase()
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
