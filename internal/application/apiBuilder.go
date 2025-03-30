package application

import (
	"fmt"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/domainservices"
	"github.com/VladSnap/gopherLoyalty/internal/features/accrual"
	"github.com/VladSnap/gopherLoyalty/internal/features/getbalance"
	"github.com/VladSnap/gopherLoyalty/internal/features/getorders"
	"github.com/VladSnap/gopherLoyalty/internal/features/getwithdrawals"
	"github.com/VladSnap/gopherLoyalty/internal/features/loginuser"
	"github.com/VladSnap/gopherLoyalty/internal/features/registrationuser"
	"github.com/VladSnap/gopherLoyalty/internal/features/uploadorder"
	"github.com/VladSnap/gopherLoyalty/internal/features/withdrawbalance"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/repositories"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
)

func CreateAPIServer(config *config.AppConfig, resMng *services.ResourceManager) (APIServer, AccrualWorker, error) {
	database, err := repositories.NewDatabaseLoyalty(config.DatabaseURI)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create DatabaseLoyalty: %w", err)
	}
	resMng.Register(database.Close)
	err = database.InitDatabase()
	if err != nil {
		return nil, nil, fmt.Errorf("failed init Database: %w", err)
	}

	userRepo := repositories.NewUserImplRepository(database)
	orderRepo := repositories.NewOrderImplRepository(database)
	withdrRepo := repositories.NewWithdrawImplRepository(database)
	bonusRepo := repositories.NewBonusCalculationImplRepository(database)
	passService := services.NewPasswordServiceImpl()
	jwtService := services.NewJWTTokenService()
	bonusAccounService := domainservices.NewBonusAccountServiceImpl(withdrRepo, bonusRepo)

	regCmd := registrationuser.NewRegistrationUserCmdHandler(userRepo, passService)
	registerUseCase := registrationuser.NewRegistrationUserUseCase(regCmd, jwtService)

	loginCmd := loginuser.NewRegistrationUserCmdHandler(userRepo, passService)
	loginUseCase := loginuser.NewLoginUserUseCase(loginCmd, jwtService)

	uploadOrderCmd := uploadorder.NewUploadOrderCmdHandler(userRepo, orderRepo, bonusRepo)
	uploadOrderUseCase := uploadorder.NewUploadOrderUseCase(uploadOrderCmd)

	getOrdersUseCase := getorders.NewGetOrdersUseCase(orderRepo)

	getBalanceUseCase := getbalance.NewGetBalanceUseCase(bonusAccounService)

	withdrawBalanceCmd := withdrawbalance.NewWithdrawBalanceCmdHandler(userRepo, orderRepo, withdrRepo, bonusAccounService)
	withdrawBalanceUseCase := withdrawbalance.NewWithdrawBalanceUseCase(withdrawBalanceCmd)

	getWithdrawalsUseCase := getwithdrawals.NewGetWithdrawalsUseCase(withdrRepo)

	server := NewAPIServer(config,
		jwtService,
		registerUseCase,
		loginUseCase,
		uploadOrderUseCase,
		getOrdersUseCase,
		getBalanceUseCase,
		withdrawBalanceUseCase,
		getWithdrawalsUseCase)

	accrualAPIClient := accrual.NewAccrualSystemClient(config.AccrualSystemAddress)
	accrualWorker := accrual.NewAccrualWorker(orderRepo, bonusRepo, bonusAccounService,
		accrualAPIClient, 1*time.Second)
	resMng.Register(accrualWorker.Close)

	return server, accrualWorker, nil
}
