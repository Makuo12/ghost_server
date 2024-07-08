package api

type AmenityItem struct {
	Tag    string `json:"tag"`
	AmType string `json:"am_type"`
	HasAm  bool   `json:"has_am"`
	ID     string `json:"id"`
}

type ListUHMAmenitiesRes struct {
	List []AmenityItem `json:"list"`
}

type CreateUpdateAmenityParams struct {
	Tag      string `json:"tag" binding:"required"`
	HasAm    bool   `json:"has_am"`
	AmType   string `json:"am_type" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

// UHMAmenityDetail would be use as get response and also for update
type UHMAmenityDetailRes struct {
	OptionID           string   `json:"option_id" binding:"required"`
	ID                 string   `json:"id"`
	LocationOption     string   `json:"location_option"`
	SizeOption         int      `json:"size_option"`
	PrivacyOption      string   `json:"privacy_option"`
	TimeOption         string   `json:"time_option"`
	StartTime          string   `json:"start_time"`
	EndTime            string   `json:"end_time"`
	AvailabilityOption string   `json:"availability_option"`
	TimeSet            bool     `json:"time_set"`
	StartMonth         string   `json:"start_month"`
	EndMonth           string   `json:"end_month"`
	TypeOption         string   `json:"type_option"`
	PriceOption        string   `json:"price_option"`
	BrandOption        string   `json:"brand_option"`
	ListOptions        []string `json:"list_options"`
}

type GetAmenityDetailParams struct {
	OptionID string `json:"option_id" binding:"required"`
	AmType   string `json:"am_type" binding:"required"`
	Tag      string `json:"tag" binding:"required"`
}

//type UpdateAmenityDetailParams struct {
//	ID string `json:"id" binding:"required"`
//}

type GetOptionLocationRes struct {
	Street               string `json:"street"`
	City                 string `json:"city"`
	State                string `json:"state"`
	Country              string `json:"country"`
	Postcode             string `json:"postcode"`
	Lat                  string `json:"lat"`
	Lng                  string `json:"lng"`
	ShowSpecificLocation bool   `json:"show_specific_location"`
}

type UpdateOptionLocationParams struct {
	Street   string `json:"street" binding:"required"`
	City     string `json:"city" binding:"required"`
	State    string `json:"state" binding:"required"`
	Country  string `json:"country" binding:"required"`
	Postcode string `json:"postcode" binding:"required"`
	Lat      string `json:"lat" binding:"required"`
	Lng      string `json:"lng" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type UpdateShowSpecificLocationParams struct {
	ShowSpecificLocation bool   `json:"show_specific_location"`
	OptionID             string `json:"option_id" binding:"required"`
}

type GetOptionDetailHighlightRes struct {
	Highlight []string `json:"highlight"`
}

type UpdateOptionDetailHighlightParams struct {
	Highlight []string `json:"highlight"`
	OptionID  string   `json:"option_id" binding:"required"`
}

type CreateUpdateWifiDetailParams struct {
	NetworkName string `json:"network_name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	OptionID    string `json:"option_id" binding:"required"`
}

type GetWifiDetailRes struct {
	NetworkName string `json:"network_name"`
	Password    string `json:"password"`
}

type GetOptionExtraInfoParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Type     string `json:"type" binding:"required,option_extra_info_type"`
}

type CreateUpdateOptionExtraInfoParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Type     string `json:"type" binding:"required,option_extra_info_type"`
	Info     string `json:"info"`
}

type OptionExtraInfoRes struct {
	Info string `json:"info"`
}

type CheckInStepRes struct {
	ID    string `json:"id"`
	Des   string `json:"des"`
	Image string `json:"image"`
}

type ListCheckInStepRes struct {
	List      []CheckInStepRes `json:"list"`
	Published bool             `json:"published"`
}

type UpdateCheckInStepParams struct {
	StepID   string `json:"step_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
	Des      string `json:"des"`
	Image    string `json:"image"`
	// This would tell us whether to use des or photo to update
	Type string `json:"type" binding:"required"`
}

type CreateCheckInStepParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Des      string `json:"des"`
	Image    string `json:"image"`
	// This would tell us whether to use des or photo to update
	Type string `json:"type" binding:"required"`
}

