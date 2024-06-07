package api

import (
	"fmt"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	//"github.com/makuo12/ghost_server/tools"
	//"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeDurationTwilioUser() time.Duration {
	// time.Now().Local().UTC() returns the time an hour before
	t_2 := time.Now().Local().UTC().Add(time.Hour * 3)
	t_1 := time.Now().Local().UTC()
	duration := t_2.Sub(t_1)
	return duration
}

func TimeTwilioUser() time.Time {
	// time.Now().Local().UTC() returns the time an hour before
	return time.Now().Local().UTC().Add(time.Hour * 3)
}

func UserIsHost(ctx *gin.Context, server *Server, user db.User) (isHost bool, hasIncomplete bool) {
	data, err := server.store.GetOptionInfoAllCount(ctx, db.GetOptionInfoAllCountParams{
		HostID:   user.ID,
		CoUserID: tools.UuidToString(user.UserID),
	})
	if err != nil {
		log.Printf("Error at GetUserIsHost in GetOptionInfoAllCount err: %v, user: %v", err, user.ID)
	} else {
		// If the user has any complete options that makes the user a host
		if data.Complete > 0 {
			isHost = true
		}
		// If the user has incomplete we would want to tell them so they can complete it
		if data.InComplete > 0 {
			hasIncomplete = true
		}
	}
	return
}

func HandleWishlist(ctx *gin.Context, server *Server, user db.User) (list ListWishlistRes) {
	emptyWish := []WishlistItem{{constants.NONE, constants.NONE, constants.NONE, constants.NONE, constants.NONE}}
	wishlistItems, err := server.store.ListWishlistItemUser(ctx, user.ID)
	if err != nil {
		log.Printf("Error at HandleWishlist in ListWishlistItemUser err: %v, user: %v", err, user.ID)
		list = ListWishlistRes{
			List:    emptyWish,
			IsEmpty: true,
		}
	} else {
		resData := []WishlistItem{}
		for _, item := range wishlistItems {
			data := WishlistItem{
				Name:           item.Name,
				WishlistID:     tools.UuidToString(item.WishlistID),
				WishlistItemID: tools.UuidToString(item.ID),
				OptionUserID:   tools.UuidToString(item.OptionUserID),
				CoverImage:     item.CoverImage,
			}
			resData = append(resData, data)
		}
		list = ListWishlistRes{
			List:    resData,
			IsEmpty: false,
		}
	}
	return

}

// useType tells me whether it is for login, createUser, etc
func HandleUserSession(ctx *gin.Context, server *Server, user db.User, useType string) (accessToken string, refreshToken string, accessPayloadStringTime string, err error) {
	// Get access token
	accessToken, accessPayload, err := server.tokenMaker.CreateTokens(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Printf("error %v at HandleUserSession for CreateTokens access token for type%v for user%v\n", err, useType, user.Username)
		err = fmt.Errorf("there was an error with the server while logging you in, try again")
		return
	}
	// Get refresh token
	refreshToken, refreshPayload, err := server.tokenMaker.CreateTokens(
		user.Username,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		log.Printf("error %v at HandleUserSession for CreateTokens refresh token for type%v for user%v\n", err, useType, user.Username)
		err = fmt.Errorf("there was an error with the server while logging you in, try again")
		return
	}
	//We delete all past sessions
	err = server.store.DeleteSessionByClientID(ctx, ctx.Request.UserAgent())
	if err != nil {
		log.Printf("error deleting session of user %v: %v", ctx.Request.UserAgent(), err)
	}

	// We also try deleting using the user username
	err = server.store.DeleteSession(ctx, user.Username)
	if err != nil {
		log.Printf("error %v at HandleUserSession for DeleteSession for type%v for user%v\n", err, useType, user.Username)
	}
	_, err = server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(), // TODO: fill it
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		log.Printf("error %v at HandleUserSession for CreateSession for type%v for user%v\n", err, useType, user.Username)
		err = fmt.Errorf("there was an error with the server while logging you in, try again")
		return
	}
	accessPayloadStringTime = tools.ConvertTimeToString(accessPayload.ExpiredAt)
	return
}

func UserUnreadMessages(ctx *gin.Context, server *Server, user db.User, funcName string) (count int) {
	contacts, err := server.store.ListMessageContactNoLimit(ctx, user.UserID)
	if err != nil || len(contacts) == 0 {
		if err != nil {
			log.Printf("Error for FuncName: %v, at UserUnreadMessages, in ListMessageContactNoLimit error: %v, userID: %v \n", funcName, err, user.ID)
			err = nil
		}

		return
	}
	for _, contact := range contacts {
		count += int(contact.UnreadMessageCount)
	}
	return
}

func UserUnreadNotifications(ctx *gin.Context, server *Server, user db.User, funcName string) (count int64) {
	count, err := server.store.CountNotificationNoLimit(ctx, user.UserID)
	if err != nil {
		log.Printf("Error for FuncName: %v, at UserUnreadMessages, in ListMessageContactNoLimit error: %v, userID: %v \n", funcName, err, user.ID)
		err = nil
		count = 0
	}
	return
}
