package api

// GuestName is the name of the guest that paid
type PayoutOptionItem struct {
	ID             string `json:"id"`
	DatePaid       string `json:"date_paid"`
	Amount         string `json:"amount"`
	GuestName      string `json:"guest_name"`
	HostOptionName string `json:"host_option_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	AccountNumber  string `json:"account_number"`
	OptionUserID   string `json:"option_user_id"`
	Currency       string `json:"currency"`
}

type ListEventPayoutParams struct {
	Offset     int  `json:"offset"`
	IsComplete bool `json:"is_complete"`
}

type ListOptionPayoutParams struct {
	Offset     int  `json:"offset"`
	IsComplete bool `json:"is_complete"`
}

type ListOptionPaymentParams struct {
	Offset int `json:"offset"`
}

type ListTicketPaymentParams struct {
	Offset int `json:"offset"`
}

type PayoutEventItem struct {
	ID              string `json:"id"`
	DatePaid        string `json:"date_paid"`
	Amount          string `json:"amount"`
	HostOptionName  string `json:"host_option_name"`
	Currency        string `json:"currency"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	AccountNumber   string `json:"account_number"`
	EventDateTimeID string `json:"event_date_time_id"`
	EventDateType   string `json:"event_date_type"` // this is recurring or single event
}

func PayoutEventItemOffset(data []PayoutEventItem, offset int, limit int) []PayoutEventItem {
	// If offset is greater than or equal to the length of data, return an empty slice.
	if offset >= len(data) {
		return []PayoutEventItem{}
	}

	// Calculate the end index based on offset and limit.
	end := offset + limit

	// If the end index is greater than the length of data, set it to the length of data.
	if end > len(data) {
		end = len(data)
	}

	// Return a subset of data starting from the offset and up to the end index.
	return data[offset:end]
}

type ListPayoutEventRes struct {
	List       []PayoutEventItem `json:"list"`
	IsComplete bool              `json:"is_complete"`
}

type ListPayoutOptionRes struct {
	List       []PayoutOptionItem `json:"list"`
	IsComplete bool               `json:"is_complete"`
}

type PaymentOptionItem struct {
	ID             string `json:"id"`
	DatePaid       string `json:"date_paid"`
	Amount         string `json:"amount"`
	HostName       string `json:"host_name"`
	HostOptionName string `json:"host_option_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	Cancelled      bool   `json:"cancelled"`
	Currency       string `json:"currency"`
}

type ListPaymentOptionRes struct {
	List       []PaymentOptionItem `json:"list"`
	IsComplete bool                `json:"is_complete"`
}

type PaymentTicketItem struct {
	ID             string `json:"id"`
	DatePaid       string `json:"date_paid"`
	Amount         string `json:"amount"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	HostOptionName string `json:"host_option_name"`
	HostName       string `json:"host_name"`
	Grade          string `json:"grade"`
	Cancelled      bool   `json:"cancelled"`
	Currency       string `json:"currency"`
}

type ListPaymentTicketRes struct {
	List []PaymentTicketItem `json:"list"`
}

type ListRefundParams struct {
	Offset     int  `json:"offset"`
	IsComplete bool `json:"is_complete"`
}

type RefundItem struct {
	Amount         string `json:"amount"`
	HostOptionName string `json:"host_option_name"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	ID             string `json:"id"`
	DatePaid       string `json:"date_paid"`
	HostName       string `json:"host_name"`
	Cancelled      bool   `json:"cancelled"`
	Currency       string `json:"currency"`
}

type ListRefundRes struct {
	List       []RefundItem `json:"list"`
	IsComplete bool         `json:"is_complete"`
}

type ListRefundPayoutRes struct {
	List       []RefundPayoutItem `json:"list"`
	IsComplete bool               `json:"is_complete"`
}

type RefundPayoutItem struct {
	ID             string `json:"id"`
	DatePaid       string `json:"date_paid"`
	Amount         string `json:"amount"`
	HostOptionName string `json:"host_option_name"`
	GuestName      string `json:"guest_name"`
	AccountNumber  string `json:"account_number"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
	Cancelled      bool   `json:"cancelled"`
	Currency       string `json:"currency"`
}
