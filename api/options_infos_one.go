package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"
	"github.com/makuo12/ghost_server/val"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// OptionType would take shortlets has it answer for main type as options
// OptionType would take events has it answer for main types as events
// MainOptionType can only be events, options, experiences

// SE MEANS SHORTLET EVENT
func (server *Server) CreateOptionSE(ctx *gin.Context) {
	var currentState string
	var previousState string
	var req CreateOptionSE
	var optionIDExist bool = true
	var optionIDReal bool = true
	var optionID uuid.UUID
	var optionData db.OptionsInfo
	var optionItemType = ""
	//var shortlet db.Shortlet
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOption in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), requestID)
		optionIDReal = false
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// validate type
	result := val.ValidateOptionType(req.OptionType)
	if !result {
		err = fmt.Errorf("an error occurred, this option type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if optionIDReal {
		// we check if the option exist
		option, err := server.store.GetOptionInfo(ctx, db.GetOptionInfoParams{
			ID:         requestID,
			HostID:     user.ID,
			IsComplete: false,
		})
		if err != nil {
			err = fmt.Errorf("this request is forbidden")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		} else {
			optionIDExist = true
			optionID = option.ID
			optionData = option

		}

	}
	if optionIDExist {
		err = server.store.RemoveCompleteOptionInfo(ctx, optionID)
		if err != nil {
			log.Printf("error at CreateOptionSE at RemoveCompleteOptionInfo: %v, optionID: %v, userID: %v\n", err.Error(), optionID, user.ID)
			err = fmt.Errorf("error occurred while performing your request")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		switch req.OptionType {
		case "shortlets":
			err = server.store.RemoveShortlet(ctx, optionID)
			if err != nil {
				log.Printf("error at RemoveOption at RemoveShortlet: %v, optionID: %v, userID: %vv\n", err.Error(), optionID, user.ID)
				err = fmt.Errorf("error occurred while taking you back, please try again")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			// Remove all default tables
			HandleRemoveOptionDefaultTables(server, ctx, optionData, user)
		case "events":
			err = server.store.RemoveEventInfo(ctx, optionID)
			if err != nil {
				log.Printf("error at RemoveOption at RemoveEventInfo: %v, optionID: %v, userID: %vv\n", err.Error(), optionID, user.ID)
				err = fmt.Errorf("error occurred while taking you back, please try again")
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return
			}
			// Remove Option Info Detail
			HandleRemoveEventDefaultTables(server, ctx, optionData, user)
		}
		err = server.store.RemoveOptionInfo(ctx, db.RemoveOptionInfoParams{
			ID:     requestID,
			HostID: user.ID,
		})
		if err != nil {
			log.Printf("error at CreateOptionSE at RemoveOptionInfo: %v, optionID: %v, userID: %v\n", err.Error(), optionID, user.ID)
			err = fmt.Errorf("error occurred while performing your request")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	arg := db.CreateOptionInfoParams{
		HostID:         user.ID,
		OptionType:     req.OptionType,
		MainOptionType: req.MainOptionType,
		Currency:       req.Currency,
		OptionImg:      req.OptionImg,
		PrimaryUserID:  user.UserID,
		TimeZone:       req.TimeZone,
	}
	option, err := server.store.CreateOptionInfo(ctx, arg)
	if err != nil {
		log.Printf("Error at CreateOption in CreateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("an error occurred while creating your option, try again or contact help")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if option.OptionType == "shortlets" {
		shortlet, err := CreateOptionShortletType(server, ctx, option, user.ID, req.ShortletType)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		currentState = utils.ShortletSpace
		previousState = utils.OptionType
		optionItemType = shortlet.TypeOfShortlet
		// Create all default tables
		HandleCreateOptionDefaultTables(server, ctx, option, user)
	}
	// we can to create an event info if it is of type event
	if option.MainOptionType == "events" {
		eventInfo, err := CreateEventType(server, ctx, option, user.ID, req.EventType)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		currentState = utils.EventSubType
		previousState = utils.OptionType
		optionItemType = eventInfo.EventType
		// Create Option detail
		HandleCreateEventDefaultTables(server, ctx, option, user)
	}
	// we create an option info detail for all of them

	argComplete := db.CreateCompleteOptionInfoParams{
		OptionID:      option.ID,
		CurrentState:  currentState,
		PreviousState: previousState,
	}
	// we create our completion
	completeOption, err := server.store.CreateCompleteOptionInfo(ctx, argComplete)
	if err != nil {
		log.Printf("Error at CreateOption in CreateCompleteOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = server.store.RemoveOptionInfo(ctx,
			db.RemoveOptionInfoParams{
				ID:     option.ID,
				HostID: user.ID,
			})
		if err != nil {
			log.Printf("Error at CreateOption in RemoveOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		}
		err = fmt.Errorf("error occurred processing your state. Try again, if it continues try restarting the app")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	userIsHost, hasIncomplete := UserIsHost(ctx, server, user)
	// We use  OptionInfoFirstResponse because the user.isHost changes here
	res := OptionInfoFirstResponse{
		OptionItemType:     optionItemType,
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
		UserIsHost:         userIsHost,
		HasIncomplete:      hasIncomplete,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateShortletSpace(ctx *gin.Context) {
	var req CreateShortletSpaceParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateShortletSpaceParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateShortletInfoParams{
		SpaceType: pgtype.Text{
			String: req.SpaceType,
			Valid:  true,
		},
		OptionID: option.ID,
	}

	shortlet, err := server.store.UpdateShortletInfo(ctx, arg)
	if err != nil {
		log.Printf("There an error at CreateShortletSpace at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

		err = fmt.Errorf("an error occurred while adding your space type. Please make sure you selected an option made available to you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(utils.ShortletInfo, utils.ShortletSpace, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at CreateShortletSpace at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	// OptionItemType will be stuff like full_place
	res := OptionInfoResponse{
		OptionItemType:     shortlet.SpaceType,
		Success:            true,
		OptionID:           tools.UuidToString(shortlet.OptionID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateShortletInfo(ctx *gin.Context) {
	var req CreateShortletInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateShortletInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.RemoveSpaceAreaAll(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Location doesn't exist %v \n", err.Error())
		} else {
			log.Println(err.Error())
			err = fmt.Errorf("error occurred while performing your request. Try again")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}

	// We check if bedroom exist
	bedroomFound, numGuests := val.ContainsBedroomAndNumGuest(req.Space)
	if !bedroomFound || numGuests == 0 {
		err = fmt.Errorf("you must have a bedroom and welcome a guest")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	argUpdate := db.UpdateShortletInfoParams{
		GuestWelcomed: pgtype.Int4{
			Int32: int32(numGuests),
			Valid: true,
		},
		AnySpaceShared: pgtype.Bool{
			Bool:  req.AnySpaceShared,
			Valid: true,
		},
		OptionID: option.ID,
	}
	shortlet, err := server.store.UpdateShortletInfo(ctx, argUpdate)
	if err != nil {
		log.Printf("There an error at CreateShortletInfo at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	for i := 0; i < len(req.Space); i++ {
		if req.Space[i] != "guest" {
			found := val.ContainsString(val.GuestAreas, req.Space[i])
			if found {
				// we want to create a new space
				err = HandleCreateGuestArea(server, ctx, option, user.ID, req.Space[i], shortlet.AnySpaceShared)
				if err != nil {
					argUp := db.UpdateShortletInfoParams{
						GuestWelcomed: pgtype.Int4{
							Int32: int32(0),
							Valid: true,
						},
						OptionID: option.ID,
					}
					_, err := server.store.UpdateShortletInfo(ctx, argUp)
					if err != nil {
						log.Printf("There an error at CreateShortletInfo at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

					}
					err = fmt.Errorf("error occurred while adding your rooms and spaces, please try again later")
					ctx.JSON(http.StatusBadRequest, errorResponse(err))
					return
				}

			}
		}
	}
	completeOption, err := HandleCompleteOption(utils.LocationView, utils.ShortletInfo, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at CreateShortletSpace at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(shortlet.OptionID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateLocation(ctx *gin.Context) {
	var req CreateLocationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateLocation in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	// check if lat and lng are good
	if !tools.CheckStringIsFloat(req.Lat) && !tools.CheckStringIsFloat(req.Lng) {
		err := fmt.Errorf("your latitude and longitude are in the wrong format. Try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.RemoveLocation(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Location doesn't exist %v \n", err.Error())
		} else {
			log.Println(err.Error())
			err = fmt.Errorf("error occurred while performing your request. Try again")
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
	}
	currentState, previousState := utils.LocationViewState(option.OptionType)
	lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
	lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
	geolocation := pgtype.Point{
		P:     pgtype.Vec2{X: lng, Y: lat},
		Valid: true,
	}
	arg := db.CreateLocationParams{
		OptionID:             option.ID,
		Street:               req.Street,
		City:                 req.City,
		State:                req.State,
		Country:              req.Country,
		Postcode:             req.Postcode,
		Geolocation:          geolocation,
		ShowSpecificLocation: req.ShowSpecificLocation,
	}

	_, err = server.store.CreateLocation(ctx, arg)
	if err != nil {
		log.Printf("Error at CreateLocation in CreateLocation: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("an error occurred while setting up location. Please select one of the space options provided")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateLocation in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) CreateAmenitiesAndSafety(ctx *gin.Context) {
	var req CreateAmenitiesAndSafety
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateAmenitiesAndSafety in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We create the amenities
	err = server.store.RemoveAllAmenity(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Amenities doesn't exist %v \n", err.Error())
		} else {
			log.Println(err.Error())
		}
	}
	// For popular am
	if len(req.PopularAm) > 0 {
		for i := 0; i < len(req.PopularAm); i++ {
			arg := db.CreateAmenityParams{
				OptionID:    option.ID,
				Tag:         req.PopularAm[i],
				AmType:      "popular",
				HasAm:       true,
				ListOptions: []string{},
			}
			_, err = server.store.CreateAmenity(ctx, arg)
			if err != nil {
				log.Printf("Error at CreateAmenitiesAndSafety in CreateOptionAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

				err = server.store.RemoveAllAmenity(ctx, option.ID)
				if err != nil {
					log.Printf("Error at CreateAmenitiesAndSafety in CreateAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				}
				err = fmt.Errorf("error occurred while adding your amenities, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}
	// For safety am
	if len(req.HomeSafetyAm) > 0 {
		for i := 0; i < len(req.HomeSafetyAm); i++ {
			arg := db.CreateAmenityParams{
				OptionID:    option.ID,
				Tag:         req.HomeSafetyAm[i],
				AmType:      "home_safety",
				HasAm:       true,
				ListOptions: []string{},
			}
			_, err = server.store.CreateAmenity(ctx, arg)
			if err != nil {
				log.Printf("Error at CreateAmenitiesAndSafety in CreateAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

				err = server.store.RemoveAllAmenity(ctx, option.ID)
				if err != nil {
					log.Printf("Error at CreateAmenitiesAndSafety in CreateOptionAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				}
				err = fmt.Errorf("error occurred while adding your amenities, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}
	completeOption, err := HandleCompleteOption(utils.Description, utils.Amenities, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateAmenitiesAndSafety in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) CreateOptionDescription(ctx *gin.Context) {
	var req CreateOptionInfoDescription
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionDescription in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID: option.ID,
		Des: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error at CreateOptionDescription in CreateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("error occurred while creating your description, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	completeOption, err := HandleCompleteOption(utils.Name, utils.Description, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateOptionDescription in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateOptionName(ctx *gin.Context) {
	var req CreateOptionInfoName
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionName in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID: option.ID,
		HostNameOption: pgtype.Text{
			String: req.HostNameOption,
			Valid:  true,
		},
	})
	if err != nil {
		log.Printf("Error at CreateOptionName in CreateOptionInfoName: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("error occurred while creating your name, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	currentState, previousState := utils.NameViewState(option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateOptionName in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

//func (server *Server) CreateOptionPhoto(ctx *gin.Context) {
//	var req CreateOptionInfoPhotoParams
//	var photoExist bool = true
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		log.Printf("Error at CreateOptionPhoto in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while taking you back")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	log.Println(req)
//	requestID, err := tools.StringToUuid(req.OptionID)
//	if err != nil {
//		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
//		err = fmt.Errorf("error occurred while processing your request")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
//	if err != nil {
//		log.Println("error something went wrong 3, ", err)
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	_, err = server.store.GetOptionInfoPhoto(ctx, option.ID)
//	if err != nil {
//		log.Println("error something went wrong 1, ", err)
//		if err == db.ErrorRecordNotFound {
//			photoExist = false
//		} else {
//			log.Println("error something went wrong, ", err)
//			err = fmt.Errorf("something went wrong while uploading your photos, try again")
//			ctx.JSON(http.StatusBadRequest, errorResponse(err))
//			return
//		}
//	}
//	if photoExist {
//		log.Println("removing photo from firebase, ", err)
//		// we want to remove the photo from firebase
//		err = RemoveAllPhoto(server, ctx, option)
//		if err != nil {
//			log.Printf("Error at CreateOptionPhoto %v, \n", err.Error())
//		}
//		err = nil
//	}
//	log.Println("error something went wrong 5, ", err)
//	_, err = server.store.CreateOptionInfoPhoto(ctx, db.CreateOptionInfoPhotoParams{
//		OptionID:   option.ID,
//		CoverImage: req.CoverImage,
//		Photo:      req.Photo,
//	})
//	if err != nil {
//		log.Printf("Error at CreateOptionPhoto in CreateOptionInfoPhoto: %v, optionID: %v, userID: %v \n for cover image: %v photos: %v\n", err.Error(), requestID, user.ID, req.CoverImage, req.Photo)

//		err = fmt.Errorf("your photos were not saved, please try uploading them again")
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//	currentState, previousState := utils.PhotoViewState(option.OptionType)
//	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
//	if err != nil {
//		log.Printf("Error at CreateOptionPhoto in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
//	}
//	res := OptionInfoResponse{
//		OptionItemType:     "",
//		Success:            true,
//		OptionID:           tools.UuidToString(option.ID),
//		UserOptionID:       tools.UuidToString(option.OptionUserID),
//		CurrentServerView:  completeOption.CurrentState,
//		PreviousServerView: completeOption.PreviousState,
//		MainOptionType:     option.MainOptionType,
//		OptionType:         option.OptionType,
//		Currency:           option.Currency,
//	}
//	ctx.JSON(http.StatusOK, res)

//}

func (server *Server) CreateOptionPhoto(ctx *gin.Context) {
	var req CreateOptionInfoPhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionPhoto in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		log.Println("error something went wrong 3, ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionData, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		// We expect that it should already exist
		log.Println("error something went wrong, ", err)
		err = fmt.Errorf("you must have a cover image with additional photos")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if tools.ServerStringEmpty(optionData.CoverImage) || tools.ServerListIsEmpty(optionData.Photo) {
		err = fmt.Errorf("you must have a cover image with additional photos")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	currentState, previousState := utils.PhotoViewState(option.OptionType)
	completeOption, err := HandleCompleteOption(currentState, previousState, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateOptionPhoto in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) UploadOptionPhoto(ctx *gin.Context) {
	var req UploadOptionInfoPhotoParams
	var optionPhotoExist bool = true
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UploadOptionPhoto in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		log.Println("error something went wrong 3, ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Println("error something went wrong 1, ", err)
		if err == db.ErrorRecordNotFound {
			optionPhotoExist = false
		} else {
			log.Println("error something went wrong, ", err)
			err = fmt.Errorf("something went wrong while uploading your photos, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	if optionPhotoExist {
		// If option photo exist we want to check if that particular photo exists
		if req.IsCover {
			_, err = server.store.UpdateOptionInfoAllPhotoCover(ctx, db.UpdateOptionInfoAllPhotoCoverParams{
				OptionID:         option.ID,
				CoverImage:       req.Photo,
				PublicCoverImage: req.PhotoUrl,
			})
			if err != nil {
				log.Println("error something went wrong, ", err)
				err = fmt.Errorf("something went wrong while uploading your photos, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		} else {
			photos := append(optionPhoto.Photo, req.Photo)
			photos = tools.HandleListReq(tools.RemoveDuplicates(photos))
			photoUrls := append(optionPhoto.PublicPhoto, req.PhotoUrl)
			photoUrls = tools.HandleListReq(tools.RemoveDuplicates(photoUrls))
			_, err = server.store.UpdateOptionInfoAllPhotoOnly(ctx, db.UpdateOptionInfoAllPhotoOnlyParams{
				OptionID:    option.ID,
				Photo:       photos,
				PublicPhoto: photoUrls,
			})
			if err != nil {
				log.Println("error something went wrong, ", err)
				err = fmt.Errorf("something went wrong while uploading your photos, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	} else {
		// We want to create the photo object
		if req.IsCover {
			_, err = server.store.CreateOptionInfoPhoto(ctx, db.CreateOptionInfoPhotoParams{
				OptionID:         option.ID,
				CoverImage:       req.Photo,
				PublicCoverImage: req.PhotoUrl,
				Photo:            []string{"none"},
				PublicPhoto:      []string{"none"},
			})
			if err != nil {
				log.Println("error something went wrong, ", err)
				err = fmt.Errorf("something went wrong while uploading your photos, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		} else {
			_, err = server.store.CreateOptionInfoPhoto(ctx, db.CreateOptionInfoPhotoParams{
				OptionID:         option.ID,
				CoverImage:       "none",
				PublicCoverImage: "none",
				Photo:            []string{req.Photo},
				PublicPhoto:      []string{req.PhotoUrl},
			})
			if err != nil {
				log.Println("error something went wrong, ", err)
				err = fmt.Errorf("something went wrong while uploading your photos, try again")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}
	res := UploadOptionInfoPhotoRes{
		Photo:    req.Photo,
		PhotoUrl: req.PhotoUrl,
		IsCover:  req.IsCover,
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) DeleteOptionPhoto(ctx *gin.Context) {
	var req db.DeleteOptionInfoPhotoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at DeleteOptionPhoto in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while taking you back")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		log.Println("error something went wrong 3, ", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Println("error something went wrong, ", err)
		err = fmt.Errorf("something went wrong while uploading your photos, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res, err := server.store.DeleteOptionPhoto(ctx, req, optionPhoto, server.Bucket)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) GetPublishData(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while adding your Question")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	switch option.MainOptionType {
	case "events":
		completeOption, err := server.store.GetCompleteOptionInfoTwo(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at GetPublishData at GetCompleteOptionInfoTwo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("there was an error while getting your data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res, err := HandleEventViewToPublish(server, ctx, option, user, completeOption)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, res)
	case "options":
		completeOption, err := server.store.GetCompleteOptionInfoTwo(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at GetPublishData at GetCompleteOptionInfoTwo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("there was an error while getting your data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res, err := HandleShortletViewToPublish(server, ctx, option, user, completeOption)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, res)
	}

}
func (server *Server) CreateOptionQuestion(ctx *gin.Context) {
	var req CreateOptionQuestionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionQuestion in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while adding your Question")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.RemoveOptionQuestion(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Option questions doesn't exist %v \n", err.Error())
		} else {
			log.Println(err.Error())
		}
	}
	lng := tools.ConvertLocationStringToFloat("0", 9)
	lat := tools.ConvertLocationStringToFloat("0", 9)
	geolocation := pgtype.Point{
		P:     pgtype.Vec2{X: lng, Y: lat},
		Valid: true,
	}
	_, err = server.store.CreateOptionQuestion(ctx, db.CreateOptionQuestionParams{
		OptionID:         option.ID,
		HostAsIndividual: req.HostAsIndividual,
		OrganizationName: req.OrganizationName,
		LegalRepresents:  []string{"none"},
		Geolocation:      geolocation,
	})
	if err != nil {
		log.Printf("There an error at CreateOptionQuestion at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

		err = fmt.Errorf("there was an error while updating your question selection please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Handle Things To Notes
	for _, note := range optionQuestionNote {
		noteData, err := server.store.GetThingToNoteByType(ctx, db.GetThingToNoteByTypeParams{
			Tag:      note.Tag,
			OptionID: option.ID,
			Type:     note.Type,
		})
		if err != nil {
			// We create
			var checked bool
			switch note.Tag {
			case "cameras_audio_devices":
				checked = req.HasSecurityCamera
			case "dangerous_animal":
				checked = req.HasDangerousAnimals
			case "weapon_on_property":
				checked = req.HasWeapons
			}
			_, err := server.store.CreateThingToNote(ctx, db.CreateThingToNoteParams{
				OptionID: option.ID,
				Tag:      note.Tag,
				Type:     note.Type,
				Checked:  checked,
			})
			if err != nil {
				log.Printf("There an error at CreateOptionQuestion at CreateThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			}
		} else {
			// We update
			var checked bool
			switch note.Tag {
			case "cameras_audio_devices":
				checked = req.HasSecurityCamera
			case "dangerous_animal":
				checked = req.HasDangerousAnimals
			case "weapon_on_property":
				checked = req.HasWeapons
			}
			_, err := server.store.UpdateThingToNote(ctx, db.UpdateThingToNoteParams{
				ID:      noteData.ID,
				Checked: checked,
			})
			if err != nil {
				log.Printf("There an error at CreateOptionQuestion at UpdateThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			}
		}
	}
	completeOption, err := HandleCompleteOption(utils.Publish, utils.HostQuestion, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at CreateOptionQuestion at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateOptionPrice(ctx *gin.Context) {
	var req CreateOptionPrice
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionPrice in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while adding your price")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.store.RemoveOptionPrice(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Option price doesn't exist %v \n", err.Error())
		} else {
			log.Println(err.Error())
		}
	}
	//priceDB, err := MoneyToDB(option.Currency, req.Price, server)
	//if err != nil {
	//	log.Printf("Error at CreateOptionPrice in MoneyToDB: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	//}
	_, err = server.store.CreateOptionPrice(ctx, db.CreateOptionPriceParams{
		OptionID:     option.ID,
		Price:        tools.MoneyStringToInt(req.Price),
		WeekendPrice: 0,
	})
	if err != nil {
		log.Printf("There an error at CreateOptionPrice at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

		err = fmt.Errorf("there was an error while processing the price you entered. Please ensure it follows the currency's standard")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(utils.Highlight, utils.Price, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("There an error at CreateOptionPrice at HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateOptionHighlight(ctx *gin.Context) {
	var req CreateOptionInfoHighlight
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateOptionHighlight in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request. Please ensure you entered the right details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, option, err := HandleGetOptionIncomplete(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	_, err = server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		OptionID:        option.ID,
		OptionHighlight: req.Highlight,
	})
	if err != nil {
		log.Printf("Error at CreateOptionHighlight in UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), requestID, user.ID)
		err = fmt.Errorf("error occurred while taking you back, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	completeOption, err := HandleCompleteOption(utils.Photo, utils.Highlight, server, ctx, option, user.ID)
	if err != nil {
		log.Printf("Error at CreateOptionHighlight in HandleCompleteOption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	res := OptionInfoResponse{
		OptionItemType:     "",
		Success:            true,
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	ctx.JSON(http.StatusOK, res)
}
