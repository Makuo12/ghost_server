package api

type GetUHMDataOptionRes struct {
	HostNameOption string   `json:"host_name_option"`
	SpaceAreas     []string `json:"space_areas"`
	SpaceType      string   `json:"space_type"`
	Category       string   `json:"category"`
	CategoryTwo    string   `json:"category_two"`
	CategoryThree  string   `json:"category_three"`
	CategoryFour   string   `json:"category_four"`
	NumOfGuest     int      `json:"num_of_guest"`
	MainImage      string   `json:"main_photo"`
	Images         []string `json:"images"`
	Price          string   `json:"price"`
	//OptionUserID automatic id we generate when creating a option because we want normal users to use this id
	OptionUserID  string `json:"option_user_id"`
	Street        string `json:"street"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Postcode      string `json:"postcode"`
	CheckInMethod string `json:"check_in_method"`
	EventType     string `json:"event_type"`
	EventSubType  string `json:"event_sub_type"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
}

type ListSpaceAreasParams struct {
	OptionID string `json:"option_id"`
}

type SpaceAreas struct {
	ID          string   `json:"id"`
	OptionID    string   `json:"option_id"`
	SharedSpace bool     `json:"shared_space"`
	SpaceType   string   `json:"space_type"`
	Image      []string `json:"images"`
	Beds        []string `json:"beds"`
	IsSuite     bool     `json:"is_suite"`
	Name        string   `json:"name"`
}

type ListSpaceAreas struct {
	List            []SpaceAreas `json:"list"`
	SharedSpaceWith []string     `json:"shared_space_with"`
}

type CreateEditSpaceAreaParams struct {
	OptionID      string   `json:"option_id"`
	Space         []string `json:"space"`
	RemoveSpaceID []string `json:"remove_space_id"`
}

type AddBedSpaceAreasParams struct {
	SpaceAreaID string   `json:"space_area_id" binding:"required"`
	OptionID    string   `json:"option_id"`
	Beds        []string `json:"beds"`
	Name        string   `json:"name" binding:"required"`
}

type AddPhotoSpaceAreasParams struct {
	SpaceAreaID string   `json:"space_area_id" binding:"required"`
	OptionID    string   `json:"option_id"`
	Photos      []string `json:"photos"`
	Name        string   `json:"name" binding:"required"`
}

type UpdateSpaceAreaParams struct {
	OptionID        string       `json:"option_id"`
	SpaceAreas      []SpaceAreas `json:"space_areas" binding:"required"`
	SharedSpaceWith []string     `json:"shared_space_with"`
}

type UpdateShortletInfoParams struct {
	SpaceType        string `json:"space_type" binding:"required"`
	TypeOfShortlet   string `json:"type_of_shortlet" binding:"required"`
	GuestWelcomed    int    `json:"guest_welcomed" binding:"required"`
	YearBuilt        int    `json:"year_built"`
	PropertySize     int    `json:"property_size"`
	PropertySizeUnit string `json:"property_size_unit" binding:"required,property_size_unit"`
	OptionID         string `json:"option_id" binding:"required"`
}

type ShortletInfoRes struct {
	SpaceType        string `json:"space_type"`
	TypeOfShortlet   string `json:"type_of_shortlet"`
	GuestWelcomed    int    `json:"guest_welcomed"`
	YearBuilt        int    `json:"year_built"`
	PropertySize     int    `json:"property_size"`
	PropertySizeUnit string `json:"property_size_unit" binding:"required,property_size_unit"`
}

type UpdateEventInfoParams struct {
	SubCategoryType string `json:"sub_category_type" binding:"required"`
	EventType       string `json:"event_type" binding:"required"`
	OptionID        string `json:"option_id" binding:"required"`
}

type EventInfoRes struct {
	SubCategoryType string `json:"sub_category_type"`
	EventType       string `json:"event_type"`
}

type UpdateOptionTitleParams struct {
	OptionID       string `json:"option_id" binding:"required"`
	HostNameOption string `json:"host_option_name" binding:"required"`
}

type UpdateOptionDesParams struct {
	OptionID string `json:"option_id" binding:"required"`
	Des      string `json:"des" binding:"required"`
	DesType  string `json:"des_type" binding:"required,des_type"`
}

type GetOptionParams struct {
	OptionID string `json:"option_id" binding:"required"`
}

