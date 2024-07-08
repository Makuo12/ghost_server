package api

import (
	"fmt"
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleUpdateCreateOptionDate(ctx *gin.Context, server *Server, optionID uuid.UUID, price string, req CreateUpdateOptionDateTimeParams, count int, userID uuid.UUID, currency string) error {

	var err error
	//count to know if data exist
	if count > 0 {
		var commonDates []string
		optionInfoTimes, err := server.store.ListAllOptionDateTime(ctx, optionID)
		if err != nil {
			log.Printf("Error at  HandleUpdateCreateOptionDate in .ListAllOptionDateTime err: %v, user: %v\n", err, userID)
			return err
		}
		for i := 0; i < len(optionInfoTimes); i++ {
			if tools.ContainsString(req.Dates, tools.ConvertDateOnlyToString(optionInfoTimes[i].Date)) {
				// Common dates
				// we next try to remove it
				//priceDB, err := MoneyToDB(currency, tools.MoneyStringToInt(price), server)
				//if err != nil {
				//	log.Printf("Error at  HandleUpdateCreateOptionDate in MoneyToDB err: %v, user: %v\n", err, userID)
				//	priceDB = price
				//}
				optionDate, err := server.store.UpdateAllOptionDateTime(ctx, db.UpdateAllOptionDateTimeParams{
					ID:        optionInfoTimes[i].ID,
					Note:      req.Note,
					Available: req.Available,
					Price:     tools.MoneyStringToInt(price),
				})
				if err != nil {
					log.Printf("Error at  HandleUpdateCreateOptionDate in .RemoveOptionDateTime err: %v, user: %v, optionID: %v\n", err, userID, optionInfoTimes[i].ID)
				} else {
					commonDates = append(commonDates, tools.ConvertDateOnlyToString(optionDate.Date))
				}
			}
		}
		// I would create for each date
		for i := 0; i < len(req.Dates); i++ {
			if !tools.ContainsString(commonDates, req.Dates[i]) {
				//priceDB, err := MoneyToDB(currency, tools.MoneyStringToInt(price), server)
				//if err != nil {
				//	log.Printf("Error at  HandleUpdateCreateOptionDate in MoneyToDB err: %v, user: %v\n", err, userID)
				//	priceDB = price
				//}
				date, err := tools.ConvertDateOnlyStringToDate(req.Dates[i])
				if err != nil {
					log.Printf("Error at  HandleUpdateCreateOptionDate in .ConvertDateOnlyStringToDate err: %v, user: %v, optionID: %v\n", err, userID, optionID)
				}
				_, err = server.store.CreateOptionDateTime(ctx, db.CreateOptionDateTimeParams{
					OptionID:  optionID,
					Note:      req.Note,
					Available: req.Available,
					Price:     tools.MoneyStringToInt(price),
					Date:      date,
				})
				if err != nil {
					log.Printf("Error at  HandleUpdateCreateOptionDate in .CreateOptionDateTime err: %v, user: %v, optionID: %v\n", err, userID, optionID)
				}
			}
		}
	} else {
		for i := 0; i < len(req.Dates); i++ {
			date, err := tools.ConvertDateOnlyStringToDate(req.Dates[i])
			if err != nil {
				log.Printf("Error at  HandleUpdateCreateOptionDate in .ConvertDateOnlyStringToDate err: %v, user: %v, optionID: %v\n", err, userID, optionID)
			}
			_, err = server.store.CreateOptionDateTime(ctx, db.CreateOptionDateTimeParams{
				OptionID:  optionID,
				Note:      req.Note,
				Available: req.Available,
				Price:     tools.MoneyStringToInt(price),
				Date:      date,
			})
			if err != nil {
				log.Printf("Error at  HandleUpdateCreateOptionDate in .CreateOptionDateTime err: %v, user: %v, optionID: %v\n", err, userID, optionID)
			}
		}

	}
	return err
}

func HandleListSpaceAreas(ctx *gin.Context, server *Server, option db.OptionsInfo, userID uuid.UUID) (res ListSpaceAreas, err error, contentFound bool) {
	contentFound = false
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("Error at ListSpaceAreas in GetShortlet: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		err = fmt.Errorf("an error %v occurred while getting your guest areas", err.Error())
		return
	}
	spaceSharedWith := tools.HandleDBList(shortlet.SharedSpacesWith)
	spaceAreas, err := server.store.ListOrderedSpaceArea(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Error at ListSpaceAreas in ListSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
			err = nil
			return
		}
		log.Printf("Error at ListSpaceAreas in ListSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		err = fmt.Errorf("an error %v occurred while getting your guest areas", err.Error())
		return
	}
	var spaceData = make(map[string]int)
	var resData []SpaceAreas
	for i := 0; i < len(spaceAreas); i++ {
		spaceData[spaceAreas[i].SpaceType] = spaceData[spaceAreas[i].SpaceType] + 1
		images := tools.HandleDBList(spaceAreas[i].Images)
		beds := tools.HandleDBList(spaceAreas[i].Beds)
		name := fmt.Sprintf("%v-%d", spaceAreas[i].SpaceType, spaceData[spaceAreas[i].SpaceType])
		data := SpaceAreas{
			ID:          tools.UuidToString(spaceAreas[i].ID),
			OptionID:    tools.UuidToString(spaceAreas[i].OptionID),
			SharedSpace: spaceAreas[i].SharedSpace,
			SpaceType:   spaceAreas[i].SpaceType,
			Image:       images,
			Beds:        beds,
			IsSuite:     spaceAreas[i].IsSuite,
			Name:        name,
		}
		resData = append(resData, data)
	}
	err = nil
	contentFound = true
	res = ListSpaceAreas{
		List:            resData,
		SharedSpaceWith: spaceSharedWith,
	}
	return
}

func HandleGetUnselectedPhotos(ctx *gin.Context, server *Server, option db.OptionsInfo, userID uuid.UUID, req AddPhotoSpaceAreasParams, reqPhotos []string, spaceAreaID uuid.UUID) (photos []string, hasPhotos bool) {
	photosDBList, err := server.store.ListSpaceAreaImages(ctx, db.ListSpaceAreaImagesParams{
		OptionID: option.ID,
		ID:       spaceAreaID,
	})
	if err != nil {
		// we don't want to through any error because it is possible not to find any photo
		log.Printf("There an error at AddPhotoSpaceAreas at ListSpaceAreaImages: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
	} else {
		fmt.Println("photoDBList", photosDBList)
	}

	// photosAll will contain all the photos in the database
	var photosAll []string
	for _, photoList := range photosDBList {
		for _, photo := range photoList {
			if photo != "none" && len(photo) > 0 {
				photosAll = append(photosAll, photo)
			}
		}
	}
	var photosNotSelected []string
	for _, photoItem := range reqPhotos {
		if !tools.ContainsString(photosAll, photoItem) {
			photosNotSelected = append(photosNotSelected, photoItem)
		}
	}
	var photosInOption []string
	if len(photosNotSelected) == 0 {
		hasPhotos = false
		return
	} else {
		// we want to check if the photos are in option info photo
		optionPhotos, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
		if err != nil {
			log.Printf("There an error at AddPhotoSpaceAreas at GetOptionInfoPhoto: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
			hasPhotos = false
			return
		}
		myOptionPhotos := optionPhotos.Images
		myOptionPhotos = append(myOptionPhotos, optionPhotos.MainImage)
		for _, photo := range photosNotSelected {
			if tools.ContainsString(myOptionPhotos, photo) {
				photosInOption = append(photosInOption, photo)
			}
		}
		hasPhotos = true
		photos = photosInOption
		return
	}
}

func HandleUpdateDes(ctx *gin.Context, server *Server, option db.OptionsInfo, userID uuid.UUID, req UpdateOptionDesParams) (res UpdateOptionDesParams, err error) {
	var description string
	switch req.DesType {
	case "des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			Des: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for Des: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.Des
		}
	case "space_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			SpaceDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for SpaceDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.SpaceDes
		}
	case "guest_access_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			GuestAccessDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for GuestAccessDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.GuestAccessDes
		}
	case "interact_with_guests_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			InteractWithGuestsDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for InteractWithGuestsDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.InteractWithGuestsDes
		}
	case "neighborhood_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			NeighborhoodDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for NeighborhoodDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.NeighborhoodDes
		}
	case "get_around_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			GetAroundDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for GetAroundDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
		} else {
			description = optionInfoDetail.GetAroundDes
		}
	case "other_des":
		optionInfoDetail, err := server.store.UpdateOptionInfoDetail(ctx, db.UpdateOptionInfoDetailParams{
			OtherDes: pgtype.Text{
				String: req.Des,
				Valid:  true,
			},
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("There an error at HandleUpdateDes at UpdateOptionInfoDetail for OtherDes: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)

		} else {
			description = optionInfoDetail.OtherDes
		}
	}
	res = UpdateOptionDesParams{
		OptionID: tools.UuidToString(option.ID),
		Des:      description,
		DesType:  req.DesType,
	}
	return
}

