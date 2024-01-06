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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nyaruka/phonenumbers"
)

func (server *Server) ForgotPasswordLogged(ctx *gin.Context) {
	var method string = strings.TrimSpace(ctx.Query("method"))
	var username string
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch method {
	case "phone":
		num, err := phonenumbers.Parse(user.PhoneNumber, user.DialCode)
		if err != nil {
			log.Printf("Error at ForgotPasswordNotLogged in method phone in phonenumbers: %v \n", err.Error())
			err = fmt.Errorf("there was an error while checking your number, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		number := fmt.Sprintf("+%d_%d", num.GetCountryCode(), num.GetNationalNumber())
		numberToSend := fmt.Sprintf("0%d", num.GetNationalNumber())
		dial_number_code := num.GetCountryCode() // +234
		dial_country := user.DialCountry         // Nigeria
		username = utils.RandomName()
		err = SendSmsOtp(server, number, username, dial_number_code, dial_country, "ForgotPasswordLogged")
		if err != nil {
			log.Printf("Error at VerifyPhoneNumber in ForgotPasswordLogged: %v \n", err.Error())
			err = fmt.Errorf("could not generate code for this phone number")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		res := VerifyPhoneNumberResponse{
			CodeSent:    true,
			Username:    username,
			PhoneNumber: numberToSend,
		}
		log.Println(res)
		ctx.JSON(http.StatusOK, res)

	case "email":
		//EMAIL
		username = utils.RandomName()
		name := tools.CapitalizeFirstCharacter(user.FirstName) + " " + tools.CapitalizeFirstCharacter(user.LastName)
		err = BrevoEmailCode(ctx, server, user.Email, name, username, "ForgotPasswordNotLogged")
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

func (server *Server) JoinVerifyEmail(ctx *gin.Context) {
	var username string
	var name string
	var emailExist bool
	var req JoinUserVerifyEmail
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("error at JoinUserVerifyEmail in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err := server.store.GetUserWithEmail(ctx, req.Email)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			emailExist = false
		} else {
			emailExist = true
		}

	} else {
		emailExist = true
	}
	if emailExist {
		err = fmt.Errorf("this email is already registered with Flizzup try logging in")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	username = utils.RandomName()
	name = ""
	err = BrevoEmailCode(ctx, server, req.Email, name, username, "ForgotPasswordNotLogged")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := VerifyEmailResponse{
		CodeSent: true,
		Username: username,
		Email:    req.Email,
	}
	log.Println(res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateVerifyEmailPhone(ctx *gin.Context) {
	var req UpdateVerifyEmailPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("error at UpdateVerifyEmailPhone in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var method string = strings.TrimSpace(ctx.Query("method"))

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch method {
	case "email":
		// We check if the email already exists
		var emailExist bool
		if len(req.Email) < 1 || req.Email == "none" {
			err := fmt.Errorf("email format wrong, please try entering your email again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, err := server.store.GetUserWithEmail(ctx, req.Email)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				emailExist = false
			} else {
				emailExist = true
			}

		} else {
			emailExist = true
		}
		if emailExist {
			log.Println("Email exist")
			err = fmt.Errorf("this email is already registered with Flizzup try logging in")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		name := tools.CapitalizeFirstCharacter(user.FirstName) + " " + tools.CapitalizeFirstCharacter(user.LastName)
		username := utils.RandomName()
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
	case "phone":
		var username string
		var numberExist bool
		if len(req.Code) == 0 || len(req.PhoneNumber) == 0 || len(req.CountryName) == 0 || req.CountryName == "none" || req.PhoneNumber == "none" || req.Code == "none" {
			err := fmt.Errorf("phone number format wrong, please try entering your phone number again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		log.Println("code ", req.Code, req.PhoneNumber)
		num, err := phonenumbers.Parse(req.PhoneNumber, req.Code)
		if err != nil {
			log.Printf("Error at VerifyPhoneNumber in phonenumbers: %v \n", err.Error())
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
		_, err = server.store.GetUserWithPhoneNum(ctx, number)
		if err != nil {
			if err == db.ErrorRecordNotFound {
				numberExist = false
			} else {
				numberExist = true
			}

		} else {
			numberExist = true
		}
		if numberExist {
			err = fmt.Errorf("this phone number already exists, try logging in or use a different phone number")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		dial_number_code := num.GetCountryCode() // +234
		dial_country := req.CountryName          // Nigeria
		username = utils.RandomName()
		err = SendSmsOtp(server, number, username, dial_number_code, dial_country, "UpdateVerifyEmailPhone")
		if err != nil {
			log.Printf("Error at VerifyPhoneNumber in SendSmsOtp(server: %v \n", err.Error())
			err = fmt.Errorf("could not generate code for this phone number")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		res := VerifyPhoneNumberResponse{
			CodeSent:    true,
			Username:    username,
			PhoneNumber: numberToSend,
		}
		log.Println(res)
		ctx.JSON(http.StatusOK, res)
	default:
		err := fmt.Errorf("method format wrong, please try selecting either email or phone as the method you want to use to reset your password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) UpdateCodeEmailPhone(ctx *gin.Context) {
	var req ConfirmCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("Error at ConfirmCode in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
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
		_, code, toPhone, _, dialCodeNumber, dialCountry, _ := split[0], split[1], split[2], split[3], split[4], split[5], split[6]
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// We want to update the dialCode, dialCountry, phoneNumber
		userUpdate, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
			DialCode: pgtype.Text{
				String: dialCodeNumber,
				Valid:  true,
			},
			DialCountry: pgtype.Text{
				String: dialCountry,
				Valid:  true,
			},
			PhoneNumber: pgtype.Text{
				String: toPhone,
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at  UpdateCodeEmailPhone in server.store.UpdateUser: %v \n", err.Error())
			err = fmt.Errorf("there was an error while updating your phone number, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "UpdateCodeEmailPhone phone")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := UpdateCodePhoneResponse{
			AccessToken:          accessToken,
			RefreshToken:         refreshToken,
			AccessTokenExpiresAt: accessPayloadStringTime,
			PhoneNumber:          tools.HandlePhoneNumber(userUpdate.PhoneNumber),
		}
		log.Printf("user phone number updated successfully (%v) \n", userUpdate.ID)
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
		_, code, _, email, _ := split[0], split[1], split[2], split[3], split[4]
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via SMS or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// We just want to update the email
		userUpdate, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
			Email: pgtype.Text{
				String: email,
				Valid:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("Error at  UpdateCodeEmailPhone in server.store.UpdateEmail: %v \n", err.Error())
			err = fmt.Errorf("there was an error while updating your email, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// We want to update the user email firebase
		params := (&auth.UserToUpdate{}).
			Email(userUpdate.Email).
			EmailVerified(true)
		u, err := server.ClientFire.UpdateUser(ctx, tools.UuidToString(userUpdate.FirebaseID), params)
		if err != nil {
			log.Printf("error updating user in firebase: %v, userID: %v\n", err, user.ID)
			err = nil
		} else {
			log.Printf("Successfully updated user in firebase: %v\n", u)
		}

		accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "UpdateCodeEmailPhone email")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := UpdateCodeEmailResponse{
			Email:                tools.HandleEmail(userUpdate.Email),
			AccessToken:          accessToken,
			RefreshToken:         refreshToken,
			AccessTokenExpiresAt: accessPayloadStringTime,
		}
		log.Printf("user email updated successfully (%v) \n", userUpdate.ID)
		ctx.JSON(http.StatusOK, res)
	default:
		err := fmt.Errorf("method format wrong, please try selecting either email or phone as the method you want to use to reset your password")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
}

func (server *Server) CreateUserAPNDetail(ctx *gin.Context) {
	var req CreateUserAPNDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at CreateUserAPNDetail in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while processing your inputs please try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	publicID, err := tools.StringToUuid(req.PublicID)
	if err != nil {
		log.Printf("Error at  CreateUserAPNDetail in tools.StringToUuid err: %v, user: %v\n", err, ctx.ClientIP)
		err = fmt.Errorf("public id not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByPD(ctx, publicID)
	if err != nil {
		log.Printf("Error at  CreateUserAPNDetail in GetUserByPD err: %v, user: %v\n", err, ctx.ClientIP)
		err = fmt.Errorf("public id not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	id, err := server.store.CreateUserAPNDetail(ctx, db.CreateUserAPNDetailParams{
		UserID:              user.ID,
		DeviceName:          req.Name,
		Model:               req.Model,
		IdentifierForVendor: req.IdentifierForVendor,
		Token:               req.Token,
	})
	if err != nil {
		log.Printf("Error at CreateUserAPNDetail in CreateUserAPNDetail err: %v, user: %v\n", err, ctx.ClientIP)
		err = fmt.Errorf("public id not in the right format")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Once we create the delete the remaining
	err = server.store.RemoveAllUserAPNDetailButOne(ctx, db.RemoveAllUserAPNDetailButOneParams{
		UserID: user.ID,
		ID:     id,
	})
	if err != nil {
		log.Printf("Error at CreateUserAPNDetail in RemoveAllUserAPNDetailButOne err: %v, user: %v\n", err, ctx.ClientIP)
	}
	res := UserResponseMsg{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

//func (server *Server) CreateUserAPNDetail(ctx *gin.Context) {
//	var req CreateUserAPNDetailParams
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("error at CreateUserAPNDetail in ShouldBindJSON: %v \n", err)
//		err = fmt.Errorf("there was an error while processing your inputs please try again later")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	publicID, err := tools.StringToUuid(req.PublicID)
//	if err != nil {
//		log.Printf("Error at  CreateUserAPNDetail in tools.StringToUuid err: %v, user: %v\n", err, ctx.ClientIP)
//		err = fmt.Errorf("public id not in the right format")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, err := server.store.GetUserByPD(ctx, publicID)
//	if err != nil {
//		log.Printf("Error at  CreateUserAPNDetail in GetUserByPD err: %v, user: %v\n", err, ctx.ClientIP)
//		err = fmt.Errorf("public id not in the right format")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	var exist bool
//	details, err := server.store.ListUidAPNDetail(ctx, user.ID)
//	if err != nil || len(details) == 0 {
//		if err != nil {
//			log.Printf("Error at  CreateUserAPNDetail in .ListUidAPNDetail err: %v, user: %v\n", err, ctx.ClientIP)
//		}
//		exist = false
//	} else {
//		exist = true
//	}

//	if exist {
//		var found bool
//		for _, d := range details {
//			// We want to check if any details match
//			if req.Name == d.DeviceName && req.IdentifierForVendor == d.IdentifierForVendor && req.Model == d.Model {
//				// Because we have found a matching device we just want to update
//				err = server.store.UpdateUserAPNDetailToken(ctx, db.UpdateUserAPNDetailTokenParams{
//					ID:    d.ID,
//					Token: req.Token,
//				})
//				if err != nil {
//					log.Printf("Error at CreateUserAPNDetail in UpdateUserAPNDetailToken err: %v, user: %v\n", err, ctx.ClientIP)
//					err = fmt.Errorf("could not update your token")
//					ctx.JSON(http.StatusBadRequest, errorResponse(err))
//					return
//				}
//				found = true
//				break
//			}
//		}
//		if !found {
//			err = server.store.CreateUserAPNDetail(ctx, db.CreateUserAPNDetailParams{
//				UserID:              user.ID,
//				DeviceName:          req.Name,
//				Model:               req.Model,
//				IdentifierForVendor: req.IdentifierForVendor,
//				Token:               req.Token,
//			})
//			if err != nil {
//				log.Printf("Error at exist CreateUserAPNDetail in CreateUserAPNDetail err: %v, user: %v\n", err, ctx.ClientIP)
//				err = fmt.Errorf("could not create your token")
//				ctx.JSON(http.StatusBadRequest, errorResponse(err))
//				return
//			}
//		}
//	} else {
//		err = server.store.CreateUserAPNDetail(ctx, db.CreateUserAPNDetailParams{
//			UserID:              user.ID,
//			DeviceName:          req.Name,
//			Model:               req.Model,
//			IdentifierForVendor: req.IdentifierForVendor,
//			Token:               req.Token,
//		})
//		if err != nil {
//			log.Printf("Error at non-exist  CreateUserAPNDetail in CreateUserAPNDetail err: %v, user: %v\n", err, ctx.ClientIP)
//			err = fmt.Errorf("public id not in the right format")
//			ctx.JSON(http.StatusBadRequest, errorResponse(err))
//			return
//		}
//	}
//	res := UserResponseMsg{
//		Success: true,
//	}
//	ctx.JSON(http.StatusOK, res)
//}
