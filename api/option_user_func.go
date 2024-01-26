package api

import (
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleListOptionExperience(ctx *gin.Context, server *Server, req ExperienceOffsetParams) (res ListExperienceOptionRes, err error, hasData bool) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	var onLastIndex bool
	hasData = true
	count, err := server.store.GetOptionExperienceCount(ctx, db.GetOptionExperienceCountParams{
		IsComplete:      true,
		IsActive:        true,
		IsActive_2:      true,
		MainOptionType:  "options",
		Category:        req.Type,
	})
	if err != nil {
		log.Printf("Error at  HandleListOptionExperience in GetOptionExperienceCount err: %v, user: %v\n", err, ctx.ClientIP())
		hasData = false
		err = fmt.Errorf("could not perform your request")
		return
	}
	if count <= int64(req.OptionOffset) {
		err = nil
		hasData = false
		return
	}
	var optionInfos []db.ListOptionExperienceByLocationRow
	optionInfos, err = server.store.ListOptionExperienceByLocation(ctx, db.ListOptionExperienceByLocationParams{
		IsComplete:      true,
		IsActive:        true,
		IsActive_2:      true,
		MainOptionType:  "options",
		Category:        req.Type,
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
		Country:         req.Country,
		State:           req.State,
		Limit:           10,
		Offset:          int32(req.OptionOffset),
	})
	if err != nil {
		if err == db.ErrorRecordNotFound {
			err = nil
			hasData = false
			return
		} else {
			log.Printf("Error at  HandleListOptionExperience in ListOptionInfo err: %v, user: %v\n", err, ctx.ClientIP())
			hasData = false
			err = fmt.Errorf("an error occurred while getting your data")
			return
		}
	}

	var resData []ExperienceOptionData
	for _, data := range optionInfos {
		basePrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.Price), data.Currency, req.Currency, dollarToNaira, dollarToCAD, data.ID)
		if err != nil {
			log.Printf("Error at  basePrice HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
			basePrice = 0.0

		}
		weekendPrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.WeekendPrice), data.Currency, req.Currency, dollarToNaira, dollarToCAD, data.ID)
		if err != nil {
			log.Printf("Error at weekendPrice HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
			weekendPrice = 0.0
		}
		newData := ExperienceOptionData{
			UserOptionID:     tools.UuidToString(data.OptionUserID),
			Name:             data.HostNameOption,
			IsVerified:       data.IsVerified,
			CoverImage:       data.CoverImage,
			HostAsIndividual: data.HostAsIndividual,
			BasePrice:        tools.ConvertFloatToString(basePrice),
			WeekendPrice:     tools.ConvertFloatToString(weekendPrice),
			Photos:           data.Photo,
			TypeOfShortlet:   data.TypeOfShortlet,
			State:            data.State,
			Country:          data.Country,
			ProfilePhoto:     data.Photo_2,
			HostName:         data.FirstName,
			HostJoined:       tools.ConvertDateOnlyToString(data.CreatedAt),
			HostVerified:     data.IsVerified_2,
			Category:         data.Category,
		}
		resData = append(resData, newData)
	}
	if err == nil && hasData {
		if count <= int64(req.OptionOffset+len(optionInfos)) {
			onLastIndex = true
		}
		res = ListExperienceOptionRes{
			List:         resData,
			OptionOffset: req.OptionOffset + len(optionInfos),
			OnLastIndex:  onLastIndex,
			Category:     req.Type,
		}
	}
	return
}

