package api

import (
	"context"
	"log"
	"sort"
	"time"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/google/uuid"
)

func HandleEventInsights(ctx context.Context, server *Server, user db.User, req GetEventInsightParams, optionUserID uuid.UUID, dollarToNaira string, dollarToCAD string, startYear time.Time, funcName string) (res GetEventInsightRes) {
	var resData []GetEventInsightItem
	var itemEmpty bool
	var startMainTime time.Time
	var eventMainDateTimeID uuid.UUID
	if req.ForOffset {
		eventDates, err := server.store.ListEventDateTimeInsight(ctx, db.ListEventDateTimeInsightParams{
			CoUserID:     tools.UuidToString(user.UserID),
			HostID:       user.ID,
			Limit:        50,
			Offset:       int32(req.Offset),
			OptionUserID: optionUserID,
		})
		if err != nil || len(eventDates) == 0 {
			if err != nil {
				log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)
				if err == db.ErrorRecordNotFound {
					empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
					resData = append(resData, empty)
					itemEmpty = true
				}
			} else {
				// If not error but data is empty
				empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
				resData = append(resData, empty)
				itemEmpty = true
			}
		} else {
			for _, d := range eventDates {
				var startTime time.Time
				var startDate string
				var endDate string
				var count int
				var earning string
				switch d.Type {
				case "single":
					startTime = d.StartDate
					startDate = tools.ConvertDateOnlyToString(d.StartDate)
					endDate = tools.ConvertDateOnlyToString(d.EndDate)
				case "recurring":
					startTime, err = tools.ConvertDateOnlyStringToDate(d.EventDate)
					if err != nil {
						log.Printf("Error at FuncName GetEventInsight tools.ConvertDateOnlyStringToDate %v, err: %v OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)
						continue
					}
					startDate = d.EventDate
					endDate = d.EventDate
				}
				// Let get count
				payouts, err := server.store.ListChargeTicketReferencePayoutInsights(ctx, db.ListChargeTicketReferencePayoutInsightsParams{
					CoUserID:        tools.UuidToString(user.UserID),
					HostID:          user.ID,
					Date:            startTime,
					EventDateID:     d.EventDateTimeID,
					OptionUserID:    optionUserID,
					PaymentComplete: true,
				})
				if err != nil {
					log.Printf("Error at FuncName %v GetEventInsight tools.CountChargeTicketReferenceOnlyByStartDate %v, OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)
					count = 0
					earning = "0.00"
				} else {
					itemCount, itemPrice, _ := GetTicketInsightCountAndPrice(payouts, dollarToNaira, dollarToCAD, req.Currency, funcName, user.ID)
					count = itemCount
					earning = tools.ConvertFloatToString(itemPrice)
				}
				data := GetEventInsightItem{
					StartDate:       startDate,
					EndDate:         endDate,
					Name:            d.Name,
					EventDateTimeID: tools.UuidToString(d.EventDateTimeID),
					Count:           count,
					Earnings:        earning,
					FakeID:          tools.UuidToString(uuid.New()),
				}
				resData = append(resData, data)
			}
			itemEmpty = false
			startMainTime = eventDates[0].StartDate
			eventMainDateTimeID = eventDates[0].EventDateTimeID
		}
	} else {
		empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
		resData = append(resData, empty)
		itemEmpty = true
	}
	if !tools.ServerStringEmpty(req.StartDate) {
		newTime, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
		if err != nil {
			log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)

		} else {
			startMainTime = newTime
		}
	}
	if !tools.ServerStringEmpty(req.EventDateTimeID) {
		newID, err := tools.StringToUuid(req.EventDateTimeID)
		if err != nil {
			log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)

		} else {
			eventMainDateTimeID = newID
		}
	}
	var count int
	var earning string
	var cancellation int
	payouts, err := server.store.ListChargeTicketReferencePayoutInsights(ctx, db.ListChargeTicketReferencePayoutInsightsParams{
		CoUserID:        tools.UuidToString(user.UserID),
		HostID:          user.ID,
		Date:            startMainTime,
		EventDateID:     eventMainDateTimeID,
		OptionUserID:    optionUserID,
		PaymentComplete: true,
	})
	if err != nil {
		log.Printf("Error at FuncName %v GetEventInsight tools.CountChargeTicketReferenceOnlyByStartDate %v, OptionUserID: %v \n", funcName, err.Error(), req.OptionUserID)
		count = 0
		cancellation = 0
		earning = "0.00"

	} else {
		itemCount, itemPrice, itemCancellation := GetTicketInsightCountAndPrice(payouts, dollarToNaira, dollarToCAD, req.Currency, funcName, user.ID)
		count = itemCount
		cancellation = itemCancellation
		earning = tools.ConvertFloatToString(itemPrice)
	}
	res = GetEventInsightRes{
		Earning:         earning,
		List:            resData,
		IsEmpty:         itemEmpty,
		ForOffset:       req.ForOffset,
		Cancellation:    cancellation,
		TicketSold:      count,
		StartYear:       startYear.Year(),
		Currency:        req.Currency,
		StartDate:       tools.ConvertDateOnlyToString(startMainTime),
		EventDateTimeID: tools.UuidToString(eventMainDateTimeID),
	}
	return
}

