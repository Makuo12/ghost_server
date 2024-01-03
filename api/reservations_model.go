package api

type ReserveHostItem struct {
	UserID      string `json:"user_id"`
	ReferenceID string `json:"reference_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	OptionID    string `json:"option_id"`
	// Host method tells if it is the host actual option or the host is co-hosting the place
	HostMethod     string `json:"host_method"`
	FirstName      string `json:"first_name"`
	HostNameOption string `json:"host_name_option"`
	UserPhoto      string `json:"user_photo"`
	ArriveAfter    string `json:"arrive_after"`
	ArriveBefore   string `json:"arrive_before"`
	LeaveBefore    string `json:"leave_before"`
	TimeZone       string `json:"time_zone"`
	CanReserve     bool   `json:"can_reserve"`
	CanScanCode    bool   `json:"can_scan_code"`
	StartTimeType  string `json:"start_time_type"`
	OptionStatus   string `json:"option_status"`
	CoverImage     string `json:"cover_image"`
}

type ListReservationDetailRes struct {
	List      []ReserveHostItem `json:"list"`
	Selection string            `json:"selection"`
}

type ListReservationDetailParams struct {
	Selection      string `json:"selection"`
	MainOptionType string `json:"main_option_type"`
	Offset         int    `json:"offset"`
}

type GetReserveHostDetailRes struct {
	IdentityVerified bool   `json:"identity_verified"`
	HasNumber        bool   `json:"has_number"`
	HasProfilePhoto  bool   `json:"has_profile_photo"`
	HasEmail         bool   `json:"has_email"`
	HasFirstName     bool   `json:"has_first_name"`
	HasPayout        bool   `json:"has_payout"`
	HasLanguage      bool   `json:"has_language"`
	HasBio           bool   `json:"has_bio"`
	UserID           string `json:"user_id"`
}

type TicketHostItem struct {
	Grade        string `json:"grade"`
	CapacityBook int    `json:"capacity_book"`
	Capacity     int    `json:"capacity"`
	TicketType   string `json:"ticket_type"`
	Type         string `json:"type"`
	IsEmpty      bool   `json:"is_empty"`
}

type DateHostItem struct {
	HostNameOption    string           `json:"host_name_option"`
	StartDate         string           `json:"start_date"`
	EndDate           string           `json:"end_date"`
	StartTime         string           `json:"start_time"`
	EndTime           string           `json:"end_time"`
	Tickets           []TicketHostItem `json:"tickets"`
	Status            string           `json:"status"`
	Timezone          string           `json:"timezone"`
	HostMethod        string           `json:"host_method"`
	EventDateTimeType string           `json:"event_date_time_type"`
	EventID           string           `json:"event_id"`
	EventDateTimeID   string           `json:"event_date_time_id"`
	CanReserve        bool             `json:"can_reserve"`
	CanScanCode       bool             `json:"can_scan_code"`
	OptionStatus      string           `json:"option_status"`
	CoverImage        string           `json:"cover_image"`
}

type ReserveEventHostItem struct {
	List      []DateHostItem `json:"list"`
	Selection string         `json:"selection"`
}

func ConcatSlicesReserveItem(slice1, slice2 []ReserveHostItem) []ReserveHostItem {
	result := make([]ReserveHostItem, len(slice1)+len(slice2))
	copy(result, slice1)
	copy(result[len(slice1):], slice2)
	return result
}

func ConcatSlicesDateItem(slice1, slice2 []DateHostItem) []DateHostItem {
	result := make([]DateHostItem, len(slice1)+len(slice2))
	copy(result, slice1)
	copy(result[len(slice1):], slice2)
	return result
}

// Charts
type OptionChartItem struct {
	//  OptionID is OptionUserID DB
	OptionID string   `json:"option_id"`
	Name     string   `json:"name"`
	Dates    []string `json:"dates"`
}

type DateItem struct {
	StartDate string `json:"start_date"`
	Count     int    `json:"count"`
}

type EventItem struct {
	//  OptionID is OptionUserID in DB
	OptionID string     `json:"option_id"`
	Name     string     `json:"name"`
	Dates    []DateItem `json:"dates"`
}

type ListEventChart struct {
	List []EventItem `json:"list"`
	// event_date this is the date the host first hosted his first event
	EventDate string `json:"event_date"`
}

type ListOptionChart struct {
	List []OptionChartItem `json:"list"`
	// option_date this is the date the host first hosted his first shortlet
	OptionDate string `json:"option_date"`
}

type ListChartDataRes struct {
	Option ListOptionChart `json:"option"`
	Event  ListEventChart  `json:"event"`
}
