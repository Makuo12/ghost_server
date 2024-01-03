package api

type GetEventDateDeepLink struct {
	EventInfoID     string `json:"event_info_id" binding:"required"`
	EventDateTimeID string `json:"event_date_time_id" binding:"required"`
}

type GetDeepLinkRes struct {
	DeepLink string `json:"deep_link"`
}

type GetDeepLinkExperienceParams struct {
	DeepLinkID string `json:"deep_link_id" binding:"required"`
	MainOption string `json:"main_option" binding:"required"`
	Currency   string `json:"currency" binding:"required"`
}

type GetEventDateDeepLinkExperienceParams struct {
	DeepLinkID  string `json:"deep_link_id" binding:"required"`
	EventLinkID string `json:"event_link_id" binding:"required"`
	MainOption  string `json:"main_option" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
}

type GetEventDateDeepLinkExperienceRes struct {
	EventDateTimeID string              `json:"event_date_time_id"`
	Data            ExperienceEventData `json:"data"`
}