func HandleAllEventInsights(ctx context.Context, server *Server, user db.User, req GetAllEventInsightParams, dollarToNaira string, dollarToCAD string, startYear time.Time, funcName string) (res GetEventInsightRes) {
	var resData []GetEventInsightItem
	var itemEmpty bool
	var startMainTime time.Time
	var eventMainDateTimeID uuid.UUID
	if req.ForOffset {
		eventDates, err := server.store.ListAllEventDateTimeInsight(ctx, db.ListAllEventDateTimeInsightParams{
			CoUserID: tools.UuidToString(user.UserID),
			HostID:   user.ID,
			Limit:    50,
			Offset:   int32(req.Offset),
		})
		if err != nil || len(eventDates) == 0 {
			if err != nil {
				log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, userID: %v \n", funcName, err.Error(), user.ID)
				if err == db.ErrorRecordNotFound {
					empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
					resData = append(resData, empty)
					itemEmpty = true
				}
			} else {
				// If not error but data is empty
				empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
				resData = append(resData, empty)
				itemEmpty = true
			}
		} else {
			for _, d := range eventDates {
				var startTime time.Time
				var startDate string
				var endDate string
				var count int
				var earning string
				switch d.Type {
				case "single":
					startTime = d.StartDate
					startDate = tools.ConvertDateOnlyToString(d.StartDate)
					endDate = tools.ConvertDateOnlyToString(d.EndDate)
				case "recurring":
					startTime, err = tools.ConvertDateOnlyStringToDate(d.EventDate)
					if err != nil {
						log.Printf("Error at FuncName GetEventInsight tools.ConvertDateOnlyStringToDate %v, err: %v, user.ID: %v \n", funcName, err.Error(), user.ID)
						continue
					}
					startDate = d.EventDate
					endDate = d.EventDate
				}
				// Let get count
				payouts, err := server.store.ListAllChargeTicketReferencePayoutInsights(ctx, db.ListAllChargeTicketReferencePayoutInsightsParams{
					CoUserID:        tools.UuidToString(user.UserID),
					HostID:          user.ID,
					Date:            startTime,
					EventDateID:     d.EventDateTimeID,
					PaymentComplete: true,
				})
				if err != nil {
					log.Printf("Error at FuncName %v GetEventInsight tools.CountChargeTicketReferenceOnlyByStartDate %v, user.ID: %v \n", funcName, err.Error(), user.ID)
					count = 0
					earning = "0.00"
				} else {
					itemCount, itemPrice, _ := GetAllTicketInsightCountAndPrice(payouts, dollarToNaira, dollarToCAD, req.Currency, funcName, user.ID)
					count = itemCount
					earning = tools.ConvertFloatToString(itemPrice)
				}
				data := GetEventInsightItem{
					StartDate:       startDate,
					EndDate:         endDate,
					Name:            d.Name,
					EventDateTimeID: tools.UuidToString(d.EventDateTimeID),
					Count:           count,
					Earnings:        earning,
					FakeID:          tools.UuidToString(uuid.New()),
				}
				resData = append(resData, data)
			}
			itemEmpty = false
			startMainTime = eventDates[0].StartDate
			eventMainDateTimeID = eventDates[0].EventDateTimeID
		}
	} else {
		empty := GetEventInsightItem{"none", "none", "none", "none", 0, "0.00", tools.UuidToString(uuid.New())}
		resData = append(resData, empty)
		itemEmpty = true
	}
	if !tools.ServerStringEmpty(req.StartDate) {
		newTime, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
		if err != nil {
			log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, user.ID: %v \n", funcName, err.Error(), user.ID)

		} else {
			startMainTime = newTime
		}
	}
	if !tools.ServerStringEmpty(req.EventDateTimeID) {
		newID, err := tools.StringToUuid(req.EventDateTimeID)
		if err != nil {
			log.Printf("Error at FuncName %v GetEventInsight ListEventDateTimeInsight %v, user.ID: %v \n", funcName, err.Error(), user.ID)

		} else {
			eventMainDateTimeID = newID
		}
	}
	var count int
	var earning string
	var cancellation int
	payouts, err := server.store.ListAllChargeTicketReferencePayoutInsights(ctx, db.ListAllChargeTicketReferencePayoutInsightsParams{
		CoUserID:        tools.UuidToString(user.UserID),
		HostID:          user.ID,
		Date:            startMainTime,
		EventDateID:     eventMainDateTimeID,
		PaymentComplete: true,
	})
	if err != nil {
		log.Printf("Error at FuncName %v GetEventInsight tools.CountChargeTicketReferenceOnlyByStartDate %v, user.ID: %v \n", funcName, err.Error(), user.ID)
		count = 0
		cancellation = 0
		earning = "0.00"

	} else {
		itemCount, itemPrice, itemCancellation := GetAllTicketInsightCountAndPrice(payouts, dollarToNaira, dollarToCAD, req.Currency, funcName, user.ID)
		count = itemCount
		cancellation = itemCancellation
		earning = tools.ConvertFloatToString(itemPrice)
	}
	res = GetEventInsightRes{
		Earning:         earning,
		List:            resData,
		IsEmpty:         itemEmpty,
		ForOffset:       req.ForOffset,
		Cancellation:    cancellation,
		TicketSold:      count,
		StartYear:       startYear.Year(),
		Currency:        req.Currency,
		StartDate:       tools.ConvertDateOnlyToString(startMainTime),
		EventDateTimeID: tools.UuidToString(eventMainDateTimeID),
	}
	return
}

