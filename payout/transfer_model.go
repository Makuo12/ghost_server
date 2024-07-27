package payout

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

type PaystackTransferResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    PaystackTransferData `json:"data"`
}

type PaystackTransferData struct {
	Amount        int                   `json:"amount"`
	CreatedAt     string                `json:"createdAt"`
	Currency      string                `json:"currency"`
	Domain        string                `json:"domain"`
	Failures      any                   `json:"failures"`
	ID            int                   `json:"id"`
	Integration   int                   `json:"integration"`
	Reason        string                `json:"reason"`
	Reference     string                `json:"reference"`
	Source        string                `json:"source"`
	SourceDetails any                   `json:"source_details"`
	Status        string                `json:"status"`
	TitanCode     string                `json:"titan_code"`
	TransferCode  string                `json:"transfer_code"`
	Request       int                   `json:"request"`
	TransferredAt string                `json:"transferred_at"`
	UpdatedAt     string                `json:"updatedAt"`
	Recipient     PaystackRecipientData `json:"recipient"`
	Session       PaystackSessionData   `json:"session"`
	FeeCharged    int                   `json:"fee_charged"`
	FeesBreakdown any                   `json:"fees_breakdown"`
}

type PaystackRecipientData struct {
	Active        bool                     `json:"active"`
	CreatedAt     string                   `json:"createdAt"`
	Currency      string                   `json:"currency"`
	Description   string                   `json:"description"`
	Domain        string                   `json:"domain"`
	Email         string                   `json:"email"`
	ID            int                      `json:"id"`
	Integration   int                      `json:"integration"`
	Metadata      any                      `json:"metadata"`
	Name          string                   `json:"name"`
	RecipientCode string                   `json:"recipient_code"`
	Type          string                   `json:"type"`
	UpdatedAt     string                   `json:"updatedAt"`
	IsDeleted     bool                     `json:"is_deleted"`
	IsDeleted2    bool                     `json:"isDeleted"`
	Details       PaystackRecipientDetails `json:"details"`
}

type PaystackRecipientDetails struct {
	AuthorizationCode string `json:"authorization_code"`
	AccountNumber     string `json:"account_number"`
	AccountName       string `json:"account_name"`
	BankCode          string `json:"bank_code"`
	BankName          string `json:"bank_name"`
}

type PaystackSessionData struct {
	Provider string `json:"provider"`
	ID       string `json:"id"`
}

type TransferResponse struct {
	Event string `json:"event"`
	Data  struct {
		Amount        int    `json:"amount"`
		Currency      string `json:"currency"`
		Domain        string `json:"domain"`
		Failures      any    `json:"failures"`
		ID            any    `json:"id"`
		Integration   any    `json:"integration"`
		Reason        string `json:"reason"`
		Reference     string `json:"reference"`
		Source        string `json:"source"`
		SourceDetails any    `json:"source_details"`
		Status        string `json:"status"`
		TitanCode     any    `json:"titan_code"`
		TransferCode  string `json:"transfer_code"`
		TransferredAt any    `json:"transferred_at"`
		Recipient     struct {
			Active        bool   `json:"active"`
			Currency      string `json:"currency"`
			Description   any    `json:"description"`
			Domain        string `json:"domain"`
			Email         string `json:"email"`
			ID            int    `json:"id"`
			Integration   int    `json:"integration"`
			Metadata      any    `json:"metadata"`
			Name          string `json:"name"`
			RecipientCode string `json:"recipient_code"`
			Type          string `json:"type"`
			IsDeleted     bool   `json:"is_deleted"`
			Details       struct {
				AuthorizationCode any    `json:"authorization_code"`
				AccountNumber     string `json:"account_number"`
				AccountName       any    `json:"account_name"`
				BankCode          string `json:"bank_code"`
				BankName          string `json:"bank_name"`
			} `json:"details"`
			CreatedAt any `json:"created_at"`
			UpdatedAt any `json:"updated_at"`
		} `json:"recipient"`
		Session struct {
			Provider string `json:"provider"`
			ID       string `json:"id"`
		} `json:"session"`
		CreatedAt any `json:"created_at"`
		UpdatedAt any `json:"updated_at"`
	} `json:"data"`
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
