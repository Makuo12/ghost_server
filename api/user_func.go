package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nyaruka/phonenumbers"
)

func (server *Server) ForgotPasswordNotLogged(ctx *gin.Context) {
	var req ForgotPasswordNotLoggedRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at ForgotPasswordNotLogged in ShouldBindJSON: %v \n", err.Error())

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var method string = strings.TrimSpace(ctx.Query("method"))

	switch method {
	case "phone":
		var username string
		var numberExist bool
		if len(req.Code) < 1 || len(req.PhoneNumber) < 1 || len(req.CountryName) < 1 || req.CountryName == "none" || req.PhoneNumber == "none" || req.Code == "none" {
			err := fmt.Errorf("phone number format wrong, please try entering your phone number again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		num, err := phonenumbers.Parse(req.PhoneNumber, req.Code)
		if err != nil {
			log.Printf("Error at ForgotPasswordNotLogged in method phone in phonenumbers: %v \n", err.Error())
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
			err = fmt.Errorf("this phone number is not registered with Flizzup, try using the number you used to create your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		if !user.IsActive {
			err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		dial_number_code := num.GetCountryCode() // +234
		dial_country := req.CountryName          // Nigeria
		username = utils.RandomName()
		err = SendSmsOtp(server, number, username, dial_number_code, dial_country, "ForgotPasswordNotLogged")
		if err != nil {
			log.Printf("Error at VerifyPhoneNumber in SendSmsOtp(server: %v \n", err.Error())
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		res := VerifyPhoneNumberResponse{
			CodeSent:    true,
			Username:    username,
			PhoneNumber: numberToSend,
		}
		ctx.JSON(http.StatusOK, res)

	case "email":
		//EMAIL
		var username string
		var name string
		var emailExist bool
		if len(req.Email) < 1 || req.Email == "none" {
			err := fmt.Errorf("email format wrong, please try entering your email again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		user, err := server.store.GetUserWithEmail(ctx, req.Email)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				emailExist = false
			} else {
				emailExist = true
			}

		} else {
			emailExist = true
		}
		if !emailExist {
			err = fmt.Errorf("this email is not registered with Flizzup, try using the email you used to create your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		if !user.IsActive {
			err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		username = utils.RandomName()
		name = tools.CapitalizeFirstCharacter(user.FirstName) + " " + tools.CapitalizeFirstCharacter(user.LastName)
		err = BrevoEmailCode(ctx, server, req.Email, name, username, "ForgotPasswordNotLogged")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := VerifyEmailResponse{
			CodeSent: true,
			Username: username,
			Email:    user.Email,
		}
		log.Println(res)
		ctx.JSON(http.StatusOK, res)

	default:
		err := fmt.Errorf("method format wrong, please try selecting either email or phone as the method you want to use to reset your password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}
func (server *Server) ConfirmCodeLogin(ctx *gin.Context) {
	var req ConfirmCodeRequest
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
	_, code, toPhone := split[0], split[1], split[2]
	if code != req.Code {
		err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserWithPhoneNum(ctx, toPhone)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("error at ConfirmCodeLogin in GetUser: %v \n", err.Error())
			err = fmt.Errorf("this account isn't registered with Flizzup")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("error at ConfirmCodeLogin in GetUser: %v \n", err.Error())
		err = fmt.Errorf("there was an error while getting your information, make sure you are sign up on Flizzup")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !user.IsActive {
		err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "ConfirmCodeLogin")
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
func (server *Server) ConfirmCode(ctx *gin.Context) {
	var req ConfirmCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at ConfirmCode in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var method string = strings.TrimSpace(ctx.Query("method"))

	switch method {
	case "phone":
		exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			if err != nil {
				log.Printf("Error at ConfirmCode in RedisClient Exists: %v \n", err.Error())
			}
			err = fmt.Errorf("details of your code was not found, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		result, err := RedisClient.Get(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			log.Printf("Error at ConfirmCode in RedisClient  RedisClient.Get: %v \n", err.Error())
			err = fmt.Errorf("details of your code was not found, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		split := strings.Split(result, "&")
		// There are six in total
		if len(split) != 7 {
			err = fmt.Errorf("details of your code was not complete, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, code, toPhone, username, dialCodeNumber, dialCountry, confirm := split[0], split[1], split[2], split[3], split[4], split[5], split[6]
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if confirm == "true" {
			err = fmt.Errorf("this code has already been used successfully")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		data := fmt.Sprintf("%v&%v&%v&%v&%v&%v&%v", "phone", code, toPhone, username, dialCodeNumber, dialCountry, "true")
		err = RedisClient.Set(RedisContext, username, data, time.Hour*4).Err()
		if err != nil {
			log.Printf("HandlePhoneVerifyWithoutTwilio SAdd(RedisContext, RedisClient.Set %v err:%v\n", username, err.Error())
			err = fmt.Errorf("something we wrong while processing your code")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		res := ConfirmCodeResponse{
			Confirmed: true,
			Username:  username,
		}
		ctx.JSON(http.StatusOK, res)
	case "email":
		exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			if err != nil {
				log.Printf("Error at ConfirmCode in RedisClient Exists: %v \n", err.Error())
			}
			err = fmt.Errorf("details of your code was not found, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		result, err := RedisClient.Get(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			log.Printf("Error at ConfirmCode in RedisClient  RedisClient.Get: %v \n", err.Error())
			err = fmt.Errorf("details of your code was not found, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		split := strings.Split(result, "&")
		// There are 4 in total
		if len(split) != 5 {
			err = fmt.Errorf("details of your code was not complete, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, code, username, email, confirm := split[0], split[1], split[2], split[3], split[4]
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if confirm == "true" {
			err = fmt.Errorf("this code has already been used successfully")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		userDetails := fmt.Sprintf("%v&%v&%v&%v&%v", "email", code, username, email, "true")
		err = RedisClient.Set(RedisContext, username, userDetails, time.Hour*4).Err()
		if err != nil {
			log.Printf("FuncName: %v, SendEmailVerifyCode RedisClient.HSe %v err:%v\n", "ConfirmCode", username, err.Error())
			err = fmt.Errorf("code not store code, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := ConfirmCodeResponse{
			Confirmed: true,
			Username:  username,
		}
		ctx.JSON(http.StatusOK, res)
	default:
		err := fmt.Errorf("method format wrong, please try selecting either email or phone as the method you want to use to reset your password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) NewPassword(ctx *gin.Context) {
	var req NewPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at NewPassword in ShouldBindJSON: %v \n", err.Error())
		err = fmt.Errorf("your password should be more than 8 characters and contain a special character($,@,&,*)")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.PasswordOne != req.PasswordTwo {
		err := fmt.Errorf("the passwords don't match. Please let the password match each other")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		if err != nil {
			log.Printf("Error at NewPassword in RedisClient Exists: %v \n", err.Error())
		}
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := RedisClient.Get(RedisContext, req.Username).Result()
	if err != nil || exist == 0 {
		log.Printf("Error at NewPassword in RedisClient  RedisClient.Get: %v \n", err.Error())
		err = fmt.Errorf("details of your code was not found, try again as it might have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var user db.User
	split := strings.Split(result, "&")
	switch split[0] {
	case "phone":
		_, _, toPhone, _, _, _, confirm := split[0], split[1], split[2], split[3], split[4], split[5], split[6]
		if confirm != "true" {
			err = fmt.Errorf("you are yet to be verified")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		user, err = server.store.GetUserWithPhoneNum(ctx, toPhone)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				log.Printf("error at NewPassword in GetUser: %v \n", err.Error())
				err = fmt.Errorf("this account isn't registered with Flizzup")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			log.Printf("error at NewPassword in GetUser: %v \n", err.Error())
			err = fmt.Errorf("there was an error while getting your information, make sure you are sign up on Flizzup")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "email":
		_, _, _, email, confirm := split[0], split[1], split[2], split[3], split[4]
		if confirm != "true" {
			err = fmt.Errorf("you are yet to be verified")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		user, err = server.store.GetUserWithEmail(ctx, email)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				log.Printf("error at NewPassword in GetUser: %v \n", err.Error())
				err = fmt.Errorf("this account isn't registered with Flizzup")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			log.Printf("error at NewPassword in GetUser: %v \n", err.Error())
			err = fmt.Errorf("there was an error while getting your information, make sure you are sign up on Flizzup")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	hashedPassword, err := utils.HashedPassword(req.PasswordOne)
	if err != nil {
		log.Printf("Error at NewPassword in HashedPassword: %v \n", err.Error())
		err = fmt.Errorf("there was an error while processing your password, make sure strong with at least 10 characters")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	updateUser, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		FirebasePassword: pgtype.Text{
			String: tools.UuidToString(uuid.New()),
			Valid:  true,
		},
		HashedPassword: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("Error at NewPassword in UpdateUser: %v \n", err.Error())
		err = fmt.Errorf("there was an error while updating your password, please check the email or phone you used to update your password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to update the user password in firebase
	params := (&auth.UserToUpdate{}).
		Password(updateUser.FirebasePassword)
	u, err := server.ClientFire.UpdateUser(ctx, tools.UuidToString(updateUser.FirebaseID), params)
	if err != nil {
		log.Printf("error updating user password in firebase: %v, userID: %v\n", err, user.ID)
		err = nil
	} else {
		log.Printf("Successfully updated user password in firebase: %v\n", u)
	}

	res := NewPasswordResponse{
		Updated: true,
	}
	log.Printf("user password updated: %v", user.ID)
	ctx.JSON(http.StatusOK, res)
}

// Not in use yet
func (server *Server) UpdatedConfirmCodeLogin(ctx *gin.Context) {
	var req ConfirmCodeRequest
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
	_, code, toPhone := split[0], split[1], split[2]
	if code != req.Code {
		err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var phoneExist bool
	var accessToken string
	var refreshToken string
	var accessPayloadStringTime string
	user, err := server.store.GetUserWithPhoneNum(ctx, toPhone)
	if err != nil {
		log.Printf("error at ConfirmCodeLogin in GetUser: %v \n", err.Error())
		phoneExist = false
	} else {
		phoneExist = true
	}
	if phoneExist {
		if !user.IsActive {
			err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}

		accessToken, refreshToken, accessPayloadStringTime, err = HandleUserSession(ctx, server, user, "ConfirmCodeLogin")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {
		// We want to create the user account
		// Because this is sign up with phone we want to create a fake email address
		email := fmt.Sprintf("%v@flizzup-email.com", uuid.New())
		password := fmt.Sprintf("%v", uuid.New())
		hashedPassword, err := utils.HashedPassword(password)
		if err != nil {
			log.Printf("error at CreateUser in HashedPassword: %v \n", err)
			err = fmt.Errorf("there was an error while processing your password, make sure strong with at least 10 characters")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		var username string = uuid.New().String()
		var firebasePassword string = uuid.New().String()
		arg := db.CreateUserParams{
			HashedPassword:   hashedPassword,
			Email:            email,
			Username:         username,
			FirebasePassword: firebasePassword,
			DateOfBirth:      time.Now(),
			FirstName:        "none",
			LastName:         "none",
			Currency:         "NGN",
		}
		user, err = server.store.CreateUser(ctx, arg)
		if err != nil {
			// when we get the error we try to convert it to pq.Err tag
			if db.ErrorCode(err) == db.UniqueViolation {
				err = fmt.Errorf("this email address already exists, try logging in or use forgot password")
				ctx.JSON(http.StatusForbidden, errorResponse(err))
			}
			log.Printf("error at CreateUser in server.store.CreateUser: %v \n", err)
			err = fmt.Errorf("there was an error in while signing you up, please try again")
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
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