func HandleDetailOptionExperience(ctx *gin.Context, server *Server, req ExperienceDetailParams) (res ExperienceOptionDetailRes, hasData bool, err error) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	optionUserID, err := tools.StringToUuid(req.OptionUserID)

	if err != nil {
		log.Printf("Error at  HandleDetailOptionExperience in StringToUuid err: %v, user: %v\n", err, ctx.ClientIP())
		hasData = false
		return
	}
	option, err := server.store.GetOptionInfoByOptionUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at  HandleDetailOptionExperience in GetOptionInfoByOptionUserID err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		hasData = false
		return
	}
	hostLanguages, hostBio, userID, cohost, book, caption := HandleExCommon(option, server, ctx)

	question, des, cancelPolicy, petsAllowed, totalReviewCount := HandleExCommonTwo(option, server, ctx, userID)

	spaceAreaDetail, spaceAreas, numOfBeds := HandleExSpaceArea(option, server, ctx)

	amenities, shortlet, location := HandleExSingle(option, server, ctx)

	review := HandleExOptionReview(option, server, ctx, userID)

	houseRules, tripLength := HandleExTripRules(option, server, ctx)

	notes, discount, addCharge := HandleExNoteDiscountCharge(option, server, ctx, req.Currency, dollarToNaira, dollarToCAD)

	checkInOut, available := HandleExQuestBookCheckAvail(option, server, ctx)

	res = ExperienceOptionDetailRes{
		SpaceAreas:           spaceAreas,
		SpaceAreaDetail:      spaceAreaDetail,
		Amenities:            amenities,
		Location:             location,
		HouseRules:           houseRules,
		Notes:                notes,
		HostLanguages:        hostLanguages,
		NumOfBeds:            numOfBeds,
		PetsAllowed:          petsAllowed,
		CoHost:               cohost,
		Review:               review,
		Des:                  des,
		TripLength:           tripLength,
		CancelPolicy:         cancelPolicy,
		Discount:             discount,
		AddCharge:            addCharge,
		ShortletDetail:       shortlet,
		Question:             question,
		CheckInOut:           checkInOut,
		Captions:             caption,
		AvailabilitySettings: available,
		BookMethod:           book,
		TotalReviewCount:     totalReviewCount,
		HostBio:              hostBio,
	}
	log.Printf("res detail %v\n", res)
	hasData = true
	err = nil
	return
}

