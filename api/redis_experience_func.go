package api

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

func ConvertSliceToExperienceOptionData(data map[string]string, dollarToNaira string, dollarToCAD string, userCurrency string, id uuid.UUID) ExperienceOptionData {
	basePrice, err := tools.ConvertPrice(getDataValue(data, constants.BASE_PRICE), data[constants.CURRENCY], userCurrency, dollarToNaira, dollarToCAD, id)
	if err != nil {
		basePrice = 0.0
	}
	weekendPrice, err := tools.ConvertPrice(getDataValue(data, constants.WEEKEND_PRICE), data[constants.CURRENCY], userCurrency, dollarToNaira, dollarToCAD, id)
	if err != nil {
		weekendPrice = 0.0
	}
	addedPrice, err := tools.ConvertPrice(getDataValue(data, constants.ADD_PRICE), data[constants.CURRENCY], userCurrency, dollarToNaira, dollarToCAD, id)
	if err != nil {
		addedPrice = 0.0
	}
	newData := ExperienceOptionData{
		UserOptionID:     getDataValue(data, constants.OPTION_USER_ID),
		Name:             getDataValue(data, constants.HOST_OPTION_NAME),
		IsVerified:       tools.ConvertStringToBool(getDataValue(data, constants.OPTION_IS_VERIFIED)),
		CoverImage:       getDataValue(data, constants.COVER_IMAGE),
		HostAsIndividual: tools.ConvertStringToBool(getDataValue(data, constants.HOST_AS_INDIVIDUAL)),
		BasePrice:        tools.ConvertFloatToString(basePrice),
		WeekendPrice:     tools.ConvertFloatToString(weekendPrice),
		Photos:           strings.Split(getDataValue(data, constants.PHOTOS), "&"),
		TypeOfShortlet:   getDataValue(data, constants.TYPE_OF_SHORTLET),
		State:            getDataValue(data, constants.STATE),
		Country:          getDataValue(data, constants.COUNTRY),
		ProfilePhoto:     getDataValue(data, constants.PROFILE_PHOTO),
		HostName:         getDataValue(data, constants.HOST_NAME),
		HostJoined:       getDataValue(data, constants.HOST_JOINED),
		HostVerified:     tools.ConvertStringToBool(getDataValue(data, constants.HOST_VERIFIED)),
		Category:         getDataValue(data, constants.CATEGORY),
		AddedPrice:       tools.ConvertFloatToString(addedPrice),
		StartDate:        getDataValue(data, constants.START_DATE),
		EndDate:          getDataValue(data, constants.END_DATE),
		AddPriceFound:    tools.ConvertStringToBool(getDataValue(data, constants.ADD_DATE_FOUND)),
		// Add more fields here in the same format
	}
	return newData
}

func ConvertSliceToExperienceEventData(data map[string]string, dollarToNaira string, dollarToCAD string, userCurrency string, id uuid.UUID) ExperienceEventData {
	lowestPrice, err := tools.ConvertPrice(data[constants.TICKET_LOWEST_PRICE], data[constants.CURRENCY], userCurrency, dollarToNaira, dollarToCAD, id)
	if err != nil {
		lowestPrice = 0.0
	}
	newData := ExperienceEventData{
		UserOptionID:      data[constants.OPTION_USER_ID],
		Name:              data[constants.HOST_OPTION_NAME],
		IsVerified:        tools.ConvertStringToBool(data[constants.OPTION_IS_VERIFIED]),
		CoverImage:        data[constants.COVER_IMAGE],
		Photos:            strings.Split(data[constants.PHOTOS], "&"),
		TicketAvailable:   tools.ConvertStringToBool(data[constants.TICKET_AVAILABLE]),
		SubEventType:      data[constants.SUB_EVENT_TYPE],
		TicketLowestPrice: tools.ConvertFloatToString(lowestPrice),
		EventStartDate:    data[constants.EVENT_START_DATE],
		EventEndDate:      data[constants.EVENT_END_DATE],
		Location:          HandleExperienceEventLocation(data[constants.LOCATION]),
		ProfilePhoto:      data[constants.PROFILE_PHOTO],
		HostAsIndividual:  tools.ConvertStringToBool(data[constants.HOST_AS_INDIVIDUAL]),
		HostName:          data[constants.HOST_NAME],
		HostJoined:        data[constants.HOST_JOINED],
		HostVerified:      tools.ConvertStringToBool(data[constants.HOST_VERIFIED]),
		Category:          data[constants.CATEGORY],
		HasFreeTicket:     tools.ConvertStringToBool(data[constants.HAS_FREE_TICKET]),
	}
	return newData
}

