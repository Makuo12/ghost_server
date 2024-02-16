package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListOptionExSearch(ctx *gin.Context) {
	var req ExControlOptionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListOptionExSearch in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	//var options
	var options []ExperienceOptionData
	if ExSearchReqHasLocation(req.Search) {
		options = HandleOptionExSearchLocation(ctx, server, req, "ListOptionExSearch")
	} else {
		options = HandleOptionExSearch(ctx, server, req, "ListOptionExSearch")
	}
	optionIndexData := GetExperienceOptionOffset(options, req.Offset, 10)
	if len(optionIndexData) == 0 {
		res := ExperienceCategoryRes{
			Category: req.Type,
		}
		ctx.JSON(http.StatusNoContent, res)
	}
	ctx.JSON(http.StatusOK, optionIndexData)
}
