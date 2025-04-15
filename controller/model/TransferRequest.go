package model

type TransferRequest struct {
	PrivateKey   string `json:"privateKey"`
	ToAddress    string `json:"toAddress"`
	TokenAddress string `json:"tokenAddress"`
	Amount       int64  `json:"amount"`
}