func HandleListOptionSelectComplete(ctx *gin.Context, server *Server, user db.User, req OptionSelectionOffsetParams) (res ListUHMOptionSelectionRes, err error, hasData bool) {
	var onLastIndex bool
	hasData = true
	count, err := server.store.CountOptionInfo(ctx, db.CountOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		log.Printf("Error at  HandleListOptionSelectComplete in CountOptionInfo err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return
	}
	// we want to get the options that are complete
	optionInfos, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		Limit:           10,
		Offset:          int32(req.OptionOffset),
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleListOptionSelectComplete in ListOptionInfo err: %v, user: %v\n", err, user.ID)
			hasData = false
			err = fmt.Errorf("an error occurred while getting your data")
			return
		}
	}
	var resData []UHMOptionSelectionRes
	for _, data := range optionInfos {
		var isActive bool
		var isCoHost bool
		var dOptionID uuid.UUID
		if data.OptionStatus == "list" || data.OptionStatus == "staged" {
			isActive = true
		}
		if data.HostType == "co_host" {
			isCoHost = true
			dOptionID = data.CoHostID

		} else {
			dOptionID = data.OptionID
		}
		newData := UHMOptionSelectionRes{
			HostNameOption: data.HostNameOption,
			MainImage:      HandleSqlNullString(data.MainImage),
			OptionID:       tools.UuidToString(dOptionID),
			MainOptionType: data.MainOptionType,
			HasName:        true,
			IsComplete:     data.IsComplete,
			IsActive:       isActive,
			IsCoHost:       isCoHost,
		}
		resData = append(resData, newData)
	}
	if hasData {
		hasData = true
		if count <= int64(req.OptionOffset+len(optionInfos)) {
			onLastIndex = true
		}
		res = ListUHMOptionSelectionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionInfos),
			OnLastIndex:  onLastIndex,
			Type:         req.Type,
		}
	}
	return
}

