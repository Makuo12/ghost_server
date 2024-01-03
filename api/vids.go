package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func (server *Server) ListVid(ctx *gin.Context) {
	var req ListVidParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListVid in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountVid(ctx)
	if err != nil {
		log.Printf("Error at  ListVid in .CountVid err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}

	vids, err := server.store.ListVid(ctx, db.ListVidParams{
		Limit:  40,
		Offset: int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at ListVid in .ListVid err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []VidItem
	for _, v := range vids {
		data := VidItem{
			Path:           v.Path,
			Filter:         v.Filter,
			OptionUserID:   tools.UuidToString(v.OptionUserID),
			MainOptionType: v.MainOptionType,
			StartDate:      v.StartDate,
			Caption:        v.Caption,
			ExtraOptionID:  tools.UuidToString(v.ExtraOptionID),
		}
		resData = append(resData, data)
	}
	res := ListVidRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}