func HandleExperienceEventLocation(data string) (res []ExperienceEventLocation) {
	// each location data is joined by  & while each individual location data is joined by /
	for _, d := range strings.Split(data, "&") {
		split := strings.Split(d, "/")
		if len(split) == 3 {
			l := ExperienceEventLocation{
				State:   split[0],
				Country: split[1],
				IsEmpty: tools.ConvertStringToBool(split[2]),
			}
			res = append(res, l)
		}
	}
	return
}

func ConvertExperienceOptionDataToSlice(ctx context.Context, server *Server, data db.ListOptionExperienceRow) []string {
	addDateFound, startDateBook, endDateBook, addPrice := HandleOptionRedisExAddPrice(ctx, server, data.ID, data.OptionUserID, data.PreparationTime, data.AvailabilityWindow, data.AdvanceNotice, data.Price, data.WeekendPrice)
	return []string{
		constants.OPTION_USER_ID,
		tools.UuidToString(data.OptionUserID),
		constants.HOST_OPTION_NAME,
		data.HostNameOption,
		constants.OPTION_IS_VERIFIED,
		tools.ConvertBoolToString(data.IsVerified),
		constants.COVER_IMAGE,
		data.CoverImage,
		constants.HOST_AS_INDIVIDUAL,
		tools.ConvertBoolToString(data.HostAsIndividual),
		constants.BASE_PRICE,
		tools.IntToMoneyString(data.Price),
		constants.WEEKEND_PRICE,
		tools.IntToMoneyString(data.WeekendPrice),
		constants.PHOTOS,
		strings.Join(data.Photo, "&"),
		constants.TYPE_OF_SHORTLET,
		data.TypeOfShortlet,
		constants.STATE,
		data.State,
		constants.COUNTRY,
		data.Country,
		constants.PROFILE_PHOTO,
		data.Photo_2,
		constants.HOST_NAME,
		data.FirstName,
		constants.HOST_JOINED,
		tools.ConvertDateOnlyToString(data.CreatedAt),
		constants.HOST_VERIFIED,
		tools.ConvertBoolToString(data.IsVerified_2),
		constants.CATEGORY,
		data.Category,
		constants.CURRENCY,
		data.Currency,
		constants.START_DATE,
		tools.ConvertDateOnlyToString(startDateBook),
		constants.END_DATE,
		tools.ConvertDateOnlyToString(endDateBook),
		constants.ADD_DATE_FOUND,
		tools.ConvertBoolToString(addDateFound),
		constants.ADD_PRICE,
		tools.IntToMoneyString(addPrice),
		// Add more fields here in the same format
	}
}

func getDataValue(data map[string]string, key string) string {
	if value, exists := data[key]; exists {
		return value
	}
	return ""
}

func CustomSort(data []ExperienceOptionData, targetState, targetCountry string) []ExperienceOptionData {
	targetStateLower := strings.ToLower(targetState)
	targetCountryLower := strings.ToLower(targetCountry)

	sort.Slice(data, func(i, j int) bool {
		// First, compare IsVerified (higher priority)
		if data[i].IsVerified != data[j].IsVerified {
			return data[i].IsVerified
		}

		// Then, compare HostVerified (second priority)
		if data[i].HostVerified != data[j].HostVerified {
			return data[i].HostVerified
		}

		// Now, compare State (third priority)
		stateI := strings.ToLower(data[i].State)
		stateJ := strings.ToLower(data[j].State)

		if stateI == targetStateLower && stateJ != targetStateLower {
			return true
		} else if stateI != targetStateLower && stateJ == targetStateLower {
			return false
		}

		// Finally, compare Country (fourth priority)
		countryI := strings.ToLower(data[i].Country)
		countryJ := strings.ToLower(data[j].Country)

		if countryI == targetCountryLower && countryJ != targetCountryLower {
			return true
		} else if countryI != targetCountryLower && countryJ == targetCountryLower {
			return false
		}

		// If all conditions are equal, maintain the original order
		return false
	})
	return data
}

func GetExperienceOptionOffset(data []ExperienceOptionData, offset int, limit int) []ExperienceOptionData {
	// If offset is greater than or equal to the length of data, return an empty slice.
	if offset >= len(data) {
		return []ExperienceOptionData{}
	}

	// Calculate the end index based on offset and limit.
	end := offset + limit

	// If the end index is greater than the length of data, set it to the length of data.
	if end > len(data) {
		end = len(data)
	}

	// Return a subset of data starting from the offset and up to the end index.
	return data[offset:end]
}

func GetExperienceEventOffset(data []ExperienceEventData, offset int, limit int) []ExperienceEventData {
	// If offset is greater than or equal to the length of data, return an empty slice.
	if offset >= len(data) {
		return []ExperienceEventData{}
	}

	// Calculate the end index based on offset and limit.
	end := offset + limit

	// If the end index is greater than the length of data, set it to the length of data.
	if end > len(data) {
		end = len(data)
	}

	// Return a subset of data starting from the offset and up to the end index.
	return data[offset:end]
}

