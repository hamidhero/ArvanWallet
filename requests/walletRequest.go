package requests

type AddTransactionRequest struct {
	Mobile int64  `json:"mobile"`
	Amount int64  `json:"amount"`
	Reason string `json:"reason"`
}
