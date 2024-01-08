package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

func (server *Server) JoinWithPhone(ctx *gin.Context) {
	var req JoinWithPhoneParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at JoinWithPhone in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var username string
	var numberExist bool
	if len(req.Code) < 1 || len(req.PhoneNumber) < 1 || len(req.CountryName) < 1 || req.CountryName == "none" || req.PhoneNumber == "none" || req.Code == "none" {
		err := fmt.Errorf("phone number format wrong, please try entering your phone number again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	num, err := phonenumbers.Parse(req.PhoneNumber, req.Code)
	if err != nil {
		log.Printf("Error at JoinWithPhone in method phone in phonenumbers: %v \n", err.Error())
		err = fmt.Errorf("there was an error while checking your number, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(num.GetCountryCode(), num.GetNationalNumber())

	number := fmt.Sprintf("+%d_%d", num.GetCountryCode(), num.GetNationalNumber())
	numberToSend := fmt.Sprintf("0%d", num.GetNationalNumber())
	if len(number) < 6 {
		err = fmt.Errorf("phone number is too short")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserWithPhoneNum(ctx, number)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			numberExist = false
		} else {
			numberExist = true
		}
	} else {
		numberExist = true
	}
	if !numberExist {
		// If number does not exist we want to actually start by creating a new account
		dial_number_code := num.GetCountryCode() // +234
		dial_country := req.CountryName          // Nigeria
		username = utils.RandomName()
		err = StoreNumber(server, number, username, dial_number_code, dial_country, "ForgotPasswordNotLogged")
		if err != nil {
			log.Printf("Error at JoinWithPhone in StoreNumber(server: %v \n", err.Error())
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := JoinWithPhoneRes{
			Type:        "sign_up",
			CodeSent:    true,
			Username:    username,
			PhoneNumber: numberToSend,
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	if !user.IsActive {
		err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if user.IsDeleted {
		err = fmt.Errorf("this account does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	dial_number_code := num.GetCountryCode() // +234
	dial_country := req.CountryName          // Nigeria
	username = utils.RandomName()
	err = SendSmsOtp(server, number, username, dial_number_code, dial_country, "ForgotPasswordNotLogged")
	if err != nil {
		log.Printf("Error at JoinWithPhone in SendSmsOtp(server: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := JoinWithPhoneRes{
		Type:        "login",
		CodeSent:    true,
		Username:    username,
		PhoneNumber: numberToSend,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ConfirmJoinSignUp(ctx *gin.Context) {
	var req ConfirmJoinSignUpParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at ConfirmCode in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		if err != nil {
			log.Printf("Error at ConfirmCodeLogin in RedisClient Exists: %v \n", err.Error())
		}
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := RedisClient.Get(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		log.Printf("Error at ConfirmCodeLogin in RedisClient  RedisClient.Get: %v \n", err.Error())
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	split := strings.Split(result, "&")
	if len(split) != 7 {
		err = fmt.Errorf("details of your code was not complete, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, _, toPhone, usernameString, dialCodeNumber, dialCountry, _ := split[0], split[1], split[2], split[3], split[4], split[5], split[6]
	err = SendSmsOtpWithName(server, toPhone, usernameString, dialCodeNumber, dialCountry, req.FirstName, req.LastName, "ConfirmJoinSignUp")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ConfirmJoinSignUpRes{
		Type:        req.Type,
		CodeSent:    true,
		Username:    usernameString,
		PhoneNumber: toPhone,
	}
	log.Printf("user logged in successfully (%v) \n", toPhone)
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) ConfirmCodeJoin(ctx *gin.Context) {
	var req ConfirmCodeJoinParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at ConfirmCode in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		if err != nil {
			log.Printf("Error at ConfirmCodeJoin in RedisClient Exists: %v \n", err.Error())
		}
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := RedisClient.Get(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		log.Printf("Error at ConfirmCodeJoin in RedisClient  RedisClient.Get: %v \n", err.Error())
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	split := strings.Split(result, "&")
	if len(split) > 6 {
		err = fmt.Errorf("details of your code was not complete, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, code, toPhone := split[0], split[1], split[2]
	if code != req.Code {
		err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var existPhone bool
	var accessToken string
	user, err := server.store.GetUserWithPhoneNum(ctx, toPhone)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			existPhone = false
		} else {
			log.Printf("error at ConfirmCodeJoin in GetUser: %v \n", err.Error())
			err = fmt.Errorf("there was an error while getting your information, make sure you are sign up on Flizzup")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		existPhone = true
	}
	if !existPhone && req.Type == "sign_up" && len(split) == 9 {
		// Here we want to create a new user
		_, _, _, _, _, _, _, firstName, lastName := split[0], split[1], split[2], split[3], split[4], split[5], split[6], split[7], split[8]
		var password string = uuid.New().String()
		hashedPassword, err := utils.HashedPassword(password)
		if err != nil {
			log.Printf("error at ConfirmCodeJoin in HashedPassword: %v \n", err)
			err = fmt.Errorf("there was an error while processing your password, make sure strong with at least 10 characters")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		email := fmt.Sprintf("%v@flizzup.com", uuid.New().String())
		var username string = uuid.New().String()
		var firebasePassword string = uuid.New().String()
		date, err := tools.ConvertDateOnlyStringToDate("1777-12-07")
		if err != nil {
			log.Printf("error at ConfirmCodeJoin in HashedPassword: %v \n", err)
			err = fmt.Errorf("Try again something went wrong with the server")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		arg := db.CreateUserParams{
			HashedPassword:   hashedPassword,
			Email:            strings.ToLower(email),
			Username:         username,
			FirebasePassword: firebasePassword,
			DateOfBirth:      date,
			FirstName:        strings.ToLower(firstName),
			LastName:         strings.ToLower(lastName),
			Currency:         req.Currency,
		}
		user, err = server.store.CreateUser(ctx, arg)
		if err != nil {
			// when we get the error we try to convert it to pq.Err tag
			if db.ErrorCode(err) == db.UniqueViolation {
				err = fmt.Errorf("this email address already exists, try logging in or use forgot password")
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
			log.Printf("error at ConfirmCodeJoin in server.store.CreateUser: %v \n", err)
			err = fmt.Errorf("there was an error in while signing you up, please try again")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		params := (&auth.UserToCreate{}).
			Email(user.Email).
			UID(tools.UuidToString(user.FirebaseID)).
			EmailVerified(true).
			Password(user.FirebasePassword).
			DisplayName(user.FirstName).
			Disabled(false)
		u, err := server.ClientFire.CreateUser(ctx, params)
		if err != nil {
			log.Printf("error creating user in firebase: %v, userID: %v\n", err, user.ID)
			err = nil
		} else {
			log.Printf("Successfully created user in firebase: %v\n", u)
		}
	} else if existPhone && req.Type == "login" {
		if !user.IsActive {
			err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		if user.IsDeleted {
			err = fmt.Errorf("this account does not exist")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

	} else {
		err = fmt.Errorf("request protocol not correct")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "ConfirmCodeJoin")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := LoginUserResponse{
		Email:                user.Email,
		FireFight:            user.FirebasePassword,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		AccessToken:          accessToken,
		ProfilePhoto:         user.Photo,
		Currency:             user.Currency,
		AccessTokenExpiresAt: accessPayloadStringTime,
		RefreshToken:         refreshToken,
		PublicID:             tools.UuidToString(user.PublicID),
	}
	log.Printf("user logged in successfully (%v) \n", user.Email)
	ctx.JSON(http.StatusOK, res)

}