func HandleListOptionSelectInProgress(ctx *gin.Context, server *Server, user db.User, req OptionSelectionOffsetParams) (res ListUHMOptionSelectionRes, err error, hasData bool) {
	var onLastIndex bool
	hasData = true
	count, err := server.store.CountOptionInfo(ctx, db.CountOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		IsActive:        true,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		log.Printf("Error at  HandleListOptionSelectInProgress in CountOptionInfo err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return

	}
	// we want to get the options that are not complete
	optionInfos, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      false,
		Limit:           10,
		Offset:          int32(req.OptionOffset),
		IsActive:        true,
		OptionStatusOne: "unlist",
		OptionStatusTwo: "unlist",
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleListOptionSelectInProgress in ListOptionInfo err: %v, user: %v\n", err, user.ID)
			hasData = false
			err = fmt.Errorf("an error occurred while getting your experiences")
			return
		}
	}
	var resData []UHMOptionSelectionRes
	for _, data := range optionInfos {
		var name string
		var hasName bool
		var isActive bool
		var isCoHost bool
		var dOptionID uuid.UUID
		if !tools.ServerStringEmpty(data.HostNameOption) {
			hasName = true
			name = data.HostNameOption
		} else {
			switch data.MainOptionType {
			case "events":
				name = HandleSqlNullString(data.EventType)
			case "options":
				name = HandleSqlNullString(data.TypeOfShortlet)
			}
		}
		if data.OptionStatus == "list" || data.OptionStatus == "staged" {
			isActive = true
		}
		if data.HostType == "co_host" {
			isCoHost = true
			dOptionID = data.CoHostID
		} else {
			dOptionID = data.OptionID
		}
		newData := UHMOptionSelectionRes{
			HostNameOption: name,
			MainImage:      HandleSqlNullString(data.MainImage),
			OptionID:       tools.UuidToString(dOptionID),
			MainOptionType: data.MainOptionType,
			HasName:        hasName,
			IsComplete:     data.IsComplete,
			IsActive:       isActive,
			IsCoHost:       isCoHost,
		}
		resData = append(resData, newData)
	}
	if hasData {
		hasData = true
		if count <= int64(req.OptionOffset+len(optionInfos)) {
			onLastIndex = true
		}
		res = ListUHMOptionSelectionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionInfos),
			OnLastIndex:  onLastIndex,
			Type:         req.Type,
		}
	}
	return
}

