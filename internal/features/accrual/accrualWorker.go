package accrual

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/features/services"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
)

type AccrualSystemClient interface {
	GetOrderDetails(orderNumber string) (*AccrualResult, error)
}

type AccrualWorkerImpl struct {
	orderRepo        domain.OrderRepository
	bonusCalcsRepo   domain.BonusCalculationRepository
	bonusAccountServ services.BonusAccountService
	apiClient        AccrualSystemClient
	pollingInterval  time.Duration
	ctxCancel        context.CancelFunc
}

func NewAccrualWorker(orderRepo domain.OrderRepository,
	bonusCalcsRepo domain.BonusCalculationRepository,
	bonusAccountServ services.BonusAccountService,
	apiClient AccrualSystemClient,
	pollingInterval time.Duration,
) *AccrualWorkerImpl {
	return &AccrualWorkerImpl{
		orderRepo:       orderRepo,
		bonusCalcsRepo:  bonusCalcsRepo,
		apiClient:       apiClient,
		pollingInterval: pollingInterval,
	}
}

func (s *AccrualWorkerImpl) StartPolling() {
	ctx, cancel := context.WithCancel(context.Background())
	s.ctxCancel = cancel

	go func() {
		ticker := time.NewTicker(s.pollingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Zap.Info("Stopping order polling...")
				return
			case <-ticker.C:
				s.processOrders(ctx)
			}
		}
	}()
}

func (s *AccrualWorkerImpl) Close() error {
	s.ctxCancel()
	return nil
}

func (s *AccrualWorkerImpl) processOrders(ctx context.Context) {
	orders, err := s.orderRepo.FindNotProcessed(ctx)
	if err != nil {
		log.Zap.Errorf("failed to get orders: %v", err)
		return
	}

	for _, order := range orders {
		accrualResp, err := s.getOrderInfoAccrual(ctx, order.GetNumber())
		if err != nil {
			log.Zap.Errorf("failed getOrderInfoAccrual for order %s: %v", order.GetNumber(), err)
			continue
		}

		if accrualResp == nil {
			log.Zap.Error("accrualResp is empty")
			continue
		}

		err = s.processOrder(ctx, &order, accrualResp)
		if err != nil {
			log.Zap.Errorf("failed processOrder %s: %w", order.GetNumber(), err)
			continue
		}
	}
}

func (s *AccrualWorkerImpl) getOrderInfoAccrual(ctx context.Context, orderNumber string) (*AccrualResult, error) {
	const retryLimit int = 10
	retryAttempts := 0
	for {
		if retryAttempts >= retryLimit {
			return nil, errors.New("failed to fetch order details retry limit exceeded")
		}
		retryAttempts++

		resp, err := s.apiClient.GetOrderDetails(orderNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch order details: %v", err)
		}

		if resp.isRetryRequired {
			time.Sleep(time.Duration(resp.retryTimeoutSec) * time.Second)
		} else {
			return resp, nil
		}
	}
}

func (s *AccrualWorkerImpl) processOrder(ctx context.Context, order *domain.Order, accrualResp *AccrualResult) error {
	if accrualResp.isNotRegisteredOrder {
		return errors.New("failed process order for the reason: accrual order NotRegistered")
	}

	status, err := domain.ParseLoyaltyStatus(accrualResp.orderInfo.Status)
	if err != nil {
		return fmt.Errorf("failed parsing accrual order status: %w", err)
	}

	accrual := "nil"
	if accrualResp.orderInfo != nil && accrualResp.orderInfo.Accrual != nil {
		accrual = strconv.FormatFloat(*accrualResp.orderInfo.Accrual, 'f', -1, 64)
	}

	log.Zap.Errorf("process order %s: accrualStatus=%s, accrual=%s",
		order.GetNumber(), status, accrual)

	switch status {
	case domain.LoyaltyStatusRegistered:
		err = s.handleRegisteredStatus(ctx, order)
	case domain.LoyaltyStatusProcessing:
		err = s.handleProcessingStatus(ctx, order, accrualResp)
	case domain.LoyaltyStatusInvalid:
		err = s.handleInvalidStatus(ctx, order, accrualResp)
	case domain.LoyaltyStatusProcessed:
		err = s.handleProcessedStatus(ctx, order, accrualResp)
	}

	if err != nil {
		return fmt.Errorf("failed handling status %s: %w", status, err)
	}

	return nil
}

func (s *AccrualWorkerImpl) handleRegisteredStatus(ctx context.Context, order *domain.Order) error {
	err := order.MarkProcessing()
	if err != nil {
		return fmt.Errorf("failed order update status: %w", err)
	}

	err = s.updateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed save order in DB: %w", err)
	}

	return nil
}

func (s *AccrualWorkerImpl) handleProcessingStatus(ctx context.Context, order *domain.Order,
	accrualResp *AccrualResult) error {
	if order.GetStatus() == domain.OrderStatusProcessed {
		log.Zap.Info("order %s still being processed", order.GetNumber())
		return nil
	}

	err := order.MarkProcessing()
	if err != nil {
		return fmt.Errorf("failed order update status: %w", err)
	}

	err = s.updateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed save order in DB: %w", err)
	}

	return nil
}

func (s *AccrualWorkerImpl) handleInvalidStatus(ctx context.Context, order *domain.Order,
	accrualResp *AccrualResult) error {
	err := order.MarkInvalid()
	if err != nil {
		return fmt.Errorf("failed order update status: %w", err)
	}

	err = s.updateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed save order in DB: %w", err)
	}

	return nil
}

func (s *AccrualWorkerImpl) handleProcessedStatus(ctx context.Context, order *domain.Order,
	accrualResp *AccrualResult) error {
	account, err := s.bonusAccountServ.GetBonusAccount(ctx, order.GetUserID())
	if err != nil {
		return fmt.Errorf("failed GetBonusAccount: %w", err)
	}

	bonusCalc, err := account.AddBonusCalc(order, domain.CurrencyFromMajorUnit(*accrualResp.orderInfo.Accrual))
	if err != nil {
		return fmt.Errorf("failed AddBonusCalc: %w", err)
	}

	err = order.MarkProcessed(*bonusCalc)
	if err != nil {
		return fmt.Errorf("failed order update status: %w", err)
	}

	err = s.bonusCalcsRepo.Create(ctx, bonusCalc)
	if err != nil {
		return fmt.Errorf("failed save new bonus calculation in DB: %w", err)
	}

	err = s.updateOrder(ctx, order)
	if err != nil {
		return fmt.Errorf("failed save order in DB: %w", err)
	}

	return nil
}

func (s *AccrualWorkerImpl) updateOrder(ctx context.Context, order *domain.Order) error {
	err := s.orderRepo.Update(ctx, order)
	if err != nil {
		return err
	}

	log.Zap.Info("order %s success updated, status: %s", order.GetNumber(), order.GetStatus().String())
	return nil
}
