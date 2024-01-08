package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/utils"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) AccountChange(ctx *gin.Context) {
	var req AccountChangeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("error at AccountChange in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var mainHeader string
	var header string
	var msgData string
	switch req.Type {
	case "delete":
		mainHeader = fmt.Sprintf("Account deletion")
		header = fmt.Sprintf("Request for account to be deleted")
		msgData = fmt.Sprintf("Hey %v, we just received a request for your account to be deleted if this was not you just ignore this email and we would handle the rest. However, if this was you just enter in the six-digit code below\nNote that this would delete your data with us and you would no more have access to this account or anything related to this account", tools.CapitalizeFirstCharacter(user.FirstName))
	default:
		err = fmt.Errorf("This page not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	username := utils.RandomName()
	err = BrevoAccountChange(ctx, server, user.Email, tools.CapitalizeFirstCharacter(user.FirstName), username, "AccountChange", mainHeader, header, msgData)
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
}

func (server *Server) VerifyAccountChange(ctx *gin.Context) {
	var req VerifyAccountChangeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//if there is an error we want to send it back to the user
		log.Printf("error at VerifyAccountChange in ShouldBindJSON: %v \n", err.Error())
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch req.Type {
	case "delete":
		exist, err := RedisClient.Exists(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			if err != nil {
				log.Printf("Error at VerifyAccountChange Delete in RedisClient Exists: %v \n", err.Error())
			}
			err = fmt.Errorf("details of your code was not found, try again as it might have expired")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		result, err := RedisClient.Get(RedisContext, req.Username).Result()
		if err != nil || exist == 0 {
			log.Printf("Error at VerifyAccountChange Delete in RedisClient  RedisClient.Get: %v \n", err.Error())
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
		_, code, _, _, _ := split[0], split[1], split[2], split[3], split[4]
		if code != req.Code {
			err = fmt.Errorf("the code was incorrect, enter the code sent to you via email or press send again to receive a new code")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
			IsDeleted: pgtype.Bool{
				Valid: true,
				Bool:  true,
			},
			ID: user.ID,
		})
		if err != nil {
			log.Printf("error deleting user account of user %v: %v", ctx.Request.UserAgent(), err)
			err = fmt.Errorf("could not delete this account try again")
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
	default:
		err = fmt.Errorf("This page not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := RemoveResponse{
		Success: true,
	}
	log.Printf("user deleted successfully (%v) \n", ctx.Request.UserAgent())
	ctx.JSON(http.StatusNoContent, res)
}
