package api

import "github.com/makuo12/ghost_server/tools"

type Bank struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Code        string      `json:"code"`
	LongCode    string      `json:"longcode"`
	Gateway     string      `json:"gateway"`
	PayWithBank bool        `json:"pay_with_bank"`
	Active      bool        `json:"active"`
	Country     string      `json:"country"`
	Currency    string      `json:"currency"`
	Type        string      `json:"type"`
	IsDeleted   bool        `json:"is_deleted"`
	CreatedAt   interface{} `json:"createdAt"`
	UpdatedAt   interface{} `json:"updatedAt"`
}

type TransferRecipientRes struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Active        bool   `json:"active"`
		CreatedAt     string `json:"createdAt"`
		Currency      string `json:"currency"`
		Domain        string `json:"domain"`
		ID            int    `json:"id"`
		Integration   int    `json:"integration"`
		Name          string `json:"name"`
		RecipientCode string `json:"recipient_code"`
		Type          string `json:"type"`
		UpdatedAt     string `json:"updatedAt"`
		IsDeleted     bool   `json:"is_deleted"`
		Details       struct {
			AuthorizationCode any    `json:"authorization_code"`
			AccountNumber     string `json:"account_number"`
			AccountName       string `json:"account_name"`
			BankCode          string `json:"bank_code"`
			BankName          string `json:"bank_name"`
		} `json:"details"`
	} `json:"data"`
}

type TransferRecipientParams struct {
	Type          string `json:"type"`
	Name          string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	Currency      string `json:"currency"`
}

type BankData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []Bank `json:"data"`
}

type AccountData struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    AccountInfo `json:"data"`
}

type AccountInfo struct {
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	BankID        int    `json:"bank_id"`
}

type BankItem struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

type ListUSSDRes struct {
	List []tools.USSDItem `json:"list"`
}



type ListBankRes struct {
	List    []BankItem `json:"list"`
	Country string     `json:"country"`
}

type ListBankParams struct {
	Country string `json:"country"`
}

type CreateAccountNumberParams struct {
	Code          string `json:"code"`
	AccountNumber string `json:"account_number"`
	Country       string `json:"country"`
	Currency      string `json:"currency"`
	BankName      string `json:"bank_name"`
}

type AccountNumberItem struct {
	AccountNumber string `json:"account_number"`
	ID            string `json:"id"`
	BankName      string `json:"bank_name"`
	Currency      string `json:"currency"`
	AccountName   string `json:"account_name"`
}

type ListAccountNumberRes struct {
	List             []AccountNumberItem `json:"list"`
	IsEmpty          bool                `json:"is_empty"`
	DefaultAccountID string              `json:"default_account_id"`
}

type SetDefaultAccountNumberParams struct {
	AccountID string `json:"account_id"`
}

type SetDefaultAccountNumberRes struct {
	Success   bool   `json:"success"`
	AccountID string `json:"account_id"`
}

type RemoveAccountNumberParams struct {
	AccountID string `json:"account_id"`
}

type RemoveAccountNumberRes struct {
	Success   bool   `json:"success"`
	AccountID string `json:"account_id"`
}
