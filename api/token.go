package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) RenewAccessToken(ctx *gin.Context) {
	var req RenewAccessTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at LoginUser in ShouldBindJSON: %v \n", err)
		err = fmt.Errorf("your session has ended, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	refreshPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		log.Printf("Error at RenewAccessToken at VerifyToken: %v \n", err)
		err = fmt.Errorf("your session has ended, try logging in")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Error at RenewAccessToken in GetSession: %v \n", err)
			err = fmt.Errorf("your session has ended, try logging in")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Printf("Error at RenewAccessToken in GetSession: %v \n", err)
		err = fmt.Errorf("your session has ended, try again")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check if session is blocked
	if session.IsBlocked {
		err := fmt.Errorf("your session is blocked, try logging in")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Check is session usernames matches
	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("incorrect session user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Check is session refreshToken matches
	if session.RefreshToken != req.RefreshToken {
		err := fmt.Errorf("your session is ended, try logging in")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("your session is ended, try logging in")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Get access token
	accessToken, accessPayload, err := server.tokenMaker.CreateTokens(
		refreshPayload.Username,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		log.Printf("Error at RenewAccessToken in CreateTokens: %v \n", err)
		err = fmt.Errorf("error while creating your session, try logging")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	accessPayloadStringTime := tools.ConvertTimeToString(accessPayload.ExpiredAt)
	res := RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayloadStringTime,
	}
	log.Printf("user access token in successfully (%v) \n", refreshPayload.Username)
	ctx.JSON(http.StatusOK, res)
}
