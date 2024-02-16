package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleHostSpecialDates(ctx *gin.Context, server *Server, option db.GetOptionInfoByOptionWithPriceUserIDRow, userCurrency string) (dates []ExOptionDateTimeItem) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	special, err := server.store.ListOptionDateTimeMore(ctx, option.ID)
	if err != nil || len(special) == 0 {
		if err != nil {
			log.Printf("Error at HandleHostSpecialDates in ListAllOptionDateTime err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		return
	}
	for _, s := range special {
		price, err := tools.ConvertPrice(tools.IntToMoneyString(s.Price), option.Currency, userCurrency, dollarToNaira, dollarToCAD, option.ID)
		if err != nil {
			log.Printf("Error at HandleHostSpecialDates in tools.ConvertPrice err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
			price = tools.ConvertStringToFloat(tools.IntToMoneyString(s.Price))
		}
		data := ExOptionDateTimeItem{
			Date:      tools.ConvertDateOnlyToString(s.Date),
			Available: s.Available,
			Price:     tools.ConvertFloatToString(price),
			IsEmpty:   false,
		}
		dates = append(dates, data)
	}
	return
}

func SortOptionDateItem(items []ExOptionDateTimeItem) []ExOptionDateTimeItem {
	// Custom sorting function based on the Date field
	sort.Slice(items, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", items[i].Date)
		dateJ, _ := time.Parse("2006-01-02", items[j].Date)
		return dateI.Before(dateJ)
	})
	return items
}

func HandleImportantDates(dates []ExOptionDateTimeItem, optionPrice string) (busyDates []ExBusyDate, priceDates []ExPriceDate, busyIsEmpty bool, priceIsEmpty bool) {
	dates = SortOptionDateItem(dates)
	start := "none"
	end := "none"
	for i, d := range dates {
		actualDate, err := time.Parse("2006-01-02", d.Date)
		if err != nil {
			continue
		}
		if time.Now().After(actualDate) {
			continue
		}
		// We check if the price is different from the option price
		if d.Price != optionPrice {
			data := ExPriceDate{
				Date:    d.Date,
				Price:   d.Price,
				IsEmpty: false,
			}
			priceDates = append(priceDates, data)
		}
		// We handle busy dates
		if !d.Available {
			if start == "none" {
				start = d.Date
				if len(dates) == i+1 {
					data := ExBusyDate{
						StartDate: start,
						EndDate:   start,
						IsEmpty:   false,
					}
					busyDates = append(busyDates, data)
				}
				continue
			} else if end == "none" {
				result, err := tools.FindDateDifference(start, d.Date)
				if err != nil {
					continue
				}
				if result <= 1 {
					end = d.Date
				} else {
					data := ExBusyDate{
						StartDate: start,
						EndDate:   start,
						IsEmpty:   false,
					}
					busyDates = append(busyDates, data)
					start = d.Date
					end = "none"
				}
				if len(dates) == i+1 {
					data := ExBusyDate{
						StartDate: start,
						EndDate:   start,
						IsEmpty:   false,
					}
					busyDates = append(busyDates, data)
				}
				continue
			} else {
				result, err := tools.FindDateDifference(end, d.Date)
				if err != nil {
					continue
				}
				if result <= 1 {
					end = d.Date
					if len(dates) == i+1 {
						data := ExBusyDate{
							StartDate: start,
							EndDate:   end,
							IsEmpty:   false,
						}
						busyDates = append(busyDates, data)
					}
					continue
				} else {
					data := ExBusyDate{
						StartDate: start,
						EndDate:   end,
						IsEmpty:   false,
					}
					busyDates = append(busyDates, data)
					start = d.Date
					end = "none"
				}
			}
		} else {
			if start != "none" {
				if end != "none" {
					if len(dates) == i+1 {
						data := ExBusyDate{
							StartDate: start,
							EndDate:   end,
							IsEmpty:   false,
						}
						busyDates = append(busyDates, data)
					}
				} else {
					if len(dates) == i+1 {
						data := ExBusyDate{
							StartDate: start,
							EndDate:   start,
							IsEmpty:   false,
						}
						busyDates = append(busyDates, data)
					}
				}
			}
		}
	}
	if len(busyDates) == 0 {
		data := ExBusyDate{"none", "none", true}
		busyIsEmpty = true
		busyDates = append(busyDates, data)
	}
	if len(priceDates) == 0 {
		data := ExPriceDate{"none", "none", true}
		priceIsEmpty = true
		priceDates = append(priceDates, data)
	}
	return
}

func HandleExAvailable(ctx *gin.Context, server *Server, optionUserID uuid.UUID, prepareTime string, autoBlock bool) (dates []ExOptionDateTimeItem) {
	charges, err := server.store.ListChargeOptionReferenceByOptionUserID(ctx, db.ListChargeOptionReferenceByOptionUserIDParams{
		OptionUserID: optionUserID,
		Cancelled:    false,
		IsComplete:   true,
	})
	if err != nil || len(charges) == 0 {
		if err != nil {
			log.Printf("Error at HandleExPrepareTime in ListChargeOptionReferenceByOptionUserID err: %v, user: %v, optionUserID: %v\n", err, ctx.ClientIP(), optionUserID)
		}
		return
	}
	if autoBlock {
		// First we handle dates that are booked
		bookDates := GetExDateTime(prepareTime, ctx.ClientIP(), charges, tools.UuidToString(optionUserID))
		if len(bookDates) != 0 {
			dates = append(dates, bookDates...)
		}
	}
	return
}