func SetupExperienceEventData(ctx context.Context, server *Server, data db.ListEventExperienceRow, dataTwo db.ListEventExperienceByLocationRow, forRedis bool, funcName string) (locationRedisList []string, locationList []ExperienceEventLocation, price float64, ticketAvailable bool, startDateData string, endDateData string, hasFreeTicket bool) {
	none := constants.NONE
	var dateTimes []db.ListEventDateTimeOnSaleRow
	var err error
	log.Println("At SetupExperienceEventData")
	if forRedis {
		dateTimes, err = server.store.ListEventDateTimeOnSale(ctx, db.ListEventDateTimeOnSaleParams{
			EventInfoID: data.ID,
			Status:      "on_sale",
		})
	} else {
		dateTimes, err = server.store.ListEventDateTimeOnSale(ctx, db.ListEventDateTimeOnSaleParams{
			EventInfoID: dataTwo.ID,
			Status:      "on_sale",
		})
	}
	if err != nil || len(dateTimes) == 0 {
		log.Println("At SetupExperienceEventData empty")
		if err != nil {
			log.Printf("Error at FuncName %v SetupExperienceEventData in ListEventDateTimeOnSale err: %v id: %v\n", funcName, err, data.ID)
		}
		startDateData = none
		endDateData = none
		ticketAvailable = false
		price = 0.00
		hasFreeTicket = false
		if forRedis {
			locationRedisList = []string{fmt.Sprintf("%v/%v/%v", "none", "none", "true")}
		} else {
			locationList = []ExperienceEventLocation{{none, none, true}}
		}
		return
	}
	// Before we call HandleExEventDates to modify dateTimes to handle recurring event we pass it in HandleExperienceRedisDateTimeLocation to store the locations in redis
	//if forRedis {
	//	HandleExperienceRedisDateTimeLocation(dateTimes)
	//}
	dateTimes = HandleExEventDates(dateTimes, "SERVER")
	if err != nil || len(dateTimes) == 0 {
		log.Println("At SetupExperienceEventData empty")
		if err != nil {
			log.Printf("Error at FuncName %v SetupExperienceEventData in ListEventDateTimeOnSale err: %v id: %v\n", funcName, err, data.ID)
		}
		startDateData = none
		endDateData = none
		ticketAvailable = false
		price = 0.00
		hasFreeTicket = false
		if forRedis {
			locationRedisList = []string{fmt.Sprintf("%v/%v/%v", "none", "none", "true")}
		} else {
			locationList = []ExperienceEventLocation{{none, none, true}}
		}
		return
	}
	startTime := dateTimes[0].StartDate
	endTime := dateTimes[0].StartDate
	// Lets store the available locations in redis

	for _, d := range dateTimes {
		if d.StartDate.Before(startTime) {
			// We want to get the min date
			startTime = d.StartDate
		}
		if d.StartDate.After(endTime) {
			// We want to get the most far date
			endTime = d.StartDate
		}
		// We only care about dates with location when setting up the location
		if !tools.ServerStringEmpty(d.Country) && !tools.ServerStringEmpty(d.State) {
			if forRedis {
				locationRedisList = append(locationRedisList, fmt.Sprintf("%v/%v/%v", d.State, d.Country, "false"))
			} else {
				locationList = append(locationList, ExperienceEventLocation{d.State, d.Country, false})
			}
		}
		// We only want to check ticket if ticketAvailable is still false
		tickets := HandleGetAvailableTicket(ctx, server, d.ID, d.StartDate, funcName)
		if len(tickets) != 0 {
			priceData := HandleGetTicketStartPrice(tickets)
			if price == 0.00 {
				price = priceData
			} else {
				if priceData < price {
					price = priceData
				}
			}
			hasFreeTicket = HandleGetFreeTicket(tickets)
			ticketAvailable = true
		}
	}
	if len(locationRedisList) == 0 && forRedis {
		locationRedisList = []string{fmt.Sprintf("%v/%v/%v", "none", "none", "true")}
	}
	if len(locationList) == 0 && !forRedis {
		locationList = []ExperienceEventLocation{{none, none, true}}
	}

	startDateData = tools.ConvertDateOnlyToString(startTime)
	endDateData = tools.ConvertDateOnlyToString(endTime)
	return
}

