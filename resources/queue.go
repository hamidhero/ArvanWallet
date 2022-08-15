package resources

type ReceiveVoucherResource struct {
	Name   string      `json:"name"`
	Amount int64       `json:"amount"`
	Data   interface{} `json:"data"`
}
