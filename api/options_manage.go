package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// This is for shortlets
func (server *Server) GetUHMOptionData(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetUHMOptionData in GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
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

	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	switch option.MainOptionType {
	case "options":
		uhmData, err := server.store.GetOptionShortletUHMData(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at GetUHMOptionData at GetOptionInfoUHMData: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not get the data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		spaceAreas, err := server.store.ListSpaceAreaType(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at GetUHMOptionData at ListSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not your space areas")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res := GetUHMDataOptionRes{
			HostNameOption: uhmData.HostNameOption,
			SpaceAreas:     spaceAreas,
			Price:          tools.IntToMoneyString(uhmData.Price),
			SpaceType:      uhmData.SpaceType,
			Category:       uhmData.Category,
			CategoryTwo:    uhmData.CategoryTwo,
			CategoryThree:  uhmData.CategoryThree,
			CategoryFour:   uhmData.CategoryFour,
			NumOfGuest:     int(uhmData.GuestWelcomed),
			CoverPhoto:     uhmData.CoverImage,
			Photos:         uhmData.Photo,
			OptionUserID:   tools.UuidToString(uhmData.OptionUserID),
			Street:         uhmData.Street,
			State:          uhmData.State,
			City:           uhmData.City,
			Country:        uhmData.Country,
			Postcode:       uhmData.Postcode,
			CheckInMethod:  uhmData.CheckInMethod,
			EventType:      "",
			EventSubType:   "",
			Currency:       uhmData.Currency,
			Status:         tools.HandleOptionStatus(uhmData.OptionStatus),
		}
		ctx.JSON(http.StatusOK, res)
	case "events":
		uhmData, err := server.store.GetOptionEventUHMData(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at GetUHMOptionData at GetOptionInfoUHMData: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not get the data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		res := GetUHMDataOptionRes{
			HostNameOption: uhmData.HostNameOption,
			SpaceAreas:     []string{"none"},
			Price:          "",
			SpaceType:      "",
			Category:       uhmData.Category,
			CategoryTwo:    uhmData.CategoryTwo,
			CategoryThree:  uhmData.CategoryThree,
			CategoryFour:   uhmData.CategoryFour,
			NumOfGuest:     0,
			CoverPhoto:     uhmData.CoverImage,
			Photos:         uhmData.Photo,
			OptionUserID:   tools.UuidToString(uhmData.OptionUserID),
			Street:         "",
			State:          "",
			City:           "",
			Country:        "",
			Postcode:       "",
			CheckInMethod:  "",
			EventType:      uhmData.EventType,
			EventSubType:   uhmData.SubCategoryType,
			Currency:       uhmData.Currency,
			Status:         tools.HandleOptionStatus(uhmData.OptionStatus),
		}
		ctx.JSON(http.StatusOK, res)

	}

}

func (server *Server) ListSpaceArea(ctx *gin.Context) {
	var req ListSpaceAreasParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListSpaceAreasParams in ShouldBindJSON: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("please make sure you select at least a start date")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionID: %v \n", err.Error(), req.OptionID)
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	res, err, contentFound := HandleListSpaceAreas(ctx, server, option, user.ID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	} else if !contentFound {
		data := "none"
		ctx.JSON(http.StatusNoContent, data)
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateEditSpaceAreas(ctx *gin.Context) {
	var req CreateEditSpaceAreaParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateEditSpaceAreaParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("try selecting at least one of the guests options made available")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(req.RemoveSpaceID) > 0 {
		for i := 0; i < len(req.RemoveSpaceID); i++ {
			removeSpaceID, err := tools.StringToUuid(req.RemoveSpaceID[i])
			if err != nil {
				log.Printf("Error at tools.StringToUuid for req.RemoveSpaceID: %v, optionID: %v \n", err.Error(), req.OptionID)
				err = fmt.Errorf("error occurred while processing your request")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
			err = server.store.RemoveSpaceArea(ctx, db.RemoveSpaceAreaParams{
				OptionID: option.ID,
				ID:       removeSpaceID,
			})
			if err != nil {
				log.Printf("There an error at CreateEditSpaceAreas at RemoveSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				err = fmt.Errorf("error occurred while processing your request")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}
	if len(req.Space) > 0 {
		for i := 0; i < len(req.Space); i++ {
			_, err := server.store.CreateSpaceArea(ctx, db.CreateSpaceAreaParams{
				OptionID:    option.ID,
				SharedSpace: false,
				SpaceType:   req.Space[i],
				Photos:      []string{"none"},
				Beds:        []string{"none"},
			})
			if err != nil {
				log.Printf("There an error at CreateEditSpaceAreas at CreateSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				err = fmt.Errorf("could not create your guest areas")
				ctx.JSON(http.StatusBadRequest, errorResponse(err))
				return
			}
		}
	}
	res, err, contentFound := HandleListSpaceAreas(ctx, server, option, user.ID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	} else if !contentFound {
		data := "none"
		ctx.JSON(http.StatusNoContent, data)
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateEditSpaceAreas", "listing space area", "create listing space area")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) AddBedSpaceAreas(ctx *gin.Context) {
	var req AddBedSpaceAreasParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at AddBedSpaceAreasParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("select the beds you to associate this room with")
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
	reqBeds := tools.HandleListReq(req.Beds)
	spaceAreaID, err := tools.StringToUuid(req.SpaceAreaID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	beds := []string{"none"}
	if !tools.ServerListIsEmpty(reqBeds) {
		beds = reqBeds
	}
	spaceArea, err := server.store.UpdateSpaceAreaBeds(ctx, db.UpdateSpaceAreaBedsParams{
		ID:   spaceAreaID,
		Beds: beds,
	})
	if err != nil {
		log.Printf("There an error at AddBedSpaceAreas at UpdateSpaceAreaBeds: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not add this bed type to this bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	photoRes := tools.HandleDBList(spaceArea.Photos)
	bedRes := tools.HandleDBList(spaceArea.Beds)
	res := SpaceAreas{
		ID:          tools.UuidToString(spaceArea.ID),
		OptionID:    tools.UuidToString(spaceArea.OptionID),
		SharedSpace: spaceArea.SharedSpace,
		SpaceType:   spaceArea.SpaceType,
		Photos:      photoRes,
		Beds:        bedRes,
		IsSuite:     spaceArea.IsSuite,
		Name:        req.Name,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "AddBedSpaceAreas", "listing bed space areas", "editing listing bed space areas")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) AddPhotoSpaceAreas(ctx *gin.Context) {
	var req AddPhotoSpaceAreasParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at AddPhotoSpaceAreasParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("select the photos you to associate this room with")
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
	reqPhotos := tools.HandleListReq(req.Photos)
	fmt.Println("reqPhotos", reqPhotos)
	spaceAreaID, err := tools.StringToUuid(req.SpaceAreaID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var photos []string
	var hasPhotos bool
	if len(reqPhotos) == 0 {
		hasPhotos = true
		photos = []string{"none"}
	} else {
		photos, hasPhotos = HandleGetUnselectedPhotos(ctx, server, option, user.ID, req, reqPhotos, spaceAreaID)
	}

	if hasPhotos {
		spaceArea, err := server.store.UpdateSpaceAreaPhotos(ctx, db.UpdateSpaceAreaPhotosParams{
			ID:     spaceAreaID,
			Photos: photos,
		})
		if err != nil {
			log.Printf("There an error at AddPhotoSpaceAreas at UpdateSpaceAreaPhotos: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not add this photo to this guest area")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		photoRes := tools.HandleDBList(spaceArea.Photos)
		bedRes := tools.HandleDBList(spaceArea.Beds)
		res := SpaceAreas{
			ID:          tools.UuidToString(spaceArea.ID),
			OptionID:    tools.UuidToString(spaceArea.OptionID),
			SharedSpace: spaceArea.SharedSpace,
			SpaceType:   spaceArea.SpaceType,
			Photos:      photoRes,
			Beds:        bedRes,
			IsSuite:     spaceArea.IsSuite,
			Name:        req.Name,
		}
		if isCoHost {
			HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "AddPhotoSpaceAreas", "listing add photo to space areas", "edit listing add photo to space areas")
		}
		ctx.JSON(http.StatusOK, res)
		return
	} else {
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

}

func (server *Server) UpdateSpaceAreas(ctx *gin.Context) {
	var req UpdateSpaceAreaParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateSpaceAreaParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("your selected options do not fit the recommended options")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, false)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Update the sharedSpaceWith
	var sharedSpaceWith []string
	if len(req.SharedSpaceWith) < 1 {
		sharedSpaceWith = []string{"none"}
	} else {
		sharedSpaceWith = req.SharedSpaceWith
	}

	_, err = server.store.UpdateShortletInfoSharedWith(ctx, db.UpdateShortletInfoSharedWithParams{
		SharedSpacesWith: sharedSpaceWith,
		OptionID:         option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateSpaceAreas in UpdateShortletInfoSharedWith: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not update your shared space with")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	for i := 0; i < len(req.SpaceAreas); i++ {
		id, err := tools.StringToUuid(req.SpaceAreas[i].ID)
		if err != nil {
			log.Printf("There an error at UpdateSpaceAreas in req.SpaceArea loop at tools.UuidToString: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not access this guest area")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		spaceAreaDB, err := server.store.GetSpaceArea(ctx, db.GetSpaceAreaParams{
			ID:       id,
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at UpdateSpaceAreas in req.SpaceArea loop at GetSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not access this guest area")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// we check to see if it needs update
		if req.SpaceAreas[i].SpaceType == "bedroom" {
			// we check to see if it is a bedroom because of isSuite is only for bedroom
			if req.SpaceAreas[i].IsSuite != spaceAreaDB.IsSuite || req.SpaceAreas[i].SharedSpace != spaceAreaDB.SharedSpace {
				_, err = server.store.UpdateSpaceAreaInfo(ctx, db.UpdateSpaceAreaInfoParams{
					IsSuite: pgtype.Bool{
						Bool:  req.SpaceAreas[i].IsSuite,
						Valid: true,
					},
					SharedSpace: pgtype.Bool{
						Bool:  req.SpaceAreas[i].SharedSpace,
						Valid: true,
					},
					ID:       spaceAreaDB.ID,
					OptionID: option.ID,
				})
				if err != nil {
					log.Printf("There an error at UpdateSpaceAreas in req.SpaceArea loop at for bedroom UpdateSpaceAreaInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
					err = fmt.Errorf("could not access this guest area")
					ctx.JSON(http.StatusBadRequest, errorResponse(err))
					return
				}

			}
		} else {
			// The rest would just check the shared space stuff
			if req.SpaceAreas[i].SharedSpace != spaceAreaDB.SharedSpace {
				_, err = server.store.UpdateSpaceAreaInfo(ctx, db.UpdateSpaceAreaInfoParams{
					SharedSpace: pgtype.Bool{
						Bool:  req.SpaceAreas[i].SharedSpace,
						Valid: true,
					},
					ID:       spaceAreaDB.ID,
					OptionID: option.ID,
				})
				if err != nil {
					log.Printf("There an error at UpdateSpaceAreas in req.SpaceArea loop at for not bedroom UpdateSpaceAreaInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
					err = fmt.Errorf("could not access this guest area")
					ctx.JSON(http.StatusBadRequest, errorResponse(err))
					return
				}
			}
		}
	}

	res, err, contentFound := HandleListSpaceAreas(ctx, server, option, user.ID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	} else if !contentFound {
		err = fmt.Errorf("could not get your guest space areas")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateSpaceArea", "listing space area", "update listing space area")
	}
	ctx.JSON(http.StatusOK, res)
}

// This updates the eventInfo
func (server *Server) UpdateEventInfo(ctx *gin.Context) {
	var req UpdateEventInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateEventInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	eventInfo, err := server.store.UpdateEventInfo(ctx, db.UpdateEventInfoParams{
		EventType: pgtype.Text{
			String: req.EventType,
			Valid:  true,
		},
		SubCategoryType: pgtype.Text{
			String: req.SubCategoryType,
			Valid:  true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateEventInfo at UpdateEventInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Because there is a change in the event type and event sub type
	CreateEventAlgo(ctx, server, option, user)

	res := EventInfoRes{
		SubCategoryType: eventInfo.SubCategoryType,
		EventType:       eventInfo.EventType,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateEventInfo", "event info", "update event information")
	}
	ctx.JSON(http.StatusOK, res)
}

// This gets the event info of the event
func (server *Server) GetEventInfo(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in GetEventInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("GetShortletInfo Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	eventInfo, err := server.store.GetEventInfo(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetEventInfo at GetEventInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := EventInfoRes{
		SubCategoryType: eventInfo.SubCategoryType,
		EventType:       eventInfo.EventType,
	}
	ctx.JSON(http.StatusOK, res)
}

// This updates the property info of the shortlet
func (server *Server) UpdateShortletInfo(ctx *gin.Context) {
	var req UpdateShortletInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateShortletInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	yearBuilt := tools.ValidateIntLessThanZero(req.YearBuilt)
	propertySize := tools.ValidateIntLessThanZero(req.PropertySize)
	shortlet, err := server.store.UpdateShortletInfo(ctx, db.UpdateShortletInfoParams{
		SpaceType: pgtype.Text{
			String: req.SpaceType,
			Valid:  true,
		},
		TypeOfShortlet: pgtype.Text{
			String: req.TypeOfShortlet,
			Valid:  true,
		},
		GuestWelcomed: pgtype.Int4{
			Int32: int32(req.GuestWelcomed),
			Valid: true,
		},
		YearBuilt: pgtype.Int4{
			Int32: int32(yearBuilt),
			Valid: true,
		},
		PropertySize: pgtype.Int4{
			Int32: int32(propertySize),
			Valid: true,
		},
		PropertySizeUnit: pgtype.Text{
			String: req.PropertySizeUnit,
			Valid:  true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateShortletInfo at UpdateShortletInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ShortletInfoRes{
		SpaceType:        shortlet.SpaceType,
		TypeOfShortlet:   shortlet.TypeOfShortlet,
		GuestWelcomed:    int(shortlet.GuestWelcomed),
		YearBuilt:        int(shortlet.YearBuilt),
		PropertySize:     int(shortlet.PropertySize),
		PropertySizeUnit: shortlet.PropertySizeUnit,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateShortletInfo", "listing information", "update listing information")
	}
	ctx.JSON(http.StatusOK, res)
}

// This gets the property info of the shortlet
func (server *Server) GetShortletInfo(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetShortletInfo in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("you must be able to welcome at least one guest and have a bedroom")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("GetShortletInfo Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetShortletInfo at GetShortlet: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ShortletInfoRes{
		SpaceType:        shortlet.SpaceType,
		TypeOfShortlet:   shortlet.TypeOfShortlet,
		GuestWelcomed:    int(shortlet.GuestWelcomed),
		YearBuilt:        int(shortlet.YearBuilt),
		PropertySize:     int(shortlet.PropertySize),
		PropertySizeUnit: shortlet.PropertySizeUnit,
	}
	ctx.JSON(http.StatusOK, res)
}

// This updates the title of the shortlet or event
func (server *Server) UpdateOptionTitle(ctx *gin.Context) {
	var req UpdateOptionTitleParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionTitleParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
		HostNameOption: pgtype.Text{
			String: req.HostNameOption,
			Valid:  true,
		},
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateOptionTitle at UpdateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateOptionTitleParams{
		OptionID:       tools.UuidToString(option.ID),
		HostNameOption: optionInfoDetail.HostNameOption,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionTitle", "listing name (title)", "update listing name")
	}
	ctx.JSON(http.StatusOK, res)
}

// This updates the optionInfoDetail for both event and option
// Note because event and option share some fields
// space_des <-> schedule
// guest_access_des <-> sponsor
// interact_with_guest_des <-> featured_guests
func (server *Server) UpdateOptionDes(ctx *gin.Context) {
	var req UpdateOptionDesParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionDesParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
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

	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	res, err := HandleUpdateDes(ctx, server, option, user.ID, req)
	if err != nil {
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionDesParams", "listing description", "update listing description")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionDes(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetDes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
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

	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	optionInfoDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionDes at UpdateOptionInfoDetail for OtherDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	res := GetOptionDesRes{
		SpaceDes:              optionInfoDetail.SpaceDes,
		GuestAccessDes:        optionInfoDetail.GuestAccessDes,
		InteractWithGuestsDes: optionInfoDetail.InteractWithGuestsDes,
		OtherDes:              optionInfoDetail.OtherDes,
		NeighborhoodDes:       optionInfoDetail.NeighborhoodDes,
		GetAroundDes:          optionInfoDetail.GetAroundDes,
		Des:                   optionInfoDetail.Des,
		OptionID:              tools.UuidToString(option.ID),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateOptionPhoto(ctx *gin.Context) {
	var req UpdateOptionPhotoParams
	photoTypes := []string{"change_cover", "create_photo", "delete_photo"}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionPhoto in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("a title is required")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	photoType := ctx.Param("photo_type")
	if !tools.ContainsString(photoTypes, photoType) {
		err := fmt.Errorf("path does not exist")
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}
	// photoType can be
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var res UpdateOptionPhotoRes
	switch photoType {
	case "change_cover":
		res, err = HandleOptionPhotoChangeCover(ctx, server, option, user, req)
	case "create_photo":
		res, err = HandleAddOptionPhotos(ctx, server, option, user, req)
	case "delete_photo":
		res, err = HandleDeleteOptionPhoto(ctx, server, option, user, req)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionPhoto", "listing photos", "update listing photos")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateOptionPhotoCaption(ctx *gin.Context) {
	var req CreateUpdateOptionPhotoCaptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateOptionPhotoCaptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	// We would try to create it however if it already exists then we update
	_, err = server.store.GetOptionPhotoCaption(ctx, db.GetOptionPhotoCaptionParams{
		OptionID: option.ID,
		PhotoID:  req.PhotoID,
	})
	var res CreateUpdateOptionPhotoCaptionRes
	if err != nil {
		photoCaption, err := server.store.CreateOptionPhotoCaption(ctx, db.CreateOptionPhotoCaptionParams{
			OptionID: option.ID,
			PhotoID:  req.PhotoID,
			Caption:  req.Caption,
		})
		if err != nil {
			log.Printf("There an error at CreateUpdateOptionPhotoCaption atCreateOptionPhotoCaption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not create your caption")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CreateUpdateOptionPhotoCaptionRes{
				PhotoID: photoCaption.PhotoID,
				Caption: photoCaption.Caption,
			}
		}
	} else {
		// we want to update the stuff instead
		photoCaption, err := server.store.UpdateOptionPhotoCaption(ctx, db.UpdateOptionPhotoCaptionParams{
			OptionID: option.ID,
			PhotoID:  req.PhotoID,
			Caption:  req.Caption,
		})
		if err != nil {
			log.Printf("There an error at CreateUpdateOptionPhotoCaption at UpdateOptionPhotoCaption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your caption")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CreateUpdateOptionPhotoCaptionRes{
				PhotoID: photoCaption.PhotoID,
				Caption: photoCaption.Caption,
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateOptionPhotoCaption", "listing photos", "create-update listing photos captions")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionPhotoCaption(ctx *gin.Context) {
	var req GetOptionPhotoCaptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionPhotoCaptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var location string = "none"
	// we want to check if it is match to any space area
	spaceAreas, err := server.store.ListOrderedSpaceArea(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at CreateUpdateOptionPhotoCaption at ListOrderedSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		location = "none"
	} else {
		var spaceData = make(map[string]int)
		for i := 0; i < len(spaceAreas); i++ {
			spaceData[spaceAreas[i].SpaceType] = spaceData[spaceAreas[i].SpaceType] + 1
			if tools.ContainsString(spaceAreas[i].Photos, req.PhotoID) {
				name := fmt.Sprintf("%v-%d", spaceAreas[i].SpaceType, spaceData[spaceAreas[i].SpaceType])
				location = name
				break
			}
		}
	}
	var photoID string
	photoCaption, err := server.store.GetOptionPhotoCaption(ctx, db.GetOptionPhotoCaptionParams{
		OptionID: option.ID,
		PhotoID:  req.PhotoID,
	})
	if err != nil {
		log.Printf("There an error at CreateUpdateOptionPhotoCaption at UpdateOptionPhotoCaption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		photoID = req.PhotoID
	} else {
		photoID = photoCaption.PhotoID
	}
	res := GetOptionPhotoCaptionRes{
		PhotoID:       photoID,
		Caption:       photoCaption.Caption,
		SpaceLocation: location,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionPhotoCaption", "listing photo captions", "create-update listing photo captions")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetUHMHighlight(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionHighlight, err := server.store.GetOptionInfoDetailHighlight(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  GetUHMHighlight in GetOptionInfoDetailHighlight err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while getting your shortlet location")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(optionHighlight) == 0 {
		optionHighlight = []string{"none"}
	}
	res := GetOptionDetailHighlightRes{
		Highlight: optionHighlight,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateUHMHighlight(ctx *gin.Context) {
	var req UpdateOptionDetailHighlightParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionDetailHighlightParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var highlight []string = tools.HandleListReq(req.Highlight)
	if len(highlight) == 0 {
		highlight = []string{"none"}
	}
	optionHighlight, err := server.store.UpdateOptionInfoDetailHighlight(ctx, db.UpdateOptionInfoDetailHighlightParams{
		OptionHighlight: highlight,
		OptionID:        option.ID,
	})
	if err != nil {
		log.Printf("Error at  UpdateUHMHighlight in UpdateOptionInfoDetailHighlight err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while getting your shortlet location")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(optionHighlight) == 0 {
		optionHighlight = []string{"none"}
	}
	res := GetOptionDetailHighlightRes{
		Highlight: optionHighlight,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateUHMHighlight", "flexr tags", "update flexr tags")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionWifiDetail(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, _, _, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	wifiDetail, err := server.store.GetWifiDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  GetOptionWifiDetail in GetWifiDetail err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	res := GetWifiDetailRes{
		NetworkName: wifiDetail.NetworkName,
		Password:    wifiDetail.Password,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateWifiDetail(ctx *gin.Context) {
	var req CreateUpdateWifiDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateWifiDetailParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	exist := true
	_, err = server.store.GetWifiDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  CreateUpdateWifiDetail in GetWifiDetail err: %v, user: %v\n", err, user.ID)
		exist = false
	}
	var res GetWifiDetailRes
	if exist {
		// we want to update the details
		wifiDetail, err := server.store.UpdateWifiDetail(ctx, db.UpdateWifiDetailParams{
			NetworkName: req.NetworkName,
			Password:    req.Password,
			OptionID:    option.ID,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateWifiDetail in UpdateWifiDetail err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not update your wifi details")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = GetWifiDetailRes{
				NetworkName: wifiDetail.NetworkName,
				Password:    wifiDetail.Password,
			}
		}
	} else {
		// we want to create it cause it doesn't exist
		wifiDetail, err := server.store.CreateWifiDetail(ctx, db.CreateWifiDetailParams{
			NetworkName: req.NetworkName,
			Password:    req.Password,
			OptionID:    option.ID,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateWifiDetail in CreateWifiDetail err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not set up your wifi details")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = GetWifiDetailRes{
				NetworkName: wifiDetail.NetworkName,
				Password:    wifiDetail.Password,
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateWifiDetail", "wifi detail", "create-update wifi details")
	}
	ctx.JSON(http.StatusOK, res)
}
