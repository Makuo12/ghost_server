package payout


type PaystackRefundResponse struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    []PaystackRefundData `json:"data"`
	Meta    PaystackMeta         `json:"meta"`
}

type PaystackRefundData struct {
	Integration          int              `json:"integration"`
	Transaction          int64            `json:"transaction"`
	Dispute              any              `json:"dispute"`
	Settlement           any              `json:"settlement"`
	ID                   int              `json:"id"`
	Domain               string           `json:"domain"`
	Currency             string           `json:"currency"`
	Amount               int              `json:"amount"`
	Status               string           `json:"status"`
	RefundedAt           string           `json:"refunded_at"`
	RefundedBy           string           `json:"refunded_by"`
	CustomerNote         string           `json:"customer_note"`
	MerchantNote         string           `json:"merchant_note"`
	DeductedAmount       int              `json:"deducted_amount"`
	FullyDeducted        int              `json:"fully_deducted"`
	CreatedAt            string           `json:"createdAt"`
	BankReference        string           `json:"bank_reference"`
	TransactionReference string           `json:"transaction_reference"`
	Reason               string           `json:"reason"`
	Customer             PaystackCustomer `json:"customer"`
	RefundType           string           `json:"refund_type"`
	TransactionAmount    int              `json:"transaction_amount"`
	InitiatedBy          string           `json:"initiated_by"`
	RefundChannel        string           `json:"refund_channel"`
	SessionID            string           `json:"session_id"`
	CollectAccountNumber bool             `json:"collect_account_number"`
}

type PaystackCustomer struct {
	ID                       int    `json:"id"`
	FirstName                any    `json:"first_name"`
	LastName                 any    `json:"last_name"`
	Email                    string `json:"email"`
	CustomerCode             string `json:"customer_code"`
	Phone                    any    `json:"phone"`
	Metadata                 any    `json:"metadata"`
	RiskAction               string `json:"risk_action"`
	InternationalFormatPhone any    `json:"international_format_phone"`
}

type PaystackMeta struct {
	Total             int `json:"total"`
	Skipped           int `json:"skipped"`
	PerPage           int `json:"perPage"`
	Page              int `json:"page"`
	PageCount         int `json:"pageCount"`
	FailedRefundCount int `json:"failedRefundCount"`
}

type RefundResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Integration          int    `json:"integration"`
		Transaction          string `json:"transaction"`
		Dispute              any    `json:"dispute"`
		Settlement           any    `json:"settlement"`
		ID                   int    `json:"id"`
		Domain               string `json:"domain"`
		Currency             string `json:"currency"`
		Amount               int    `json:"amount"`
		Status               string `json:"status"`
		RefundedAt           string `json:"refunded_at"`
		RefundedBy           string `json:"refunded_by"`
		CustomerNote         string `json:"customer_note"`
		MerchantNote         string `json:"merchant_note"`
		DeductedAmount       int    `json:"deducted_amount"`
		FullyDeducted        int    `json:"fully_deducted"`
		CreatedAt            string `json:"createdAt"`
		BankReference        string `json:"bank_reference"`
		TransactionReference string `json:"transaction_reference"`
		Reason               string `json:"reason"`
		Customer             struct {
			ID                       int    `json:"id"`
			FirstName                any    `json:"first_name"`
			LastName                 any    `json:"last_name"`
			Email                    string `json:"email"`
			CustomerCode             string `json:"customer_code"`
			Phone                    any    `json:"phone"`
			Metadata                 any    `json:"metadata"`
			RiskAction               string `json:"risk_action"`
			InternationalFormatPhone any    `json:"international_format_phone"`
		} `json:"customer"`
		RefundType           string `json:"refund_type"`
		TransactionAmount    int    `json:"transaction_amount"`
		InitiatedBy          string `json:"initiated_by"`
		RefundChannel        string `json:"refund_channel"`
		SessionID            string `json:"session_id"`
		CollectAccountNumber bool   `json:"collect_account_number"`
	} `json:"data"`
}
