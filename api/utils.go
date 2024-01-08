package api

import (
	"context"
	db "flex_server/db/sqlc"
	"flex_server/token"
	"flex_server/tools"

	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleGetUser(ctx *gin.Context, server *Server) (user db.User, err error) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		log.Printf("Error at ctx.GET does not exist")
		err = fmt.Errorf("your credentials are wrong try logging in again, try again")
		return
	}
	data := payload.(*token.Payload)

	var username string = *&data.Username
	user, err = server.store.GetUserWithUsername(ctx, username)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Error at HandleGetOptionIncomplete in GetUserWithUsername: %v, userID: %v \n", err.Error(), user.ID)
			err = fmt.Errorf("this account isn't registered with Flexr")
			return
		}
		log.Printf("Error at HandleGetOptionIncomplete in GetUserWithUsername: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("an error occurred, make sure you are signed up on Flexr")
		return
	}
	if !user.IsActive {
		err = fmt.Errorf("your request is forbidden as your account is deactivated. Contact our support team to know how to activate your account")
		return
	}
	if user.IsDeleted {
		err = fmt.Errorf("this account does not exist")
		return
	}
	return
}

func HandleAmenityLocationData(server *Server, ctx *gin.Context, option db.OptionsInfo, completeOption db.CompleteOptionInfo) (res CreateLocationParams) {
	location, err := server.store.GetLocation(ctx, option.ID)
	if err != nil {
		res = CreateLocationParams{
			OptionID:             tools.UuidToString(option.ID),
			Street:               "none",
			City:                 "none",
			State:                "none",
			Country:              "none",
			Postcode:             "none",
			Lat:                  "none",
			Lng:                  "none",
			ShowSpecificLocation: false,
			UserOptionID:         tools.UuidToString(option.OptionUserID),
			CurrentServerView:    completeOption.CurrentState,
			PreviousServerView:   completeOption.PreviousState,
			MainOptionType:       option.MainOptionType,
			Currency:             option.Currency,
			OptionType:           option.OptionType,
		}
		return
	}

	res = CreateLocationParams{
		OptionID:             tools.UuidToString(option.ID),
		Street:               location.Street,
		City:                 location.City,
		State:                location.State,
		Country:              location.Country,
		Postcode:             location.Postcode,
		Lat:                  tools.ConvertFloatToLocationString(location.Geolocation.P.Y, 9),
		Lng:                  tools.ConvertFloatToLocationString(location.Geolocation.P.X, 9),
		ShowSpecificLocation: location.ShowSpecificLocation,
		UserOptionID:         tools.UuidToString(option.OptionUserID),
		CurrentServerView:    completeOption.CurrentState,
		PreviousServerView:   completeOption.PreviousState,
		MainOptionType:       option.MainOptionType,
		Currency:             option.Currency,
		OptionType:           option.OptionType,
	}
	return
}

// This creates all the default tables for shortlets
func HandleRemoveOptionDefaultTables(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) {
	// Remove Option Info Detail
	err := server.store.RemoveOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveOptionInfoRemove: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option detail removed successfully")
	}

	// Remove Option Trip Lengths
	err = server.store.RemoveOptionTripLength(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveOptionTripLength: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option trip length removed successfully")
	}

	// Remove Cancel Policies
	err = server.store.RemoveCancelPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveCancelPolicy: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Cancel policy removed successfully")
	}

	// Remove Available settings
	err = server.store.RemoveOptionAvailabilitySetting(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveOptionAvailabilitySetting: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option availability settings removed")
	}

	// Remove Check In Out Details
	err = server.store.RemoveCheckInOutDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveCheckInOutDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Check in out detail removed successfully")
	}

	// Remove Option book methods
	err = server.store.RemoveOptionBookMethod(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book method removed successfully")
	}

	// Remove Book requirements
	err = server.store.RemoveBookRequirement(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book requirements removed successfully")
	}

	// Remove Option Info Status
	err = server.store.RemoveOptionInfoStatus(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveOptionDefaultTables in RemoveOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Remove option info status removed successfully")
	}

}

