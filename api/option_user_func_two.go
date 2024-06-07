package api

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

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

func HandleExOptionReview(option db.OptionsInfo, server *Server, ctx *gin.Context, hostID uuid.UUID) (review UserExReview) {
	var environment float64
	var accuracy float64
	var communication float64
	var location float64
	var checkIn float64
	var general float64
	var resData []UserExReviewItem
	var isEmpty bool
	var fiveCount int
	var fourCount int
	var threeCount int
	var twoCount int
	var oneCount int
	reviewData, err := server.store.ListChargeOptionReview(ctx, option.OptionUserID)
	if err != nil || len(reviewData) == 0 {
		if err != nil {
			log.Printf("Error atHandleExOptionReview in ListOptionInfoReview err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		isEmpty = true
		data := UserExReviewItem{
			ID:               tools.UuidToString(uuid.New()),
			General:          tools.ConvertFloatToString(0.0),
			Environment:      tools.ConvertFloatToString(0.0),
			Accuracy:         tools.ConvertFloatToString(0.0),
			CheckIn:          tools.ConvertFloatToString(0.0),
			Communication:    tools.ConvertFloatToString(0.0),
			Location:         tools.ConvertFloatToString(0.0),
			PublicNote:       "none",
			HostPublicNote:   "none",
			Average:          tools.ConvertFloatToString(0.0),
			YearJoined:       "none",
			DateBooked:       "none",
			DateHostResponse: "none",
			ProfilePhoto:     "none",
			FirstName:        "none",
		}
		resData = append(resData, data)

	} else {
		log.Printf("review data good %v\n", reviewData)
		for _, rev := range reviewData {
			environment += float64(rev.Environment)
			accuracy += float64(rev.Accuracy)
			communication += float64(rev.Communication)
			location += float64(rev.Location)
			checkIn += float64(rev.CheckIn)
			general += float64(rev.General)
			average := (float64(rev.Environment) + float64(rev.Accuracy) + float64(rev.Communication) + float64(rev.Location) + float64(rev.CheckIn) + float64(rev.General)) / 6
			switch int(average) {
			case 5:
				fiveCount += 1
			case 4:
				fourCount += 1
			case 3:
				threeCount += 1
			case 2:
				twoCount += 1
			case 1:
				oneCount += 1
			}
			data := UserExReviewItem{
				ID:               tools.UuidToString(uuid.New()),
				General:          tools.ConvertInt32ToString(rev.General),
				Environment:      tools.ConvertInt32ToString(rev.Environment),
				Accuracy:         tools.ConvertInt32ToString(rev.Accuracy),
				CheckIn:          tools.ConvertInt32ToString(rev.CheckIn),
				Communication:    tools.ConvertInt32ToString(rev.Communication),
				Location:         tools.ConvertInt32ToString(rev.Location),
				PublicNote:       rev.PublicNote,
				HostPublicNote:   "none",
				Average:          tools.ConvertFloatToString(average),
				YearJoined:       tools.ConvertDateOnlyToString(rev.UserJoined),
				DateBooked:       tools.ConvertDateOnlyToString(rev.DateBooked),
				DateHostResponse: "none",
				ProfilePhoto:     rev.Photo,
				FirstName:        rev.FirstName,
			}
			resData = append(resData, data)
		}
		// We divide it to get an average
		environment = environment / float64(len(reviewData))
		accuracy = accuracy / float64(len(reviewData))
		communication = communication / float64(len(reviewData))
		location = location / float64(len(reviewData))
		checkIn = checkIn / float64(len(reviewData))
		general = general / float64(len(reviewData))
	}
	average := (environment + accuracy + communication + location + checkIn + general) / 6
	review = UserExReview{
		Total:         len(reviewData),
		Five:          fiveCount,
		Four:          fourCount,
		Three:         threeCount,
		Two:           twoCount,
		One:           oneCount,
		Count:         len(resData),
		Environment:   tools.ConvertFloatToString(environment),
		Accuracy:      tools.ConvertFloatToString(accuracy),
		Communication: tools.ConvertFloatToString(communication),
		Location:      tools.ConvertFloatToString(location),
		CheckIn:       tools.ConvertFloatToString(checkIn),
		General:       tools.ConvertFloatToString(general),
		Average:       tools.ConvertFloatToString(average),
		List:          resData,
		IsEmpty:       isEmpty,
	}
	return
}

func HandleExEventReview(option db.OptionsInfo, server *Server, ctx *gin.Context, hostID uuid.UUID) (review UserExReview) {
	var environment float64
	var accuracy float64
	var communication float64
	var location float64
	var checkIn float64
	var general float64
	var resData []UserExReviewItem
	var isEmpty bool = true
	var fiveCount int
	var fourCount int
	var threeCount int
	var twoCount int
	var oneCount int
	var total int
	data := UserExReviewItem{
		ID:               tools.UuidToString(uuid.New()),
		General:          tools.ConvertFloatToString(0.0),
		Environment:      tools.ConvertFloatToString(0.0),
		Accuracy:         tools.ConvertFloatToString(0.0),
		CheckIn:          tools.ConvertFloatToString(0.0),
		Communication:    tools.ConvertFloatToString(0.0),
		Location:         tools.ConvertFloatToString(0.0),
		PublicNote:       "none",
		HostPublicNote:   "none",
		Average:          tools.ConvertFloatToString(0.0),
		YearJoined:       "none",
		DateBooked:       "none",
		DateHostResponse: "none",
		ProfilePhoto:     "none",
		FirstName:        "none",
	}
	resData = append(resData, data)
	average := (environment + accuracy + communication + location + checkIn + general) / 6
	review = UserExReview{
		Total:         total,
		Five:          fiveCount,
		Four:          fourCount,
		Three:         threeCount,
		Two:           twoCount,
		One:           oneCount,
		Count:         len(resData),
		Environment:   tools.ConvertFloatToString(environment),
		Accuracy:      tools.ConvertFloatToString(accuracy),
		Communication: tools.ConvertFloatToString(communication),
		Location:      tools.ConvertFloatToString(location),
		CheckIn:       tools.ConvertFloatToString(checkIn),
		General:       tools.ConvertFloatToString(general),
		Average:       tools.ConvertFloatToString(average),
		List:          resData,
		IsEmpty:       isEmpty,
	}
	return
}

func HandleListOptionExReview(ctx context.Context, server *Server, req ListExReviewDetailReq) (res ListExReviewDetailRes, hasData bool, err error) {
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at  HandleListOptionExReview in tools.StringToUuid err: %v, user: %v\n", err, req.OptionUserID)
		hasData = false
		return
	}
	count, err := server.store.CountChargeOptionReviewIndex(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at  HandleListOptionExReview in CountChargeOptionReviewIndex err: %v, user: %v\n", err, optionUserID)
		hasData = false
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		hasData = false
		return
	}
	reviews, err := server.store.ListChargeOptionReviewIndex(ctx, db.ListChargeOptionReviewIndexParams{
		OptionUserID: optionUserID,
		Limit:        10,
		Offset:       int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error atHandleListOptionExReview in ListChargeOptionReviewIndex err: %v, user: %v\n", err, optionUserID)
		hasData = false
		err = nil
		return
	}
	var resData []UserExReviewItem

	for _, rev := range reviews {
		average := (float64(rev.Environment) + float64(rev.Accuracy) + float64(rev.Communication) + float64(rev.Location) + float64(rev.CheckIn) + float64(rev.General)) / 6
		data := UserExReviewItem{
			ID:               tools.UuidToString(uuid.New()),
			General:          tools.ConvertInt32ToString(rev.General),
			Environment:      tools.ConvertInt32ToString(rev.Environment),
			Accuracy:         tools.ConvertInt32ToString(rev.Accuracy),
			CheckIn:          tools.ConvertInt32ToString(rev.CheckIn),
			Communication:    tools.ConvertInt32ToString(rev.Communication),
			Location:         tools.ConvertInt32ToString(rev.Location),
			PublicNote:       rev.PublicNote,
			HostPublicNote:   "none",
			Average:          tools.ConvertFloatToString(average),
			YearJoined:       tools.ConvertDateOnlyToString(rev.UserJoined),
			DateBooked:       tools.ConvertDateOnlyToString(rev.DateBooked),
			DateHostResponse: "none",
			ProfilePhoto:     rev.Photo,
			FirstName:        rev.FirstName,
		}
		resData = append(resData, data)
	}
	onLastIndex := false
	hasData = true
	if count <= int64(req.Offset+len(reviews)) {
		onLastIndex = true
	}
	res = ListExReviewDetailRes{
		List:        resData,
		Offset:      req.Offset + len(reviews),
		OnLastIndex: onLastIndex,
	}
	return
}

func (server *Server) ListExReviewDetail(ctx *gin.Context) {
	var req ListExReviewDetailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListReserveUserItem in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ListExReviewDetailRes
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleListOptionExReview(ctx, server, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resOption
		}
	case "events":
		data := "none"
		ctx.JSON(http.StatusNoContent, data)
		return
	}

	ctx.JSON(http.StatusOK, res)

}
