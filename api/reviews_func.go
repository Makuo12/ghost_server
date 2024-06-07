package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleGetCharge(ctx *gin.Context, server *Server, chargeIDString string, mainOption string, funcName string) (chargeID uuid.UUID, mainOptionType string, chargeType string, user db.User, err error) {
	user, err = HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	chargeID, err = tools.StringToUuid(chargeIDString)
	if err != nil {
		log.Printf("Error FuncName: %v occurred at HandleGetCharge in tools.StringToUuid %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeIDString, user.ID)
		return
	}
	switch mainOption {
	case "options":
		charge, errOption := server.store.GetChargeOptionReference(ctx, db.GetChargeOptionReferenceParams{
			ID:               chargeID,
			UserID:           user.UserID,
			PaymentCompleted: true,
			ChargeCancelled:  false,
			RequestApproved:  true,
		})
		if errOption != nil {
			log.Printf("Error FuncName: %v occurred at HandleGetCharge in GetChargeOptionReference %v, chargeID: %v, userID: %v \n", funcName, errOption.Error(), chargeIDString, user.ID)
			err = fmt.Errorf("your stay was not found")
			return
		}
		// We need to make sure it is a 5 hours after the user checkout date that he writes the review
		if !time.Now().After(charge.EndDate.Add(time.Hour * 5)) {
			// If it not true then we send an error
			err = fmt.Errorf("you can only write a review 5 hours after your checkout date")
			return
		}
		chargeID = charge.ID
		chargeType = constants.CHARGE_OPTION_REFERENCE
		mainOptionType = mainOption
	case "events":
		charge, errEvent := server.store.GetChargeTicketReference(ctx, db.GetChargeTicketReferenceParams{
			ID:         chargeID,
			UserID:     user.UserID,
			Cancelled:  false,
			IsComplete: true,
		})
		if errEvent != nil {
			log.Printf("Error FuncName: %v occurred at HandleGetCharge in GetChargeTicketReference %v, chargeID: %v, userID: %v \n", funcName, errEvent.Error(), chargeIDString, user.ID)
			err = fmt.Errorf("your ticket was not found")
			return
		}
		// We need to make sure it is a 5 hours after the user checkout date that he writes the review
		if !time.Now().After(charge.EndDate.Add(time.Hour * 5)) {
			// If it not true then we send an error
			err = fmt.Errorf("you can only write a review 5 hours after the event ends")
			return
		}
		chargeID = charge.ID
		chargeType = constants.CHARGE_TICKET_REFERENCE
		mainOptionType = mainOption
	default:
		err = fmt.Errorf("type not found")
		return
	}
	return
}

func HandleOptionGetCharge(ctx *gin.Context, server *Server, chargeIDString string, mainOption string, funcName string) (user db.User, optionReview db.GetOptionChargeReviewRow, chargeID uuid.UUID, err error) {
	user, err = HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	chargeID, err = tools.StringToUuid(chargeIDString)
	if err != nil {
		log.Printf("Error FuncName: %v occurred at HandleGetCharge in tools.StringToUuid %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeIDString, user.ID)
		return
	}
	optionReview, err = server.store.GetOptionChargeReview(ctx, db.GetOptionChargeReviewParams{
		ChargeID: chargeID,
		UserID:   user.UserID,
	})
	if err != nil {
		log.Printf("Error FuncName: %v occurred at HandleGetCharge in GetOptionChargeReview %v, chargeID: %v, userID: %v \n", funcName, err.Error(), chargeIDString, user.ID)
		return
	}
	return
}

func HandleAddAmenityReviewItem(amenities []string, req CreateAmenityReviewItem) (result []string) {

	for _, a := range amenities {
		// First we split to get the tag
		if !tools.ServerStringEmpty(a) {
			split := strings.Split(a, "&")
			if len(split) == 2 {
				if split[0] != req.Tag {
					// If the amenity is not equal to the one selected then we add it to the list
					result = append(result, a)
				}
			}
		}
	}
	// Then we add the requested item
	value := fmt.Sprintf("%v&%v", req.Tag, req.Answer)
	result = append(result, value)
	return
}

func HandleRemoveAmenityReviewItem(amenities []string, req RemoveAmenityReviewItem) (result []string) {

	for _, a := range amenities {
		// First we split to get the tag
		if !tools.ServerStringEmpty(a) {
			split := strings.Split(a, "&")
			if len(split) == 2 {
				if split[0] != req.Tag {
					// If the amenity is not equal to the one selected then we add it to the list
					result = append(result, a)
				}
			}
		}
	}
	if len(result) == 0 {
		result = append(result, "none")
	}
	return
}

func HandleAmenityReviewRes(amenities []string, chargeID uuid.UUID) (res ListAmenityReviewRes) {
	chargeIDString := tools.UuidToString(chargeID)
	var resData []CreateAmenityReviewItem
	var isEmpty bool
	result := tools.HandleDBList(amenities)
	for _, a := range result {
		if !tools.ServerStringEmpty(a) {
			split := strings.Split(a, "&")
			if len(split) == 2 {
				data := CreateAmenityReviewItem{
					Tag:      split[0],
					Answer:   split[1],
					ChargeID: chargeIDString,
				}
				resData = append(resData, data)
			}
		}
	}
	if len(resData) == 0 {
		isEmpty = true
		resData = []CreateAmenityReviewItem{{"none", "none", "none"}}
	}
	res = ListAmenityReviewRes{
		IsEmpty: isEmpty,
		List:    resData,
	}
	return
}
