package api

import (
	"context"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type ExFilterRangeReq struct {
	Type     string `json:"type"`
	Currency string `json:"currency"`
}

type ExFilterRangeRes struct {
	MaxPrice        string `json:"max_price"`
	MinPrice        string `json:"min_price"`
	AveragePrice    string `json:"average_price"`
	AddMaxPrice     string `json:"add_max_price"`
	AddMinPrice     string `json:"add_min_price"`
	AddDayCount     int    `json:"add_day_count"`
	AverageAddPrice string `json:"average_add_price"`
}

type ExFilterOptionRequest struct {
	MaxPrice    string `json:"max_price"`
	MinPrice    string `json:"min_price"`
	AddMaxPrice string `json:"add_max_price"`
	AddMinPrice string `json:"add_min_price"`
	// OnAddPrice would tell us what price it is on
	OnAddPrice        bool     `json:"on_add_price"`
	Currency          string   `json:"currency"`
	ShortletSpaceType []string `json:"shortlet_space_type"`
	Bedrooms          int      `json:"bedrooms"`
	Beds              int      `json:"beds"`
	Bathrooms         int      `json:"bathrooms"`
	CategoryTypes     []string `json:"category_types"`
	Amenities         []string `json:"amenities"`
	CanInstantBook    bool     `json:"can_instant_book"`
	CanSelfCheck      bool     `json:"can_self_check"`
}

type ExFilterEventRequest struct {
	MaxPrice    string `json:"max_price"`
	MinPrice    string `json:"min_price"`
	AddMaxPrice string `json:"add_max_price"`
	AddMinPrice string `json:"add_min_price"`
	// OnAddPrice would tell us what price it is on
	OnAddPrice      bool     `json:"on_add_price"`
	Currency        string   `json:"currency"`
	TicketAvailable bool     `json:"ticket_available"`
	CategoryType    []string `json:"category_type"`
	TicketType      []string `json:"ticket_type"`
	SubCategory     []string `json:"sub_category"`
}

type ExSearchRequest struct {
	PeriodType     string        `json:"period_type"`
	PeriodDaySpace int           `json:"period_day_space"`
	StayType       string        `json:"stay_type"`
	Months         []ExMonthItem `json:"months"`
	GuestTypes     []string      `json:"guest_types"`
	Street         string        `json:"street"`
	City           string        `json:"city"`
	State          string        `json:"state"`
	Country        string        `json:"country"`
	Lat            string        `json:"lat"`
	Lng            string        `json:"lng"`
	PostCode       string        `json:"post_code"`
	StartDate      string        `json:"start_date"`
	EndDate        string        `json:"end_date"`
	Currency       string        `json:"currency"`
}

type ExControlOptionRequest struct {
	Offset int `json:"offset"`
	// Type is category type
	Type          string                `json:"type"`
	Filter        ExFilterOptionRequest `json:"filter"`
	IsFilterEmpty bool                  `json:"is_filter_empty"`
	IsSearchEmpty bool                  `json:"is_search_empty"`
	Search        ExSearchRequest       `json:"search"`
}

type ExControlEventRequest struct {
	Offset int `json:"offset"`
	// Type is category type
	Type          string               `json:"type"`
	Filter        ExFilterEventRequest `json:"filter"`
	IsFilterEmpty bool                 `json:"is_filter_empty"`
	IsSearchEmpty bool                 `json:"is_search_empty"`
	Search        ExSearchRequest      `json:"search"`
}

type ExMonthItem struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func ExSearchReqHasLocation(data ExSearchRequest) bool {
	var exist bool
	list := []string{data.City, data.State, data.Country, data.PostCode, data.Street}
	for _, s := range list {
		if !tools.ServerStringEmpty(s) {
			exist = true
		}
	}
	if !exist {
		lat := tools.ConvertStringToFloat(data.Lat)
		lng := tools.ConvertStringToFloat(data.Lng)
		if lat != 0.0 || lng != 0.0 {
			exist = true
		}
	}
	return exist
}

func HandleOptionExSearchLocation(ctx context.Context, server *Server, req ExControlOptionRequest, funcName string) (list []ExperienceOptionData) {
	funcName = "HandleOptionExSearchLocation"
	lat := tools.ConvertLocationStringToFloat(req.Search.Lat, 9)
	lng := tools.ConvertLocationStringToFloat(req.Search.Lng, 9)
	options, err := server.store.ListOptionInfoSearchLocation(ctx, db.ListOptionInfoSearchLocationParams{
		City:        req.Search.City,
		Street:      req.Search.Street,
		LlToEarth:   lat,
		LlToEarth_2: lng,
		Column5:     100.0,
	})
	if err != nil || len(options) == 0 {
		if err != nil {
			log.Printf("Error at FuncName %v, HandleOptionExSearchLocation ListOptionInfoSearchLocation err: %v \n", funcName, err.Error())
		}
		return
	}
	var guestData = make(map[string]int)
	for _, g := range req.Search.GuestTypes {
		guestData[g] += 1
	}
	totalGuests := guestData["children"] + guestData["adult"]
	log.Println("total guests: ", totalGuests)
	for _, o := range options {
		var addPrice string
		//var priceFloat float64
		var basePrice float64
		var weekendPrice float64
		var startDateBook time.Time
		var addPriceInt int64
		var endDateBook time.Time
		if !req.IsSearchEmpty {
			var confirmBook bool
			dateTimes := GetOptionDateTimes(ctx, server, o.ID, funcName)
			if int(o.GuestWelcomed) < totalGuests || (!o.PetsAllowed && guestData["pet"] > 0) {
				log.Println("failed guests and pets", o.GuestWelcomed)
				continue
			}

			switch req.Search.PeriodType {
			case "flexible":
				startDateBook, endDateBook, confirmBook = HandleOptionSearchFlexible(ctx, server, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, req.Search, funcName, dateTimes, o.ID)
				if !confirmBook {
					continue
				}
			case "choose_date":
				startDate, endDate, errDate := OptionSearchChooseDate(ctx, server, req.Search, funcName)
				if errDate != nil {
					continue
				}
				startDateBook, endDateBook, confirmBook = HandleOptionSearchChooseDate(ctx, server, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, req.Search, funcName, startDate, endDate, dateTimes, o.ID)
				if !confirmBook {
					continue
				}
			}

			addPrice, _, basePrice, weekendPrice, err = HandleOptionSearchPrice(ctx, server, req.Search, o.ID, o.OptionUserID, o.Currency, o.Price, o.WeekendPrice, dateTimes, startDateBook, endDateBook, funcName)
			if err != nil {
				continue
			}
		} else {
			_, startDateBook, endDateBook, addPriceInt = HandleOptionRedisExAddPrice(ctx, server, o.ID, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, o.Price, o.WeekendPrice)
			addPriceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(addPriceInt), o.Currency, req.Search.Currency, server.config.DollarToNaira, server.config.DollarToCAD, o.ID)
			if err != nil {
				log.Printf("Error at addedPrice HandleOptionExSearchLocation in ConvertPrice err: %v, user: %v\n", err, o.ID)
				addPriceFloat = 0.0
			}
			addPrice = tools.ConvertFloatToString(addPriceFloat)
		}
		if !req.IsFilterEmpty {
			optionCategory := []string{o.Category, o.CategoryTwo, o.CategoryThree, o.CategoryFour}
			confirm := HandleOptionFilter(ctx, server, o.ID, o.OptionUserID, req.Filter, o.Price, o.Currency, funcName, o.SpaceType, optionCategory, o.InstantBook, o.CheckInMethod)
			if !confirm {
				continue
			}
		}
		hostJoined := tools.ConvertDateOnlyToString(o.CreatedAt)
		data := ExperienceOptionData{
			UserOptionID:     tools.UuidToString(o.OptionUserID),
			Name:             o.HostNameOption,
			IsVerified:       o.OptionIsVerified,
			CoverImage:       o.CoverImage,
			HostAsIndividual: o.HostAsIndividual,
			BasePrice:        tools.ConvertFloatToString(basePrice),
			WeekendPrice:     tools.ConvertFloatToString(weekendPrice),
			AddedPrice:       addPrice,
			AddPriceFound:    true,
			StartDate:        tools.ConvertDateOnlyToString(startDateBook),
			EndDate:          tools.ConvertDateOnlyToString(endDateBook),
			Photos:           o.Photo,
			TypeOfShortlet:   o.TypeOfShortlet,
			State:            o.State,
			Country:          o.Country,
			ProfilePhoto:     o.ProfilePhoto,
			HostName:         o.HostName,
			HostJoined:       hostJoined,
			HostVerified:     o.HostVerified,
			Category:         o.Category,
		}
		list = append(list, data)
	}
	return
}

