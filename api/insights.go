package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListOptionInsight(ctx *gin.Context) {
	var req ListOptionInsightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListOptionInsight in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountOptionInfoInsight(ctx, db.CountOptionInfoInsightParams{
		HostID:         user.ID,
		IsComplete:     true,
		IsActive:       true,
		CoUserID:       tools.UuidToString(user.UserID),
		MainOptionType: req.MainOption,
	})
	if err != nil {
		log.Printf("Error at  HandleListOptionSelectComplete in CountOptionInfoInsight err: %v, user: %v\n", err, user.ID)
		if err == db.ErrorRecordNotFound {
			ctx.JSON(http.StatusNoContent, "none")
			return
		}
		err = fmt.Errorf("could not perform your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if count <= int64(req.Offset) {
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	// we want to get the options that are complete
	optionInfos, err := server.store.ListOptionInfoInsight(ctx, db.ListOptionInfoInsightParams{
		HostID:         user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		IsComplete:     true,
		Limit:          10,
		MainOptionType: req.MainOption,
		Offset:         int32(req.Offset),
		IsActive:       true,
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			ctx.JSON(http.StatusNoContent, "none")
			return
		} else {
			err = fmt.Errorf("could not perform your request")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var resData []OptionInsightItem
	for _, data := range optionInfos {
		var isCoHost bool
		if data.HostType == "co_host" {
			isCoHost = true
		}
		newData := OptionInsightItem{
			HostNameOption: data.HostNameOption,
			CoverImage:     data.CoverImage,
			OptionUserID:   tools.UuidToString(data.OptionUserID),
			MainOptionType: data.MainOptionType,
			HasName:        true,
			IsCoHost:       isCoHost,
		}
		resData = append(resData, newData)
	}
	res := ListOptionInsightRes{
		List:       resData,
		MainOption: req.MainOption,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionInsight(ctx *gin.Context) {
	var req GetOptionInsightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionInsight in ShouldBindJSON: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startYear, err := server.store.GetOptionInfoStartYear(ctx, db.GetOptionInfoStartYearParams{
		HostID:         user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		IsComplete:     true,
		IsActive:       true,
		MainOptionType: "options",
	})
	if err != nil {
		log.Printf("Error at GetOptionInfoStartYear tools.StringToUuid: %v, ClientID: %v \n", err.Error(), ctx.ClientIP())
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payouts, err := server.store.ListOptionMainPayoutInsights(ctx, db.ListOptionMainPayoutInsightsParams{
		CoUserID:     tools.UuidToString(user.UserID),
		HostID:       user.ID,
		OptionUserID: optionUserID,
		Year:         int32(req.Year),
	})
	if err != nil || len(payouts) == 0 {
		log.Printf("Error at ListOptionMainPayoutInsights tools.StringToUuid: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				ctx.JSON(http.StatusNoContent, "none")
				return
			} else {
				err = fmt.Errorf("could not perform your request")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusNoContent, "none")
		return
	}

	var resData []GetOptionInsightItem
	var earning float64
	var booking int
	var cancellation int
	var monthExist bool
	var newMonth string
	// We start by grouping the data
	grouped := GroupChargeOptionInsightByMonth(payouts)
	for month, data := range grouped {
		if len(data) == 0 {
			continue
		}
		if len(newMonth) == 0 {
			// We set a new month just incase the month we have doesn't have any data
			newMonth = month
		}
		itemCount, itemPrice := GetChargeOptionInsightCountAndPrice(data, dollarToNaira, dollarToCAD, user.ID, req.Currency, "GetOptionInsight")
		item := GetOptionInsightItem{
			Month:   month,
			Count:   itemCount,
			Earning: tools.ConvertFloatToString(itemPrice),
		}
		resData = append(resData, item)
		if month == tools.CapitalizeFirstCharacter(req.Month) {
			monthExist = true
			for _, v := range data {
				if v.Cancelled {
					cancellation += 1
				} else {
					price := tools.IntToMoneyString(v.Amount)
					priceFloat, err := tools.ConvertPrice(price, v.Currency, req.Currency, dollarToNaira, dollarToCAD, user.ID)
					if err != nil {
						log.Printf("Error at tools.ConvertPrice(: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
						continue
					}
					earning += priceFloat
					booking += 1
				}
			}
		}
	}
	var setMonth string
	if req.FromMonth {
		setMonth = req.Month
	} else {
		if monthExist {
			setMonth = req.Month

		} else {
			setMonth = newMonth
		}
	}
	res := GetOptionInsightRes{
		Earning:      tools.ConvertFloatToString(earning),
		Booking:      booking,
		Cancellation: cancellation,
		List:         resData,
		Currency:     req.Currency,
		StartYear:    startYear.Year(),
		Month:        setMonth,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetAllOptionInsight(ctx *gin.Context) {
	var req GetAllOptionInsightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionInsight in ShouldBindJSON: %v, Client: %v \n", err.Error(), ctx.ClientIP())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startYear, err := server.store.GetOptionInfoStartYear(ctx, db.GetOptionInfoStartYearParams{
		HostID:         user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		IsComplete:     true,
		IsActive:       true,
		MainOptionType: "options",
	})
	if err != nil {
		log.Printf("Error at GetOptionInfoStartYear tools.StringToUuid: %v, ClientID: %v \n", err.Error(), ctx.ClientIP())
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	payouts, err := server.store.ListAllOptionMainPayoutInsights(ctx, db.ListAllOptionMainPayoutInsightsParams{
		CoUserID: tools.UuidToString(user.UserID),
		HostID:   user.ID,
		Year:     int32(req.Year),
	})
	if err != nil || len(payouts) == 0 {
		if err != nil {
			if err == db.ErrorRecordNotFound {
				ctx.JSON(http.StatusNoContent, "none")
				return
			} else {
				err = fmt.Errorf("could not perform your request")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusNoContent, "none")
		return
	}

	var resData []GetOptionInsightItem
	var earning float64
	var booking int
	var cancellation int
	var monthExist bool
	var newMonth string
	// We start by grouping the data
	grouped := GroupAllChargeOptionInsightByMonth(payouts)
	for month, data := range grouped {

		if len(data) == 0 {

			continue
		}
		if len(newMonth) == 0 {
			// We set a new month just incase the month we have doesn't have any data
			newMonth = month
		}
		itemCount, itemPrice := GetAllChargeOptionInsightCountAndPrice(data, dollarToNaira, dollarToCAD, user.ID, req.Currency, "GetOptionInsight")
		item := GetOptionInsightItem{
			Month:   month,
			Count:   itemCount,
			Earning: tools.ConvertFloatToString(itemPrice),
		}
		resData = append(resData, item)
		if month == tools.CapitalizeFirstCharacter(req.Month) {
			monthExist = true
			for _, v := range data {
				if v.Cancelled {
					cancellation += 1
				} else {
					price := tools.IntToMoneyString(v.Amount)
					priceFloat, err := tools.ConvertPrice(price, v.Currency, req.Currency, dollarToNaira, dollarToCAD, user.ID)
					if err != nil {
						log.Printf("Error at tools.ConvertPrice(: %v, clientID: %v \n", err.Error(), ctx.ClientIP())
						continue
					}
					earning += priceFloat
					booking += 1
				}
			}
		}
	}
	var setMonth string
	if req.FromMonth {
		setMonth = req.Month
	} else {
		if monthExist {
			setMonth = req.Month

		} else {
			setMonth = newMonth
		}
	}
	res := GetOptionInsightRes{
		Earning:      tools.ConvertFloatToString(earning),
		Booking:      booking,
		Cancellation: cancellation,
		List:         resData,
		Currency:     req.Currency,
		StartYear:    startYear.Year(),
		Month:        setMonth,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventInsight(ctx *gin.Context) {
	var req GetEventInsightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionInsight in ShouldBindJSON: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startYear, err := server.store.GetOptionInfoStartYear(ctx, db.GetOptionInfoStartYearParams{
		HostID:         user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		IsComplete:     true,
		IsActive:       true,
		MainOptionType: "events",
	})
	if err != nil {
		log.Printf("Error at GetEventInsight GetOptionInfoStartYear %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	optionUserID, err := tools.StringToUuid(req.OptionUserID)
	if err != nil {
		log.Printf("Error at GetEventInsight GetOptionInfoStartYear tools.StringToUuid: %v, OptionUserID: %v \n", err.Error(), req.OptionUserID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := HandleEventInsights(ctx, server, user, req, optionUserID, dollarToNaira, dollarToCAD, startYear, "GetEventInsight")

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetAllEventInsight(ctx *gin.Context) {
	var req GetAllEventInsightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionInsight in ShouldBindJSON: %v, ClientIP: %v \n", err.Error(), ctx.ClientIP())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	log.Println(req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startYear, err := server.store.GetOptionInfoStartYear(ctx, db.GetOptionInfoStartYearParams{
		HostID:         user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		IsComplete:     true,
		IsActive:       true,
		MainOptionType: "events",
	})
	if err != nil {
		log.Printf("Error at GetEventInsight GetOptionInfoStartYear %v, UserID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := HandleAllEventInsights(ctx, server, user, req, dollarToNaira, dollarToCAD, startYear, "GetAllEventInsight")

	ctx.JSON(http.StatusOK, res)
}
