package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) GetReserveHostDetail(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userData, err := server.store.GetUserVerify(ctx, user.ID)
	if err != nil {
		log.Printf("Error at  GetReserveHostDetail in GetUserVerify err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not access your profile")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var identityVerified bool
	if userData.Status != "not_started" {
		identityVerified = true
	}
	res := GetReserveHostDetailRes{
		IdentityVerified: identityVerified,
		HasNumber:        !tools.ServerStringEmpty(userData.PhoneNumber),
		HasProfilePhoto:  !tools.ServerStringEmpty(userData.Photo),
		HasEmail:         !tools.ServerStringEmpty(userData.Email),
		HasFirstName:     !tools.ServerStringEmpty(userData.FirstName),
		HasPayout:        !tools.ServerStringEmpty(userData.DefaultAccountID),
		HasLanguage:      !tools.ServerListIsEmpty(userData.Languages),
		HasBio:           !tools.ServerStringEmpty(userData.Bio),
		UserID:           tools.UuidToString(user.UserID),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListReservationDetail(ctx *gin.Context) {
	var req ListReservationDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at FinalOptionReserveDetail in ShouldBindJSON: %v, Selection: %v \n", err.Error(), req.Selection)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch req.MainOptionType {
	case "events":
		resData, hasData, err := HandleReserveEventHost(req.Selection, server, ctx, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		log.Println("reserve host event count, ", resData)
		resDataOffset := ReserveEventHostItemOffset(resData, req.Offset, 15)
		log.Println("reserve host event offset count, ", resDataOffset)
		if !hasData || len(resDataOffset) == 0 {
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
			return
		}
		res := ReserveEventHostItem{
			List:      resDataOffset,
			Selection: req.Selection,
		}
		ctx.JSON(http.StatusOK, res)
		return
	case "options":
		resData, hasData, err := HandleReserveOptionHost(req.Selection, server, ctx, user)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		resDataOffset := ListReservationDetailResOffset(resData, req.Offset, 15)
		if !hasData || len(resDataOffset) == 0 {
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
			return
		}
		res := ListReservationDetailRes{
			List:      resDataOffset,
			Selection: req.Selection,
		}
		ctx.JSON(http.StatusOK, res)
		return
	default:
		err = fmt.Errorf("nothing here")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return

	}

}