func GetTicketInsightCountAndPrice(data []db.ListChargeTicketReferencePayoutInsightsRow, dollarToNaira string, dollarToCAD string, userCurrency string, funcName string, userID uuid.UUID) (count int, price float64, cancelled int) {
	for _, d := range data {
		if d.Cancelled {
			// We don't want to count cancelled stuff
			cancelled += 1
			continue
		}
		priceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(d.Amount), d.Currency, userCurrency, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("Error at FuncName %v GetTicketInsightCountAndPrice ConvertPrice: %v, d.ID: %v \n", funcName, err.Error(), d.ChargeID)
			continue
		}
		count += 1
		price += priceFloat
	}
	return
}

func GetAllTicketInsightCountAndPrice(data []db.ListAllChargeTicketReferencePayoutInsightsRow, dollarToNaira string, dollarToCAD string, userCurrency string, funcName string, userID uuid.UUID) (count int, price float64, cancelled int) {
	for _, d := range data {
		if d.Cancelled {
			// We don't want to count cancelled stuff
			cancelled += 1
			continue
		}
		priceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(d.Amount), d.Currency, userCurrency, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("Error at FuncName %v GetTicketInsightCountAndPrice ConvertPrice: %v, d.ID: %v \n", funcName, err.Error(), d.ChargeID)
			continue
		}
		count += 1
		price += priceFloat
	}
	return
}

func GroupAllChargeOptionInsightByMonth(rows []db.ListAllOptionMainPayoutInsightsRow) map[string][]db.ListAllOptionMainPayoutInsightsRow {
	grouped := make(map[string][]db.ListAllOptionMainPayoutInsightsRow)

	// Sort rows by StartDate before grouping
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].StartDate.Before(rows[j].StartDate)
	})

	for _, row := range rows {
		// Format the month without the year
		month := row.StartDate.Format("January")
		grouped[month] = append(grouped[month], row)
	}

	return grouped
}

func GroupChargeOptionInsightByMonth(rows []db.ListOptionMainPayoutInsightsRow) map[string][]db.ListOptionMainPayoutInsightsRow {
	grouped := make(map[string][]db.ListOptionMainPayoutInsightsRow)

	// Sort rows by StartDate before grouping
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].StartDate.Before(rows[j].StartDate)
	})

	for _, row := range rows {
		// Format the month without the year
		month := row.StartDate.Format("January")
		grouped[month] = append(grouped[month], row)
	}

	return grouped
}

func GetChargeOptionInsightCountAndPrice(data []db.ListOptionMainPayoutInsightsRow, dollarToNaira string, dollarToCAD string, userID uuid.UUID, userCurrency string, funcName string) (count int, price float64) {
	for _, d := range data {
		if d.Cancelled {
			continue
		}
		priceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(d.Amount), d.Currency, userCurrency, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("Error at FuncName %v GetChargeOptionInsightCountAndPrice ConvertPrice: %v, d.ID: %v \n", funcName, err.Error(), d.ChargeID)
			continue
		}
		count += 1
		price += priceFloat
	}
	return
}

func GetAllChargeOptionInsightCountAndPrice(data []db.ListAllOptionMainPayoutInsightsRow, dollarToNaira string, dollarToCAD string, userID uuid.UUID, userCurrency string, funcName string) (count int, price float64) {
	for _, d := range data {
		if d.Cancelled {
			continue
		}
		priceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(d.Amount), d.Currency, userCurrency, dollarToNaira, dollarToCAD, userID)
		if err != nil {
			log.Printf("Error at FuncName %v GetChargeOptionInsightCountAndPrice ConvertPrice: %v, d.ID: %v \n", funcName, err.Error(), d.ChargeID)
			continue
		}
		count += 1
		price += priceFloat
	}
	return
}
