package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) GetChargeCode(ctx *gin.Context) {
	var req GetChargeCodeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetChargeCode in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userDetails, err := HandleCodeExist(user.UserID, req.ID, "GetChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	chargeID, optionUserID, scanned, grade, _, scannedTime, scannedByName, scannedUserImage, ticketType, err := GetChargeForScanned(ctx, server, req, user, "GetChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if scanned {
		// If the code has been scanned we don't want to generate a code for the user
		var msg string
		switch req.MainOption {
		case "options":
			msg = "Your code for this listing has was successfully scanned."
		case "events":
			msg = "Your code for this event has was successfully scanned"
		}
		res := GetChargeCodeScannedRes{
			Message:               msg,
			ID:                    req.ID,
			ScannedByName:         scannedByName,
			ScannedTime:           tools.ConvertTimeToString(scannedTime),
			ScannedUserImage: scannedUserImage,
		}

		ctx.JSON(http.StatusAlreadyReported, res)
		return
	}
	value := fmt.Sprintf("%v&%v&%v&%v", user.UserID, chargeID, optionUserID, user.FirstName)
	key := tools.UuidToString(uuid.New())
	// We store the key with the user details to make sure we generate a new code and delete the previous ones
	err = RedisClient.Set(RedisContext, userDetails, key, time.Hour*1).Err()
	if err != nil {
		log.Printf("Error at  GetChargeCode in .GetChargeOptionReferenceScanned err: %v, user: %v\n", err, user.ID)
		err = errors.New("could not find your booking, try again later. Or try contacting us")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// This id would only be valid for 1hr
	err = RedisClient.Set(RedisContext, key, value, time.Hour*1).Err()
	if err != nil {
		log.Printf("Error at  GetChargeCode in .RedisClient.Set err: %v, user: %v\n", err, user.ID)
		err = errors.New("could not generate a code for this listing")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	code, err := HandleCodeEncrypt(server, key, user, "GetChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetChargeCodeRes{
		Code:       code,
		ID:         req.ID,
		TicketType: ticketType,
		Grade:      grade,
		MainOption: req.MainOption,
	}
	ctx.JSON(http.StatusOK, res)
}

// This deletes the code that was stored in redis
// This maybe used if the guest feels their code was stolen
func (server *Server) DeleteChargeCode(ctx *gin.Context) {
	var req GetChargeCodeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at DeleteChargeCode in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = HandleCodeExist(user.UserID, req.ID, "DeleteChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, _, scanned, _, _, _, _, _, _, err := GetChargeForScanned(ctx, server, req, user, "DeleteChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := DeleteChargeCodeRes{
		WasScanned: scanned,
		Success:    true,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ValidateChargeCode(ctx *gin.Context) {
	var req ValidateChargeCodeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ValidateChargeCode in ShouldBindJSON: %v, OptionID: %v \n", err.Error(), req.OptionID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println("req validate charge code, ", req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at  ValidateChargeCode in .GetChargeOptionReferenceScanned err: %v, optionID: %v\n", err, req.OptionID)
		err = errors.New("your unable to access this listing")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	code, err := HandleCodeDecrypt(server, req.Code, ctx.ClientIP(), "ValidateChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to get the details in the code
	guestID, chargeID, optionUserID, _, err := GetCodeDetails(code, "ValidateChargeCode")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch req.MainOption {
	case "events":
		chargeMain, used, err := CreateAndValidateChargeCodeTicket(ctx, server, req, "ValidateChargeCode", user, guestID, chargeID, optionUserID, optionID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		title := "Checked in successful"
		msg := fmt.Sprintf("Hey %v, you were successfully checked in by %v. Enjoy the event!!!", tools.CapitalizeFirstCharacter(chargeMain.GuestFirstName), tools.CapitalizeFirstCharacter(user.FirstName))
		// Send an apn
		HandleUIDApn(ctx, server, chargeMain.GuestID, title, msg)
		res := ValidateChargeCodeRes{
			FirstName:      chargeMain.GuestFirstName,
			HostOptionName: chargeMain.HostNameOption,
			Grade:          chargeMain.TicketGrade,
			TicketType:     chargeMain.TicketType,
			MainOption:     req.MainOption,
			Success:        true,
			StartDate:      tools.ConvertDateOnlyToString(chargeMain.EventStartDate),
			Used:           used,
			EndDate:        tools.ConvertDateOnlyToString(chargeMain.EventEndDate),
		}
		ctx.JSON(http.StatusOK, res)
	case "options":
		chargeMain, used, err := CreateAndValidateChargeCodeOption(ctx, server, req, "ValidateChargeCode", user, guestID, chargeID, optionUserID, optionID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		title := "Checked in successful"
		msg := fmt.Sprintf("Hey %v, you were successfully checked in by %v. Enjoy your stay!!!", tools.CapitalizeFirstCharacter(chargeMain.GuestFirstName), tools.CapitalizeFirstCharacter(user.FirstName))
		// Send an apn
		HandleUIDApn(ctx, server, chargeMain.GuestID, title, msg)
		res := ValidateChargeCodeRes{
			FirstName:      chargeMain.GuestFirstName,
			HostOptionName: chargeMain.HostNameOption,
			Grade:          "none",
			TicketType:     "none",
			MainOption:     req.MainOption,
			Success:        true,
			Used:           used,
			StartDate:      tools.ConvertDateOnlyToString(chargeMain.StartDate),
			EndDate:        tools.ConvertDateOnlyToString(chargeMain.EndDate),
		}
		ctx.JSON(http.StatusOK, res)
	default:
		err = fmt.Errorf("option not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
