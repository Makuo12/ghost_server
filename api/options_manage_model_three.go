package api

type OptionCoHostItem struct {
	ID             string `json:"id"`
	MainHostName   string `json:"main_host_name"`
	MainOptionName string `json:"main_option_name"`
	OptionCoHostID string `json:"option_co_host_id"`
	MainImage      string `json:"main_image"`
	MainOption     string `json:"main_option"`
	HostImage      string `json:"host_image"`
	IsPrimaryHost  bool   `json:"is_primary_host"`
}

type ListOptionCoHostItemRes struct {
	List []OptionCoHostItem `json:"list"`
}

type ListOptionCoHostItemParams struct {
	Offset int `json:"offset"`
}

type DeactivateOptionCOHostParams struct {
	ID string `json:"id"`
}

type OptionCoHostItemDetailParams struct {
	ID string `json:"id"`
}

type DeactivateOptionCOHostRes struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type OptionCoHostItemDetailRes struct {
	Reservations        bool `json:"reservations"`
	Post                bool `json:"post"`
	ScanCode            bool `json:"scan_code"`
	Calender            bool `json:"calender"`
	EditOptionInfo      bool `json:"edit_option_info"`
	EditEventDatesTimes bool `json:"edit_event_dates_times"`
	EditCoHosts         bool `json:"edit_co_hosts"`
	Insights            bool `json:"insights"`
}

type ValidateOptionCOHostParams struct {
	Code string `json:"code" binding:"required,len=6"`
}

type UpdatePublishCheckInStepParams struct {
	OptionID string `json:"option_id" binding:"required"`
}

type UpdatePublishCheckInStepRes struct {
	Published bool `json:"published" binding:"required"`
}