// This creates all the default tables for shortlets
func HandleCreateOptionDefaultTables(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) {
	// Create Option Info Detail
	_, err := server.store.CreateOptionInfoDetail(ctx, db.CreateOptionInfoDetailParams{
		OptionID:        option.ID,
		OptionHighlight: []string{""},
	})
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option detail created")
	}

	_, err = server.store.CreateOptionAvailabilitySetting(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateOptionAvailabilitySetting: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option availability settings created")
	}

	startDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
	if err != nil {
		log.Printf("Error at start HandleCreateOptionDefaultTables in ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		startDate = time.Now()
	}
	endDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
	if err != nil {
		log.Printf("Error at end HandleCreateOptionDefaultTables in ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		endDate = time.Now()
	}
	_, err = server.store.CreateOptionInfoStatus(ctx, db.CreateOptionInfoStatusParams{
		OptionID:        option.ID,
		SnoozeStartDate: startDate,
		SnoozeEndDate:   endDate,
	})
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

	// Create Option Trip Lengths
	_, err = server.store.CreateOptionTripLength(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateOptionTripLength: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option trip length created")
	}

	// Create Cancel Policies
	_, err = server.store.CreateCancelPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateCancelPolicy: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Cancel policy created")
	}

	// Create Check In Out Details
	_, err = server.store.CreateCheckInOutDetail(ctx, db.CreateCheckInOutDetailParams{
		OptionID:               option.ID,
		RestrictedCheckInDays:  []string{"none"},
		RestrictedCheckOutDays: []string{"none"},
	})
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateCheckInOutDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Check in out detail created")
	}

	// Create Option book methods
	// We want instant book to be off for shortlets
	_, err = server.store.CreateOptionBookMethod(ctx, db.CreateOptionBookMethodParams{
		OptionID:    option.ID,
		InstantBook: false,
	})
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book method created")
	}

	// Create Book requirements
	_, err = server.store.CreateBookRequirement(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateOptionDefaultTables in CreateBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book requirements created")
	}

}

func HandleEventDatesList(list []string) (dates []string, err error) {
	for _, d := range list {
		_, err = tools.ConvertDateOnlyStringToDate(d)
		if err != nil {
			return
		}
		dates = append(dates, d)
	}
	return
}

func DescriptionHandleLocationParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateLocationParams) {
	location, err := server.store.GetLocation(ctx, option.ID)
	if err != nil {
		log.Printf("Error at DescriptionHandleLocationParams in GetLocation: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, userID)
		res = CreateLocationParams{
			OptionID:             tools.UuidToString(option.ID),
			Street:               "",
			City:                 "",
			State:                "",
			Country:              "",
			Postcode:             "",
			Lng:                  "",
			Lat:                  "",
			ShowSpecificLocation: true,
			UserOptionID:         tools.UuidToString(option.OptionUserID),
			CurrentServerView:    completeOption.CurrentState,
			PreviousServerView:   completeOption.PreviousState,
			MainOptionType:       option.MainOptionType,
			Currency:             option.Currency,
			OptionType:           option.OptionType,
		}
		return
	}
	res = CreateLocationParams{
		OptionID:             tools.UuidToString(option.ID),
		Street:               location.Street,
		City:                 location.City,
		State:                location.State,
		Country:              location.Country,
		Postcode:             location.Postcode,
		Lat:                  tools.ConvertFloatToLocationString(location.Geolocation.P.Y, 9),
		Lng:                  tools.ConvertFloatToLocationString(location.Geolocation.P.X, 9),
		ShowSpecificLocation: true,
		UserOptionID:         tools.UuidToString(option.OptionUserID),
		CurrentServerView:    completeOption.CurrentState,
		PreviousServerView:   completeOption.PreviousState,
		MainOptionType:       option.MainOptionType,
		Currency:             option.Currency,
		OptionType:           option.OptionType,
	}
	return
}

func DescriptionHandleAmenitiesParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateAmenitiesAndSafety) {
	var popularAm []string
	var homeSafetyAm []string
	arg := db.ListAmenitiesParams{
		OptionID: option.ID,
		HasAm:    true,
	}
	amenityData, err := server.store.ListAmenities(ctx, arg)
	fmt.Println("amenityData", amenityData)
	if err != nil {
		log.Printf("Error at DescriptionHandleAmenitiesParams in ListAmenities: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, userID)
		res = CreateAmenitiesAndSafety{
			OptionID:           tools.UuidToString(option.ID),
			PopularAm:          popularAm,
			HomeSafetyAm:       homeSafetyAm,
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			Currency:           option.Currency,
			OptionType:         option.OptionType,
		}
		return
	}

	for i := 0; i < len(amenityData); i++ {
		switch amenityData[i].AmType {
		case "popular":
			popularAm = append(popularAm, amenityData[i].Tag)
		case "home_safety":
			homeSafetyAm = append(homeSafetyAm, amenityData[i].Tag)
		}
	}
	fmt.Println("popular_am", popularAm)
	fmt.Println("home_safety_am", homeSafetyAm)
	if len(homeSafetyAm) == 0 {
		homeSafetyAm = []string{""}
	}
	if len(popularAm) == 0 {
		popularAm = []string{""}
	}
	res = CreateAmenitiesAndSafety{
		OptionID:           tools.UuidToString(option.ID),
		PopularAm:          popularAm,
		HomeSafetyAm:       homeSafetyAm,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		Currency:           option.Currency,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
	}
	return
}

func HighlightHandleNameParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionInfoName) {
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HighlightHandleNameParams in GetOptionInfoDetail: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, userID)
		res = CreateOptionInfoName{
			OptionID:           tools.UuidToString(option.ID),
			HostNameOption:     "",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			MainOptionType:     option.MainOptionType,
			PreviousServerView: completeOption.PreviousState,
			Currency:           option.Currency,
			OptionType:         option.OptionType,
		}
		return
	}
	res = CreateOptionInfoName{
		OptionID:           tools.UuidToString(option.ID),
		HostNameOption:     optionDetail.HostNameOption,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		MainOptionType:     option.MainOptionType,
		PreviousServerView: completeOption.PreviousState,
		Currency:           option.Currency,
		OptionType:         option.OptionType,
	}
	return
}

func HighlightHandlePriceParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionPrice) {
	optionPrice, err := server.store.GetOptionPrice(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HighlightHandlePriceParams in GetOptionPrice: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, userID)
		res = CreateOptionPrice{
			OptionID:           tools.UuidToString(option.ID),
			Price:              "",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			MainOptionType:     option.MainOptionType,
			PreviousServerView: completeOption.PreviousState,
			Currency:           option.Currency,
			OptionType:         option.OptionType,
		}
		return
	}
	var priceData string
	//priceRes, err := MoneyDBToApp(option.Currency, optionPrice.Price, server)
	//if err != nil {
	//	log.Printf("Error at HighlightHandlePriceParams in MoneyDBToApp: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, userID)
	//	priceRes = optionPrice.Price
	//}
	price, err := tools.ConvertStringToInt64(tools.IntToMoneyString(optionPrice.Price))
	if err != nil {
		priceData = tools.IntToMoneyString(optionPrice.Price)
	} else {
		priceData = tools.ConvertInt64ToString(price)
	}
	res = CreateOptionPrice{
		OptionID:           tools.UuidToString(option.ID),
		Price:              priceData,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		MainOptionType:     option.MainOptionType,
		PreviousServerView: completeOption.PreviousState,
		Currency:           option.Currency,
		OptionType:         option.OptionType,
	}
	return
}

