package api

//import (
//	"database/sql"
//	"flex_server/constants"
//	db "flex_server/db/sqlc"
//	"flex_server/token"
//	"flex_server/utils"
//	"fmt"
//	"log"
//	"net/http"
//	"strings"
//	"time"

//	"github.com/gin-gonic/gin"
//	"github.com/nyaruka/phonenumbers"
//)

//func (server *Server) VerifyPhoneNumber(ctx *gin.Context) {
//	var req VerifyPhoneNumberRequest
//	var username string
//	var numberExist bool
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		//if there is an error we want to send it back to the user
//		log.Printf("Error at VerifyPhoneNumber in ShouldBindJSON: %v \n", err.Error())
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	if len(req.Code) == 0 || len(req.PhoneNumber) == 0 || len(req.CountryName) == 0 || req.CountryName == "none" || req.PhoneNumber == "none" || req.Code == "none" {
//		err := fmt.Errorf("phone number format wrong, please try entering your phone number again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	num, err := phonenumbers.Parse(req.PhoneNumber, req.Code)
//	if err != nil {
//		log.Printf("Error at VerifyPhoneNumber in phonenumbers: %v \n", err.Error())
//		err = fmt.Errorf("there was an error while checking your number, try again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(num.GetCountryCode(), num.GetNationalNumber())

//	number := fmt.Sprintf("+%d_%d", num.GetCountryCode(), num.GetNationalNumber())
//	numberToSend := fmt.Sprintf("0%d", num.GetNationalNumber())
//	if len(number) < 6 {
//		err = fmt.Errorf("phone number is too short")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	_, err = server.store.GetUserWithPhoneNum(ctx, number)
//	if err != nil {
//		if err ==  {
//			numberExist = false
//		} else {
//			numberExist = true
//		}

//	} else {
//		numberExist = true
//	}
//	if numberExist {
//		err = fmt.Errorf("this phone number already exists, try logging in or use a different phone number")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	dial_number_code := num.GetCountryCode() // +234
//	dial_country := req.CountryName          // Nigeria
//	username = utils.RandomName()
//	err = HandlePhoneVerifyWithoutTwilio(server, number, username, dial_number_code, dial_country)
//	if err != nil {
//		log.Printf("Error at VerifyPhoneNumber in HandlePhoneVerifyWithoutTwilio: %v \n", err.Error())
//		err = fmt.Errorf("could not generate code for this phone number")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}
//	res := VerifyPhoneNumberResponse{
//		CodeSent:    true,
//		Username:    username,
//		PhoneNumber: numberToSend,
//	}
//	log.Println(res)
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) ConfirmNumber(ctx *gin.Context) {
//	var req ConfirmPhoneNumberRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		//if there is an error we want to send it back to the user
//		log.Printf("Error at ConfirmNumber in ShouldBindJSON: %v \n", err.Error())
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
//	if err != nil || exist == 0 {
//		if err != nil {
//			log.Printf("Error at ConfirmNumber in RedisClient Exists: %v \n", err.Error())
//		}
//		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	result, err := RedisClient.Get(RedisContext, req.Username).Result()
//	if err != nil || exist == 0 {
//		log.Printf("Error at ConfirmNumber in RedisClient  RedisClient.Get: %v \n", err.Error())
//		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	split := strings.Split(result, "&")
//	if len(split) != 6 {
//		err = fmt.Errorf("details of your code was not complete, try again as it might have expired")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	code, toPhone, username, dialCodeNumber, dialCountry, confirm := split[0], split[1], split[2], split[3], split[4], split[5]
//	if code != req.Code {
//		err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	if confirm == "true" {
//		err = fmt.Errorf("this code has already been used successfully")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	data := fmt.Sprintf("%v&%v&%v&%v&%v&%v", code, toPhone, username, dialCodeNumber, dialCountry, "true")
//	err = RedisClient.Set(RedisContext, username, data, time.Hour*4).Err()
//	if err != nil {
//		log.Printf("HandlePhoneVerifyWithoutTwilio SAdd(RedisContext, RedisClient.Set %v err:%v\n", username, err.Error())
//		err = fmt.Errorf("something we wrong while processing your code")
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	res := ConfirmPhoneNumberResponse{
//		Confirmed: true,
//		Username:  username,
//	}
//	ctx.JSON(http.StatusOK, res)
//}

//func (server *Server) UpdateUserPhoneNumber(ctx *gin.Context) {
//	var req ChangePhoneRequest

//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		//if there is an error we want to send it back to the user
//		log.Printf("Error at UpdateUserPhoneNumber in ShouldBindJSON: %v \n", err.Error())
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	payload, exists := ctx.Get(authorizationPayloadKey)
//	if !exists {
//		log.Printf("Error at ctx.GET does not exist")
//		err := fmt.Errorf("there was an error while getting your information, make sure you are logged in")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	data := payload.(*token.Payload)

//	var username string = *&data.Username
//	user, err := server.store.GetUserWithUsername(ctx, username)

//	if err != nil {
//		if err ==  {
//			log.Printf("Error at UpdateUserPhoneNumber in GetUserIDWithUsername: %v \n", err.Error())
//			err = fmt.Errorf("this account isn't registered with Flexr")
//			ctx.JSON(http.StatusNotFound, errorResponse(err))
//			return
//		}
//		log.Printf("Error at UpdateUserPhoneNumber in GetUserIDWithUsername: %v \n", err.Error())
//		err = fmt.Errorf("there was an error while getting your information, make sure you are sign up on Flexr")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	if !user.IsActive {
//		err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
//		ctx.JSON(http.StatusForbidden, errorResponse(err))
//		return
//	}

//	userData, err := RedisClient.HGetAll(RedisContext, req.Username).Result()

//	if err != nil {
//		log.Printf("Error at  UpdateUserPhoneNumber in RedisClient HGet: %v \n", err.Error())
//		err = fmt.Errorf("there was an error while confirming your code, try again")
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}

//	// We want to update the dialCode, dialCountry, and phone number
//	userUpdate, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
//		DialCode: sql.NullString{
//			String: userData[constants.DIAL_CODE_NUMBER],
//			Valid:  true,
//		},
//		DialCountry: sql.NullString{
//			String: userData[constants.DIAL_COUNTRY],
//			Valid:  true,
//		},
//		PhoneNumber: sql.NullString{
//			String: userData[constants.PHONE_NUMBER],
//			Valid:  true,
//		},
//		ID: user.ID,
//	})
//	if err != nil {
//		log.Printf("Error occurred in UpdateUser %v\n", err.Error())
//		err = fmt.Errorf("an error occurred while updating your phone number, try again")
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//	res := ChangePhoneResponse{
//		Updated:     true,
//		PhoneNumber: userUpdate.PhoneNumber,
//		Code:        userUpdate.DialCode,
//		CountryName: userUpdate.DialCountry,
//	}
//	ctx.JSON(http.StatusOK, res)
//}
