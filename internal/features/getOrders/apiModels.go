package getOrders

type OrderResponse struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"` // Enum: "NEW", "PROCESSING", "INVALID", "PROCESSED"
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"` // Format: date-time
}

type OrderListResponse []OrderResponse
