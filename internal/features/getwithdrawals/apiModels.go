package getwithdrawals

import "time"

type WithdrawalResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"` // Format: date-time
}

type WithdrawalListResponse []WithdrawalResponse

func (r *WithdrawalListResponse) NoContent() bool {
	return len(*r) == 0
}

func (r *WithdrawalListResponse) HTTPStatus() int {
	if len(*r) == 0 {
		return 0
	}
	return 200
}

func (r *WithdrawalListResponse) ExpectedHTTPStatuses() []int {
	return []int{}
}
