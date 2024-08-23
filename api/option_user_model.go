package api

type ExperienceOptionData struct {
	UserOptionID     string `json:"user_option_id"`
	Name             string `json:"name"`
	IsVerified       bool   `json:"is_verified"`
	HostAsIndividual bool   `json:"host_as_individual"`
	BasePrice        string `json:"base_price"`
	WeekendPrice     string `json:"weekend_price"`
	// AddedPrice for when we are calculating based on more than one night
	AddedPrice     string   `json:"added_price"`
	AddPriceFound  bool     `json:"add_price_found"`
	StartDate      string   `json:"start_date"`
	EndDate        string   `json:"end_date"`
	TypeOfShortlet string   `json:"type_of_shortlet"`
	State          string   `json:"state"`
	Country        string   `json:"country"`
	Street         string   `json:"street"`
	City           string   `json:"city"`
	HostName       string   `json:"host_name"`
	HostJoined     string   `json:"host_joined"`
	HostVerified   bool     `json:"host_verified"`
	Category       string   `json:"category"`
	MainUrl        string   `json:"main_url"`
	Urls           []string `json:"urls"`
	HostUrl        string   `json:"host_url"`
}

type ListExperienceOptionRes struct {
	List         []ExperienceOptionData `json:"list"`
	OptionOffset int                    `json:"option_offset"`
	OnLastIndex  bool                   `json:"on_last_index"`
	Category     string                 `json:"category"`
}

type ExperienceEventLocation struct {
	State   string `json:"state"`
	Country string `json:"country"`
	IsEmpty bool   `json:"is_empty"`
}

type ExperienceEventData struct {
	UserOptionID      string                    `json:"user_option_id"`
	Name              string                    `json:"name"`
	IsVerified        bool                      `json:"is_verified"`
	Category          string                    `json:"category"`
	HostName          string                    `json:"host_name"`
	TicketAvailable   bool                      `json:"ticket_available"`
	SubEventType      string                    `json:"sub_event_type"`
	HostAsIndividual  bool                      `json:"host_as_individual"`
	TicketLowestPrice string                    `json:"ticket_lowest_price"`
	EventStartDate    string                    `json:"event_start_date"`
	EventEndDate      string                    `json:"event_end_date"`
	HasFreeTicket     bool                      `json:"has_free_ticket"`
	Location          []ExperienceEventLocation `json:"location"`
	HostJoined        string                    `json:"host_joined"`
	HostVerified      bool                      `json:"host_verified"`
	MainUrl           string                    `json:"main_url"`
	Urls              []string                  `json:"urls"`
	HostUrl           string                    `json:"host_url"`
}

type ListExperienceEventRes struct {
	List         []ExperienceEventData `json:"list"`
	OptionOffset int                   `json:"option_offset"`
	OnLastIndex  bool                  `json:"on_last_index"`
	Category     string                `json:"category"`
}

type ExperienceOffsetParams struct {
	OptionOffset   int    `json:"option_offset"`
	MainOptionType string `json:"main_option_type"`
	Type           string `json:"type" binding:"required"`
	Currency       string `json:"currency" binding:"required"`
	Country        string `json:"country"`
	State          string `json:"state"`
}

type ExperienceSpaceArea struct {
	Name        string   `json:"name"`
	AreaType    string   `json:"area_type"`
	SharedSpace bool     `json:"shared_space"`
	Images      []string `json:"images"`
	Beds        []string `json:"beds"`
	IsSuite     bool     `json:"is_suite"`
	IsEmpty     bool     `json:"is_empty"`
}

type ExperienceSm struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Checked bool   `json:"checked"`
	IsEmpty bool   `json:"is_empty"`
}

type ExperienceDetailCoHost struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	HostImage string `json:"host_image"`
	IsEmpty   bool   `json:"is_empty"`
}