type RemoveCheckInStepParams struct {
	StepID   string `json:"step_id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type UpdateCheckInMethodParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Des      string `json:"des"`
	Method   string `json:"method"`
	// This would tell us whether to use des or method to update
	Type string `json:"type" binding:"required"`
}

type GetShortletCheckInMethodRes struct {
	Des    string `json:"des"`
	Method string `json:"method"`
}

type ThingToNoteItem struct {
	Tag     string `json:"tag"`
	Checked bool   `json:"checked"`
	Type    string `json:"type"`
	ID      string `json:"id"`
}

type ListThingToNoteRes struct {
	List []ThingToNoteItem `json:"list"`
}

// CU means CreateUpdate
type CUThingToNoteParams struct {
	Tag      string `json:"tag" binding:"required"`
	Checked  bool   `json:"checked"`
	Type     string `json:"type" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type UThingToNoteDetailRes struct {
	ID        string `json:"id" binding:"required"`
	Des       string `json:"des" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time"`
}

type UThingToNoteDetailReq struct {
	OptionID  string `json:"option_id" binding:"required"`
	ID        string `json:"id" binding:"required"`
	Des       string `json:"des" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time"`
}

type GetThingToNoteDetailParams struct {
	OptionID string `json:"option_id" binding:"required"`
	ID       string `json:"id" binding:"required"`
}

type OptionRuleItem struct {
	Tag     string `json:"tag"`
	Checked bool   `json:"checked"`
	Type    string `json:"type"`
	ID      string `json:"id"`
	Des     string `json:"des"`
}

type ListOptionRuleRes struct {
	List []OptionRuleItem `json:"list"`
}

// CU means CreateUpdate
type CUOptionRuleParams struct {
	Tag      string `json:"tag" binding:"required"`
	Checked  bool   `json:"checked"`
	Type     string `json:"type" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type UOptionRuleDetailRes struct {
	ID        string `json:"id" binding:"required"`
	Des       string `json:"des"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type UOptionRuleDetailReq struct {
	OptionID  string `json:"option_id" binding:"required"`
	ID        string `json:"id" binding:"required"`
	Des       string `json:"des"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type GetOptionRuleDetailParams struct {
	OptionID string `json:"option_id" binding:"required"`
	ID       string `json:"id" binding:"required"`
}

type UOptionAvailabilitySettingRes struct {
	AdvanceNotice          string `json:"advance_notice"`
	AdvanceNoticeCondition string `json:"advance_notice_condition"`
	PreparationTime        string `json:"preparation_time"`
	AvailabilityWindow     string `json:"availability_window"`
}

type UOptionAvailabilitySettingReq struct {
	OptionID               string `json:"option_id" binding:"required"`
	AdvanceNotice          string `json:"advance_notice"`
	AdvanceNoticeCondition string `json:"advance_notice_condition"`
	PreparationTime        string `json:"preparation_time"`
	AvailabilityWindow     string `json:"availability_window"`
}

type UOptionTripLengthRes struct {
	MinStayDay                  int  `json:"min_stay_day" binding:"required,gte=0"`
	MaxStayNight                int  `json:"max_stay_night" binding:"required,gte=0"`
	ManualApproveRequestPassMax bool `json:"manual_approve_request_pass_max"`
	AllowReservationRequest     bool `json:"allow_reservation_request"`
}

type UOptionTripLengthReq struct {
	OptionID                    string `json:"option_id" binding:"required"`
	MinStayDay                  int    `json:"min_stay_day" binding:"required,gte=0"`
	MaxStayNight                int    `json:"max_stay_night" binding:"required,gte=0"`
	ManualApproveRequestPassMax bool   `json:"manual_approve_request_pass_max"`
	AllowReservationRequest     bool   `json:"allow_reservation_request"`
}

type UCheckInOutDetailRes struct {
	ArriveAfter            string   `json:"arrive_after"`
	ArriveBefore           string   `json:"arrive_before"`
	LeaveBefore            string   `json:"leave_before"`
	RestrictedCheckInDays  []string `json:"restricted_check_in_days"`
	RestrictedCheckOutDays []string `json:"restricted_check_out_days"`
}

type UCheckInOutDetailReq struct {
	OptionID               string   `json:"option_id" binding:"required"`
	ArriveAfter            string   `json:"arrive_after"`
	ArriveBefore           string   `json:"arrive_before"`
	LeaveBefore            string   `json:"leave_before"`
	RestrictedCheckInDays  []string `json:"restricted_check_in_days"`
	RestrictedCheckOutDays []string `json:"restricted_check_out_days"`
}

// This would be used to the update response
type UpdateCancelPolicyReq struct {
	OptionID      string `json:"option_id" binding:"required"`
	TypeOne       string `json:"type_one"`
	TypeTwo       string `json:"type_two"`
	RequestRefund bool   `json:"request_refund"`
}

type UpdateCancelPolicyRes struct {
	TypeOne       string `json:"type_one"`
	TypeTwo       string `json:"type_two"`
	RequestRefund bool   `json:"request_refund"`
}

type GetOptionBookMethodRes struct {
	InstantBook      bool   `json:"instant_book"`
	IdentityVerified bool   `json:"identity_verified"`
	GoodTrackRecord  bool   `json:"good_track_record"`
	PreBookMsg       string `json:"pre_book_msg"`
}

// This would be used to the get response and the update response
type UOptionBookMethodRes struct {
	InstantBook      bool `json:"instant_book"`
	IdentityVerified bool `json:"identity_verified"`
	GoodTrackRecord  bool `json:"good_track_record"`
}

type UOptionBookMethodMsgRes struct {
	PreBookMsg string `json:"pre_book_msg"`
}

// This would update te request
type UOptionBookMethodReq struct {
	OptionID         string `json:"option_id" binding:"required"`
	InstantBook      bool   `json:"instant_book"`
	IdentityVerified bool   `json:"identity_verified"`
	GoodTrackRecord  bool   `json:"good_track_record"`
}

type UOptionBookMethodMsgReq struct {
	OptionID   string `json:"option_id" binding:"required"`
	PreBookMsg string `json:"pre_book_msg"`
}

// This would be used to for the update response
type UBookRequirementRes struct {
	ProfilePhoto bool `json:"profile_photo"`
}

// This would be used to for the update req
type UBookRequirementReq struct {
	OptionID     string `json:"option_id" binding:"required"`
	ProfilePhoto bool   `json:"profile_photo"`
}

// This would be used to for the get response
type GetBookRequirementRes struct {
	Email        bool `json:"email"`
	PhoneNumber  bool `json:"phone_number"`
	Rules        bool `json:"rules"`
	PaymentInfo  bool `json:"payment_info"`
	ProfilePhoto bool `json:"profile_photo"`
}

type OptionCOHostItem struct {
	ID            string `json:"id" binding:"required"`
	Email         string `json:"email" binding:"required"`
	Accepted      bool   `json:"accepted"`
	FirstName     string `json:"first_name" binding:"required"`
	HostImage     string `json:"host_image" binding:"required"`
	Date          string `json:"date" binding:"required"`
	IsPrimaryHost bool   `json:"is_primary_host"`
	IsMainHost    bool   `json:"is_main_host"`
}

type ListOptionCoHostParams struct {
	OptionID string `json:"option_id"`
	Offset   int    `json:"offset"`
}

type ListOptionCOHostRes struct {
	List []OptionCOHostItem `json:"list" binding:"required"`
}

type CreateOptionCOHostParams struct {
	OptionID           string `json:"option_id" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Reservations       bool   `json:"reservations"`
	Post               bool   `json:"post"`
	ScanCode           bool   `json:"scan_code"`
	Calender           bool   `json:"calender"`
	EditOptionInfo     bool   `json:"edit_option_info"`
	Insights           bool   `json:"insights"`
	EditEventDateTimes bool   `json:"edit_event_date_times"`
	EditCoHosts        bool   `json:"edit_co_hosts"`
}

type UpdateOptionCOHostParams struct {
	ID                 string `json:"id" binding:"required"`
	OptionID           string `json:"option_id" binding:"required"`
	Email              string `json:"email" binding:"required,email"`
	Reservations       bool   `json:"reservations"`
	Post               bool   `json:"post"`
	Insights           bool   `json:"insights"`
	ScanCode           bool   `json:"scan_code"`
	Calender           bool   `json:"calender"`
	EditOptionInfo     bool   `json:"edit_option_info"`
	EditEventDateTimes bool   `json:"edit_event_date_times"`
	EditCoHosts        bool   `json:"edit_co_hosts"`
}

type GetOptionCOHostRes struct {
	Reservations       bool `json:"reservations"`
	Post               bool `json:"post"`
	ScanCode           bool `json:"scan_code"`
	Calender           bool `json:"calender"`
	Insights           bool `json:"insights"`
	EditOptionInfo     bool `json:"edit_option_info"`
	EditEventDateTimes bool `json:"edit_event_date_times"`
	EditCoHosts        bool `json:"edit_co_hosts"`
}

type GetOptionCOHostParams struct {
	ID       string `json:"id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
}

type RemoveCOHostParams struct {
	ID       string `json:"id" binding:"required"`
	OptionID string `json:"option_id" binding:"required"`
	Message  string `json:"message"`
}

type RemoveItemRes struct {
	Success bool `json:"success" binding:"required"`
}

type ResendInviteRes struct {
	Success bool `json:"success" binding:"required"`
}

// U UPDATE
type UOptionInfoStatusRes struct {
	Status          string `json:"status"`
	StatusReason    string `json:"status_reason"`
	SnoozeStartDate string `json:"snooze_start_date"`
	SnoozeEndDate   string `json:"snooze_end_date"`
	UnlistReason    string `json:"unlist_reason"`
	UnlistDes       string `json:"unlist_des"`
}

type UOptionInfoStatusReq struct {
	OptionID        string `json:"option_id" binding:"required"`
	Status          string `json:"status"`
	StatusReason    string `json:"status_reason"`
	SnoozeStartDate string `json:"snooze_start_date"`
	SnoozeEndDate   string `json:"snooze_end_date"`
	UnlistReason    string `json:"unlist_reason"`
	UnlistDes       string `json:"unlist_des"`
}

type GetOptionQuestionRes struct {
	OptionID          string   `json:"option_id" binding:"required"`
	HostAsIndividual  bool     `json:"host_as_individual"`
	OrganizationName  string   `json:"organization_name"`
	OrganizationEmail string   `json:"organization_email"`
	LegalRepresents   []string `json:"legal_represents"`
	Street            string   `json:"street"`
	State             string   `json:"state"`
	Country           string   `json:"country"`
	Postcode          string   `json:"postcode"`
	Lat               string   `json:"lat"`
	Lng               string   `json:"lng"`
}

type UpdateOptionQuestionParams struct {
	OptionID          string `json:"option_id" binding:"required"`
	HostAsIndividual  bool   `json:"host_as_individual"`
	OrganizationName  string `json:"organization_name"`
	OrganizationEmail string `json:"organization_email"`
}

type UpdateOptionQuestionRes struct {
	HostAsIndividual  bool   `json:"host_as_individual"`
	OrganizationName  string `json:"organization_name"`
	OrganizationEmail string `json:"organization_email"`
}

type AddLegalRepresentParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LegalRepresentRes struct {
	LegalRepresents []string `json:"legal_represents"`
}

type RemoveLegalRepresentParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type GetMainOptionQuestionRes struct {
	HostAsIndividual  bool   `json:"host_as_individual"`
	OrganizationName  string `json:"organization_name"`
	OrganizationEmail string `json:"organization_email"`
}
