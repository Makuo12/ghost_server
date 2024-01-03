package api

type ListOptionInsightParams struct {
	MainOption string `json:"main_option"`
	Offset     int    `json:"offset"`
}

type OptionInsightItem struct {
	HostNameOption string `json:"host_name_option"`
	CoverImage     string `json:"cover_image"`
	OptionUserID   string `json:"option_user_id"`
	HasName        bool   `json:"has_name"`
	MainOptionType string `json:"main_option_type"`
	IsCoHost       bool   `json:"is_co_host"`
}

type ListOptionInsightRes struct {
	List       []OptionInsightItem `json:"list"`
	MainOption string              `json:"main_option"`
}

type GetOptionInsightItem struct {
	Month   string `json:"month"`
	Count   int    `json:"count"`
	Earning string `json:"earning"`
}

type GetOptionInsightParams struct {
	OptionUserID string `json:"option_user_id"`
	Currency     string `json:"currency"`
	Year         int    `json:"year"`
	Month        string `json:"month"`
	FromMonth    bool   `json:"from_month"`
}

type GetEventInsightParams struct {
	OptionUserID    string `json:"option_user_id"`
	Offset          int    `json:"offset"`
	ForOffset       bool   `json:"for_offset"`
	Currency        string `json:"currency"`
	EventDateTimeID string `json:"event_date_time_id"`
	StartDate       string `json:"start_date"`
}

type GetAllEventInsightParams struct {
	Offset          int    `json:"offset"`
	ForOffset       bool   `json:"for_offset"`
	Currency        string `json:"currency"`
	EventDateTimeID string `json:"event_date_time_id"`
	StartDate       string `json:"start_date"`
}

type GetAllOptionInsightParams struct {
	Currency  string `json:"currency"`
	Year      int    `json:"year"`
	Month     string `json:"month"`
	FromMonth bool   `json:"from_month"`
}

type GetOptionInsightRes struct {
	Earning      string                 `json:"earning"`
	Currency     string                 `json:"currency"`
	Booking      int                    `json:"booking"`
	StartYear    int                    `json:"start_year"`
	Cancellation int                    `json:"cancellation"`
	Month        string                 `json:"month"`
	List         []GetOptionInsightItem `json:"list"`
}

type GetEventInsightItem struct {
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	Name            string `json:"name"`
	EventDateTimeID string `json:"event_date_time_id"`
	Count           int    `json:"count"`
	Earnings        string `json:"earnings"`
	FakeID          string `json:"fake_id"`
}



type GetEventInsightRes struct {
	Earning         string                `json:"earning"`
	Currency        string                `json:"currency"`
	TicketSold      int                   `json:"ticket_sold"`
	StartYear       int                   `json:"start_year"`
	ForOffset       bool                  `json:"for_offset"`
	Cancellation    int                   `json:"cancellation"`
	List            []GetEventInsightItem `json:"list"`
	IsEmpty         bool                  `json:"is_empty"`
	StartDate       string                `json:"start_date"`
	EventDateTimeID string                `json:"event_date_time_id"`
}

