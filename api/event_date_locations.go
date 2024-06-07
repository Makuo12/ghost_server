package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) CreateUpdateEventDateLocation(ctx *gin.Context) {
	var req CreateUpdateEventDateLocationReq
	var dataExist bool
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateEventDateLocationParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("please make sure yuo select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at CreateUpdateEventDateLocation at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set location for this event date, please try again using the format on the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateLocation(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at  GetEventDateLocation in GetEventDateLocation err: %v, user: %v\n", err, user.ID)
		dataExist = false
	} else {
		dataExist = true
	}
	var res CreateUpdateEventDateLocationRes
	var eventDateLocation db.EventDateLocation
	if dataExist {
		// If data exist we want to update
		lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
		lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
		geolocation := pgtype.Point{
			P:     pgtype.Vec2{X: lng, Y: lat},
			Valid: true,
		}
		eventDateLocation, err = server.store.UpdateEventDateLocationTwo(ctx, db.UpdateEventDateLocationTwoParams{
			Street:          req.Street,
			City:            req.City,
			State:           req.State,
			Country:         req.Country,
			Postcode:        req.Postcode,
			Geolocation:     geolocation,
			EventDateTimeID: eventDateTimeID,
		})
		if err != nil {
			log.Printf("Error at  CreateUpdateEventDateLocation in UpdateEventDateLocationTw err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("an error occurred while updating your event date location")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	} else {

		// If data does not exist we want to update
		lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
		lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
		geolocation := pgtype.Point{
			P:     pgtype.Vec2{X: lng, Y: lat},
			Valid: true,
		}
		eventDateLocation, err = server.store.CreateEventDateLocation(ctx, db.CreateEventDateLocationParams{
			Street:          req.Street,
			City:            req.City,
			State:           req.State,
			Country:         req.Country,
			Postcode:        req.Postcode,
			Geolocation:     geolocation,
			EventDateTimeID: eventDateTimeID,
		})
		if err != nil {
			log.Printf("Error at  CreateUpdateEventDateLocation in .CreateEventDateLocation err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("an error occurred while creating your event date location")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	res = CreateUpdateEventDateLocationRes{
		EventDateTimeID: tools.UuidToString(eventDateLocation.EventDateTimeID),
		Street:          eventDateLocation.Street,
		City:            eventDateLocation.City,
		State:           eventDateLocation.State,
		Country:         eventDateLocation.Country,
		Postcode:        eventDateLocation.Postcode,
		Lat:             tools.ConvertFloatToLocationString(eventDateLocation.Geolocation.P.Y, 9),
		Lng:             tools.ConvertFloatToLocationString(eventDateLocation.Geolocation.P.X, 9),
	}
	fmt.Println(res)
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateEventDateLocation", "event date location", "update event date location")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetEventDateLocation(ctx *gin.Context) {
	var req GetEventDateLocationReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetEventDateLocationParams in ShouldBindJSON: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("some required inputs were not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	eventID, err := tools.StringToUuid(req.EventInfoID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventInfoID: %v \n", err.Error(), req.EventInfoID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, EventDateTimeID: %v \n", err.Error(), req.EventDateTimeID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditEventDateTimes(eventID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.GetEventDateTimeByOption(ctx, db.GetEventDateTimeByOptionParams{
		ID:         eventDateTimeID,
		ID_2:       user.ID,
		ID_3:       option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("There an error at GetEventDateLocation at GetEventDateTimeByOption: %v, EventInfoID: %v, userID: %v \n", err.Error(), eventDateTimeID, user.ID)
		err = fmt.Errorf("could not set location for this event date, please try again using the format on the app")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	eventDateLocation, err := server.store.GetEventDateLocation(ctx, eventDateTimeID)
	if err != nil {
		log.Printf("Error at  GetEventDateLocation in GetEventDateLocation err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return

	}
	res := CreateUpdateEventDateLocationRes{
		EventDateTimeID: tools.UuidToString(eventDateLocation.EventDateTimeID),
		Street:          eventDateLocation.Street,
		City:            eventDateLocation.City,
		State:           eventDateLocation.State,
		Country:         eventDateLocation.Country,
		Postcode:        eventDateLocation.Postcode,
		Lat:             tools.ConvertFloatToLocationString(eventDateLocation.Geolocation.P.Y, 9),
		Lng:             tools.ConvertFloatToLocationString(eventDateLocation.Geolocation.P.X, 9),
	}
	fmt.Println(res)

	ctx.JSON(http.StatusOK, res)
}