type ExperienceHouseRules struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Checked   bool   `json:"checked"`
	Des       string `json:"des"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	IsEmpty   bool   `json:"is_empty"`
}

type ExperienceDetailLocation struct {
	Street               string `json:"street"`
	City                 string `json:"city"`
	ShowSpecificLocation bool   `json:"show_specific_location"`
	Lat                  string `json:"lat"`
	Lng                  string `json:"lng"`

	IsEmpty bool `json:"is_empty"`
}

type ExperienceDetailDes struct {
	Des                  string `json:"des"`
	GetAroundDes         string `json:"get_around_des"`
	InteractWithGuestDes string `json:"interact_with_guest_des"`
	SpaceDes             string `json:"space_des"`
	NeighborhoodDes      string `json:"neighborhood_des"`
	GuestAccessDes       string `json:"guest_access_des"`
	OtherDes             string `json:"other_des"`
	IsEmpty              bool   `json:"is_empty"`
}

type ExOptionTripLength struct {
	MinStayDays                 int  `json:"min_stay_days"`
	MaxStayDays                 int  `json:"max_stay_days"`
	ManualApproveRequestPassMax bool `json:"manual_approve_request_pass_max"`
	AllowReservationRequest     bool `json:"allow_reservation_request"`
	IsEmpty                     bool `json:"is_empty"`
}

type ExCancelPolicy struct {
	TypeOne string `json:"type_one"`
	TypeTwo string `json:"type_two"`
	IsEmpty bool   `json:"is_empty"`
}

type ExOptionDiscount struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Percent int    `json:"percent"`
	Des     string `json:"des"`
	Name    string `json:"name"`
	IsEmpty bool   `json:"is_empty"`
}
type ExOptionAddCharge struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	MainFee     string `json:"main_fee"`
	ExtraFee    string `json:"extra_fee"`
	NumOfGuests int    `json:"num_of_guests"`
	IsEmpty     bool   `json:"is_empty"`
}

type ExShortletDetail struct {
	CheckInMethod    string   `json:"check_in_method"`
	AnySpaceShared   bool     `json:"any_space_shared"`
	SpaceType        string   `json:"space_type"`
	NumOfGuests      int      `json:"num_of_guests"`
	YearBuilt        int      `json:"year_built"`
	PropertySize     int      `json:"property_size"`
	PropertySizeUnit string   `json:"property_size_unit"`
	SharedSpacesWith []string `json:"shared_spaces_with"`
	TimeZone         string   `json:"time_zone"`
	IsEmpty          bool     `json:"is_empty"`
}

type ExOptionQuestions struct {
	OrganizationName  string   `json:"organization_name"`
	OrganizationEmail string   `json:"organization_email"`
	LegalRepresents   []string `json:"legal_represents"`
	Street            string   `json:"street"`
	State             string   `json:"state"`
	City              string   `json:"city"`
	Country           string   `json:"country"`
	Postcode          string   `json:"postcode"`
	Lat               string   `json:"lat"`
	Lng               string   `json:"lng"`
	IsEmpty           bool     `json:"is_empty"`
}

type ExCheckInOutDetails struct {
	ArriveAfter            string   `json:"arrive_after"`
	ArriveBefore           string   `json:"arrive_before"`
	LeaveBefore            string   `json:"leave_before"`
	RestrictedCheckInDays  []string `json:"restricted_check_in_days"`
	RestrictedCheckOutDays []string `json:"restricted_check_out_days"`
	IsEmpty                bool     `json:"is_empty"`
}

type ExOptionAvailabilitySettings struct {
	AdvanceNotice          string `json:"advance_notice"`
	AdvanceNoticeCondition string `json:"advance_notice_condition"`
	PreparationTime        string `json:"preparation_time"`
	AvailabilityWindow     string `json:"availability_window"`
	IsEmpty                bool   `json:"is_empty"`
}

type ExOptionBookMethod struct {
	InstantBook bool   `json:"instant_book"`
	PreBookMsg  string `json:"pre_book_msg"`
	IsEmpty     bool   `json:"is_empty"`
}

type ExOptionPhotoCaptions struct {
	PhotoID string `json:"photo_id"`
	Caption string `json:"caption"`
	IsEmpty bool   `json:"is_empty"`
}

type ExperienceDetailParams struct {
	OptionUserID   string `json:"option_user_id"`
	MainOptionType string `json:"main_option_type"`
	Currency       string `json:"currency"`
}

type ExperienceDetailAmParams struct {
	OptionUserID string `json:"option_user_id"`
	Tag          string `json:"tag"`
}

type ExperienceOptionDetailRes struct {
	SpaceAreas           []string                     `json:"space_areas"`
	SpaceAreaDetail      []ExperienceSpaceArea        `json:"space_area_detail"`
	Amenities            []string                     `json:"amenities"`
	Location             ExperienceDetailLocation     `json:"location"`
	HouseRules           []ExperienceHouseRules       `json:"house_rules"`
	Notes                []ExperienceSm               `json:"notes"`
	HostLanguages        []string                     `json:"host_languages"`
	NumOfBeds            int                          `json:"num_of_beds"`
	PetsAllowed          bool                         `json:"pets_allowed"`
	CoHost               []ExperienceDetailCoHost     `json:"co_host"`
	Review               UserExReview                 `json:"review"`
	Des                  ExperienceDetailDes          `json:"des"`
	TripLength           ExOptionTripLength           `json:"trip_length"`
	CancelPolicy         ExCancelPolicy               `json:"cancel_policy"`
	Discount             []ExOptionDiscount           `json:"discount"`
	AddCharge            []ExOptionAddCharge          `json:"add_charge"`
	ShortletDetail       ExShortletDetail             `json:"shortlet_detail"`
	Question             ExOptionQuestions            `json:"question"`
	CheckInOut           ExCheckInOutDetails          `json:"check_in_out"`
	Captions             []ExOptionPhotoCaptions      `json:"captions"`
	AvailabilitySettings ExOptionAvailabilitySettings `json:"availability_settings"`
	BookMethod           ExOptionBookMethod           `json:"book_method"`
	TotalReviewCount     int                          `json:"total_review_count"`
	HostBio              string                       `json:"host_bio"`
}

type ExOptionAmDetail struct {
	TimeSet            bool     `json:"time_set"`
	LocationOption     string   `json:"location_option"`
	SizeOption         int      `json:"size_option"`
	PrivacyOption      string   `json:"privacy_option"`
	TimeOption         string   `json:"time_option"`
	StartTime          string   `json:"start_time"`
	EndTime            string   `json:"end_time"`
	AvailabilityOption string   `json:"availability_option"`
	StartMonth         string   `json:"start_month"`
	EndMonth           string   `json:"end_month"`
	TypeOption         string   `json:"type_option"`
	PriceOption        string   `json:"price_option"`
	BrandOption        string   `json:"brand_option"`
	ListOptions        []string `json:"list_options"`
}

// This is for dates that are not available
type ExBusyDate struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsEmpty   bool   `json:"is_empty"`
}

// This is for dates that have special prices
type ExPriceDate struct {
	Date    string `json:"date"`
	Price   string `json:"price"`
	IsEmpty bool   `json:"is_empty"`
}

type ExOptionDateTimeItem struct {
	Date      string `json:"date"`
	Available bool   `json:"available"`
	Price     string `json:"price"`
	IsEmpty   bool   `json:"is_empty"`
}

type ListExOptionDateTimeRes struct {
	BusyDates    []ExBusyDate  `json:"busy_dates"`
	PriceDates   []ExPriceDate `json:"price_dates"`
	BusyIsEmpty  bool          `json:"busy_is_empty"`
	PriceIsEmpty bool          `json:"price_is_empty"`
}

type ExEventDateTimes struct {
	Name      string                   `json:"name"`
	StartDate string                   `json:"start_date"`
	EndDate   string                   `json:"end_date"`
	Type      string                   `json:"type"`
	StartTime string                   `json:"start_time"`
	EndTime   string                   `json:"end_time"`
	MainID    string                   `json:"main_id"`
	Timezone  string                   `json:"timezone"`
	Location  ExEventDateTimesLocation `json:"location"`
	RandomID  string                   `json:"random_id"`
	Status    string                   `json:"status"`
	IsEmpty   bool                     `json:"is_empty"`
}

type ExEventDateTimesLocation struct {
	State   string `json:"state"`
	Country string `json:"country"`
	Street  string `json:"street"`
	City    string `json:"city"`
	Lat     string `json:"lat"`
	Lng     string `json:"lng"`
	IsEmpty bool   `json:"is_empty"`
}

type ExperienceEventDetailRes struct {
	HostLanguages    []string                 `json:"host_languages"`
	CoHost           []ExperienceDetailCoHost `json:"co_host"`
	Des              ExperienceDetailDes      `json:"des"`
	CancelPolicy     ExCancelPolicy           `json:"cancel_policy"`
	Captions         []ExOptionPhotoCaptions  `json:"captions"`
	Review           UserExReview             `json:"review"`
	BookMethod       ExOptionBookMethod       `json:"book_method"`
	EventDateTimes   []ExEventDateTimes       `json:"event_date_times"`
	Question         ExOptionQuestions        `json:"question"`
	TotalReviewCount int                      `json:"total_review_count"`
	HostBio          string                   `json:"host_bio"`
}

type ExperienceCategoryRes struct {
	Category string `json:"category"`
}

type ExEventTicketData struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Price            string `json:"price"`
	AbsorbFees       bool   `json:"absorb_fees"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Level            string `json:"level"`
	TicketType       string `json:"ticket_type"`
	NumOfSeats       int    `json:"num_of_seats"`
	FreeRefreshments bool   `json:"free_refreshments"`
	IsEmpty          bool   `json:"is_empty"`
}