func ConfirmExperience(option db.OptionsInfo, server *Server, ctx *gin.Context) (err error) {
	switch option.OptionType {
	case "shortlets":
		_, err = server.store.GetShortlet(ctx, option.ID)
		if err != nil {
			log.Printf("Error at ConfirmRoutes at GetShortlet: %v", err.Error())
			err = fmt.Errorf("please make you filled all the detail to get this shortlet registered on Flexr")
			return
		}
	}
	return

}

func CreateOptionShortletType(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, shortletType string) (shortlet db.Shortlet, err error) {
	err = server.store.RemoveShortlet(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Shortlet doesn't exist %v \n", err.Error())
			err = nil
		} else {
			log.Println(err.Error())
			err = fmt.Errorf("error occurred while performing your request. Try again")

			return
		}
	}
	argShortlet := db.CreateShortletParams{
		OptionID:         option.ID,
		TypeOfShortlet:   shortletType,
		GuestWelcomed:    0,
		YearBuilt:        0,
		PropertySize:     0,
		PropertySizeUnit: "ft",
		SharedSpacesWith: []string{"none"},
	}

	shortlet, err = server.store.CreateShortlet(ctx, argShortlet)
	if err != nil {
		log.Printf("Error at CreateOptionShortletType in CreateShortlet: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)

		err = server.store.RemoveOptionInfo(ctx, db.RemoveOptionInfoParams{
			ID:     option.ID,
			HostID: userID,
		})
		if err != nil {
			log.Printf("Error at CreateOptionShortletType in RemoveOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)

		}
		err = fmt.Errorf("an error occurred while creating your shortlet, try again or contact help")
		return
	}

	return
}

func RemoveAllPhoto(server *Server, ctx *gin.Context, option db.OptionsInfo) (err error) {
	photoData, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		log.Printf("an error at RemoveAllPhoto in GetOptionInfoPhoto err: %v for optionID: %q", err.Error(), option.ID)
		err = fmt.Errorf("an error occurred while removing your photos")
		return
	}
	photos := []string{}
	deletedPhotos := []string{}
	notDeletedPhotos := []string{}
	photos = append(photos, photoData.Photo...)
	// delete cover photo
	err = RemoveFirebasePhoto(server, ctx, photoData.CoverImage)
	if err != nil {
		log.Printf("An error at RemoveAllPhoto in RemoveFirebasePhoto err: %v for optionID: %q", err.Error(), option.ID)
		err = fmt.Errorf("an error occurred while removing your photos")
		return

	}
	argCover := db.UpdateOptionInfoPhotoCoverParams{
		CoverImage: "none",
		OptionID:   option.ID,
	}
	_, err = server.store.UpdateOptionInfoPhotoCover(ctx, argCover)
	if err != nil {
		log.Printf("An error at RemoveAllPhoto in UpdateOptionInfoPhoto err: %v for optionID: %q", err.Error(), option.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		return
	}
	// delete all the photos in the array
	for i := 0; i < len(photos); i++ {
		object := photos[i]
		err = RemoveFirebasePhoto(server, ctx, object)
		if err != nil {
			log.Printf("An error at RemoveAllPhoto in RemoveFirebasePhoto err: %v for optionID: %q", err.Error(), option.ID)

		} else {
			deletedPhotos = append(deletedPhotos, object)
		}

		if i == len(photos)-1 {
			if len(deletedPhotos) != len(photos) {
				for j := 0; j < len(photos); j++ {
					found := false
					for k := 0; k < len(deletedPhotos); k++ {
						if deletedPhotos[j] == photos[k] {
							found = true
						}
					}
					if !found {
						notDeletedPhotos = append(notDeletedPhotos, deletedPhotos[j])
					}
				}

			}

		}

	}
	if len(notDeletedPhotos) < 1 {
		err = server.store.RemoveOptionInfoPhoto(ctx, option.ID)
		if err != nil {
			log.Printf("An error at RemoveAllPhoto in RemoveOptionInfoPhoto err: %v for optionID: %q", err.Error(), option.ID)
			err = fmt.Errorf("an error occurred while removing your past photos some might be deleted")
			return
		}
		err = nil
		return
	}
	// we would just update it
	argPhotos := db.UpdateOptionInfoPhotoOnlyParams{
		Photo:    notDeletedPhotos,
		OptionID: option.ID,
	}
	_, err = server.store.UpdateOptionInfoPhotoOnly(ctx, argPhotos)
	if err != nil {
		log.Printf("An error at RemoveAllPhoto in UpdateOptionInfoPhoto err: %v for optionID: %q", err.Error(), option.ID)
		err = fmt.Errorf("error occurred while performing your request, please try again")
		return
	}
	err = nil
	return

}

func RemoveFirebasePhoto(server *Server, ctx *gin.Context, object string) (err error) {
	// First we delete cover photo
	if object == "none" || len(object) < 1 {
		err = fmt.Errorf("no object found here try again")
		return
	}
	contextOne, cancel := context.WithTimeout(ctx, time.Second*20)
	defer cancel()

	o := server.Bucket.Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		err = fmt.Errorf("object.Attrs: %v", err)
		return
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err = o.Delete(contextOne); err != nil {
		err = fmt.Errorf("Object(%q).Delete: %v", object, err)
		return
	}
	fmt.Printf("Object %v was deleted", object)
	return nil
}

func HandleCreateGuestArea(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, spaceArea string, sharedSpace bool) (err error) {
	arg := db.CreateSpaceAreaParams{
		OptionID:    option.ID,
		SharedSpace: sharedSpace,
		SpaceType:   spaceArea,
		Photos:      []string{"none"},
		Beds:        []string{"none"},
	}
	_, err = server.store.CreateSpaceArea(ctx, arg)
	if err != nil {
		log.Printf("An error at HandleCreateGuestArea in CreateSpaceArea err: %v for optionID: %q", err.Error(), option.ID)
	}
	if err != nil {
		err = server.store.RemoveSpaceAreaAll(ctx, option.ID)
		if err != nil {
			log.Printf("An error at HandleCreateGuestArea in RemoveSpaceAreaAll err: %v for optionID: %q", err.Error(), option.ID)
		}
		err = fmt.Errorf("error occurred while adding your rooms and spaces, please try again later")
	}
	return
}

// Get Publish data for shortlet to publish view
func HandleShortletViewToPublish(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User, completeOption db.CompleteOptionInfo) (publishOptionParams PublishShortletOptionParams, err error) {
	var popularAm []string
	var homeSafetyAm []string
	publishData, err := server.store.GetShortletPublishData(ctx, option.ID)
	if err != nil {
		log.Printf("An error at HandleHostQuestionToPublish in GetShortletPublishData err: %v for optionID: %q", err.Error(), option.ID)
		return
	}
	spaceArea, err := server.store.GetSpaceAreaType(ctx, option.ID)
	if err != nil {
		log.Printf("An error at HandleHostQuestionToPublish in GetSpaceAreaType err: %v for optionID: %q", err.Error(), option.ID)
		return
	}
	if len(spaceArea) == 0 {
		spaceArea = []string{""}
	}
	arg := db.ListAmenitiesParams{
		OptionID: option.ID,
		HasAm:    true,
	}
	amenityData, err := server.store.ListAmenities(ctx, arg)
	if err != nil {
		log.Printf("An error at HandleHostQuestionToPublish in ListAmenities err: %v for optionID: %q", err.Error(), option.ID)
		return
	}
	for i := 0; i < len(amenityData); i++ {
		switch amenityData[i].AmType {
		case "popular":
			popularAm = append(popularAm, amenityData[i].Tag)
		case "home_safety":
			homeSafetyAm = append(homeSafetyAm, amenityData[i].Tag)
		}
	}
	if len(homeSafetyAm) == 0 {
		homeSafetyAm = []string{""}
	}
	if len(popularAm) == 0 {
		popularAm = []string{""}
	}
	publishOptionParams = PublishShortletOptionParams{
		OptionID:             tools.UuidToString(option.ID),
		UserOptionID:         tools.UuidToString(option.OptionUserID),
		CurrentServerView:    completeOption.CurrentState,
		MainOptionType:       option.MainOptionType,
		PreviousServerView:   completeOption.PreviousState,
		Currency:             option.Currency,
		OptionType:           option.OptionType,
		Name:                 publishData.HostNameOption,
		OptionMainType:       publishData.TypeOfShortlet,
		NumOfGuest:           int(publishData.GuestWelcomed),
		Space:                spaceArea,
		PopularAm:            popularAm,
		CoverImage:           publishData.CoverImage,
		FirstName:            user.FirstName,
		HomeSafetyAm:         homeSafetyAm,
		Street:               publishData.Street,
		State:                publishData.State,
		City:                 publishData.City,
		Country:              publishData.Country,
		Postcode:             publishData.Postcode,
		Description:          publishData.Des,
		ShowSpecificLocation: publishData.ShowSpecificLocation,
	}
	return

}

// Get Publish data for event to publish view
func HandleEventViewToPublish(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User, completeOption db.CompleteOptionInfo) (publishOptionParams PublishEventOptionParams, err error) {
	publishData, err := server.store.GetEventPublishData(ctx, option.ID)
	if err != nil {
		log.Printf("An error at HandleEventViewToPublish in GetEventPublishData err: %v for optionID: %q", err.Error(), option.ID)
		return
	}
	publishOptionParams = PublishEventOptionParams{
		OptionID:           tools.UuidToString(option.ID),
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		MainOptionType:     option.MainOptionType,
		PreviousServerView: completeOption.PreviousState,
		Currency:           option.Currency,
		OptionType:         option.OptionType,
		Name:               publishData.HostNameOption,
		OptionMainType:     publishData.EventType,
		CoverImage:         publishData.CoverImage,
		FirstName:          user.FirstName,
		Description:        publishData.Des,
	}
	return

}

// Get the host current option data for shortlet
func HandleShortletHostCurrentOption(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) (currentOption HostCurrentOptionParams, err error) {
	currentOptionData, err := server.store.GetShortletCurrentOptionData(ctx, db.GetShortletCurrentOptionDataParams{
		ID:         option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("An error at HandleShortletHostCurrentOption in GetShortletCurrentOptionData err: %v for optionID: %q", err.Error(), option.ID)
		return
	}

	currentOption = HostCurrentOptionParams{
		OptionID:       tools.UuidToString(option.ID),
		MainOptionType: option.MainOptionType,
		Currency:       option.Currency,
		OptionType:     option.OptionType,
		HostNameOption: currentOptionData.HostNameOption,
		OptionMainType: currentOptionData.TypeOfShortlet,
		CoverImage:     currentOptionData.CoverImage,
		State:          currentOptionData.State,
		Country:        currentOptionData.Country,
	}
	return

}

// Get the host current option data for Event
func HandleEventHostCurrentOption(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) (currentOption HostCurrentOptionParams, err error) {
	currentOptionData, err := server.store.GetEventCurrentOptionData(ctx, db.GetEventCurrentOptionDataParams{
		ID:         option.ID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("An error at HandleEventHostCurrentOption in GetEventCurrentOptionData err: %v for optionID: %q", err.Error(), option.ID)
		return
	}
	currentOption = HostCurrentOptionParams{
		OptionID:       tools.UuidToString(option.ID),
		MainOptionType: option.MainOptionType,
		Currency:       option.Currency,
		OptionType:     option.OptionType,
		HostNameOption: currentOptionData.HostNameOption,
		OptionMainType: currentOptionData.EventType,
		CoverImage:     currentOptionData.CoverImage,
		State:          "none",
		Country:        "none",
	}
	return

}
