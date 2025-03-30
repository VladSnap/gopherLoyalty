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
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	oapi "github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/swgui"
	swguiv5 "github.com/swaggest/swgui/v5emb"
)

type APIServer interface {
	Start() error
}

type AccrualWorker interface {
	StartPolling()
}

type SwaggestAPIServer struct {
	config                *config.AppConfig
	jwtService            services.JWTTokenService
	registerUseCase       RegistrationUserUseCase
	loginUseCase          LoginUserUseCase
	getBalanceUseCase     GetBalanceUseCase
	getOrdersUseCase      GetOrdersUseCase
	getWithdrawalsUseCase GetWithdrawalsUseCase
	uploadOrderUseCase    UploadOrderUseCase
	wthdrawBalanceUseCase WithdrawBalanceUseCase
}

func NewAPIServer(config *config.AppConfig,
	jwtService services.JWTTokenService,
	registerUseCase RegistrationUserUseCase,
	loginUseCase LoginUserUseCase,
	uploadOrderUseCase UploadOrderUseCase,
	getOrdersUseCase GetOrdersUseCase,
	getBalanceUseCase GetBalanceUseCase,
	wthdrawBalanceUseCase WithdrawBalanceUseCase,
	getWithdrawalsUseCase GetWithdrawalsUseCase) *SwaggestAPIServer {
	server := new(SwaggestAPIServer)
	server.config = config
	server.jwtService = jwtService
	server.registerUseCase = registerUseCase
	server.loginUseCase = loginUseCase
	server.uploadOrderUseCase = uploadOrderUseCase
	server.getOrdersUseCase = getOrdersUseCase
	server.getBalanceUseCase = getBalanceUseCase
	server.wthdrawBalanceUseCase = wthdrawBalanceUseCase
	server.getWithdrawalsUseCase = getWithdrawalsUseCase
	return server
}

func (server *SwaggestAPIServer) Start() error {
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

	err := server.registerRoutes(service)
	if err != nil {
		return fmt.Errorf("failed registerRoutes: %w", err)
	}

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
	err = serv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed server listen: %w", err)
	}
	return nil
}

func (sr *SwaggestAPIServer) registerRoutes(sv *web.Service) error {
	// Middleware для аутентификации
	authMiddleware := middlewares.NewAuthMiddleware(sr.jwtService)
	apiDocAuthDoc := nethttp.APIKeySecurityMiddleware(
		sv.OpenAPICollector, "UserAuthToken", "Authorization", oapi.InHeader, "Authorization token.",
	)

	// Роуты без аутентификации
	regRoute(sv.Router, sr.registerUseCase.Execute, http.MethodPost, "/api/user/register", http.StatusOK,
		"Регистрация пользователя.", "Регистрация нового пользователя.", "Not Auth")
	regRoute(sv.Router, sr.loginUseCase.Execute, http.MethodPost, "/api/user/login", http.StatusOK,
		"Аутентификация пользователя.", "Вход пользователя в систему.", "Not Auth")

	// Роуты с аутентификацией
	sv.Route("/api/user", func(r chi.Router) {
		r.Group(func(rg chi.Router) {
			rg.Use(authMiddleware, apiDocAuthDoc) // Применяем middleware для аутентификации

			regRoute(rg, sr.uploadOrderUseCase.Execute, http.MethodPost, "/orders", 0,
				"Загрузка номера заказа.", "Требуется аутентификация.", "Required Auth")
			regRoute(rg, sr.getOrdersUseCase.Execute, http.MethodGet, "/orders", 0,
				"Получение списка загруженных номеров заказов.", "Требуется аутентификация.", "Required Auth")
			regRoute(rg, sr.getBalanceUseCase.Execute, http.MethodGet, "/balance", 0,
				"Получение текущего баланса.", "Требуется аутентификация.", "Required Auth")
			regRoute(rg, sr.wthdrawBalanceUseCase.Execute, http.MethodPost, "/balance/withdraw", http.StatusOK,
				"Запрос на списание баллов.", "Требуется аутентификация.", "Required Auth")
			regRoute(rg, sr.getWithdrawalsUseCase.Execute, http.MethodGet, "/withdrawals", 0,
				"Получение информации о выводе средств.", "Требуется аутентификация.", "Required Auth")
		})
	})

	return nil
}
