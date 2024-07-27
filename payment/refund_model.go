package payment




type PaystackFullRefundParams struct {
	Transaction string `json:"transaction"`
}

type PaystackPartialRefundParams struct {
	Transaction string `json:"transaction"`
	Amount      int    `json:"amount"`
}

type PaystackCreateRefundData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Transaction struct {
			ID            int64  `json:"id"`
			Domain        string `json:"domain"`
			Reference     string `json:"reference"`
			ReceiptNumber any    `json:"receipt_number"`
			Amount        int    `json:"amount"`
			PaidAt        string `json:"paid_at"`
			Channel       string `json:"channel"`
			Currency      string `json:"currency"`
			Authorization struct {
				ExpMonth    any `json:"exp_month"`
				ExpYear     any `json:"exp_year"`
				AccountName any `json:"account_name"`
			} `json:"authorization"`
			Customer struct {
				InternationalFormatPhone any `json:"international_format_phone"`
			} `json:"customer"`
			Plan struct {
			} `json:"plan"`
			Subaccount struct {
				Currency any `json:"currency"`
			} `json:"subaccount"`
			Split struct {
			} `json:"split"`
			OrderID            any    `json:"order_id"`
			PaidAt0            string `json:"paidAt"`
			PosTransactionData any    `json:"pos_transaction_data"`
			Source             any    `json:"source"`
			FeesBreakdown      any    `json:"fees_breakdown"`
		} `json:"transaction"`
		Integration    int    `json:"integration"`
		DeductedAmount int    `json:"deducted_amount"`
		Channel        any    `json:"channel"`
		MerchantNote   string `json:"merchant_note"`
		CustomerNote   string `json:"customer_note"`
		Status         string `json:"status"`
		RefundedBy     string `json:"refunded_by"`
		ExpectedAt     string `json:"expected_at"`
		Currency       string `json:"currency"`
		Domain         string `json:"domain"`
		Amount         int    `json:"amount"`
		FullyDeducted  bool   `json:"fully_deducted"`
		ID             int    `json:"id"`
		CreatedAt      string `json:"createdAt"`
		UpdatedAt      string `json:"updatedAt"`
	} `json:"data"`
}