type CreateReportOptionParams struct {
	OptionUserID string `json:"option_user_id" binding:"required"`
	TypeOne      string `json:"type_one" binding:"required,report_option"`
	TypeTwo      string `json:"type_two"`
	TypeThree    string `json:"type_three"`
	Description  string `json:"description"`
}

type ListExperienceEventTicketsParams struct {
	EventDateTimeID string `json:"event_date_time_id"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	Currency        string `json:"currency"`
}

type ListExperienceEventTicketsRes struct {
	List            []ExEventTicketData `json:"list"`
	EventDateTimeID string              `json:"event_date_time_id"`
	IsEmpty         bool                `json:"is_empty"`
}

type UserExReviewItem struct {
	ID               string `json:"id"`
	General          string `json:"general"`
	Environment      string `json:"environment"`
	Accuracy         string `json:"accuracy"`
	CheckIn          string `json:"check_in"`
	Communication    string `json:"communication"`
	Location         string `json:"location"`
	PublicNote       string `json:"public_note"`
	HostPublicNote   string `json:"host_public_note"`
	Average          string `json:"average"`
	YearJoined       string `json:"year_joined"`
	DateBooked       string `json:"date_booked"`
	DateHostResponse string `json:"date_host_response"`
	HostImage        string `json:"host_image"`
	FirstName        string `json:"first_name"`
}

type UserExReview struct {
	Five          int                `json:"five"`
	Four          int                `json:"four"`
	Three         int                `json:"three"`
	Two           int                `json:"two"`
	One           int                `json:"one"`
	Total         int                `json:"total"`
	Average       string             `json:"average"`
	Count         int                `json:"count"`
	General       string             `json:"general"`
	Environment   string             `json:"environment"`
	Communication string             `json:"communication"`
	Accuracy      string             `json:"accuracy"`
	CheckIn       string             `json:"check_in"`
	Location      string             `json:"location"`
	List          []UserExReviewItem `json:"list"`
	IsEmpty       bool               `json:"is_empty"`
}

type ListExReviewDetailRes struct {
	List        []UserExReviewItem `json:"list"`
	Offset      int                `json:"offset"`
	OnLastIndex bool               `json:"on_last_index"`
}

type ListExReviewDetailReq struct {
	OptionUserID string `json:"option_user_id"`
	Offset       int    `json:"offset"`
	MainOption   string `json:"main_option"`
}
