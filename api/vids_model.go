package api

type ListVidParams struct {
	Offset int `json:"offset"`
}

type CreateVidParams struct {
	Path           string `json:"path"`
	Filter         string `json:"filter"`
	OptionUserID   string `json:"option_user_id"`
	MainOptionType string `json:"main_option_type"`
	StartDate      string `json:"start_date"`
	Caption        string `json:"caption"`
	ExtraOptionID  string `json:"extra_option_id"`
}

type VidItem struct {
	Path           string `json:"path"`
	Filter         string `json:"filter"`
	OptionUserID   string `json:"option_user_id"`
	MainOptionType string `json:"main_option_type"`
	StartDate      string `json:"start_date"`
	Caption        string `json:"caption"`
	ExtraOptionID  string `json:"extra_option_id"`
}

type ListVidRes struct {
	List []VidItem `json:"list"`
}
