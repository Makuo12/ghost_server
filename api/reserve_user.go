package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListReserveUserItem(ctx *gin.Context) {
	var req ListReserveUserItemParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListReserveUserItem in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Printf("reserves req %v \n", req)
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ListReserveUserItemRes
	if req.Type == "current" {
		switch req.MainOption {
		case "options":
			resOption, hasData, err := HandleCurrentListReserveUserOptionItem(server, ctx, user, req)
			if err != nil || !hasData {
				data := "none"
				ctx.JSON(http.StatusNoContent, data)
				return
			} else {
				res = resOption
			}
		case "events":
			resEvent, hasData, err := HandleCurrentListReserveUserEventItem(server, ctx, user, req)
			if err != nil || !hasData {
				data := "none"
				ctx.JSON(http.StatusNoContent, data)
				return
			} else {
				res = resEvent
			}
		}
	} else if req.Type == "visited" {
		switch req.MainOption {
		case "options":
			resOption, hasData, err := HandleVisitedListReserveUserOptionItem(server, ctx, user, req)
			if err != nil || !hasData {
				data := "none"
				ctx.JSON(http.StatusNoContent, data)
				return
			} else {
				res = resOption
			}
		case "events":
			resEvent, hasData, err := HandleVisitedListReserveUserEventItem(server, ctx, user, req)
			if err != nil || !hasData {
				data := "none"
				ctx.JSON(http.StatusNoContent, data)
				return
			} else {
				res = resEvent
			}
		}
	} else {
		err = fmt.Errorf("page not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetReserveUserDirection(ctx *gin.Context) {
	var req ReserveUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetReserveUserDirection in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ReserveUserDirectionRes
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleGetRUOptionDirection(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resOption
		}
	case "events":
		resEvent, hasData, err := HandleGetRUEventDirection(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resEvent
		}
	}

	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetRUCheckInStep(ctx *gin.Context) {
	var req ReserveUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetRUOptionDirection in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res RUListCheckInStepRes
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleGetRUOptionCheckInStep(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resOption
		}
	case "events":
		resEvent, hasData, err := HandleGetRUEventCheckInStep(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resEvent
		}
	}

	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetRUHelp(ctx *gin.Context) {
	var req ReserveUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetRUOptionDirection in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res RUHelpManualRes
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleGetRUOptionHelp(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resOption
		}
	case "events":
		resEvent, hasData, err := HandleGetRUEventHelp(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resEvent
		}
	}

	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetRUWifi(ctx *gin.Context) {
	var req ReserveUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetRUOptionDirection in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ReserveUserWifiRes
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleGetRUOptionWifi(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resOption
		}
	case "events":
		resEvent, hasData, err := HandleGetRUEventWifi(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			res = resEvent
		}
	}

	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetRUReceipt(ctx *gin.Context) {
	var req ReserveUserInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetRUOptionDirection in ShouldBindJSON: %v, ID: %v \n", err.Error(), req.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch req.MainOption {
	case "options":
		resOption, hasData, err := HandleGetRUOptionReceipt(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			ctx.JSON(http.StatusOK, resOption)
			return
		}
	case "events":
		resEvent, hasData, err := HandleGetRUEventReceipt(server, ctx, user, req)
		if err != nil || !hasData {
			data := "none"
			ctx.JSON(http.StatusNoContent, data)
			return
		} else {
			ctx.JSON(http.StatusOK, resEvent)
			return
		}
	}

}