func HandleListOptionSelectInActive(ctx *gin.Context, server *Server, user db.User, req OptionSelectionOffsetParams) (res ListUHMOptionSelectionRes, err error, hasData bool) {
	var onLastIndex bool
	hasData = true
	count, err := server.store.CountOptionInfo(ctx, db.CountOptionInfoParams{
		HostID:          user.ID,
		IsComplete:      true,
		IsActive:        true,
		OptionStatusOne: "unlist",
		OptionStatusTwo: "snooze",
	})
	if err != nil {
		log.Printf("Error at  HandleListOptionSelectInActive in CountOptionInfo err: %v, user: %v\n", err, user.ID)
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return
	}
	// we want to get the options that are not inactive

	optionInfos, err := server.store.ListOptionInfo(ctx, db.ListOptionInfoParams{
		HostID:          user.ID,
		CoUserID:        tools.UuidToString(user.UserID),
		IsComplete:      true,
		Limit:           10,
		Offset:          int32(req.OptionOffset),
		IsActive:        true,
		OptionStatusOne: "unlist",
		OptionStatusTwo: "snooze",
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleListOptionSelectInActive in ListOptionInfo err: %v, user: %v\n", err, user.ID)
			hasData = false
			err = fmt.Errorf("an error occurred while getting your experiences")
			return
		}
	}
	var resData []UHMOptionSelectionRes
	for _, data := range optionInfos {
		var isActive bool
		var isCoHost bool
		var dOptionID uuid.UUID
		if data.OptionStatus == "list" || data.OptionStatus == "staged" {
			isActive = true
		}
		if data.HostType == "co_host" {
			isCoHost = true
			dOptionID = data.CoHostID
		} else {
			dOptionID = data.OptionID
		}
		newData := UHMOptionSelectionRes{
			HostNameOption: data.HostNameOption,
			MainImage:      HandleSqlNullString(data.MainImage),
			OptionID:       tools.UuidToString(dOptionID),
			MainOptionType: data.MainOptionType,
			HasName:        true,
			IsComplete:     data.IsComplete,
			IsActive:       isActive,
			IsCoHost:       isCoHost,
		}
		resData = append(resData, newData)
	}
	if hasData {
		hasData = true
		if count <= int64(req.OptionOffset+len(optionInfos)) {
			onLastIndex = true
		}
		res = ListUHMOptionSelectionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionInfos),
			OnLastIndex:  onLastIndex,
			Type:         req.Type,
		}
	}
	return
}

