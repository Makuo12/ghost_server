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
		log.Println("used location")
		options = HandleOptionExSearchLocation(ctx, server, req, "ListOptionExSearch")
	} else {
		log.Println("not used location")
		options = HandleOptionExSearch(ctx, server, req, "ListOptionExSearch")
	}
	optionIndexData := GetExperienceOptionOffset(options, req.Offset, 10)
	if len(optionIndexData) == 0 {
		res := ExperienceCategoryRes{
			Category: req.Type,
		}
		ctx.JSON(http.StatusNoContent, res)
		return
	}

	res := ListExperienceOptionRes{
		List:         optionIndexData,
		OptionOffset: req.Offset + len(optionIndexData),
		OnLastIndex:  false,
		Category:     req.Type,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListEventExSearch(ctx *gin.Context) {
	var req ExControlEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListEventExSearch in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	events, err := server.store.ListEvent(ctx)
	if err != nil || len(events) == 0 {
		if err != nil {
			log.Printf("Error at FuncName %v, HandleEventSearch ListEvent err: %v \n", "ListEventExSearch", err.Error())
		}
		ctx.JSON(http.StatusNoContent, errorResponse(err))
		return
	}
	startPrice, endPrice := HandleExEventFilterPrice(ctx, server, "ListEventExSearch", req)
	//var eventData
	var eventData []ExperienceEventData
	if ExSearchReqHasLocation(req.Search) {
		eventData = HandleExEventLocation(ctx, server, "ListEventExSearch", req, events, startPrice, endPrice)
	} else {
		eventData = HandleExEvent(ctx, server, "ListEventExSearch", req, events, startPrice, endPrice)
	}
	eventIndexData := GetExperienceEventOffset(eventData, req.Offset, 10)
	if len(eventIndexData) == 0 {
		res := ExperienceCategoryRes{
			Category: req.Type,
		}
		ctx.JSON(http.StatusNoContent, res)
		return

	}
	res := ListExperienceEventRes{
		List:         eventIndexData,
		OptionOffset: req.Offset + len(eventIndexData),
		OnLastIndex:  false,
		Category:     req.Type,
	}
	ctx.JSON(http.StatusOK, res)
}
