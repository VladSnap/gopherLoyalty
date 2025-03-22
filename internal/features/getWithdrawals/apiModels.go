package getWithdrawals

type WithdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"` // Format: date-time
}

type WithdrawalListResponse []WithdrawalResponse