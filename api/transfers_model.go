package api

import "github.com/google/uuid"

type PayoutData struct {
	ChargeIDs            []uuid.UUID
	HostID               uuid.UUID
	HostUserID           uuid.UUID
	HostDefaultAccountID string
	HostName             string
	Amounts              []float64
	PayoutID             string
}

type TransferQueueResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []struct {
		Reference    string `json:"reference"`
		Recipient    string `json:"recipient"`
		Amount       int    `json:"amount"`
		TransferCode string `json:"transfer_code"`
		Currency     string `json:"currency"`
		Status       string `json:"status"`
	} `json:"data"`
}

type Transfer struct {
	Amount    int    `json:"amount"`
	Reference string `json:"reference"`
	Reason    string `json:"reason"`
	Recipient string `json:"recipient"`
}

type BulkTransferRequest struct {
	Currency  string     `json:"currency"`
	Source    string     `json:"source"`
	Transfers []Transfer `json:"transfers"`
}

