package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func HandleOptionSearchFlexible(ctx context.Context, server *Server, optionUserID uuid.UUID, prepareTime string, window string, advanceNotice string, req ExSearchRequest, funcName string, dateTimes []db.OptionDateTime, optionID uuid.UUID) (startDateBook time.Time, endDateBook time.Time, confirmBook bool) {
	switch req.StayType {
	case "weekend":
		weekends := tools.ListWeekends()
		for _, w := range weekends {
			startDate := w.Friday
			endDate := w.Sunday
			confirm := OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, startDate, endDate, funcName, dateTimes, optionID)
			if confirm {
				startDateBook = startDate
				endDateBook = endDate
				confirmBook = true
				break
			}

		}
	case "week":
		weeks := tools.ListWeeks()
		for _, w := range weeks {
			startDate := w.Monday
			endDate := w.Sunday
			confirm := OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, startDate, endDate, funcName, dateTimes, optionID)
			if confirm {
				startDateBook = startDate
				endDateBook = endDate
				confirmBook = true
				break
			}

		}
	case "month":
		months := tools.ListMonths()
		for _, m := range months {
			startDate := m.StartDateOfMonth
			endDate := m.EndDateOfMonth
			confirm := OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, startDate, endDate, funcName, dateTimes, optionID)
			if confirm {
				startDateBook = startDate
				endDateBook = endDate
				confirmBook = true
				break
			}

		}
	}
	return
}

func HandleOptionSearchChooseDate(ctx context.Context, server *Server, optionUserID uuid.UUID, prepareTime string, window string, advanceNotice string, req ExSearchRequest, funcName string, startDate time.Time, endDate time.Time, dateTimes []db.OptionDateTime, optionID uuid.UUID) (startDateBook time.Time, endDateBook time.Time, confirmBook bool) {

	confirm := OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, startDate, endDate, funcName, dateTimes, optionID)
	if confirm {
		startDateBook = startDate
		endDateBook = endDate
		confirmBook = true
		return
	}
	// if the first one fails and then we can try the next
	if req.PeriodDaySpace != 0 {
		newStartDate := startDate.AddDate(0, 0, -req.PeriodDaySpace)
		confirm = OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, newStartDate, endDate, funcName, dateTimes, optionID)
		if confirm {
			startDateBook = newStartDate
			endDateBook = endDate
			confirmBook = true
			return
		}
		// We try we end date
		newEndDate := endDate.AddDate(0, 0, req.PeriodDaySpace)
		confirm = OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, startDate, newEndDate, funcName, dateTimes, optionID)
		if confirm {
			startDateBook = startDate
			endDateBook = newEndDate
			confirmBook = true
			return
		}
		// We try both
		confirm = OptionSearchMain(ctx, server, req, prepareTime, window, advanceNotice, optionUserID, newStartDate, newEndDate, funcName, dateTimes, optionID)
		if confirm {
			startDateBook = newStartDate
			endDateBook = newEndDate
			confirmBook = true
			return
		}
	}
	return
}

func OptionSearchChooseDate(ctx context.Context, server *Server, req ExSearchRequest, funcName string) (startDate time.Time, endDate time.Time, err error) {
	startDate, err = tools.ConvertDateOnlyStringToDate(req.StartDate)
	if err != nil {
		log.Printf("Error start at FuncName %v, OptionSearchChooseDate ConvertDateOnlyStringToDate err: %v \n", funcName, err.Error())
		err = fmt.Errorf("your start date was not selected")
		return
	}
	endDate, err = tools.ConvertDateOnlyStringToDate(req.EndDate)
	if err != nil {
		log.Printf("Error end at FuncName %v, OptionSearchChooseDate ConvertDateOnlyStringToDate err: %v \n", funcName, err.Error())
		err = fmt.Errorf("your end date was not selected")
		return
	}

	return
}

