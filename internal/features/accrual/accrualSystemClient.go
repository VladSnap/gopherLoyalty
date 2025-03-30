package accrual

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
)

type AccrualResult struct {
	orderInfo            *AccrualResponse
	isNotRegisteredOrder bool
	isRetryRequired      bool
	retryTimeoutSec      int
}

// Response from external API
type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

type AccrualSystemClientImpl struct {
	client  *http.Client
	baseURL string
}

func NewAccrualSystemClient(baseURL string) *AccrualSystemClientImpl {
	return &AccrualSystemClientImpl{
		client:  &http.Client{},
		baseURL: baseURL,
	}
}

func (c *AccrualSystemClientImpl) GetOrderDetails(orderNumber string) (*AccrualResult, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, orderNumber)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	// Чтение тела ответа в буфер
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Zap.Info("AccrualAPI GetOrderDetails %s, status: %s", orderNumber, resp.StatusCode)
	log.Zap.Info("Body: %s", string(body))

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		seconds, err := strconv.Atoi(retryAfter)
		if err != nil {
			return nil, errors.New("invalid Retry-After header")
		}

		return &AccrualResult{isRetryRequired: true, retryTimeoutSec: seconds}, nil
	}

	if resp.StatusCode == http.StatusNoContent {
		return &AccrualResult{isNotRegisteredOrder: true}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var accrualResp AccrualResponse
	if err := json.Unmarshal(body, &accrualResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &AccrualResult{orderInfo: &accrualResp}, nil
}
