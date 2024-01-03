package api

type OptionSelectionOffsetParams struct {
	OptionOffset int    `json:"option_offset"`
	Type         string `json:"type" binding:"required"`
}

type UHMOptionSelectionRes struct {
	HostNameOption string `json:"host_name_option"`
	CoverImage     string `json:"cover_image"`
	OptionID       string `json:"option_id"`
	HasName        bool   `json:"has_name"`
	MainOptionType string `json:"main_option_type"`
	IsComplete     bool   `json:"is_complete"`
	IsActive       bool   `json:"is_active"`
	IsCoHost       bool   `json:"is_co_host"`
}

type ListUHMOptionSelectionRes struct {
	List         []UHMOptionSelectionRes `json:"list"`
	OptionOffset int                     `json:"option_offset"`
	Type         string                  `json:"type"`
	OnLastIndex  bool                    `json:"on_last_index"`
}