func HandleExperienceRedisDateTimeLocation(dateTimes []db.ListEventDateTimeOnSaleRow) {
	for _, d := range dateTimes {
		if !tools.ServerStringEmpty(d.Country) && !tools.ServerStringEmpty(d.State) {
			log.Println("Enter locations")
			locationKey := fmt.Sprintf("%v&%v", constants.EXPERIENCE_EVENT_LOCATION, d.ID)
			lat := tools.ConvertFloatToLocationString(d.Geolocation.P.Y, 9)
			lng := tools.ConvertFloatToLocationString(d.Geolocation.P.X, 9)
			location := &redis.GeoLocation{
				Latitude:  tools.ConvertLocationStringToFloat(lat, 9),
				Longitude: tools.ConvertLocationStringToFloat(lng, 9),
				Name:      locationKey,
			}
			err := RedisClient.GeoAdd(RedisContext, constants.ALL_EXPERIENCE_LOCATION, location).Err()
			if err != nil {
				log.Printf("Error at HandleExperienceRedisDateTimeLocation in RedisClient.GeoAdd(ctx, constants.ALL_EXPERIENCE_LOCATION err: %v\n", err)
				continue
			}
		}
	}
}

func HandleGetAvailableTicket(ctx context.Context, server *Server, eventDateTimeID uuid.UUID, startDate time.Time, funcName string) (at []db.EventDateTicket) {
	log.Println("eventDateID ticket ", eventDateTimeID)
	tickets, err := server.store.ListEventDateTicket(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at FuncName %v SetupExperienceEventData in server.store.ListEventDateTicket err: %v id: %v\n", funcName, err, eventDateTimeID)
		return
	}
	for _, t := range tickets {
		count, err := server.store.CountChargeTicketReferenceByStartDate(ctx, db.CountChargeTicketReferenceByStartDateParams{
			Date:        startDate,
			EventDateID: eventDateTimeID,
			Cancelled:   false,
			IsComplete:  true,
			TicketID:    t.ID,
		})
		log.Println("tickets, ", count)
		if err != nil {
			log.Printf("Error at FuncName %v SetupExperienceEventData in CountChargeTicketReferenceByStartDate err: %v id: %v ticketID: %v\n", funcName, err, eventDateTimeID, t.ID)
			// If there was an error we say the ticket is available
			at = append(at, t)
			continue
		}
		if int32(count) < t.Capacity {
			at = append(at, t)
		}
	}
	return
}

func CustomEventExperienceDataSort(data []ExperienceEventData, state, country string) []ExperienceEventData {
	// Define a custom sorting function using sort.Slice
	sort.Slice(data, func(i, j int) bool {
		// Compare state, country, HostVerified, and IsVerified for the two elements
		a := data[i]
		b := data[j]

		// Compare state and country (case-insensitive)
		if strings.EqualFold(a.Location[0].State, state) && strings.EqualFold(a.Location[0].Country, country) {
			// If state and country match, prioritize HostVerified and IsVerified
			return a.HostVerified && a.IsVerified && (!b.HostVerified || !b.IsVerified)
		}

		// If state and country do not match, prioritize other elements
		return false
	})

	return data
}

func HandleOptionRedisExAddPrice(ctx context.Context, server *Server, optionID uuid.UUID, optionUserID uuid.UUID, prepareTime string, window string, advanceNotice string, basePrice int64, weekendPrice int64) (addDateFound bool, startDateBook time.Time, endDateBook time.Time, addPrice int64) {
	dateTimes := GetOptionDateTimes(ctx, server, optionID, "HandleOptionRedisExAddPrice")
	intervals := tools.ListDateIntervals(5)
	for _, v := range intervals {
		confirm := OptionRedisExMain(ctx, server, prepareTime, window, window, optionUserID, v.StartDate, v.EndDate, "HandleOptionRedisExAddPrice", dateTimes, optionID)
		if confirm {
			addPrice = HandleOptionPrice(basePrice, weekendPrice, dateTimes, v.StartDate, v.EndDate, "HandleOptionRedisExAddPrice")
			startDateBook = v.StartDate
			endDateBook = v.EndDate
			addDateFound = true
			break
		}
	}
	return
}

func OptionRedisExMain(ctx context.Context, server *Server, prepareTime string, window string, advanceNotice string, optionUserID uuid.UUID, startDate time.Time, endDate time.Time, funcName string, dateTimes []db.OptionDateTime, optionID uuid.UUID) bool {
	// We check the available settings
	settingGood := HandleOptionSearchSetting(startDate, endDate, advanceNotice, prepareTime, window, optionUserID, funcName, optionID)
	// We check search charges
	chargeGood := HandleOptionSearchCharges(ctx, server, optionUserID, funcName, prepareTime, startDate, endDate)
	// We check dates
	dateTimeGood := HandleOptionSearchDates(ctx, server, optionUserID, funcName, startDate, endDate, dateTimes)
	return settingGood && chargeGood && dateTimeGood
}
