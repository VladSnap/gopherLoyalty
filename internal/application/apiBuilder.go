package application

import (
	"fmt"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/domain/domainServices"
	"github.com/VladSnap/gopherLoyalty/internal/features/accrual"
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
	orderRepo := repositories.NewOrderImplRepository(database)
	withdrRepo := repositories.NewWithdrawImplRepository(database)
	bonusRepo := repositories.NewBonusCalculationImplRepository(database)
	passService := services.NewPasswordServiceImpl()
	jwtService := services.NewJWTTokenService()
	bonusAccounService := domainServices.NewBonusAccountServiceImpl(withdrRepo, bonusRepo)

	regCmd := registrationUser.NewRegistrationUserCmdHandler(userRepo, passService)
	registerUseCase := registrationUser.NewRegistrationUserUseCase(regCmd, jwtService)

	loginCmd := loginUser.NewRegistrationUserCmdHandler(userRepo, passService)
	loginUseCase := loginUser.NewLoginUserUseCase(loginCmd, jwtService)

	uploadOrderCmd := uploadOrder.NewUploadOrderCmdHandler(userRepo, orderRepo, bonusRepo)
	uploadOrderUseCase := uploadOrder.NewUploadOrderUseCase(uploadOrderCmd)

	getOrdersUseCase := getOrders.NewGetOrdersUseCase(orderRepo)

	getBalanceUseCase := getBalance.NewGetBalanceUseCase(bonusAccounService)

	withdrawBalanceCmd := withdrawBalance.NewWithdrawBalanceCmdHandler(userRepo, orderRepo, withdrRepo, bonusAccounService)
	withdrawBalanceUseCase := withdrawBalance.NewWithdrawBalanceUseCase(withdrawBalanceCmd)

	getWithdrawalsUseCase := getWithdrawals.NewGetWithdrawalsUseCase(withdrRepo)

	server := NewApiServer(config,
		jwtService,
		registerUseCase,
		loginUseCase,
		uploadOrderUseCase,
		getOrdersUseCase,
		getBalanceUseCase,
		withdrawBalanceUseCase,
		getWithdrawalsUseCase)

	accrualApiClient := accrual.NewAccrualSystemClient(config.AccrualSystemAddress)
	accrualWorker := accrual.NewAccrualWorker(orderRepo, bonusRepo, bonusAccounService,
		accrualApiClient, 1*time.Second)
	resMng.Register(accrualWorker.Close)

	return server, nil
}