func HandleOptionExSearch(ctx context.Context, server *Server, req ExControlOptionRequest, funcName string) (list []ExperienceOptionData) {
	funcName = "HandleOptionExSearch"
	options, err := server.store.ListOptionInfoSearch(ctx)
	if err != nil || len(options) == 0 {
		if err != nil {
			log.Printf("Error at FuncName %v, HandleOptionExSearch ListOptionInfoSearch err: %v \n", funcName, err.Error())
		}
		return
	}
	var guestData = make(map[string]int)
	for _, g := range req.Search.GuestTypes {
		guestData[g] += 1
	}
	totalGuests := guestData["children"] + guestData["adult"]
	for _, o := range options {
		var addPrice string
		//var priceFloat float64
		var basePrice float64
		var weekendPrice float64
		var startDateBook time.Time
		var addPriceInt int64
		var endDateBook time.Time
		if !req.IsSearchEmpty {
			var confirmBook bool
			dateTimes := GetOptionDateTimes(ctx, server, o.ID, funcName)
			if int(o.GuestWelcomed) < totalGuests || (!o.PetsAllowed && guestData["pet"] > 0) {
				continue
			}
			switch req.Search.PeriodType {
			case "flexible":
				startDateBook, endDateBook, confirmBook = HandleOptionSearchFlexible(ctx, server, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, req.Search, funcName, dateTimes, o.ID)
				if !confirmBook {
					continue
				}
			case "choose_date":
				startDate, endDate, errDate := OptionSearchChooseDate(ctx, server, req.Search, funcName)
				if errDate != nil {
					continue
				}
				startDateBook, endDateBook, confirmBook = HandleOptionSearchChooseDate(ctx, server, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, req.Search, funcName, startDate, endDate, dateTimes, o.ID)
				if !confirmBook {
					continue
				}
			}

			addPrice, _, basePrice, weekendPrice, err = HandleOptionSearchPrice(ctx, server, req.Search, o.ID, o.OptionUserID, o.Currency, o.Price, o.WeekendPrice, dateTimes, startDateBook, endDateBook, funcName)
			if err != nil {
				continue
			}
		} else {
			_, startDateBook, endDateBook, addPriceInt = HandleOptionRedisExAddPrice(ctx, server, o.ID, o.OptionUserID, o.PreparationTime, o.AvailabilityWindow, o.AdvanceNotice, o.Price, o.WeekendPrice)
			addPriceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(addPriceInt), o.Currency, req.Search.Currency, server.config.DollarToNaira, server.config.DollarToCAD, o.ID)
			if err != nil {
				log.Printf("Error at addedPrice HandleOptionExSearch in ConvertPrice err: %v, user: %v\n", err, o.ID)
				addPriceFloat = 0.0
			}
			addPrice = tools.ConvertFloatToString(addPriceFloat)
		}
		if !req.IsFilterEmpty {
			optionCategory := []string{o.Category, o.CategoryTwo, o.CategoryThree, o.CategoryFour}
			confirm := HandleOptionFilter(ctx, server, o.ID, o.OptionUserID, req.Filter, o.Price, o.Currency, funcName, o.SpaceType, optionCategory, o.InstantBook, o.CheckInMethod)
			if !confirm {
				continue
			}
		}
		hostJoined := tools.ConvertDateOnlyToString(o.CreatedAt)
		data := ExperienceOptionData{
			UserOptionID:     tools.UuidToString(o.OptionUserID),
			Name:             o.HostNameOption,
			IsVerified:       o.OptionIsVerified,
			CoverImage:       o.CoverImage,
			HostAsIndividual: o.HostAsIndividual,
			BasePrice:        tools.ConvertFloatToString(basePrice),
			WeekendPrice:     tools.ConvertFloatToString(weekendPrice),
			AddedPrice:       addPrice,
			AddPriceFound:    true,
			StartDate:        tools.ConvertDateOnlyToString(startDateBook),
			EndDate:          tools.ConvertDateOnlyToString(endDateBook),
			Photos:           o.Photo,
			TypeOfShortlet:   o.TypeOfShortlet,
			State:            o.State,
			Country:          o.Country,
			ProfilePhoto:     o.ProfilePhoto,
			HostName:         o.HostName,
			HostJoined:       hostJoined,
			HostVerified:     o.HostVerified,
			Category:         o.Category,
		}
		list = append(list, data)
	}
	return
}

func GetOptionFilterMaxMinPrice(ctx context.Context, server *Server, req ExFilterRangeReq, clientIP string) (minPrice float64, maxPrice float64, err error) {
	optionPrices, err := server.store.ListOptionInfoPrice(ctx)
	if err != nil || len(optionPrices) == 0 {
		if err != nil {
			log.Printf("Error at GetOptionFilterRange in ListOptionInfoPrice err: %v, user: %v\n", err, clientIP)
		}
		minPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(server.config.OptionMinPriceNaira)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error 0 at min FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", "GetOptionFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
		maxPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(server.config.OptionMaxPriceNaira)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error 0 max at FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", "GetOptionFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
	} else {
		minPriceInt := optionPrices[0].Price
		maxPriceInt := optionPrices[len(optionPrices)-1].Price
		if minPriceInt == maxPriceInt {
			// We want to ensure it is not equal
			maxPriceInt = int64(server.config.OptionMaxPriceNaira)
		}
		minPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(minPriceInt)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error at min FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", "GetOptionFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
		maxPrice, err = tools.ConvertPrice(tools.IntToMoneyString(maxPriceInt), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error at max FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", "GetOptionFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}

	}
	return
}
