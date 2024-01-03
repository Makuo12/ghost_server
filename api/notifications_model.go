package api

type NotificationItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Header    string `json:"header"`
	Message   string `json:"message"`
	Handled   bool   `json:"handled"`
	CreatedAt string `json:"created_at"`
}

type ListNotificationRes struct {
	List        []NotificationItem `json:"list"`
	Offset      int                `json:"offset"`
	OnLastIndex bool               `json:"on_last_index"`
	UserID      string             `json:"user_id"`
	Time        string             `json:"time"`
}

type ListNotificationParams struct {
	Offset int `json:"offset"`
}

type GetNotificationDetailParams struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type NotificationOptionReserveDetailParams struct {
	Currency string `json:"currency"`
	ID       string `json:"id"`
}

type NotificationOptionReserveDetailRes struct {
	ReserveData   ExperienceReserveOModel   `json:"reserve_data"`
	DefaultCardID string                    `json:"default_card_id"`
	HasCard       bool                      `json:"has_card"`
	CardDetail    CardDetailResponse        `json:"card_detail"`
	OptionDetail  ExperienceOptionDetailRes `json:"option_detail"`
	Option        ExperienceOptionData      `json:"option"`
}

type ListNotificationListenRes struct {
	List        []NotificationItem `json:"list"`
	CurrentTime string             `json:"current_time"`
}


