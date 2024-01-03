package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/val"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// This has the server functions for booking mostly

// This gets the option price
func (server *Server) GetOptionPrice(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in GetOptionPrice in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("GetOptionPrice Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	price, err := server.store.GetOptionPrice(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionPrice at GetOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := OptionPriceRes{
		Price:        tools.IntToMoneyString(price.Price),
		WeekendPrice: tools.IntToMoneyString(price.WeekendPrice),
	}
	ctx.JSON(http.StatusOK, res)
}

// Updates option price
func (server *Server) UpdateOptionPrice(ctx *gin.Context) {
	var req UpdateOptionPriceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionPrice in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	price, err := server.store.UpdateOptionPrice(ctx, db.UpdateOptionPriceParams{
		Price: pgtype.Int8{
			Int64: tools.MoneyStringToInt(req.Price),
			Valid: true,
		},
		WeekendPrice: pgtype.Int8{
			Int64: tools.MoneyStringToInt(req.WeekendPrice),
			Valid: true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateOptionPrice at UpdateOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := OptionPriceRes{
		Price:        tools.IntToMoneyString(price.Price),
		WeekendPrice: tools.IntToMoneyString(price.WeekendPrice),
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionPrice", "listing price", "update listing price")
	}
	ctx.JSON(http.StatusOK, res)
}

// This list option additional charge
func (server *Server) ListOptionAddCharge(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ListOptionAddCharge in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("ListOptionAddCharge Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	charges, err := server.store.ListOptionAddCharge(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at ListOptionAddCharge at ListOptionAddCharge: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		// This error can be expected
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	petsAllowed, err := server.store.GetOptionInfoDetailPetsAllow(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at ListOptionAddCharge at ListOptionAddCharge: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		petsAllowed = false
	}
	var resData []OptionAddChargeItem
	for _, charge := range charges {
		data := OptionAddChargeItem{
			ID:         tools.UuidToString(charge.ID),
			MainFee:    tools.IntToMoneyString(charge.MainFee),
			ExtraFee:   tools.IntToMoneyString(charge.ExtraFee),
			NumOfGuest: int(charge.NumOfGuest),
			Type:       charge.Type,
		}
		resData = append(resData, data)
	}
	res := ListOptionAddChargeRes{
		List:        resData,
		PetsAllowed: petsAllowed,
	}
	ctx.JSON(http.StatusOK, res)
}

// This would create a charge
func (server *Server) CreateUpdateOptionAddCharge(ctx *gin.Context) {
	var req CreateUpdateOptionAddChargeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionAddChargeParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	_, err = server.store.GetOptionAddCharge(ctx, db.GetOptionAddChargeParams{
		OptionID: option.ID,
		Type:     req.Type,
	})

	var res OptionAddChargeItem
	if err != nil {
		// We expect there should be an error here
		// Because we believe it should not be created
		// We would want to create the charge
		res, err = HandleCreateOptionAddCharge(ctx, server, option, user, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		// Since there is no err we believe it is to update
		res, err = HandleUpdateOptionAddCharge(ctx, server, option, user, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateOptionAddCharge", "listing additional charge", "create-update listing additional charge")
	}
	ctx.JSON(http.StatusOK, res)
}

// Updates option detail pets allowed
func (server *Server) UpdatePetsAllowed(ctx *gin.Context) {
	var req UpdatePetsAllowedParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdatePetsAllowed in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	optionDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		PetsAllowed: pgtype.Bool{
			Bool:  req.PetsAllowed,
			Valid: true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdatePetsAllowed at UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdatePetsAllowedRes{
		PetsAllowed: optionDetail.PetsAllowed,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdatePetsAllowed", "pets allowed", "update pets allowed")
	}
	ctx.JSON(http.StatusOK, res)
}

// This would create and or update a option discounts
// LOT means length of stay
func (server *Server) LOTCreateUpdateOptionDiscount(ctx *gin.Context) {
	var req LOTCreateUpdateOptionDiscountParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at LOTCreateUpdateOptionDiscount in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	res, err, hasData := HandleLOTCreateUpdateOptionDiscount(ctx, server, option, user, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !hasData {
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "LOTCreateUpdateOptionDiscount", "listing discount", "create-update listing discount")
	}
	ctx.JSON(http.StatusOK, res)
}

// This would list all the option discounts if any
// LOT means length of stay
func (server *Server) LOTListOptionDiscount(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at LOTListOptionDiscount in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {

		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	discounts, err := server.store.ListOptionDiscountByMainType(ctx, db.ListOptionDiscountByMainTypeParams{
		OptionID: option.ID,
		MainType: "length_of_stay",
	})
	if err != nil {
		log.Printf("There an error at LOTListOptionDiscount at ListOptionDiscountByMainType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	var resData []OptionDiscountItem
	for _, d := range discounts {
		data := OptionDiscountItem{
			ID:        tools.UuidToString(d.ID),
			Type:      d.Type,
			MainType:  d.MainType,
			Percent:   int(d.Percent),
			Name:      d.Name,
			ExtraType: d.ExtraType,
			Des:       d.Des,
		}
		resData = append(resData, data)
	}
	res := ListOptionDiscountRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateUnlistedOptionCurrency(ctx *gin.Context) {
	var req UpdateOptionCurrencyParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionCurrencyParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	currency, err := server.store.UpdateOptionInfoCurrency(ctx, db.UpdateOptionInfoCurrencyParams{
		ID:       option.ID,
		Currency: req.Currency,
	})
	if err != nil {
		log.Printf("There an error at UpdateUnlistedOptionCurrency at UpdateOptionInfoCurrency: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not update the currency ")
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}
	res := UpdateUnlistedOptionCurrencyRes{
		Currency: currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateOptionCurrency(ctx *gin.Context) {
	var req UpdateOptionCurrencyParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionCurrencyParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	oldCurrency := option.Currency
	currency, err := server.store.UpdateOptionInfoCurrency(ctx, db.UpdateOptionInfoCurrencyParams{
		ID:       option.ID,
		Currency: req.Currency,
	})
	if err != nil {
		log.Printf("There an error at UpdateOptionCurrency at UpdateOptionInfoCurrency: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not update the currency ")
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}
	log.Printf("New currency, %v. old currency %v\n", currency, oldCurrency)
	// we want to update the price to the local currency
	switch option.MainOptionType {
	case "options":
		UpdateOptionDataCurrency(currency, ctx, server, option, oldCurrency, user.ID)
		UpdateOptionAddChargeCurrency(currency, ctx, server, option, oldCurrency, user.ID)
	case "events":
		eventDateTimes, err := server.store.ListEventDateTimeNoLimit(ctx, option.ID)
		log.Println("eventDateTR", eventDateTimes)
		if err != nil {
			log.Printf("There an error at UpdateOptionCurrency at ListEventDateTimeNoLimit: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		} else {
			for _, eventDT := range eventDateTimes {
				UpdateEventCurrency(currency, ctx, server, option, oldCurrency, eventDT.ID, user.ID)
			}
		}
	}
	var price = ""
	if option.MainOptionType == "options" {
		optionPrice, err := server.store.GetOptionPrice(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at UpdateOptionCurrency at GetOptionPrice: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		} else {
			price = tools.IntToMoneyString(optionPrice.Price)
		}
	}
	res := UpdateOptionCurrencyRes{
		Currency: currency,
		Price:    price,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "LOTCreateUpdateOptionDiscount", "listing discount", "create-update listing discount")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListOptionRule(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var res ListOptionRuleRes
	var resData []OptionRuleItem
	// lastly we want to setup resData for pets_allowed
	petsAllowed, err := server.store.GetOptionInfoDetailPetsAllow(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at ListOptionRule at GetOptionInfoDetailPetsAllow: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	petsRuleData := OptionRuleItem{Tag: "pets_allowed", Checked: petsAllowed, Type: "house_rule", ID: "none", Des: "none"}
	rules, err := server.store.ListOptionRuleOne(ctx, option.ID)
	if err != nil {
		// This error is here because there is no data
		log.Printf("There an error at ListOptionRule at ListOptionRuleOne: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		resData = append(resData, petsRuleData)
		res = ListOptionRuleRes{
			List: resData,
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	for _, note := range rules {
		data := OptionRuleItem{
			Tag:     note.Tag,
			Checked: note.Checked,
			Type:    note.Type,
			ID:      tools.UuidToString(note.ID),
			Des:     note.Des,
		}
		resData = append(resData, data)
	}

	resData = append(resData, petsRuleData)
	res = ListOptionRuleRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

// CU means CreateUpdate
func (server *Server) CUOptionRule(ctx *gin.Context) {
	var req CUOptionRuleParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CUOptionRuleParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// We want to confirm the tag
	if !val.ValidateHouseRuleTag(req.Tag) {
		err = fmt.Errorf("tag not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res OptionRuleItem
	// We want to check if the tag is for pets
	if req.Tag == "pets_allowed" {
		data, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			PetsAllowed: pgtype.Bool{
				Bool:  req.Checked,
				Valid: true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at CUOptionRule at UpdateOptionRule: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update this house rule")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if !data.PetsAllowed {
			// We want to remove it from additional charge
			_, err = server.store.UpdateOptionAddChargeByType(ctx, db.UpdateOptionAddChargeByTypeParams{
				OptionID:   option.ID,
				Type:       "pet_fee",
				MainFee:    0,
				ExtraFee:   0,
				NumOfGuest: 0,
			})
			if err != nil {
				// Error can occur because we don't know if this pet_fee charge has been created
				log.Printf("There an error at CUOptionRule at UpdateOptionAddChargeByType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			}

		}
		res = OptionRuleItem{
			Tag:     req.Tag,
			Checked: data.PetsAllowed,
			Type:    "house_rule",
			ID:      "",
			Des:     "none",
		}
	} else {
		// This is when the tag is not for pets
		// We would try to create it however if it already exists then we update
		rule, err := server.store.GetOptionRuleByType(ctx, db.GetOptionRuleByTypeParams{
			OptionID: option.ID,
			Type:     req.Type,
			Tag:      req.Tag,
		})
		var exists bool
		if err != nil {
			log.Printf("There an error at CUOptionRule at GetOptionRuleByType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			if err == db.ErrorRecordNotFound {
				exists = false
			} else {
				exists = false
				//err = fmt.Errorf("could not add this house rule to your option")
				//ctx.JSON(http.StatusBadRequest, errorResponse(err))
				//return
			}
		} else {
			exists = true
		}

		if exists {
			// We want to update the house rule
			log.Println("at house rules update")
			update, err := server.store.UpdateOptionRule(ctx, db.UpdateOptionRuleParams{
				Checked: req.Checked,
				ID:      rule.ID,
			})
			if err != nil {
				log.Printf("There an error at CUOptionRule at UpdateOptionRule: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				err = fmt.Errorf("could not update this house rule")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			} else {
				res = OptionRuleItem{
					Tag:     update.Tag,
					Checked: update.Checked,
					Type:    update.Type,
					ID:      tools.UuidToString(update.ID),
					Des:     update.Des,
				}
			}
		} else {
			// We want to create
			log.Println("at house rules create")
			create, err := server.store.CreateOptionRule(ctx, db.CreateOptionRuleParams{
				OptionID: option.ID,
				Tag:      req.Tag,
				Type:     req.Type,
				Checked:  req.Checked,
			})
			if err != nil {
				log.Printf("There an error at CUOptionRule at CreateOptionRule: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				err = fmt.Errorf("could not add this house rule to your option")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			} else {
				res = OptionRuleItem{
					Tag:     create.Tag,
					Checked: create.Checked,
					Type:    create.Type,
					ID:      tools.UuidToString(create.ID),
					Des:     create.Des,
				}
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CUOptionRule", "listing rules", "create-update listing rules")
	}
	fmt.Println("house rules", res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionRuleDetail(ctx *gin.Context) {
	var req GetOptionRuleDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionRuleDetailParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this house rule detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	houseRuleID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, houseRuleID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	ruleDetail, err := server.store.GetOptionRule(ctx, db.GetOptionRuleParams{
		ID:       houseRuleID,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at GetOptionRuleDetail at GetOptionRule: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := fmt.Errorf("none")
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	res := UOptionRuleDetailRes{
		Des:       ruleDetail.Des,
		ID:        tools.UuidToString(ruleDetail.ID),
		StartTime: tools.ConvertTimeOnlyToString(ruleDetail.StartTime),
		EndTime:   tools.ConvertTimeOnlyToString(ruleDetail.EndTime),
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UOptionRuleDetail(ctx *gin.Context) {
	var req UOptionRuleDetailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionRuleDetailRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this house rule detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ruleID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, ruleID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// get rule
	rule, err := server.store.GetOptionRule(ctx, db.GetOptionRuleParams{
		ID:       ruleID,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UOptionRuleDetail start time at tools.GetOptionRule: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to find the resource")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var ruleDetail db.UpdateOptionRuleDetailRow
	if rule.Tag == "quit_hours" {
		startTime, err := tools.ConvertStringToTimeOnly(req.StartTime)
		if err != nil {
			log.Printf("There an error at UOptionRuleDetail start time at tools.ConvertStringToTimeOnly: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("start time is in the wrong format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endTime, err := tools.ConvertStringToTimeOnly(req.EndTime)
		if err != nil {
			log.Printf("There an error at end time UOptionRuleDetail at tools.ConvertStringToTimeOnly: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end time is in the wrong format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ruleDetail, err = server.store.UpdateOptionRuleDetail(ctx, db.UpdateOptionRuleDetailParams{
			StartTime: pgtype.Time{
				Microseconds: tools.TimeToMicroseconds(startTime),
				Valid:        true,
			},
			EndTime: pgtype.Time{
				Microseconds: tools.TimeToMicroseconds(endTime),
				Valid:        true,
			},
			OptionID: option.ID,
			ID:       ruleID,
		})
		if err != nil {
			log.Printf("There an error at time UOptionRuleDetail at UpdateOptionRuleDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end time is in the wrong format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		ruleDetail, err = server.store.UpdateOptionRuleDetail(ctx, db.UpdateOptionRuleDetailParams{
			Des: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
			ID:       ruleID,
		})
		if err != nil {
			log.Printf("There an error at des UOptionRuleDetail at UpdateOptionRuleDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end time is in the wrong format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res := UOptionRuleDetailRes{
		Des:       ruleDetail.Des,
		ID:        tools.UuidToString(ruleDetail.ID),
		StartTime: tools.ConvertTimeOnlyToString(ruleDetail.StartTime),
		EndTime:   tools.ConvertTimeOnlyToString(ruleDetail.EndTime),
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UOptionRuleDetail", "listing rules details", "update listing rules details")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionAvailabilitySetting(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	setting, err := server.store.GetOptionAvailabilitySetting(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionAvailabilitySetting at GetOptionAvailabilitySetting: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your availability settings")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionAvailabilitySettingRes{
		AdvanceNotice:          setting.AdvanceNotice,
		AdvanceNoticeCondition: setting.AdvanceNoticeCondition,
		PreparationTime:        setting.PreparationTime,
		AvailabilityWindow:     setting.AvailabilityWindow,
	}
	fmt.Println(res)
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UOptionAvailabilitySetting(ctx *gin.Context) {
	var req UOptionAvailabilitySettingReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionAvailabilitySettingRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// get rule
	var advanceNotice string
	var preparationTime string
	var availabilityWindow string
	// Advance Notice
	if len(strings.TrimSpace(req.AdvanceNotice)) == 0 || req.AdvanceNotice == "none" {
		// We know advance notice is empty
		advanceNotice = "none"
	} else {
		// We check to make sure it matches
		if !val.ValidateAvailability(req.AdvanceNotice, "advance_notice") {
			err = fmt.Errorf("invalid request data sent")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		advanceNotice = req.AdvanceNotice
	}
	// Preparation Time
	if len(strings.TrimSpace(req.PreparationTime)) == 0 || req.PreparationTime == "none" {
		// We know preparation time is empty
		preparationTime = "none"
	} else {
		// We check to make sure it matches
		if !val.ValidateAvailability(req.PreparationTime, "preparation_time") {
			err = fmt.Errorf("invalid request data sent")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		preparationTime = req.PreparationTime
	}
	// Availability Window
	if len(strings.TrimSpace(req.AvailabilityWindow)) == 0 || req.AvailabilityWindow == "none" {
		// We know availability window is empty
		availabilityWindow = "none"
	} else {
		// We check to make sure it matches
		if !val.ValidateAvailability(req.AvailabilityWindow, "availability_window") {
			err = fmt.Errorf("invalid request data sent")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		availabilityWindow = req.AvailabilityWindow
	}

	setting, err := server.store.UpdateOptionAvailabilitySetting(ctx, db.UpdateOptionAvailabilitySettingParams{
		AdvanceNotice:          advanceNotice,
		AdvanceNoticeCondition: req.AdvanceNoticeCondition,
		PreparationTime:        preparationTime,
		AvailabilityWindow:     availabilityWindow,
		OptionID:               option.ID,
	})
	if err != nil {
		log.Printf("There an error at UOptionAvailabilitySetting at UpdateOptionAvailabilitySetting: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update details for this amenity")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionAvailabilitySettingRes{
		AdvanceNotice:          setting.AdvanceNotice,
		AdvanceNoticeCondition: setting.AdvanceNoticeCondition,
		PreparationTime:        setting.PreparationTime,
		AvailabilityWindow:     setting.AvailabilityWindow,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UOptionAvailabilitySetting", "listing availability settings", "update listing availability settings")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionTripLength(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	trip, err := server.store.GetOptionTripLength(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionTripLength at GetOptionTripLength: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your availability settings")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionTripLengthRes{
		MinStayDay:                  int(trip.MinStayDay),
		MaxStayNight:                int(trip.MaxStayNight),
		ManualApproveRequestPassMax: trip.ManualApproveRequestPassMax,
		AllowReservationRequest:     trip.AllowReservationRequest,
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UOptionTripLength(ctx *gin.Context) {
	var req UOptionTripLengthReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionTripLengthRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if req.MaxStayNight < req.MinStayDay {
		err = fmt.Errorf("your maximum stay period cannot be less than the minimum amount of nights guests can stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trip, err := server.store.UpdateOptionTripLength(ctx, db.UpdateOptionTripLengthParams{
		OptionID:                    option.ID,
		MinStayDay:                  int32(req.MinStayDay),
		MaxStayNight:                int32(req.MaxStayNight),
		ManualApproveRequestPassMax: req.ManualApproveRequestPassMax,
		AllowReservationRequest:     req.AllowReservationRequest,
	})
	if err != nil {
		log.Printf("There an error at UOptionTripLength at UpdateOptionTripLength: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update trip length detail for this stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionTripLengthRes{
		MinStayDay:                  int(trip.MinStayDay),
		MaxStayNight:                int(trip.MaxStayNight),
		ManualApproveRequestPassMax: trip.ManualApproveRequestPassMax,
		AllowReservationRequest:     trip.AllowReservationRequest,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UOptionTripLength", "listing trip length", "update listing trip length")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetCheckInOutDetail(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	checkDetail, err := server.store.GetCheckInOutDetail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetCheckInOutDetail at GetCheckInOutDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your check in and out details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UCheckInOutDetailRes{
		ArriveAfter:            checkDetail.ArriveAfter,
		ArriveBefore:           checkDetail.ArriveBefore,
		LeaveBefore:            checkDetail.LeaveBefore,
		RestrictedCheckInDays:  checkDetail.RestrictedCheckInDays,
		RestrictedCheckOutDays: checkDetail.RestrictedCheckOutDays,
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UCheckInOutDetail(ctx *gin.Context) {
	var req UCheckInOutDetailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UCheckInOutDetailRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var arriveAfter string
	var arriveBefore string
	var leaveBefore string
	// lets handle arrive after
	if tools.ServerStringEmpty(req.ArriveAfter) {
		arriveAfter = "none"
	} else {
		arriveAfter = req.ArriveAfter
	}
	// lets handle arrive before
	if tools.ServerStringEmpty(req.ArriveBefore) {
		arriveBefore = "none"
	} else {
		arriveBefore = req.ArriveBefore
	}
	// lets handle leave before
	if tools.ServerStringEmpty(req.LeaveBefore) {
		leaveBefore = "none"
	} else {
		leaveBefore = req.LeaveBefore
	}
	checkDetail, err := server.store.UpdateCheckInOutDetail(ctx, db.UpdateCheckInOutDetailParams{
		ArriveAfter:            arriveAfter,
		ArriveBefore:           arriveBefore,
		LeaveBefore:            leaveBefore,
		RestrictedCheckInDays:  tools.ServerListToDB(req.RestrictedCheckInDays),
		RestrictedCheckOutDays: tools.ServerListToDB(req.RestrictedCheckOutDays),
		OptionID:               option.ID,
	})
	if err != nil {
		log.Printf("There an error at UCheckInOutDetail at UpdateCheckInOutDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your check in out details for this stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UCheckInOutDetailRes{
		ArriveAfter:            checkDetail.ArriveAfter,
		ArriveBefore:           checkDetail.ArriveBefore,
		LeaveBefore:            checkDetail.LeaveBefore,
		RestrictedCheckInDays:  checkDetail.RestrictedCheckInDays,
		RestrictedCheckOutDays: checkDetail.RestrictedCheckOutDays,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UCheckInOutDetail", "listing check in out details", "update listing check in out details")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionBookMethod(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	bookMethod, err := server.store.GetOptionBookMethod(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionBookMethod at GetOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your availability settings")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetOptionBookMethodRes{
		InstantBook:      bookMethod.InstantBook,
		IdentityVerified: bookMethod.IdentityVerified,
		GoodTrackRecord:  bookMethod.GoodTrackRecord,
		PreBookMsg:       bookMethod.PreBookMsg,
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UOptionBookMethod(ctx *gin.Context) {
	var req UOptionBookMethodReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionBookMethodReq in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var canInstantBook bool
	if option.MainOptionType == "events" {
		canInstantBook = true
	} else {
		canInstantBook = req.InstantBook
	}
	bookMethod, err := server.store.UpdateOptionBookMethod(ctx, db.UpdateOptionBookMethodParams{
		InstantBook:      canInstantBook,
		IdentityVerified: req.IdentityVerified,
		GoodTrackRecord:  req.GoodTrackRecord,
		OptionID:         option.ID,
	})
	if err != nil {
		log.Printf("There an error at UOptionBookMethod at UpdateOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your check in out details for this stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionBookMethodRes{
		InstantBook:      bookMethod.InstantBook,
		IdentityVerified: bookMethod.IdentityVerified,
		GoodTrackRecord:  bookMethod.GoodTrackRecord,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UOptionBookMethod", "listing book method", "update listing book method")
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UOptionBookMethodMsg(ctx *gin.Context) {
	var req UOptionBookMethodMsgReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionBookMethodMsgReq in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var preBookMsg string
	if tools.ServerStringEmpty(req.PreBookMsg) {
		preBookMsg = "none"
	} else {
		preBookMsg = req.PreBookMsg

	}
	msg, err := server.store.UpdateOptionBookMethodMsg(ctx, db.UpdateOptionBookMethodMsgParams{
		PreBookMsg: preBookMsg,
		OptionID:   option.ID,
	})
	if err != nil {
		log.Printf("There an error at UOptionBookMethodMsg at UpdateOptionBookMethodMsg: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your check in out details for this stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionBookMethodMsgRes{
		PreBookMsg: msg,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UOptionBookMethodMsg", "listing book method", "update listing book method message")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetBookRequirement(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	bookRequirement, err := server.store.GetBookRequirement(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetBookRequirement at GetBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your availability settings")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetBookRequirementRes{
		Email:        bookRequirement.Email,
		PhoneNumber:  bookRequirement.PhoneNumber,
		Rules:        bookRequirement.Rules,
		PaymentInfo:  bookRequirement.PaymentInfo,
		ProfilePhoto: bookRequirement.ProfilePhoto,
	}
	log.Println(res)
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UBookRequirement(ctx *gin.Context) {
	var req UBookRequirementReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UBookRequirementRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	bookRequirement, err := server.store.UpdateBookRequirement(ctx, db.UpdateBookRequirementParams{
		ProfilePhoto: pgtype.Bool{
			Bool:  req.ProfilePhoto,
			Valid: true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UBookRequirement at UpdateBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your check in out details for this stay")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UBookRequirementRes{
		ProfilePhoto: bookRequirement.ProfilePhoto,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UBookRequirement", "listing booking requirement", "update listing booking requirement")
	}
	ctx.JSON(http.StatusOK, res)
}
