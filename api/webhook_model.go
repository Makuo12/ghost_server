package api

import "strings"

type WebhookEvent struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
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

type RefundEvent struct {
	Event string `json:"event"`
	Data  struct {
		Status               string `json:"status"`
		TransactionReference string `json:"transaction_reference"`
		RefundReference      any    `json:"refund_reference"`
		Amount               int    `json:"amount"`
		Currency             string `json:"currency"`
		Processor            string `json:"processor"`
		Customer             struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
		} `json:"customer"`
		Integration int    `json:"integration"`
		Domain      string `json:"domain"`
	} `json:"data"`
}

func IsRefund(input string) bool {
	lowerInput := strings.ToLower(input)
	return strings.Contains(lowerInput, "refund")
}

func IsTransfer(input string) bool {
	lowerInput := strings.ToLower(input)
	return strings.Contains(lowerInput, "transfer")
}