func HandleAddOptionPhotos(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req UpdateOptionPhotoParams) (res UpdateOptionPhotoRes, err error) {
	var images = []string{}
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleAddOptionImages in GetOptionInfoPhotoOnly err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not update your images")
		return
	}
	if len(req.CreateImages) == 0 {
		err = fmt.Errorf("you must at least add a photo")
		return
	}
	images = append(images, optionPhoto.Images...)
	for i := 0; i < len(req.CreateImages); i++ {
		images = append(images, req.CreateImages[i])
	}

	// We update the database

	optionInfoPhoto, err := server.store.UpdateOptionInfoImages(ctx, db.UpdateOptionInfoImagesParams{
		Images:   images,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("Error at  HandleAddOptionPhotos in UpdateOptionInfoImages err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not update your photos")
		return
	}
	res = UpdateOptionPhotoRes{
		MainImage: optionInfoPhoto.MainImage,
		Images:    optionInfoPhoto.Images,
	}
	return
}

func HandleOptionPhotoChangeCover(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req UpdateOptionPhotoParams) (res UpdateOptionPhotoRes, err error) {
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleOptionPhotoChangeCover in GetOptionInfoPhoto err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not set up your new cover photo")
		return
	}
	// compare the cover photo paths to make sure if they match we send a error
	if optionPhoto.MainImage == req.ChangeMainImage {
		err = fmt.Errorf("change not made because the paths are the same")
		return
	}
	// we make sure the path already exist
	if !tools.ContainsString(optionPhoto.Images, req.ChangeMainImage) {
		err = fmt.Errorf("could not change photo because path doesn't match any of your saved photo paths")
		return
	}
	// we store the current option photos inside photos
	var images []string

	for _, photo := range optionPhoto.Images {
		if photo != req.ChangeMainImage {
			images = append(images, photo)
		}
	}
	// we append the cover image so that the cover image still remain path of the option images
	images = append(images, optionPhoto.MainImage)
	// we update the new cover image
	// we update the Photo to our new images
	optionInfoPhoto, err := server.store.UpdateOptionInfoPhoto(ctx, db.UpdateOptionInfoPhotoParams{
		MainImage: req.ChangeMainImage,
		Images:    images,
		OptionID:  option.ID,
	})
	if err != nil {
		log.Printf("Error at  HandleOptionPhotoChangeCover in UpdateOptionInfoPhoto err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not set up your new cover photo")
		return
	}
	res = UpdateOptionPhotoRes{
		MainImage: optionInfoPhoto.MainImage,
		Images:    optionInfoPhoto.Images,
	}
	return
}

