package api

import (
	"context"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) GetEventFilterRange(ctx *gin.Context) {
	var req ExFilterRangeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionFilterRange in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Type)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var maxPrice float64
	var minPrice float64
	var averagePrice float64
	var addMaxPrice float64
	var addMinPrice float64
	var averageAddPrice float64
	minPrice, maxPrice, err := GetEventFilterMaxMinPrice(ctx, server, req, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We calculate the average
	averagePrice = (maxPrice + minPrice) / 2

	// We handle it when it is 5 days
	addMaxPrice = maxPrice * float64(server.config.OptionExDayCount)
	addMinPrice = minPrice * float64(server.config.OptionExDayCount)
	averageAddPrice = (addMaxPrice + addMinPrice) / 2
	res := ExFilterRangeRes{
		MaxPrice:        tools.ConvertFloatToString(maxPrice),
		MinPrice:        tools.ConvertFloatToString(minPrice),
		AveragePrice:    tools.ConvertFloatToString(averagePrice),
		AddMaxPrice:     tools.ConvertFloatToString(addMaxPrice),
		AddMinPrice:     tools.ConvertFloatToString(addMinPrice),
		AddDayCount:     server.config.OptionExDayCount,
		AverageAddPrice: tools.ConvertFloatToString(averageAddPrice),
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetOptionFilterRange(ctx *gin.Context) {
	var req ExFilterRangeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionFilterRange in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Type)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var maxPrice float64
	var minPrice float64
	var averagePrice float64
	var addMaxPrice float64
	var addMinPrice float64
	var averageAddPrice float64
	minPrice, maxPrice, err := GetOptionFilterMaxMinPrice(ctx, server, req, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We calculate the average
	averagePrice = (maxPrice + minPrice) / 2

	// We handle it when it is 5 days
	addMaxPrice = maxPrice * float64(server.config.OptionExDayCount)
	addMinPrice = minPrice * float64(server.config.OptionExDayCount)
	averageAddPrice = (addMaxPrice + addMinPrice) / 2
	res := ExFilterRangeRes{
		MaxPrice:        tools.ConvertFloatToString(maxPrice),
		MinPrice:        tools.ConvertFloatToString(minPrice),
		AveragePrice:    tools.ConvertFloatToString(averagePrice),
		AddMaxPrice:     tools.ConvertFloatToString(addMaxPrice),
		AddMinPrice:     tools.ConvertFloatToString(addMinPrice),
		AddDayCount:     server.config.OptionExDayCount,
		AverageAddPrice: tools.ConvertFloatToString(averageAddPrice),
	}
	ctx.JSON(http.StatusOK, res)

}

func HandleOptionFilter(ctx context.Context, server *Server, optionID uuid.UUID, optionUserID uuid.UUID, req ExFilterOptionRequest, basePrice int64, optionCurrency string, funcName string, shortletSpaceType string, optionCategory []string, optionCanInstantBook bool, checkInMethod string) (confirm bool) {
	basePriceFloat, err := tools.ConvertPrice(tools.IntToMoneyString(basePrice), optionCurrency, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, optionUserID)
	if err != nil {
		log.Printf("Error at FuncName %v, HandleOptionFilter tools.ConvertPrice err: %v \n", funcName, err.Error())
		confirm = false
		return
	}
	// First we deal with the price range
	var confirmPrice bool
	if req.OnAddPrice {
		addMaxPrice := tools.ConvertStringToFloat(req.AddMaxPrice)
		addMinPrice := tools.ConvertStringToFloat(req.AddMinPrice)
		if addMaxPrice == 0 && addMinPrice == 0 {

			confirmPrice = true
		} else if addMinPrice <= basePriceFloat && addMaxPrice == 0 {
			confirmPrice = true
		}  else if addMinPrice <= basePriceFloat && basePriceFloat <= addMaxPrice {
			confirmPrice = true
		}
	} else {
		maxPrice := tools.ConvertStringToFloat(req.MaxPrice)
		minPrice := tools.ConvertStringToFloat(req.MinPrice)
		if minPrice == 0 && maxPrice == 0 {
			log.Printf("1 passed price maxPrice: %v, minPrice: %v, basePrice:%v \n", maxPrice, minPrice, basePriceFloat)
			confirmPrice = true
		} else if minPrice <= basePriceFloat && maxPrice == 0 {
			log.Printf("2 passed price maxPrice: %v, minPrice: %v, basePrice:%v \n", maxPrice, minPrice, basePriceFloat)
			confirmPrice = true
		}  else if minPrice <= basePriceFloat && basePriceFloat <= maxPrice {
			log.Printf("3 passed price maxPrice: %v, minPrice: %v, basePrice:%v \n", maxPrice, minPrice, basePriceFloat)
			confirmPrice = true
		}
	}
	// Next we deal with shortlet space type that is things like entire space, private room
	var shortletSpaceConfirm bool
	if len(req.ShortletSpaceType) == 0 {
		shortletSpaceConfirm = true
	} else {
		for _, s := range req.ShortletSpaceType {
			if s == shortletSpaceType {
				shortletSpaceConfirm = true
				break
			}
		}
	}

	// Next we handle amenities
	var amenityConfirm bool
	if len(req.Amenities) == 0 {
		amenityConfirm = true
	} else {
		amenities, err := server.store.ListAmenitiesTag(ctx, db.ListAmenitiesTagParams{
			OptionID: optionID,
			HasAm:    true,
		})
		if err != nil || len(amenities) == 0 {
			if err != nil {
				log.Printf("Error at FuncName %v, HandleOptionFilter .ListAmenities err: %v \n", funcName, err.Error())
			}
			amenityConfirm = true
		} else {
			var exist = true
			for _, ma := range req.Amenities {
				if !tools.IsInList(amenities, ma) {
					exist = false
				}
			}
			if exist {
				amenityConfirm = true
			}
		}
	}

	// Next we handle category types
	var categoryConfirm bool
	if len(req.CategoryTypes) == 0 {
		categoryConfirm = true
	} else {
		var exist = false
		for _, ca := range req.CategoryTypes {
			if tools.IsInList(optionCategory, ca) {
				exist = true
			}
		}
		if exist {
			categoryConfirm = true
		}
	}

	// Next we handle space areas
	var spaceConfirm = OptionFilterRooms(ctx, server, req, optionID, funcName)

	// Next we handle instant check in
	var instantBookConfirm bool
	if req.CanInstantBook {
		// We only want to check if the user switch it on
		instantBookConfirm = req.CanInstantBook == optionCanInstantBook
	} else {
		instantBookConfirm = true
	}

	// Next we handle self-check in
	var selfCheckIn bool
	if req.CanSelfCheck {
		// We only want to check if the user switch it on
		selfCheckIn = "self_check_in" == checkInMethod
	} else {
		selfCheckIn = true
	}
	return confirmPrice && shortletSpaceConfirm && amenityConfirm && categoryConfirm && spaceConfirm && instantBookConfirm && selfCheckIn
}

func OptionFilterRooms(ctx context.Context, server *Server, req ExFilterOptionRequest, optionID uuid.UUID, funcName string) bool {
	var bathroomConfirm bool
	var bedroomConfirm bool
	var bedConfirm bool
	if req.Bathrooms == 0 && req.Bedrooms == 0 && req.Beds == 0 {
		return true
	} else {
		spaceAreas, err := server.store.ListSpaceAreaRooms(ctx, optionID)
		if err != nil || len(spaceAreas) == 0 {
			if err != nil {
				log.Printf("Error at FuncName %v, HandleOptionFilter .ListSpaceAreaRooms err: %v \n", funcName, err.Error())
			}
			return true
		} else {
			var spaceData = make(map[string]int)
			for _, sa := range spaceAreas {
				spaceData[sa.SpaceType] += 1
				spaceData["bed"] += len(sa.Beds)
			}
			// Bathrooms
			if req.Bathrooms == 0 {
				bathroomConfirm = true
			} else {
				if req.Bathrooms == 8 && spaceData["full_bathroom"] >= 8 {
					bathroomConfirm = true
				} else if req.Bathrooms == spaceData["full_bathroom"] {
					bathroomConfirm = true
				}
			}
			// Bedrooms
			if req.Bedrooms == 0 {
				bedroomConfirm = true
			} else {
				if req.Bedrooms == 8 && spaceData["bedroom"] >= 8 {
					bedroomConfirm = true
				} else if req.Bedrooms == spaceData["bedroom"] {
					bedroomConfirm = true
				}
			}
			// Beds
			if req.Beds == 0 {
				bedConfirm = true
			} else {
				if req.Beds == 8 && spaceData["bed"] >= 8 {
					bedConfirm = true
				} else if req.Beds == spaceData["bed"] {
					bedConfirm = true
				}
			}
		}
	}
	log.Println(bedConfirm)
	return bedroomConfirm && bathroomConfirm && bedConfirm
}

func GetEventFilterMaxMinPrice(ctx context.Context, server *Server, req ExFilterRangeReq, clientIP string) (minPrice float64, maxPrice float64, err error) {
	eventPrices, err := server.store.ListTicketForRange(ctx)
	if err != nil || len(eventPrices) == 0 {
		if err != nil {
			log.Printf("Error at GetEventFilterRange in ListTicketForRange err: %v, user: %v\n", err, clientIP)
		}
		minPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(server.config.EventMinPriceNaira)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error 0 at min FuncName %v, HandleEventSearchPrice  tools.ConvertPrice err: %v \n", "GetEventFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
		maxPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(server.config.EventMaxPriceNaira)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error 0 max at FuncName %v, HandleEventSearchPrice  tools.ConvertPrice err: %v \n", "GetEventFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
	} else {
		minPriceInt := eventPrices[0].Price
		maxPriceInt := eventPrices[len(eventPrices)-1].Price
		if minPriceInt == maxPriceInt {
			// We want to ensure it is not equal
			maxPriceInt = int64(server.config.EventMaxPriceNaira)
		}
		minPrice, err = tools.ConvertPrice(tools.IntToMoneyString(int64(minPriceInt)), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error at min FuncName %v, HandleEventSearchPrice  tools.ConvertPrice err: %v \n", "GetEventFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}
		maxPrice, err = tools.ConvertPrice(tools.IntToMoneyString(maxPriceInt), utils.NGN, req.Currency, server.config.DollarToNaira, server.config.DollarToCAD, uuid.New())
		if err != nil {
			log.Printf("Error at max FuncName %v, HandleEventSearchPrice  tools.ConvertPrice err: %v \n", "GetEventFilterRange", err.Error())
			err = fmt.Errorf("range not found")
			return
		}

	}
	return
}
