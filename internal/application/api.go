package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/VladSnap/gopherLoyalty/internal/application/config"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api/middlewares"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	oapi "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui"
	swguiv5 "github.com/swaggest/swgui/v5emb"
	"github.com/swaggest/usecase"
)

type ApiServer interface {
	Start() error
}

type SwaggestApiServer struct {
	config                *config.AppConfig
	registerUseCase       RegistrationUserUseCase
	loginUseCase          LoginUserUseCase
	getBalanceUseCase     GetBalanceUseCase
	getOrdersUseCase      GetOrdersUseCase
	getWithdrawalsUseCase GetWithdrawalsUseCase
	uploadOrderUseCase    UploadOrderUseCase
	wthdrawBalanceUseCase WithdrawBalanceUseCase
}

func NewApiServer(config *config.AppConfig,
	registerUseCase RegistrationUserUseCase,
	loginUseCase LoginUserUseCase,
	uploadOrderUseCase UploadOrderUseCase,
	getOrdersUseCase GetOrdersUseCase,
	getBalanceUseCase GetBalanceUseCase,
	wthdrawBalanceUseCase WithdrawBalanceUseCase,
	getWithdrawalsUseCase GetWithdrawalsUseCase) *SwaggestApiServer {
	server := new(SwaggestApiServer)
	server.config = config
	server.registerUseCase = registerUseCase
	server.loginUseCase = loginUseCase
	server.uploadOrderUseCase = uploadOrderUseCase
	server.getOrdersUseCase = getOrdersUseCase
	server.getBalanceUseCase = getBalanceUseCase
	server.wthdrawBalanceUseCase = wthdrawBalanceUseCase
	server.getWithdrawalsUseCase = getWithdrawalsUseCase
	return server
}

func (server *SwaggestApiServer) Start() error {
	reflector := openapi31.NewReflector()
	service := web.NewService(reflector)

	// It allows OpenAPI configuration.
	service.OpenAPISchema().SetTitle("GopherLoyalty API")
	service.OpenAPISchema().SetDescription("This service provides API to Gopher Loyalty System.")
	service.OpenAPISchema().SetVersion("v1.0.0")

	// Additional middlewares can be added.
	service.Use(
		middlewares.LogMiddleware,
		middleware.Compress(1, "/*"),
		middleware.Recoverer,
	)

	server.registerRoutes(service)

	serv := &http.Server{
		Addr:    server.config.RunAddress,
		Handler: service,
	}

	// Горутина для прослушивания сигналов завершения.
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Zap.Info("Termination signal received. Stopping server....")
		if err := serv.Shutdown(context.Background()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Zap.Error("Error while stopping the server: %v\n", err)
		}
	}()

	swaggeConfig := swgui.Config{
		SettingsUI: map[string]string{
			"withCredentials": "true",
		},
	}
	service.Docs("/docs", swguiv5.NewWithConfig(swaggeConfig))

	// Запускаем прослушивание запросов.
	err := serv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed server listen: %w", err)
	}
	return nil
}

func (server *SwaggestApiServer) registerRoutes(service *web.Service) error {
	// Подготовка и настройка Interactor для корректной работы swagger
	registerInter := usecase.NewInteractor(server.registerUseCase.Execute)
	registerInter.SetTitle("Регистрация пользователя.")
	registerInter.SetName("registerUseCase")
	registerInter.SetTags("Not Auth")

	loginInter := usecase.NewInteractor(server.loginUseCase.Execute)
	loginInter.SetTitle("Аутентификация пользователя.")
	loginInter.SetName("loginUseCase")
	loginInter.SetTags("Not Auth")

	uploadOrderInter := usecase.NewInteractor(server.uploadOrderUseCase.Execute)
	uploadOrderInter.SetTitle("Загрузка номера заказа. Требуется аутентификация.")
	uploadOrderInter.SetName("uploadOrderUseCase")
	uploadOrderInter.SetTags("Required Auth")

	getOrdersInter := usecase.NewInteractor(server.getOrdersUseCase.Execute)
	getOrdersInter.SetTitle("Получение списка загруженных номеров заказов. Требуется аутентификация и авторизация.")
	getOrdersInter.SetName("getOrdersUseCase")
	getOrdersInter.SetTags("Required Auth")

	getBalanceInter := usecase.NewInteractor(server.getBalanceUseCase.Execute)
	getBalanceInter.SetTitle("Получение текущего баланса счёта баллов лояльности пользователя. Требуется аутентификация и авторизация.")
	getBalanceInter.SetName("getBalanceUseCase")
	getBalanceInter.SetTags("Required Auth")

	wthdrawBalanceInter := usecase.NewInteractor(server.wthdrawBalanceUseCase.Execute)
	wthdrawBalanceInter.SetTitle("Запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа. Требуется аутентификация и авторизация.")
	wthdrawBalanceInter.SetName("wthdrawBalanceUseCase")
	wthdrawBalanceInter.SetTags("Required Auth")

	getWithdrawalsInter := usecase.NewInteractor(server.getWithdrawalsUseCase.Execute)
	getWithdrawalsInter.SetTitle("Получение информации о выводе средств с накопительного счёта пользователем. Требуется аутентификация и авторизация.")
	getWithdrawalsInter.SetName("getWithdrawalsUseCase")
	getWithdrawalsInter.SetTags("Required Auth")

	apiDocAuthDoc := nethttp.APIKeySecurityMiddleware(service.OpenAPICollector, "UserAuthToken",
		"Authorization", oapi.InHeader, "Authorization token.")

	service.Post("/api/user/register", registerInter, nethttp.SuccessStatus(http.StatusOK))
	service.Post("/api/user/login", loginInter, nethttp.SuccessStatus(http.StatusOK))

	// Регистрация защищенных аутентификацией и авторизацией маршрутов
	service.Route("/api/user", func(r chi.Router) {
		r.Group(func(rg chi.Router) {
			rg.Use(middlewares.AuthMiddleware, apiDocAuthDoc)
			rg.Method(http.MethodPost, "/orders", nethttp.NewHandler(uploadOrderInter, nethttp.SuccessStatus(http.StatusOK)))
			rg.Method(http.MethodGet, "/orders", nethttp.NewHandler(getOrdersInter))
			rg.Method(http.MethodGet, "/balance", nethttp.NewHandler(getBalanceInter))
			rg.Method(http.MethodPost, "/balance/withdraw", nethttp.NewHandler(wthdrawBalanceInter, nethttp.SuccessStatus(http.StatusOK)))
			rg.Method(http.MethodGet, "/withdrawals", nethttp.NewHandler(getWithdrawalsInter))
		})
	})

	return nil
}