func HandleDeleteOptionPhoto(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req UpdateOptionPhotoParams) (res UpdateOptionPhotoRes, err error) {
	if len(req.DeleteImage) == 0 {
		err = fmt.Errorf("photo path is empty")
		return
	}
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleDeleteOptionPhoto in GetOptionInfoPhoto err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not delete photo")
		return
	}
	// First confirm path exist
	if !tools.ContainsString(optionPhoto.Images, req.DeleteImage) {
		err = fmt.Errorf("could not remove this photo because path doesn't match any of your saved photo paths")
		return
	}
	// we want to go to firebase and delete the photo
	err = RemoveFirebasePhoto(server, ctx, req.DeleteImage)
	if err != nil {
		return
	}

	// We want to remove it from spaceAreas
	spaceAreas, err := server.store.ListOrderedSpaceArea(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at HandleDeleteOptionPhoto atListOrderedSpaceArea: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		for i := 0; i < len(spaceAreas); i++ {
			if tools.ContainsString(spaceAreas[i].Images, req.DeleteImage) {
				var spaceImages []string
				for j := 0; j < len(spaceAreas[i].Images); j++ {
					if spaceAreas[i].Images[i] != req.DeleteImage {
						spaceImages = append(spaceImages, spaceAreas[i].Images[i])
					}
				}
				_, err = server.store.UpdateSpaceAreaImages(ctx, db.UpdateSpaceAreaImagesParams{
					ID:     spaceAreas[i].ID,
					Images: spaceImages,
				})
				if err != nil {
					log.Printf("There an error at HandleDeleteOptionPhoto at UpdateSpaceAreaImages: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				}
				break
			}
		}
	}
	// We want to check if the photo has a caption and then delete it
	err = server.store.RemoveOptionPhotoCaption(ctx, db.RemoveOptionPhotoCaptionParams{
		OptionID: option.ID,
		PhotoID:  req.DeleteImage,
	})
	if err != nil {
		log.Printf("There an error at HandleDeleteOptionPhoto at RemoveOptionPhotoCaption: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}
	// if the photo was successfully deleted from firebase now we need to update our photos
	var newPhotos []string
	for _, item := range optionPhoto.Images {
		// we want to add the current photos in the db to the newPhotos
		// however we would not add the photo path that was deleted
		if item != req.DeleteImage {
			newPhotos = append(newPhotos, item)
		}
	}

	// We update the database
	// WE USE newPhotos because it contains the new photos
	optionInfoPhoto, err := server.store.UpdateOptionInfoImages(ctx, db.UpdateOptionInfoImagesParams{
		Images:   newPhotos,
		OptionID: option.ID,
	})
	if err != nil {
		log.Printf("Error at  HandleDeleteOptionPhoto in UpdateOptionInfoImages err: %v, user: %v, optionID: %v , photo path to delete: %v\n", err, user.ID, option.ID, req.DeleteImage)
		err = fmt.Errorf("could not update your photos")
		return
	}
	res = UpdateOptionPhotoRes{
		MainImage: optionInfoPhoto.MainImage,
		Images:    optionInfoPhoto.Images,
	}
	return
}

func HandleCreateOptionAddCharge(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req CreateUpdateOptionAddChargeParams) (res OptionAddChargeItem, err error) {
	var extraFee string
	var numOfGuest int
	switch req.Type {
	case "cleaning_fee":
		if tools.ConvertStringToFloat(req.ExtraFee) < 1 {
			extraFee = "0.0"
		} else {
			extraFee = req.ExtraFee
		}
		numOfGuest = 0
	case "pet_fee":
		extraFee = "0.0"
		numOfGuest = 0
	case "extra_guest_fee":
		extraFee = "0.0"
		numOfGuest = req.NumOfGuest
	default:
		err = fmt.Errorf("could not find this type")
		return OptionAddChargeItem{}, err
	}
	//priceDB, err := MoneyToDB(option.Currency, tools.MoneyStringToInt(req.MainFee), server)
	//if err != nil {
	//	log.Printf("Error at  HandleCreateOptionAddCharge in in mainFee MoneyToDB err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
	//	priceDB = tools.MoneyStringToInt(req.MainFee)
	//}
	//extraPriceDB, err := MoneyToDB(option.Currency, extraFee, server)
	//if err != nil {
	//	log.Printf("Error at  HandleCreateOptionAddCharge in in extraFee MoneyToDB err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
	//	extraPriceDB = extraFee
	//}
	charge, err := server.store.CreateOptionAddCharge(ctx, db.CreateOptionAddChargeParams{
		OptionID:   option.ID,
		Type:       req.Type,
		MainFee:    tools.MoneyStringToInt(req.MainFee),
		ExtraFee:   tools.MoneyStringToInt(extraFee),
		NumOfGuest: int32(numOfGuest),
	})
	if err != nil {
		log.Printf("Error at  HandleCreateOptionAddCharge in GetOptionInfoPhotoOnly err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not create this additional charge")
		return
	}
	err = nil
	res = OptionAddChargeItem{
		ID:         tools.UuidToString(charge.ID),
		MainFee:    tools.IntToMoneyString(charge.MainFee),
		Type:       req.Type,
		ExtraFee:   tools.IntToMoneyString(charge.ExtraFee),
		NumOfGuest: int(charge.NumOfGuest),
	}
	return
}

func HandleUpdateOptionAddCharge(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req CreateUpdateOptionAddChargeParams) (res OptionAddChargeItem, err error) {
	var extraFee string
	var numOfGuest int
	switch req.Type {
	case "cleaning_fee":
		if tools.ConvertStringToFloat(req.ExtraFee) < 1 {
			extraFee = "0.0"
		} else {
			extraFee = req.ExtraFee
		}
		numOfGuest = 0
	case "pet_fee":
		extraFee = "0.0"
		numOfGuest = 0
	case "extra_guest_fee":
		extraFee = "0.0"
		if req.NumOfGuest < 1 {
			numOfGuest = 1
		} else {
			numOfGuest = req.NumOfGuest
		}

	default:
		err = fmt.Errorf("could not find this type")
		return OptionAddChargeItem{}, err
	}
	charge, err := server.store.UpdateOptionAddChargeByType(ctx, db.UpdateOptionAddChargeByTypeParams{
		OptionID:   option.ID,
		Type:       req.Type,
		MainFee:    tools.MoneyStringToInt(req.MainFee),
		ExtraFee:   tools.MoneyStringToInt(extraFee),
		NumOfGuest: int32(numOfGuest),
	})
	if err != nil {
		log.Printf("Error at  HandleUpdateOptionAddCharge in GetOptionInfoPhotoOnly err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
		err = fmt.Errorf("could not create this additional charge")
		return
	}
	err = nil
	res = OptionAddChargeItem{
		ID:         tools.UuidToString(charge.ID),
		MainFee:    tools.IntToMoneyString(charge.MainFee),
		Type:       req.Type,
		ExtraFee:   tools.IntToMoneyString(charge.ExtraFee),
		NumOfGuest: int(charge.NumOfGuest),
	}
	return
}

func HandleLOTCreateUpdateOptionDiscount(ctx *gin.Context, server *Server, option db.OptionsInfo, user db.User, req LOTCreateUpdateOptionDiscountParams) (res ListOptionDiscountRes, err error, hasData bool) {
	var exists = false
	var resData []OptionDiscountItem
	discounts, err := server.store.ListOptionDiscountByMainType(ctx, db.ListOptionDiscountByMainTypeParams{
		OptionID: option.ID,
		MainType: "length_of_stay",
	})
	// if there an error that means no data
	if err != nil {
		log.Printf("Error at  HandleLOTCreateUpdateOptionDiscount in ListOptionDiscountByMainType err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
	} else {
		exists = true
	}

	if exists && len(discounts) > 0 {
		// if exists then we remove all of them
		err = server.store.RemoveOptionDiscountByMainType(ctx, db.RemoveOptionDiscountByMainTypeParams{
			OptionID: option.ID,
			MainType: "length_of_stay",
		})
		if err != nil {
			log.Printf("Error at  HandleLOTCreateUpdateOptionDiscount in .RemoveOptionDiscountByMainType err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
			err = fmt.Errorf("could make changes to your stay discount")
			hasData = false
			return
		}
	}
	// If there are not discounts in db
	if len(req.List) == 0 {
		// If length is zero we can't do anything
		hasData = false
		return
	} else {
		for _, dis := range req.List {
			discount, err := server.store.CreateOptionDiscount(ctx, db.CreateOptionDiscountParams{
				OptionID:  option.ID,
				Type:      dis.Type,
				MainType:  "length_of_stay",
				Percent:   int32(dis.Percent),
				Name:      "none",
				ExtraType: "none",
				Des:       "none",
			})
			if err != nil {
				log.Printf("Error at  HandleLOTCreateUpdateOptionDiscount in CreateOptionDiscount err: %v, user: %v, optionID: %v\n", err, user.ID, option.ID)
			} else {
				resData = append(resData, OptionDiscountItem{
					ID:        tools.UuidToString(discount.ID),
					Type:      discount.Type,
					MainType:  discount.MainType,
					Percent:   int(discount.Percent),
					Name:      discount.Name,
					ExtraType: discount.ExtraType,
					Des:       discount.Des,
				})
			}
		} //: End of loop
		res = ListOptionDiscountRes{
			List: resData,
		}
		if len(resData) > 0 {
			hasData = true
		}
		return
	}
}