func OptionSearchMain(ctx context.Context, server *Server, req ExSearchRequest, prepareTime string, window string, advanceNotice string, optionUserID uuid.UUID, startDate time.Time, endDate time.Time, funcName string, dateTimes []db.OptionDateTime, optionID uuid.UUID) bool {
	// We check if the weekend is in the month
	var monthGood bool
	if req.PeriodType == "flexible" {
		monthGood = OptionSearchMonthGood(req, startDate, endDate)
	} else {
		monthGood = true
	}
	// We check the available settings
	log.Println("at available")
	settingGood := HandleOptionSearchSetting(startDate, endDate, advanceNotice, prepareTime, window, optionUserID, funcName, optionID)
	// We check search charges
	log.Println("at charges")
	chargeGood := HandleOptionSearchCharges(ctx, server, optionUserID, funcName, prepareTime, startDate, endDate)
	// We check dates
	log.Println("at dates")
	dateTimeGood := HandleOptionSearchDates(ctx, server, optionUserID, funcName, startDate, endDate, dateTimes)
	log.Println("was it good: ", monthGood && settingGood && chargeGood && dateTimeGood)
	return monthGood && settingGood && chargeGood && dateTimeGood
}

func OptionSearchMonthGood(req ExSearchRequest, startDate time.Time, endDate time.Time) (exist bool) {
	if len(req.Months) == 0 {
		exist = true
		return
	}
	for _, month := range req.Months {
		if (int(startDate.Month()) == month.Month && startDate.Year() == month.Year) && (int(endDate.Month()) == month.Month && endDate.Year() == month.Year) {
			exist = true
			break
		}
	}

	return
}

func HandleOptionSearchSetting(startDate time.Time, endDate time.Time, advanceNotice string, prepare string, window string, optionUserID uuid.UUID, funcName string, optionID uuid.UUID) (confirm bool) {
	// Check the available settings
	confirm, err := ReserveAvailableSetting(tools.ConvertDateOnlyToString(startDate), tools.ConvertDateOnlyToString(endDate), optionID, advanceNotice, prepare, window)
	if err != nil {
		log.Printf("Error at FuncName %v, HandleOptionExSearchLocation ListChargeOptionReferenceDatesMore err: %v optionUserID: %v\n", funcName, err.Error(), optionUserID)
	}
	return
}

func HandleOptionSearchCharges(ctx context.Context, server *Server, optionUserID uuid.UUID, funcName string, prepareTime string, startMainDate time.Time, endMainDate time.Time) (confirm bool) {
	confirm = true
	var dayChange int
	if prepareTime == constants.PREPARE_ONE_NIGHT {
		dayChange = 1
	} else if prepareTime == constants.PREPARE_TWO_NIGHT {
		dayChange = 2
	}
	charges, err := server.store.ListChargeOptionReferenceDatesMore(ctx, db.ListChargeOptionReferenceDatesMoreParams{
		OptionUserID: optionUserID,
		Cancelled:    false,
		IsComplete:   true,
	})
	if err != nil {
		log.Printf("Error at FuncName %v, HandleOptionExSearchLocation ListChargeOptionReferenceDatesMore err: %v optionUserID: %v\n", funcName, err.Error(), optionUserID)
		err = nil
	}
	for _, charge := range charges {
		var startDate time.Time
		var endDate time.Time
		if dayChange != 0 {
			startDate = charge.StartDate.AddDate(0, 0, -dayChange)
			endDate = charge.EndDate.AddDate(0, 0, dayChange)
		} else {
			startDate = charge.StartDate
			endDate = charge.EndDate
		}
		if !((startMainDate.Before(startDate) && endMainDate.Before(startDate)) || (startMainDate.After(endDate) && endMainDate.After(endDate))) {
			confirm = false
			break
		}
	}
	return
}

func GetOptionDateTimes(ctx context.Context, server *Server, optionID uuid.UUID, funcName string) (dateTimes []db.OptionDateTime) {
	dateTimes, err := server.store.ListOptionDateTimeMore(ctx, optionID)
	if err != nil {
		log.Printf("Error at FuncName %v, GetOptionDateTimes ListOptionDateTimeMore err: %v optionUserID: %v\n", funcName, err.Error(), optionID)
		err = nil
	}
	return
}

