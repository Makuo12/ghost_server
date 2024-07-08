package api

type DateOptionBookedItem struct {
	ID        string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	UserImage string `json:"user_image"`
	Booked    bool   `json:"booked"`
	FirstName string `json:"first_name"`
	UserID    string `json:"user_id"`
	IsEmpty   bool   `json:"is_empty"`
}

type DateEventBookedItem struct {
	FakeID    string `json:"fake_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Booked    bool   `json:"booked"`
	Count     int    `json:"count"`
}

type ListDateEventBookedRes struct {
	List []DateEventBookedItem `json:"list"`
}

type UpdateHostEventDateBookedParams struct {
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	NewStartDate string `json:"new_start_date"`
	NewEndDate   string `json:"new_end_date"`
	EventDateID  string `json:"event_date_id"`
	EventID      string `json:"event_id"`
	ReasonOne    string `json:"reason_one"`
	Message      string `json:"message"`
}

type UpdateHostEventDateBookedRes struct {
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	EventDateID string `json:"event_date_id"`
}