func HandleExCommon(option db.OptionsInfo, server *Server, ctx *gin.Context) (hostLanguages []string, hostBio string, userID uuid.UUID, cohost []ExperienceDetailCoHost, book ExOptionBookMethod, caption []ExOptionPhotoCaptions) {
	optionData, err := server.store.GetOptionHost(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleExCommon in GetOptionHost err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		hostLanguages = []string{"none"}
		hostBio = "none"
	} else {
		hostLanguages = optionData.Languages
		hostBio = optionData.Bio
		userID = optionData.HostID
	}
	coHostData, err := server.store.ListOptionCOHostUser(ctx, option.ID)
	if err != nil || len(coHostData) == 0 {
		if err != nil {
			log.Printf("Error at HandleExCommon in .ListOptionCOHostUser err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		cohost = []ExperienceDetailCoHost{{"none", "none", "none", true}}
	} else {
		for _, host := range coHostData {
			data := ExperienceDetailCoHost{
				UserID:       tools.UuidToString(host.UserID),
				Name:         host.FirstName,
				ProfilePhoto: host.Photo,
				IsEmpty:      false,
			}
			cohost = append(cohost, data)
		}
	}
	b, err := server.store.GetOptionBookMethod(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionAvailabilitySetting err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		book = ExOptionBookMethod{false, "none", true}
	} else {
		book = ExOptionBookMethod{
			InstantBook: b.InstantBook,
			PreBookMsg:  b.PreBookMsg,
			IsEmpty:     false,
		}
	}
	ps, err := server.store.ListOptionPhotoCaption(ctx, option.ID)
	if err != nil || len(ps) == 0 {
		if err != nil {
			log.Printf("Error at HandleExNotes in ListOptionPhotoCaption err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		caption = []ExOptionPhotoCaptions{{"none", "none", true}}
	} else {
		for _, p := range ps {
			data := ExOptionPhotoCaptions{
				PhotoID: p.PhotoID,
				Caption: p.Caption,
				IsEmpty: false,
			}
			caption = append(caption, data)
		}
	}
	return

}

func HandleExCommonTwo(option db.OptionsInfo, server *Server, ctx *gin.Context, hostID uuid.UUID) (question ExOptionQuestions, des ExperienceDetailDes, cancelPolicy ExCancelPolicy, petsAllowed bool, totalReviewCount int) {
	//count, err := server.store.GetOptionInfoReviewUserCount(ctx, hostID)
	//if err != nil {
	//	log.Printf("Error atHandleExOptionReview in GetOptionInfoReviewUserCount err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
	//} else {
	//	totalReviewCount = int(count)
	//}
	q, err := server.store.GetOptionQuestion(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionQuestion err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		question = ExOptionQuestions{"none", "none", []string{"none"}, "none", "none", "none", "none", "none", "none", "none", true}
	} else {
		question = ExOptionQuestions{
			OrganizationName:  q.OrganizationName,
			OrganizationEmail: q.OrganizationEmail,
			LegalRepresents:   q.LegalRepresents,
			Street:            q.Street,
			State:             q.State,
			City:              q.City,
			Country:           q.Country,
			Postcode:          q.Postcode,
			Lat:               tools.ConvertFloatToLocationString(q.Geolocation.P.Y, 9),
			Lng:               tools.ConvertFloatToLocationString(q.Geolocation.P.X, 9),
			IsEmpty:           false,
		}
	}
	d, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionInfoDetail err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		des = ExperienceDetailDes{"none", "none", "none", "none", "none", "none", "none", true}
	} else {
		des = ExperienceDetailDes{
			Des:                  d.Des,
			GetAroundDes:         d.GetAroundDes,
			InteractWithGuestDes: d.InteractWithGuestsDes,
			SpaceDes:             d.SpaceDes,
			NeighborhoodDes:      d.NeighborhoodDes,
			GuestAccessDes:       d.GuestAccessDes,
			OtherDes:             d.OtherDes,
			IsEmpty:              false,
		}
		petsAllowed = d.PetsAllowed
	}

	cancel, err := server.store.GetCancelPolicy(ctx, option.ID)
	if err != nil || (tools.ServerStringEmpty(cancel.TypeOne) && tools.ServerStringEmpty(cancel.TypeTwo)) {
		if err != nil {
			log.Printf("Error at  HandleExDesTripPolicy in GetCancelPolicy err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		cancelPolicy = ExCancelPolicy{"none", "none", true}
	} else {
		cancelPolicy = ExCancelPolicy{
			TypeOne: cancel.TypeOne,
			TypeTwo: cancel.TypeTwo,
			IsEmpty: false,
		}
	}
	return
}

func HandleExSpaceArea(option db.OptionsInfo, server *Server, ctx *gin.Context) (spaceAreaDetail []ExperienceSpaceArea, spaceAreas []string, numOfBeds int) {
	spaceAreasDB, err := server.store.ListOrderedSpaceArea(ctx, option.ID)
	if err != nil || len(spaceAreasDB) == 0 {
		if err != nil {
			log.Printf("Error at  HandleExSpaceArea in StringToUuid err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		spaceAreas = []string{"none"}
		spaceAreaDetail = []ExperienceSpaceArea{{"none", "none", false, []string{"none"}, []string{"none"}, false, true}}
		return
	} else {
		numOfBeds = 0
		var spaceData = make(map[string]int)
		for i := 0; i < len(spaceAreasDB); i++ {
			spaceData[spaceAreasDB[i].SpaceType] = spaceData[spaceAreasDB[i].SpaceType] + 1
			photos := tools.HandleDBList(spaceAreasDB[i].Photos)
			beds := tools.HandleDBList(spaceAreasDB[i].Beds)
			if !tools.ServerListIsEmpty(beds) {
				numOfBeds += len(beds)
			}
			name := fmt.Sprintf("%v-%d", spaceAreasDB[i].SpaceType, spaceData[spaceAreasDB[i].SpaceType])
			data := ExperienceSpaceArea{
				AreaType:    spaceAreasDB[i].SpaceType,
				SharedSpace: spaceAreasDB[i].SharedSpace,
				Photos:      photos,
				Beds:        beds,
				IsSuite:     spaceAreasDB[i].IsSuite,
				Name:        name,
				IsEmpty:     false,
			}
			spaceAreaDetail = append(spaceAreaDetail, data)
			spaceAreas = append(spaceAreas, spaceAreasDB[i].SpaceType)
		}
	}

	return
}

func HandleExSingle(option db.OptionsInfo, server *Server, ctx *gin.Context) (amenities []string, shortlet ExShortletDetail, location ExperienceDetailLocation) {
	amenities, err := server.store.ListAmenitiesTag(ctx, db.ListAmenitiesTagParams{
		OptionID: option.ID,
		HasAm:    true,
	})
	if err != nil {
		log.Printf("Error at HandleExSingle in ListAmenitiesTag err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		amenities = []string{"none"}

	}

	shortletData, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleExSingle in GetShortlet err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		shortlet = ExShortletDetail{"none", false, "none", 0, 0, 0, "none", []string{"none"}, "none", true}
	} else {
		shortlet = ExShortletDetail{
			CheckInMethod:    shortletData.CheckInMethod,
			AnySpaceShared:   shortletData.AnySpaceShared,
			SpaceType:        shortletData.SpaceType,
			NumOfGuests:      int(shortletData.GuestWelcomed),
			YearBuilt:        int(shortletData.YearBuilt),
			PropertySize:     int(shortletData.PropertySize),
			PropertySizeUnit: shortletData.PropertySizeUnit,
			SharedSpacesWith: shortletData.SharedSpacesWith,
			TimeZone:         option.TimeZone,
			IsEmpty:          false,
		}
	}
	locationData, err := server.store.GetLocation(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleExPacked in GetLocation err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		location = ExperienceDetailLocation{"none", "none", false, "none", "none", true}
	} else {
		location = ExperienceDetailLocation{
			Street:               locationData.Street,
			City:                 locationData.City,
			ShowSpecificLocation: locationData.ShowSpecificLocation,
			Lat:                  tools.ConvertFloatToLocationString(locationData.Geolocation.P.Y, 9),
			Lng:                  tools.ConvertFloatToLocationString(locationData.Geolocation.P.X, 9),
			IsEmpty:              false,
		}
	}

	return
}

func HandleExOptionReview(option db.OptionsInfo, server *Server, ctx *gin.Context, hostID uuid.UUID) (review ExperienceReviewData) {
	reviewData, err := server.store.ListChargeReview(ctx, option.OptionUserID)
	environment := 0.0
	accuracy := 0.0
	communication := 0.0
	location := 0.0
	checkIn := 0.0
	general := 0.0

	if err != nil || len(reviewData) == 0 {
		log.Printf("review data bad %v\n", reviewData)
		if err != nil {
			log.Printf("Error atHandleExOptionReview in ListOptionInfoReview err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		rate := (environment + accuracy + communication + location + checkIn + general) / 6
		review = ExperienceReviewData{
			Environment:   tools.ConvertFloatToString(environment),
			Accuracy:      tools.ConvertFloatToString(accuracy),
			Communication: tools.ConvertFloatToString(communication),
			Location:      tools.ConvertFloatToString(location),
			CheckIn:       tools.ConvertFloatToString(checkIn),
			General:       tools.ConvertFloatToString(general),
			Rate:          tools.ConvertFloatToString(rate),
			ReviewCount:   0,
			IsEmpty:       true,
		}
	} else {
		log.Printf("review data good %v\n", reviewData)
		for _, rev := range reviewData {
			environment += float64(rev.Environment)
			accuracy += float64(rev.Accuracy)
			communication += float64(rev.Communication)
			location += float64(rev.Location)
			checkIn += float64(rev.CheckIn)
			general += float64(rev.General)
		}
		environment = environment / float64(len(reviewData))
		accuracy = accuracy / float64(len(reviewData))
		communication = communication / float64(len(reviewData))
		location = location / float64(len(reviewData))
		checkIn = checkIn / float64(len(reviewData))
		general = general / float64(len(reviewData))
		rate := (environment + accuracy + communication + location + checkIn + general) / 6
		review = ExperienceReviewData{
			Environment:   tools.ConvertFloatToString(environment),
			Accuracy:      tools.ConvertFloatToString(accuracy),
			Communication: tools.ConvertFloatToString(communication),
			Location:      tools.ConvertFloatToString(location),
			CheckIn:       tools.ConvertFloatToString(checkIn),
			General:       tools.ConvertFloatToString(general),
			Rate:          tools.ConvertFloatToString(rate),
			ReviewCount:   len(reviewData),
			IsEmpty:       false,
		}
	}
	return
}

func HandleExTripRules(option db.OptionsInfo, server *Server, ctx *gin.Context) (houseRules []ExperienceHouseRules, tripLength ExOptionTripLength) {
	rules, err := server.store.ListAllOptionRule(ctx, option.ID)
	if err != nil || len(rules) == 0 {
		if err != nil {
			log.Printf("Error at  HandleExDesTripPolicy in ListAllOptionRule err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		houseRules = []ExperienceHouseRules{{"none", "none", false, "none", "none", "none", true}}
	} else {
		for _, rule := range rules {
			data := ExperienceHouseRules{
				ID:        tools.UuidToString(rule.ID),
				Type:      rule.Type,
				Checked:   rule.Checked,
				Des:       rule.Des,
				StartTime: tools.ConvertTimeOnlyToString(rule.StartTime),
				EndTime:   tools.ConvertTimeOnlyToString(rule.EndTime),
				IsEmpty:   false,
			}
			houseRules = append(houseRules, data)
		}
	}
	t, err := server.store.GetOptionTripLength(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionTripLength err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		tripLength = ExOptionTripLength{0, 0, false, false, true}
	} else {
		tripLength = ExOptionTripLength{
			MinStayDays:                 int(t.MinStayDay),
			MaxStayDays:                 int(t.MaxStayNight),
			ManualApproveRequestPassMax: t.ManualApproveRequestPassMax,
			AllowReservationRequest:     t.AllowReservationRequest,
		}
	}
	return
}

func HandleExNoteDiscountCharge(option db.OptionsInfo, server *Server, ctx *gin.Context, userCurrency string, dollarToNaira string, dollarToCAD string) (notes []ExperienceSm, discount []ExOptionDiscount, addCharge []ExOptionAddCharge) {

	noteData, err := server.store.ListThingToNoteOne(ctx, option.ID)
	if err != nil || len(noteData) == 0 {
		log.Printf("Error at HandleExNotes in ListThingToNoteOne err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		notes = []ExperienceSm{{"none", "none", false, true}}
	} else {
		for _, rule := range noteData {
			data := ExperienceSm{
				ID:      tools.UuidToString(rule.ID),
				Type:    rule.Type,
				Checked: rule.Checked,
				IsEmpty: false,
			}
			notes = append(notes, data)
		}
	}

	ds, err := server.store.ListOptionDiscount(ctx, option.ID)
	if err != nil || len(ds) == 0 {
		if err != nil {
			log.Printf("Error at HandleExNotes in ListOptionDiscount err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		discount = []ExOptionDiscount{{"none", "none", 0, "none", "none", true}}
	} else {
		for _, d := range ds {
			data := ExOptionDiscount{
				ID:      tools.UuidToString(d.ID),
				Type:    d.Type,
				Percent: int(d.Percent),
				Des:     d.Des,
				Name:    d.Name,
				IsEmpty: false,
			}
			discount = append(discount, data)
		}
	}

	cs, err := server.store.ListOptionAddCharge(ctx, option.ID)
	if err != nil || len(cs) == 0 {
		if err != nil {
			log.Printf("Error at HandleExNotes in ListOptionAddCharge err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		}
		addCharge = []ExOptionAddCharge{{"none", "none", "none", "none", 0, true}}
	} else {
		for _, char := range cs {
			mainFee, err := tools.ConvertPrice(tools.IntToMoneyString(char.MainFee), option.Currency, userCurrency, dollarToNaira, dollarToCAD, char.ID)
			if err != nil {
				log.Printf("Error at  mainFee HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
				mainFee = 0.0

			}
			extraFee, err := tools.ConvertPrice(tools.IntToMoneyString(char.ExtraFee), option.Currency, userCurrency, dollarToNaira, dollarToCAD, char.ID)
			if err != nil {
				log.Printf("Error at extraFee HandleWishlistOptionExperience in ConvertPrice err: %v, user: %v\n", err, ctx.ClientIP())
				extraFee = 0.0
			}
			data := ExOptionAddCharge{
				ID:          tools.UuidToString(char.ID),
				Type:        char.Type,
				MainFee:     tools.ConvertFloatToString(mainFee),
				ExtraFee:    tools.ConvertFloatToString(extraFee),
				NumOfGuests: int(char.NumOfGuest),
				IsEmpty:     false,
			}
			addCharge = append(addCharge, data)
		}
	}
	return
}

func HandleExQuestBookCheckAvail(option db.OptionsInfo, server *Server, ctx *gin.Context) (checkInOut ExCheckInOutDetails, available ExOptionAvailabilitySettings) {

	a, err := server.store.GetOptionAvailabilitySetting(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionAvailabilitySetting err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		available = ExOptionAvailabilitySettings{"none", "none", "none", "none", true}
	} else {
		available = ExOptionAvailabilitySettings{
			AdvanceNotice:          a.AdvanceNotice,
			AdvanceNoticeCondition: a.AdvanceNoticeCondition,
			PreparationTime:        a.PreparationTime,
			AvailabilityWindow:     a.AvailabilityWindow,
			IsEmpty:                false,
		}
	}

	check, err := server.store.GetCheckInOutDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  HandleExDesTripPolicy in GetOptionAvailabilitySetting err: %v, user: %v, optionID: %v\n", err, ctx.ClientIP(), option.ID)
		checkInOut = ExCheckInOutDetails{"none", "none", "none", []string{"none"}, []string{"none"}, true}
	} else {
		arriveBefore, arriveAfter, leaveBefore := HandleCheckInTime(check.ArriveBefore, check.ArriveAfter, check.LeaveBefore)
		checkInOut = ExCheckInOutDetails{
			ArriveAfter:            arriveAfter,
			ArriveBefore:           arriveBefore,
			LeaveBefore:            leaveBefore,
			RestrictedCheckInDays:  check.RestrictedCheckInDays,
			RestrictedCheckOutDays: check.RestrictedCheckOutDays,
			IsEmpty:                false,
		}
	}
	return
}

func HandleCheckInTime(arriveBefore string, arriveAfter string, leaveBefore string) (arriveBeforeRes string, arriveAfterRes string, leaveBeforeRes string) {
	if tools.ServerStringEmpty(arriveBefore) {
		arriveBeforeRes = "20:30"
	} else {
		arriveBeforeRes = arriveBefore
	}
	if tools.ServerStringEmpty(arriveAfter) {
		arriveAfterRes = "08:30"
	} else {
		arriveAfterRes = arriveAfter
	}
	if tools.ServerStringEmpty(leaveBefore) {
		leaveBeforeRes = "19:30"
	} else {
		leaveBeforeRes = leaveBefore
	}
	return
}

func GetExDateTime(prepareTime string, id string, charges []db.ListChargeOptionReferenceByOptionUserIDRow, optionUserID string) (dates []ExOptionDateTimeItem) {
	for _, ch := range charges {
		// First we do the prepare time
		startDate := tools.ConvertDateOnlyToString(ch.StartDate)
		endDate := tools.ConvertDateOnlyToString(ch.EndDate)
		if prepareTime == constants.PREPARE_ONE_NIGHT || prepareTime == constants.PREPARE_TWO_NIGHT {
			dates = append(dates, HandlePrepareDates(optionUserID, prepareTime, startDate, endDate)...)
		}
		// We want to get the dates from startDate to endDate
		datesString, err := tools.GenerateDateListString(startDate, endDate)
		if err != nil {
			log.Printf("Error at HandleExPrepareTime in tools.DateByAddOrSubtractDays err: %v, user: %v, prepareOne, optionID: %v\n", err, prepareTime, optionUserID)
		} else {
			for _, d := range datesString {
				dates = append(dates, ExOptionDateTimeItem{d, false, "0.00", false})
			}

		}
	}

	return
}

func HandlePrepareDates(optionUserID, prepareTime, startDate, endDate string) (dates []ExOptionDateTimeItem) {
	var general int
	if prepareTime == constants.PREPARE_ONE_NIGHT {
		general = 1
	} else if prepareTime == constants.PREPARE_TWO_NIGHT {
		general = 2
	}
	for i := 1; i <= general; i++ {
		prepareOne, err := tools.DateByAddOrSubtractDays(startDate, -i)
		if err != nil {
			log.Printf("Error at HandleExPrepareTime in tools.DateByAddOrSubtractDays err: %v, user: %v, prepareOne, optionID: %v\n", err, prepareTime, optionUserID)
		} else {
			dates = append(dates, ExOptionDateTimeItem{prepareOne, false, "0.00", false})
		}
		prepareTwo, err := tools.DateByAddOrSubtractDays(endDate, i)
		if err != nil {
			log.Printf("Error at HandleExPrepareTime in tools.DateByAddOrSubtractDays err: %v, user: %v, prepareOne, optionID: %v\n", err, prepareTime, optionUserID)
		} else {
			dates = append(dates, ExOptionDateTimeItem{prepareTwo, false, "0.00", false})
		}
	}
	return
}
