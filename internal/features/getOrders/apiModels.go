package getOrders

import "time"

type OrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"` // Enum: "NEW", "PROCESSING", "INVALID", "PROCESSED"
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"` // Format: date-time
}

type OrderListResponse []OrderResponse

func (r *OrderListResponse) NoContent() bool {
	return len(*r) == 0
}

func (r *OrderListResponse) HTTPStatus() int {
	if len(*r) == 0 {
		return 0
	}
	return 200
}

func (r *OrderListResponse) ExpectedHTTPStatuses() []int {
	return []int{}
}