type GetOptionDesRes struct {
	SpaceDes              string `json:"space_des"`
	GuestAccessDes        string `json:"guest_access_des"`
	InteractWithGuestsDes string `json:"interact_with_guests_des"`
	OtherDes              string `json:"other_des"`
	NeighborhoodDes       string `json:"neighborhood_des"`
	GetAroundDes          string `json:"get_around_des"`
	Des                   string `json:"des"`
	OptionID              string `json:"option_id"`
}

type UpdateOptionPhotoParams struct {
	OptionID        string   `json:"option_id" binding:"required"`
	CreateImages    []string `json:"create_images"`
	DeleteImage     string   `json:"delete_image"`
	ChangeMainImage string   `json:"change_main_image"`
}

type UpdateOptionPhotoRes struct {
	MainImage string   `json:"main_image"`
	Images    []string `json:"images"`
}

type CreateUpdateOptionPhotoCaptionParams struct {
	OptionID string `json:"option_id" binding:"required"`
	PhotoID  string `json:"photo_id" binding:"required"`
	Caption  string `json:"caption"`
}

type GetOptionPhotoCaptionParams struct {
	OptionID string `json:"option_id" binding:"required"`
	PhotoID  string `json:"photo_id" binding:"required"`
}

type GetOptionPhotoCaptionRes struct {
	PhotoID       string `json:"photo_id"`
	Caption       string `json:"caption"`
	SpaceLocation string `json:"space_location"`
}

type CreateUpdateOptionPhotoCaptionRes struct {
	PhotoID string `json:"photo_id"`
	Caption string `json:"caption"`
}

// User Host Manage for Booking

type UpdateOptionPriceParams struct {
	OptionID     string `json:"option_id" binding:"required"`
	Price        string `json:"price" binding:"required"`
	WeekendPrice string `json:"weekend_price"`
}

type OptionPriceRes struct {
	Price        string `json:"price"`
	WeekendPrice string `json:"weekend_price"`
}

type OptionAddChargeItem struct {
	ID         string `json:"id"`
	MainFee    string `json:"main_fee"`
	Type       string `json:"type"`
	ExtraFee   string `json:"extra_fee"`
	NumOfGuest int    `json:"num_of_guest"`
}

//type UpdateOptionAddChargeParams struct {
//	ID         string `json:"id"`
//	MainFee    string `json:"main_fee"`
//	Type       string `json:"type"`
//	ExtraFee   string `json:"extra_fee"`
//	NumOfGuest int    `json:"num_of_guest"`
//}

type CreateUpdateOptionAddChargeParams struct {
	OptionID   string `json:"option_id" binding:"required"`
	MainFee    string `json:"main_fee" binding:"required"`
	Type       string `json:"type" binding:"required"`
	ExtraFee   string `json:"extra_fee"`
	NumOfGuest int    `json:"num_of_guest"`
}

type ListOptionAddChargeRes struct {
	List        []OptionAddChargeItem `json:"list"`
	PetsAllowed bool                  `json:"pets_allowed"`
}

type UpdatePetsAllowedParams struct {
	PetsAllowed bool   `json:"pets_allowed"`
	OptionID    string `json:"option_id" binding:"required"`
}

type UpdatePetsAllowedRes struct {
	PetsAllowed bool `json:"pets_allowed"`
}

type OptionDiscountItem struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	MainType  string `json:"main_type"`
	Percent   int    `json:"percent"`
	Name      string `json:"name"`
	ExtraType string `json:"extra_type"`
	Des       string `json:"des"`
}

type ListOptionDiscountRes struct {
	List []OptionDiscountItem `json:"list"`
}

type CreateUpdateOptionDiscountItem struct {
	Type      string `json:"type"`
	MainType  string `json:"main_type"`
	Percent   int    `json:"percent"`
	Name      string `json:"name"`
	ExtraType string `json:"extra_type"`
	Des       string `json:"des"`
}

// LOT means Length of stay
type LOTCreateUpdateOptionDiscountParams struct {
	List     []CreateUpdateOptionDiscountItem `json:"list"`
	OptionID string                           `json:"option_id" binding:"required"`
}

type UpdateOptionCurrencyParams struct {
	Currency string `json:"currency" binding:"required,currency"`
	OptionID string `json:"option_id" binding:"required"`
}

type UpdateOptionCurrencyRes struct {
	Currency string `json:"currency"`
	Price    string `json:"price"`
}

type UpdateUnlistedOptionCurrencyRes struct {
	Currency string `json:"currency"`
}
