package api

import (
	// "database/sql"
	// "errors"
	// "log"
	// "strings"

	db "flex_server/db/sqlc"
	"flex_server/token"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	//"github.com/lib/pq"
)

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("error at CreateUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("your password should be more than 8 characters and contain a special character($,@,&,*)")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		log.Printf("error at CreateUser in HashedPassword: %v \n", err)
		err = fmt.Errorf("there was an error while processing your password, make sure strong with at least 10 characters")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	date, err := tools.ConvertDateOnlyStringToDate(req.DateOfBirth)
	if err != nil {
		log.Printf("error at CreateUser in ConvertDateOnlyStringToDate: %v \n", err)
		err = fmt.Errorf("there was an error in while verifying your birthday, try confirming the birthday entered")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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
	split := strings.Split(result, "&")
	_, _, _, email, confirm := split[0], split[1], split[2], split[3], split[4]
	if confirm != "true" {
		err = fmt.Errorf("you are yet to be verified")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var username string = uuid.New().String()
	var firebasePassword string = uuid.New().String()
	arg := db.CreateUserParams{
		HashedPassword:   hashedPassword,
		Email:            strings.ToLower(email),
		Username:         username,
		FirebasePassword: firebasePassword,
		DateOfBirth:      date,
		FirstName:        strings.ToLower(strings.TrimSpace(req.FirstName)),
		LastName:         strings.ToLower(strings.TrimSpace(req.LastName)),
		Currency:         req.Currency,
	}
	user, err := server.store.CreateUser(ctx, arg)
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
	log.Println("email is: ", user.Email)
	// We want to create the user in firebase
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

	err = server.store.CreateIdentity(ctx, db.CreateIdentityParams{
		UserID:          user.ID,
		IDPhotoList:     []string{"none"},
		FacialPhotoList: []string{"none"},
		IDBackPhotoList: []string{"none"},
	})
	if err != nil {
		log.Printf("Error at CreateUser at CreateIdentity %v for user: %v \n", err, user.ID)
	}

	err = server.store.CreateUserProfile(ctx, db.CreateUserProfileParams{
		UserID:    user.ID,
		Languages: []string{"none"},
	})
	if err != nil {
		log.Printf("Error at CreateUser at CreateUserProfile %v for user: %v \n", err, user.ID)

	}

	// We create a user account for the user
	// We create dollar account
	err = server.store.CreateAccount(ctx, db.CreateAccountParams{
		UserID:   user.ID,
		Currency: utils.USD,
	})
	if err != nil {
		log.Printf("error at CreateUser USD in server.store.CreateAccount user account was not created %v \n", err.Error())
	}
	// We create naira account
	err = server.store.CreateAccount(ctx, db.CreateAccountParams{
		UserID:   user.ID,
		Currency: utils.NGN,
	})
	if err != nil {
		log.Printf("error at CreateUser NGN in server.store.CreateAccount user account was not created %v \n", err.Error())
	}
	err = RedisClient.Del(RedisContext, req.Username).Err()
	if err != nil {
		log.Printf("error at CreateUser on DEL function for user (%s) error message: %v \n", req.Username, err)
	}
	accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "CreateUser")
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

	log.Printf("user created successfully (%v) \n", user.Email)
	ctx.JSON(http.StatusCreated, res)
}

func (server *Server) LoginUser(ctx *gin.Context) {
	var req LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at LoginUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("your password should be more than 8 characters and contain a special character($,@,&,*)")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserWithEmail(ctx, req.Email)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("error at LoginUser in GetUserWithEmail: %v \n", err)
			err = fmt.Errorf("this email isn't registered with Flizzup, try entering the email again or signing up")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("error at LoginUser in GetUserWithEmail: %v \n", err)
		err = fmt.Errorf("there was an error while logging you in, try again")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		log.Printf("error at LoginUser at CheckPassword: %v \n", err)
		err = fmt.Errorf("this password is incorrect, try entering it again or forgot password")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "login_user")
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

func (server *Server) UpdateCurrency(ctx *gin.Context) {
	var req UpdateCurrencyParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateCurrency in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while updating your currency, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		log.Printf("error at ctx.GET does not exist")
		err := fmt.Errorf("there was an error while updating your currency, make sure you are logged in")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	data := payload.(*token.Payload)

	var username string = *&data.Username
	user, err := server.store.GetUserWithUsername(ctx, username)

	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("error at UpdateCurrency in GetUserWithUsername: %v \n", err)
			err = fmt.Errorf("this account isn't registered with Flizzup")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("error at UpdateCurrency in GetUserWithUsername: %v \n", err)
		err = fmt.Errorf("there was an error while updating your currency, make sure you are sign up on Flizzup")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !user.IsActive {
		err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// We just want to update the currency
	userUpdate, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		Currency: pgtype.Text{
			String: req.Currency,
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("error at UpdateCurrency in UpdateUser HGet: %v \n", err)
		err = fmt.Errorf("there was an error while updating your currency, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	res := UpdateUserResponse{
		Email:         userUpdate.Email,
		FireFight:     userUpdate.FirebasePassword,
		FirstName:     userUpdate.FirstName,
		LastName:      userUpdate.LastName,
		IsHost:        userIsHost,
		Currency:      userUpdate.Currency,
		ProfilePhoto:  userUpdate.Photo,
		HasIncomplete: hasIncomplete,
	}
	log.Printf("user currency was updated in UpdateCurrency successfully (%v) \n", user.Email)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateUserInfo(ctx *gin.Context) {
	var req UpdateUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateUserInfo in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("there was an error while updating your currency, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	birthTime, err := tools.ConvertDateOnlyStringToDate(req.DateOfBirth)
	if err != nil {
		log.Printf("error at UpdateUserInfo in ConvertDateOnlyStringToDate: %v \n", err)
		err = fmt.Errorf("there was an error in while verifying your birthday, try confirming the birthday entered")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// We want to update the first name, last name, and date of birth
	userUpdateData, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		DateOfBirth: pgtype.Date{
			Time:  birthTime,
			Valid: true,
		},
		FirstName: pgtype.Text{
			String: strings.ToLower(req.FirstName),
			Valid:  true,
		},
		LastName: pgtype.Text{
			String: strings.ToLower(req.LastName),
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("error at UpdateUserInfo in UpdateUser: %v \n", err)
		err = fmt.Errorf("there was an error in while verifying your birthday, try confirming the birthday entered")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// We want to update the user display name in firebase
	params := (&auth.UserToUpdate{}).
		DisplayName(userUpdateData.FirstName)
	u, err := server.ClientFire.UpdateUser(ctx, tools.UuidToString(userUpdateData.FirebaseID), params)
	if err != nil {
		log.Printf("error updating user display name in firebase: %v, userID: %v\n", err, user.ID)
		err = nil
	} else {
		log.Printf("Successfully updated user display name in firebase: %v\n", u)
	}

	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	userUpdateRes := UpdateUserResponse{
		Email:         userUpdateData.Email,
		FireFight:     userUpdateData.FirebasePassword,
		FirstName:     userUpdateData.FirstName,
		LastName:      userUpdateData.LastName,
		ProfilePhoto:  userUpdateData.Photo,
		Currency:      userUpdateData.Currency,
		IsHost:        userIsHost,
		DateOfBirth:   tools.ConvertDateOnlyToString(userUpdateData.DateOfBirth),
		HasIncomplete: hasIncomplete,
	}
	log.Printf("user info data successfully updated %v\n", userUpdateData.Email)
	ctx.JSON(http.StatusOK, userUpdateRes)
}

func (server *Server) UpdateUserPassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at UpdateUserPassword in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("your password should be more than 8 characters and contain a special character($,@,&,*)")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Password == req.CurrentPassword {
		err := fmt.Errorf("the new password and current password cannot be the same")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Password != req.PasswordTwo {
		err := fmt.Errorf("the new password and confirm password must be the same")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.CurrentPassword, user.HashedPassword)
	if err != nil {
		err = fmt.Errorf("the current password entered was incorrect. Try entering the password again or press forget password")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	hashedPassword, err := utils.HashedPassword(req.Password)
	if err != nil {
		log.Printf("error occurred while hashing your password in HashedPassword %v\n", err.Error())
		err = fmt.Errorf("an error occurred while processing your password, try again")
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

	accessToken, refreshToken, accessPayloadStringTime, err := HandleUserSession(ctx, server, user, "update_password")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ChangePasswordResponse{
		Updated:              true,
		AccessToken:          accessToken,
		RefreshToken:         refreshToken,
		AccessTokenExpiresAt: accessPayloadStringTime,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetUser(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	wishlist := HandleWishlist(ctx, server, user)
	res := GetUserParams{
		Email:         user.Email,
		FireFight:     user.FirebasePassword,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		IsHost:        userIsHost,
		ProfilePhoto:  user.Photo,
		Currency:      user.Currency,
		WishlistList:  wishlist,
		HasIncomplete: hasIncomplete,
	}
	log.Printf("user logged in successfully (%v) \n", user.Email)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) LogoutUser(ctx *gin.Context) {
	//We delete all past sessions
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.DeleteSessionByClientID(ctx, ctx.Request.UserAgent())
	if err != nil {
		log.Printf("error deleting session of user %v: %v", ctx.Request.UserAgent(), err)
	}
	// We also try deleting using the user username
	err = server.store.DeleteSession(ctx, user.Username)
	if err != nil {
		log.Printf("error %v at HandleUserSession for DeleteSession for type%v for user%v\n", err, "logout", user.Username)
	}
	// We also delete all the apns
	err = server.store.RemoveAllUserAPNDetail(ctx, user.ID)
	if err != nil {
		log.Printf("error %v at HandleUserSession for RemoveAllUserAPNDetail for type%v for user%v\n", err, "logout", user.Username)
	}
	res := RemoveResponse{
		Success: true,
	}
	log.Printf("user logged out successfully (%v) \n", ctx.Request.UserAgent())
	ctx.JSON(http.StatusNoContent, res)
}

func (server *Server) GetAppPolicy(ctx *gin.Context) {
	var req GetAppPolicyParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("error at GetAppPolicy in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("Try again later")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var link string
	switch req.Type {
	case "term_of_service":
		link = server.config.TermsOfService
	case "privacy_policy":
		link = server.config.PrivacyPolicy
	}
	res := GetAppPolicyRes{
		Type: req.Type,
		Link: link,
	}
	log.Printf("user got app service in successfully (%v) \n", ctx.ClientIP())
	ctx.JSON(http.StatusOK, res)
}