func HandleOptionSearchDates(ctx context.Context, server *Server, optionUserID uuid.UUID, funcName string, startMainDate time.Time, endMainDate time.Time, dateTimes []db.OptionDateTime) (confirm bool) {
	confirm = true
	mainDates := tools.GenerateDateListStringFromTime(startMainDate, endMainDate)
	for _, m := range mainDates {
		for _, d := range dateTimes {
			if tools.ConvertDateOnlyToString(d.Date) == m && !d.Available {
				confirm = false
				return
			}
		}
	}

	return
}

func HandleOptionSearchPrice(ctx context.Context, server *Server, req ExSearchRequest, optionID uuid.UUID, optionUserID uuid.UUID, optionCurrency string, optionPrice int64, optionWeekendPrice int64, dateTimes []db.OptionDateTime, startDateBook time.Time, endDateBook time.Time, funcName string) (addPrice string, priceFloat float64, basePrice float64, weekendPrice float64, err error) {
	basePrice, err = tools.ConvertPrice(tools.IntToMoneyString(optionPrice), optionCurrency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, optionUserID)
	if err != nil {
		log.Printf("Error at FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", funcName, err.Error())
		return
	}
	weekendPrice, err = tools.ConvertPrice(tools.IntToMoneyString(optionWeekendPrice), optionCurrency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, optionUserID)
	if err != nil {
		log.Printf("Error at FuncName %v, HandleOptionSearchPrice  tools.ConvertPrice err: %v \n", funcName, err.Error())
		return
	}
	dates := tools.GenerateDateListStringFromTime(startDateBook, endDateBook)
	for _, d := range dates {
		var exist bool
		var datePrice int64

		for _, dt := range dateTimes {
			if tools.ConvertDateOnlyToString(dt.Date) == d && dt.Available {
				exist = true
				datePrice = dt.Price
				break
			}
		}
		if exist && datePrice > 0 {
			datePriceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(datePrice), optionCurrency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, optionUserID)
			if err != nil {
				log.Printf("Error at FuncName %v, HandleOptionSearchPrice tools.ConvertPrice err: %v \n", funcName, err.Error())
				err = nil
			} else {
				priceFloat += datePriceFloat
				continue
			}
		}
		// If price is not in special dateTimes or an error happen we just calculate using standard price format
		// We check if the date is a weekend and there is a weekend price
		isWeekend, err := tools.IsWeekend(d)
		if err == nil {
			if isWeekend && optionWeekendPrice > 0 {
				priceFloat += weekendPrice
				continue
			}
		} else {
			log.Printf("Error at FuncName %v, HandleOptionSearchPrice tools.IsWeekend err: %v \n", funcName, err.Error())
			err = nil
		}
		// We isWeekend goes wrong or it is not a weekend
		priceFloat += basePrice
	}
	log.Println("price for a night: ", priceFloat)
	addPrice = tools.ConvertFloatToString(priceFloat)
	return
}

// We make this function possible to be used by anyone other function
func HandleOptionPrice(basePrice int64, weekendPrice int64, dateTimes []db.OptionDateTime, startDateBook time.Time, endDateBook time.Time, funcName string) (addPrice int64) {
	dates := tools.GenerateDateListStringFromTime(startDateBook, endDateBook)
	for _, d := range dates {
		var exist bool
		var datePrice int64

		for _, dt := range dateTimes {
			if tools.ConvertDateOnlyToString(dt.Date) == d && dt.Available {
				exist = true
				datePrice = dt.Price
				break
			}
		}
		if exist && datePrice > 0 {
			addPrice += datePrice
			continue
		} else {
			// If price is not in special dateTimes or an error happen we just calculate using standard price format
			// We check if the date is a weekend and there is a weekend price
			isWeekend, err := tools.IsWeekend(d)
			if err == nil {
				if isWeekend && weekendPrice > 0 {
					addPrice += weekendPrice
					continue
				}
			} else {
				log.Printf("Error at FuncName %v, HandleOptionSearchPrice tools.IsWeekend err: %v \n", funcName, err.Error())
				err = nil
			}
			// We isWeekend goes wrong or it is not a weekend
			addPrice += basePrice
		}
	}
	return
}
