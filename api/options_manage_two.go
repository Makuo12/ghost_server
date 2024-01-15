package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/val"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (server *Server) ListUHMAmenities(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
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
	amenities, err := server.store.ListAmenitiesOne(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at ListUHMAmenities at ListAmenitiesOne: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := "none"
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	var res ListUHMAmenitiesRes
	var resData []AmenityItem
	for _, amenity := range amenities {
		data := AmenityItem{
			Tag:    amenity.Tag,
			AmType: amenity.AmType,
			HasAm:  amenity.HasAm,
			ID:     tools.UuidToString(amenity.ID),
		}
		resData = append(resData, data)
	}

	res = ListUHMAmenitiesRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateAmenity(ctx *gin.Context) {
	var req CreateUpdateAmenityParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateAmenityParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
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
	// We want to confirm the tag
	if !val.ValidateAmTag(req.AmType, req.Tag) {
		err = fmt.Errorf("this amenity does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We would try to create it however if it already exists then we update
	amenity, err := server.store.GetAmenityByType(ctx, db.GetAmenityByTypeParams{
		OptionID: option.ID,
		AmType:   req.AmType,
		Tag:      req.Tag,
	})
	var exists bool
	if err != nil {
		log.Printf("There an error at CreateUpdateAmenity at GetAmenityByType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		if err == db.ErrorRecordNotFound {
			exists = false
		} else {
			//err = fmt.Errorf("could not add this amenity to your option")
			//ctx.JSON(http.StatusBadRequest, errorResponse(err))
			//return
			exists = false
		}
	} else {
		exists = true
	}
	var res AmenityItem
	if exists {
		// We want to update the amenity
		update, err := server.store.UpdateAmenity(ctx, db.UpdateAmenityParams{
			HasAm: req.HasAm,
			ID:    amenity.ID,
		})
		if err != nil {
			log.Printf("There an error at CreateUpdateAmenity at UpdateAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update this amenity")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = AmenityItem{
				Tag:    update.Tag,
				AmType: update.AmType,
				HasAm:  update.HasAm,
				ID:     tools.UuidToString(update.ID),
			}
		}
	} else {
		// We want to create
		create, err := server.store.CreateAmenity(ctx, db.CreateAmenityParams{
			OptionID:    option.ID,
			Tag:         req.Tag,
			AmType:      req.AmType,
			HasAm:       true,
			ListOptions: []string{},
		})
		if err != nil {
			log.Printf("There an error at CreateUpdateAmenity at CreateAmenity: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not add this amenity to your option")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = AmenityItem{
				Tag:    create.Tag,
				AmType: create.AmType,
				HasAm:  create.HasAm,
				ID:     tools.UuidToString(create.ID),
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateAmenity", "listing amenities", "create-update listing amenities")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetAmenityDetail(ctx *gin.Context) {
	var req GetAmenityDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetAmenityDetailParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
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
	amDetail, err := server.store.GetAmenityDetail(ctx, db.GetAmenityDetailParams{
		OptionID: option.ID,
		AmType:   req.AmType,
		Tag:      req.Tag,
	})
	if err != nil {
		log.Printf("There an error at GetAmenityDetail at GetAmenityDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to get details for this amenity")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startTime := tools.ConvertTimeOnlyToString(amDetail.StartTime)
	endTime := tools.ConvertTimeOnlyToString(amDetail.EndTime)
	res := UHMAmenityDetailRes{
		ID:                 tools.UuidToString(amDetail.ID),
		LocationOption:     amDetail.LocationOption,
		SizeOption:         int(amDetail.SizeOption),
		PrivacyOption:      amDetail.PrivacyOption,
		TimeOption:         amDetail.TimeOption,
		StartTime:          startTime,
		TimeSet:            amDetail.TimeSet,
		EndTime:            endTime,
		AvailabilityOption: amDetail.AvailabilityOption,
		StartMonth:         amDetail.StartMonth,
		EndMonth:           amDetail.EndMonth,
		TypeOption:         amDetail.TypeOption,
		PriceOption:        amDetail.PriceOption,
		BrandOption:        amDetail.BrandOption,
		ListOptions:        amDetail.ListOptions,
		OptionID:           tools.UuidToString(option.ID),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateAmenityDetail(ctx *gin.Context) {
	var req UHMAmenityDetailRes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UHMAmenityDetailRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, ID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var timeSet bool = false
	var timeOne string
	var timeTwo string
	if req.TimeOption == "Open specific hours" {
		timeSet = true
		timeOne = req.StartTime
		timeTwo = req.EndTime
	} else {
		timeOne = "00:00"
		timeTwo = "00:00"
	}
	if req.AvailabilityOption == "Available seasonally" {
		if !val.ValidateMonth(req.StartMonth) || !val.ValidateMonth(req.EndMonth) {
			err = fmt.Errorf("your start month or end month is not set properly")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	startTime, err := tools.ConvertStringToTimeOnly(timeOne)
	if err != nil {
		err = fmt.Errorf("your start time is set wrongly")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	endTime, err := tools.ConvertStringToTimeOnly(timeTwo)
	if err != nil {
		err = fmt.Errorf("your end time is set wrongly")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	amDetail, err := server.store.UpdateAmenityDetail(ctx, db.UpdateAmenityDetailParams{
		LocationOption:     req.LocationOption,
		SizeOption:         int32(req.SizeOption),
		PrivacyOption:      req.PrivacyOption,
		TimeSet:            timeSet,
		TimeOption:         req.TimeOption,
		StartTime:          startTime,
		EndTime:            endTime,
		AvailabilityOption: req.AvailabilityOption,
		StartMonth:         req.StartMonth,
		EndMonth:           req.EndMonth,
		TypeOption:         req.TypeOption,
		PriceOption:        req.PriceOption,
		ListOptions:        req.ListOptions,
		ID:                 id,
		OptionID:           option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateAmenityDetail at UpdateAmenityDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update details for this amenity")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	startTimeString := tools.ConvertTimeOnlyToString(amDetail.StartTime)
	endTimeString := tools.ConvertTimeOnlyToString(amDetail.EndTime)
	res := UHMAmenityDetailRes{
		ID:                 tools.UuidToString(amDetail.ID),
		OptionID:           tools.UuidToString(option.ID),
		LocationOption:     amDetail.LocationOption,
		SizeOption:         int(amDetail.SizeOption),
		PrivacyOption:      amDetail.PrivacyOption,
		TimeOption:         amDetail.TimeOption,
		StartTime:          startTimeString,
		TimeSet:            amDetail.TimeSet,
		EndTime:            endTimeString,
		AvailabilityOption: amDetail.AvailabilityOption,
		StartMonth:         amDetail.StartMonth,
		EndMonth:           amDetail.EndMonth,
		TypeOption:         amDetail.TypeOption,
		PriceOption:        amDetail.PriceOption,
		BrandOption:        amDetail.BrandOption,
		ListOptions:        amDetail.ListOptions,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateAmenityDetail", "listing amenity details", "update listing amenity details")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateOptionLocation(ctx *gin.Context) {
	var req UpdateOptionLocationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionLocationParams in ShouldBindJSON: %v,OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v,OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	lng := tools.ConvertLocationStringToFloat(req.Lng, 9)
	lat := tools.ConvertLocationStringToFloat(req.Lat, 9)
	geolocation := pgtype.Point{
		P:     pgtype.Vec2{X: lng, Y: lat},
		Valid: true,
	}
	// We would update the location data
	location, err := server.store.UpdateLocationTwo(ctx, db.UpdateLocationTwoParams{
		Street:      req.Street,
		City:        req.City,
		State:       req.State,
		Country:     req.Country,
		Postcode:    req.Postcode,
		Geolocation: geolocation,
		OptionID:    option.ID,
	})
	if err != nil {
		log.Printf("Error at  UpdateOptionLocation in UpdateLocationTwo err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while updating your shortlet location")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetOptionLocationRes{
		Street:               location.Street,
		City:                 location.City,
		State:                location.State,
		Country:              location.Country,
		Postcode:             location.Postcode,
		Lat:                  tools.ConvertFloatToLocationString(location.Geolocation.P.Y, 9),
		Lng:                  tools.ConvertFloatToLocationString(location.Geolocation.P.X, 9),
		ShowSpecificLocation: location.ShowSpecificLocation,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionLocation", "listing location", "update listing location")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateShowSpecificLocation(ctx *gin.Context) {
	var req UpdateShowSpecificLocationParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateShowSpecificLocationParams in ShouldBindJSON: %v,OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	requestID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v,OptionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We would update the location data
	location, err := server.store.UpdateSpecificLocation(ctx, db.UpdateSpecificLocationParams{
		ShowSpecificLocation: req.ShowSpecificLocation,
		OptionID:             option.ID,
	})
	if err != nil {
		log.Printf("Error at  UpdateShowSpecificLocation in UpdateLocationTwo err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while updating your shortlet location")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateShowSpecificLocationParams{
		ShowSpecificLocation: location.ShowSpecificLocation,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateHowSpecificLocation", "listing location", "update listing specific location")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionLocation(ctx *gin.Context) {
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
	location, err := server.store.GetLocation(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  GetOptionLocation in GetLocation err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("an error occurred while getting your shortlet location")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetOptionLocationRes{
		Street:               location.Street,
		City:                 location.City,
		State:                location.State,
		Country:              location.Country,
		Postcode:             location.Postcode,
		Lat:                  tools.ConvertFloatToLocationString(location.Geolocation.P.Y, 9),
		Lng:                  tools.ConvertFloatToLocationString(location.Geolocation.P.X, 9),
		ShowSpecificLocation: location.ShowSpecificLocation,
	}
	fmt.Println(res)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionExtraInfo(ctx *gin.Context) {
	var req GetOptionExtraInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionExtraInfoParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
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
	optionExtraInfo, err := server.store.GetOptionExtraInfo(ctx, db.GetOptionExtraInfoParams{
		OptionID: option.ID,
		Type:     req.Type,
	})
	if err != nil {
		log.Printf("Error at  GetOptionExtraInfo in GetOptionExtraInfo err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	res := OptionExtraInfoRes{
		Info: optionExtraInfo.Info,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateUpdateOptionExtraInfo(ctx *gin.Context) {
	var req CreateUpdateOptionExtraInfoParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateUpdateOptionExtraInfoParams in ShouldBindJSON: %v, OptionTimeID: %v \n", err.Error(), req.OptionID)
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
	_, err = server.store.GetOptionExtraInfo(ctx, db.GetOptionExtraInfoParams{
		OptionID: option.ID,
		Type:     req.Type,
	})
	if err != nil {
		log.Printf("Error at  CreateUpdateOptionExtraInfo in GetOptionExtraInfo err: %v, user: %v\n", err, user.ID)
		exist = false
	}
	var res OptionExtraInfoRes
	if exist {
		// we want to update the details
		optionExtraInfo, err := server.store.UpdateOptionExtraInfo(ctx, db.UpdateOptionExtraInfoParams{
			OptionID: option.ID,
			Type:     req.Type,
			Info:     req.Info,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateOptionExtraInfo in UpdateOptionExtraInfo err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could not update the data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = OptionExtraInfoRes{
				Info: optionExtraInfo.Info,
			}
		}
	} else {
		// we want to create it cause it doesn't exist
		optionExtraInfo, err := server.store.CreateOptionExtraInfo(ctx, db.CreateOptionExtraInfoParams{
			OptionID: option.ID,
			Type:     req.Type,
			Info:     req.Info,
		})
		if err != nil {
			log.Printf("Error at CreateUpdateOptionExtraInfo in CreateOptionExtraInfo err: %v, user: %v\n", err, user.ID)
			err = fmt.Errorf("could add this data")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = OptionExtraInfoRes{
				Info: optionExtraInfo.Info,
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateUpdateOptionExtraInfo", "listing extra info", "create-update listing extra info")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListCheckInStep(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
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
	steps, err := server.store.ListCheckInStepOrdered(ctx, option.ID)
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("There an error at ListCheckInStep at ListCheckInStepOrdered: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		}
		res := "none"
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	var res ListCheckInStepRes
	var resData []CheckInStepRes
	for i := 0; i < len(steps); i++ {
		data := CheckInStepRes{
			ID:    tools.UuidToString(steps[i].ID),
			Des:   tools.HandleString(steps[i].Des),
			Photo: steps[i].Photo,
		}
		resData = append(resData, data)
	}
	var published bool
	if len(steps) > 0 {
		published = steps[0].PublishCheckInSteps
	}
	res = ListCheckInStepRes{
		List:      resData,
		Published: published,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveCheckInStepPhoto(ctx *gin.Context) {
	var req RemoveCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveCheckInStepPhotoParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}

	// We want to remove photo from fire base
	stepDetail, err := server.store.GetCheckInStep(ctx, db.GetCheckInStepParams{
		OptionID: option.ID,
		ID:       stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCheckInStepPhoto at GetCheckInStep: %v, optionID: %v, userID: %v, stepID: %v \n", err.Error(), option.ID, user.ID, stepID)
		err = fmt.Errorf("could not find this step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = RemoveFirebasePhoto(server, ctx, stepDetail.Photo)
	if err != nil {
		log.Printf("There an error at RemoveCheckInStepPhoto at RemoveFirebasePhoto: %v, optionID: %v, userID: %v, stepID: %v \n", err.Error(), option.ID, user.ID, stepID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	step, err := server.store.UpdateCheckInStepPhoto(ctx, db.UpdateCheckInStepPhotoParams{
		Photo:    "none",
		OptionID: option.ID,
		ID:       stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCheckInStepPhoto at UpdateCheckInStepPhoto: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("your photo was deleted but not updated on the database, please if anything feels wrong just connect us")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := CheckInStepRes{
		ID:    tools.UuidToString(step.ID),
		Des:   step.Des,
		Photo: step.Photo,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "RemoveCheckInStepPhoto", "listing check in step photo", "remove listing check in step photo")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveCheckInStep(ctx *gin.Context) {
	var req RemoveCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveCheckInStepParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	stepDetail, err := server.store.GetCheckInStep(ctx, db.GetCheckInStepParams{
		OptionID: option.ID,
		ID:       stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCheckInStep at GetCheckInStep: %v, optionID: %v, userID: %v, stepID: %v \n", err.Error(), option.ID, user.ID, stepID)
		err = fmt.Errorf("could not find this step")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(stepDetail.Photo) != 0 && stepDetail.Photo != "none" {
		err = RemoveFirebasePhoto(server, ctx, stepDetail.Photo)
		if err != nil {
			log.Printf("There an error at RemoveCheckInStep at RemoveFirebasePhoto: %v, optionID: %v, userID: %v, stepID: %v \n", err.Error(), option.ID, user.ID, stepID)
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}

	err = server.store.RemoveCheckInStep(ctx, db.RemoveCheckInStepParams{
		OptionID: option.ID,
		ID:       stepID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCheckInStep at RemoveCheckInStep: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("the photo for this step was removed however something went wrong while updating it on the database. please refresh then connect us if anything feels wrong")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	steps, err := server.store.ListCheckInStepOrdered(ctx, option.ID)
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("There an error at RemoveCheckInStep at ListCheckInStepOrdered: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		}
		res := "none"
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	var res ListCheckInStepRes
	var resData []CheckInStepRes
	for i := 0; i < len(steps); i++ {
		data := CheckInStepRes{
			ID:    tools.UuidToString(steps[i].ID),
			Des:   steps[i].Des,
			Photo: steps[i].Photo,
		}
		resData = append(resData, data)
	}
	res = ListCheckInStepRes{
		List: resData,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "RemoveCheckInStep", "listing check in step", "remove listing check in step")
	}
	ctx.JSON(http.StatusOK, res)
}
func (server *Server) UpdateCheckInStep(ctx *gin.Context) {
	var req UpdateCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateCheckInStepParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	stepID, err := tools.StringToUuid(req.StepID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, StepID: %v \n", err.Error(), req.StepID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var res CheckInStepRes
	switch req.Type {
	case "photo":
		step, err := server.store.UpdateCheckInStepPhoto(ctx, db.UpdateCheckInStepPhotoParams{
			Photo:    req.Photo,
			OptionID: option.ID,
			ID:       stepID,
		})
		if err != nil {
			log.Printf("There an error at UpdateCheckInStep at UpdateCheckInStepPhoto: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   step.Des,
				Photo: step.Photo,
			}
		}
	case "des":
		step, err := server.store.UpdateCheckInStepDes(ctx, db.UpdateCheckInStepDesParams{
			Des:      req.Des,
			OptionID: option.ID,
			ID:       stepID,
		})
		if err != nil {
			log.Printf("There an error at UpdateCheckInStep at UpdateCheckInStepDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your des in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   step.Des,
				Photo: step.Photo,
			}
		}
	default:
		err = fmt.Errorf("type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateCheckInStep", "listing check in step description", "create-update listing check in step description")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateCheckInStep(ctx *gin.Context) {
	var req CreateCheckInStepParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateCheckInStepParams in CreateCheckInStep in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var res CheckInStepRes
	switch req.Type {
	case "photo":
		step, err := server.store.CreateCheckInStep(ctx, db.CreateCheckInStepParams{
			OptionID: option.ID,
			Photo:    req.Photo,
			Des:      "none",
		})
		if err != nil {
			log.Printf("There an error at CreateCheckInStep at CreateCheckInStep for photo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   step.Des,
				Photo: step.Photo,
			}
		}
	case "des":
		step, err := server.store.CreateCheckInStep(ctx, db.CreateCheckInStepParams{
			OptionID: option.ID,
			Photo:    "none",
			Des:      req.Des,
		})
		if err != nil {
			log.Printf("There an error at CreateCheckInStep at CreateCheckInStep for des: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = CheckInStepRes{
				ID:    tools.UuidToString(step.ID),
				Des:   step.Des,
				Photo: step.Photo,
			}
		}
	default:
		err = fmt.Errorf("type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateCheckInStep", "listing check in step", "create listing check in step")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetShortletCheckInMethod(ctx *gin.Context) {
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
	checkMethod, err := server.store.GetShortletCheckInMethod(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetShortletCheckInMethod at GetShortletCheckInMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("an error occurred while getting your check in method")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := GetShortletCheckInMethodRes{
		Des:    checkMethod.CheckInMethodDes,
		Method: checkMethod.CheckInMethod,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) UpdateShortletCheckInMethod(ctx *gin.Context) {
	var req UpdateCheckInMethodParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateCheckInMethodParams in UpdateCheckInMethod in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	optionID, err := tools.StringToUuid(req.OptionID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("error occurred while processing your request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(optionID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var res GetShortletCheckInMethodRes
	switch req.Type {
	case "method":
		checkMethod, err := server.store.UpdateShortletCheckInMethod(ctx, db.UpdateShortletCheckInMethodParams{
			CheckInMethod: pgtype.Text{
				String: req.Method,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at UpdateCheckInMethod at UpdateShortletInfo for photo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = GetShortletCheckInMethodRes{
				Method: checkMethod.CheckInMethod,
				Des:    checkMethod.CheckInMethodDes,
			}
		}
	case "des":
		checkMethod, err := server.store.UpdateShortletCheckInMethod(ctx, db.UpdateShortletCheckInMethodParams{
			CheckInMethodDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at UpdateCheckInMethod at UpdateShortletInfo for photo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update your photo in server")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = GetShortletCheckInMethodRes{
				Method: checkMethod.CheckInMethod,
				Des:    checkMethod.CheckInMethodDes,
			}
		}
	default:
		err = fmt.Errorf("type not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateCheckInMethod", "listing check in method", "update listing check in method")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListThingToNote(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
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
	notes, err := server.store.ListThingToNoteOne(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at ListThingToNote at ListThingToNoteOne: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := "none"
		ctx.JSON(http.StatusNotFound, res)
		return
	}
	var res ListThingToNoteRes
	var resData []ThingToNoteItem
	for _, note := range notes {
		data := ThingToNoteItem{
			Tag:     note.Tag,
			Checked: note.Checked,
			Type:    note.Type,
			ID:      tools.UuidToString(note.ID),
		}
		resData = append(resData, data)
	}

	res = ListThingToNoteRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

// CU means CreateUpdate
// ThingToNote is related to sm
// sm is related to UHMSafety
func (server *Server) CUThingToNote(ctx *gin.Context) {
	var req CUThingToNoteParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CUThingToNoteParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while setting up this amenity, try again")
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
	// We want to confirm the tag
	if !val.ValidateSmTag(req.Type, req.Tag) {
		err = fmt.Errorf("tag not found")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// We would try to create it however if it already exists then we update
	note, err := server.store.GetThingToNoteByType(ctx, db.GetThingToNoteByTypeParams{
		OptionID: option.ID,
		Type:     req.Type,
		Tag:      req.Tag,
	})
	var exists bool
	if err != nil {
		log.Printf("There an error at CUThingToNote at GetThingToNoteByType: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		if err == db.ErrorRecordNotFound {
			exists = false
		} else {
			exists = false
			//err = fmt.Errorf("could not add this amenity to your option")
			//ctx.JSON(http.StatusBadRequest, errorResponse(err))
			//return
		}
	} else {
		exists = true
	}
	var res ThingToNoteItem
	if exists {
		// We want to update the amenity
		update, err := server.store.UpdateThingToNote(ctx, db.UpdateThingToNoteParams{
			Checked: req.Checked,
			ID:      note.ID,
		})
		if err != nil {
			log.Printf("There an error at CUThingToNote at UpdateThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not update this amenity")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = ThingToNoteItem{
				Tag:     update.Tag,
				Checked: update.Checked,
				Type:    update.Type,
				ID:      tools.UuidToString(update.ID),
			}
		}
	} else {
		// We want to create
		create, err := server.store.CreateThingToNote(ctx, db.CreateThingToNoteParams{
			OptionID: option.ID,
			Tag:      req.Tag,
			Type:     req.Type,
			Checked:  req.Checked,
		})
		if err != nil {
			log.Printf("There an error at CUThingToNote at CreateThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("could not add this amenity to your option")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			res = ThingToNoteItem{
				Tag:     create.Tag,
				Checked: create.Checked,
				Type:    create.Type,
				ID:      tools.UuidToString(create.ID),
			}
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CUThingToNote", "listing things to note", "create-update listing things to note")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetThingToNoteDetail(ctx *gin.Context) {
	var req GetThingToNoteDetailParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetThingToNoteDetailParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
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
	noteID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, noteID: %v \n", err.Error(), req.ID)
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
	smDetail, err := server.store.GetThingToNote(ctx, db.GetThingToNoteParams{
		ID:       noteID,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("There an error at GetThingToNoteDetail at GetThingToNote: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		res := fmt.Errorf("none")
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	res := UThingToNoteDetailRes{
		Des: smDetail.Des,
		ID:  tools.UuidToString(smDetail.ID),
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UThingToNoteDetail(ctx *gin.Context) {
	var req UThingToNoteDetailReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UThingToNoteDetailRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
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
	noteID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, noteID: %v \n", err.Error(), req.ID)
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
	smDetail, err := server.store.UpdateThingToNoteDetail(ctx, db.UpdateThingToNoteDetailParams{
		OptionID: option.ID,
		ID:       noteID,
		Des:      req.Des,
	})
	if err != nil {
		log.Printf("There an error at UThingToNoteDetail at UpdateThingToNoteDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your details")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UThingToNoteDetailRes{
		Des: smDetail.Des,
		ID:  tools.UuidToString(smDetail.ID),
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UThingToNoteDetail", "listing thing to note detail", "update listing thing to note detail")
	}
	ctx.JSON(http.StatusOK, res)
}
