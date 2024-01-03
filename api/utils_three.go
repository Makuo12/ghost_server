package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleSqlNullString(s pgtype.Text) string {
	log.Println("string check good", s.Valid)
	if s.Valid {

		log.Println("string all good", s)
		return s.String
	}
	log.Println("string not all good", s)
	return ""
}

func HandleSqlNullTime(s pgtype.Date) time.Time {
	log.Println("Time check good", s.Valid)
	if s.Valid {
		log.Println("Time all good", s)
		return s.Time
	}
	log.Println("Time not all good", s)
	return time.Now()
}

func HandleSqlNullTimestamp(s pgtype.Timestamptz) time.Time {
	log.Println("Time check good", s.Valid)
	if s.Valid {
		log.Println("Time all good", s)
		return s.Time
	}
	log.Println("Time not all good", s)
	return time.Now()
}

func HandleSqlNullBool(s pgtype.Bool) bool {
	log.Println("bool check good", s.Valid)
	if s.Valid {
		log.Println("bool all good", s)
		return s.Bool
	}
	log.Println("bool not all good", s)
	return false
}

func HandleSqlNullUUID(s pgtype.UUID) uuid.UUID {
	log.Println("bool check good", s.Valid)
	if s.Valid {
		log.Println("bool all good", s)
		return uuid.UUID(s.Bytes)
	}
	log.Println("bool not all good", s)
	return uuid.New()
}

func FindMinMaxDates(dateStrings []string) (time.Time, time.Time, error) {
	if len(dateStrings) == 0 {
		return time.Time{}, time.Time{}, fmt.Errorf("empty list of dates")
	}

	var minDate, maxDate time.Time

	for i, dateString := range dateStrings {
		date, err := time.Parse("2006-01-02", dateString) // Assuming the date format is "YYYY-MM-DD"
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("error parsing date at index %d: %v", i, err)
		}

		if i == 0 {
			minDate = date
			maxDate = date
		} else {
			if date.Before(minDate) {
				minDate = date
			}
			if date.After(maxDate) {
				maxDate = date
			}
		}
	}

	return minDate, maxDate, nil
}

// If the user selects a date that was already selected we remove that date
func HandleEventDatesUpdateList(server *Server, ctx *gin.Context, requestID uuid.UUID, reqDates []string) (dates []string, err error) {
	log.Println("At recurring d")
	var mainDates []string
	// first we would make sure the event dates are in the right format
	for _, d := range reqDates {
		_, err = tools.ConvertDateOnlyStringToDate(d)
		if err != nil {
			return
		}
		mainDates = append(mainDates, d)
	}
	log.Println("At recurring dd", mainDates)
	eventDates, err := server.store.GetEventDateTimeDates(ctx, requestID)
	log.Println("At recurring dates", eventDates)
	if err != nil {
		log.Println("At recurring dates err", err)
		err = fmt.Errorf("could not get your event dates")
		return
	}
	// we want to append the dates that were in the request sent that are not in the eventDates(SERVER).
	// we append it to dates []string
	for _, d := range mainDates {
		// we only want to add the dates that do not exist
		log.Println("c", d)
		if !tools.ContainsString(eventDates, d) {
			log.Println("cs", d)
			dates = append(dates, d)
		}
	}
	// we want to append the eventDates(SERVER) that are not in the sent request (mainDates).
	// we append it to dates []string
	for _, d := range eventDates {
		// we only want to add the dates that do not exist
		log.Println("h", d)
		if !tools.ContainsString(mainDates, d) {
			log.Println("hs", d)
			dates = append(dates, d)
		} else {
			fmt.Printf("contains")
		}
	}
	return
}

func HandleCompleteOption(currentState string, previousState string, server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID) (completeOption db.CompleteOptionInfo, err error) {
	completeOption, err = server.store.UpdateCompleteOptionInfo(ctx, db.UpdateCompleteOptionInfoParams{
		OptionID:      option.ID,
		CurrentState:  currentState,
		PreviousState: previousState,
	})
	if err != nil {
		log.Printf("Error at HandleCompleteOption in UpdateCompleteOptionInfo: %v, option.ID: %v, userID: %v  \n", err.Error(), option.ID, userID)
		err = fmt.Errorf("an error occurred while updating your state. try again")
		return
	}
	return
}
func CreateEventType(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, eventType string) (eventInfo db.EventInfo, err error) {
	err = server.store.RemoveEventInfo(ctx, option.ID)
	if err != nil {
		if err == db.ErrorRecordNotFound {
			log.Printf("Event doesn't exist %v \n", err.Error())
			err = nil
		} else {
			log.Println(err.Error())
			err = fmt.Errorf("error occurred while performing your request. Try again")
			return
		}
	}
	argEvent := db.CreateEventInfoParams{
		OptionID:  option.ID,
		EventType: eventType,
	}
	eventInfo, err = server.store.CreateEventInfo(ctx, argEvent)
	if err != nil {
		log.Printf("Error at CreateEventType in CreateEvent: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)

		err = server.store.RemoveOptionInfo(ctx, db.RemoveOptionInfoParams{
			ID:     option.ID,
			HostID: userID,
		})
		if err != nil {
			log.Printf("Error at CreateEventType in RemoveOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)

		}
		err = fmt.Errorf("an error occurred while creating your event, try again or contact help")
		return
	}

	return
}

func HandleRemoveEventInfo(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID) (err error) {
	err = server.store.RemoveEventInfo(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleEventInfo for optionID: %v, userID: %v, err is: %v", option.ID, userID, err.Error())
		err = fmt.Errorf("an error occurred while removing your event, try again")
		return
	}
	return

}

// This creates all the default tables for events
func HandleCreateEventDefaultTables(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) {
	// Create Option Info Detail
	_, err := server.store.CreateOptionInfoDetail(ctx, db.CreateOptionInfoDetailParams{
		OptionID:        option.ID,
		OptionHighlight: []string{""},
	})
	if err != nil {
		log.Printf("Error at HandleCreateEventDefaultTables in CreateOptionInfoDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option detail created")
	}

	startDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
	if err != nil {
		log.Printf("Error at start HandleCreateEventDefaultTables in ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		startDate = time.Now()
	}
	endDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
	if err != nil {
		log.Printf("Error at end HandleCreateEventDefaultTables in ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		endDate = time.Now()
	}

	// Create Option status
	_, err = server.store.CreateOptionInfoStatus(ctx, db.CreateOptionInfoStatusParams{
		OptionID:        option.ID,
		SnoozeStartDate: startDate,
		SnoozeEndDate:   endDate,
	})
	if err != nil {
		log.Printf("Error at HandleCreateEventDefaultTables in CreateOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	}

	// Create Cancel Policies
	_, err = server.store.CreateCancelPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateEventDefaultTables in CreateCancelPolicy: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Cancel policy created")
	}

	// Create Option book methods

	// We want instant book to be on for events
	_, err = server.store.CreateOptionBookMethod(ctx, db.CreateOptionBookMethodParams{
		OptionID:    option.ID,
		InstantBook: true,
	})
	if err != nil {
		log.Printf("Error at HandleCreateEventDefaultTables in CreateOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book method created")
	}

	// Create Book requirements
	_, err = server.store.CreateBookRequirement(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleCreateEventDefaultTables in CreateBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book requirements created")
	}

}

// This removes all the default tables for shortlets
func HandleRemoveEventDefaultTables(server *Server, ctx *gin.Context, option db.OptionsInfo, user db.User) {
	// Remove Option Info Detail
	err := server.store.RemoveOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveEventDefaultTables in RemoveOptionInfoRemove: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option detail removed successfully")
	}

	// Remove Cancel Policies
	err = server.store.RemoveCancelPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveEventDefaultTables in RemoveCancelPolicy: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Cancel policy removed successfully")
	}

	// Remove Option book methods
	err = server.store.RemoveOptionBookMethod(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveEventDefaultTables in RemoveOptionBookMethod: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book method removed successfully")
	}

	// Remove Book requirements
	err = server.store.RemoveBookRequirement(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveEventDefaultTables in RemoveBookRequirement: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Option book requirements removed successfully")
	}

	// Remove Option Info Status
	err = server.store.RemoveOptionInfoStatus(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleRemoveEventDefaultTables in RemoveOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		log.Println("Remove option info status removed successfully")
	}

}

func PublishHandleHostQuestion(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionQuestionParams) {
	var hasSecurityCamera bool
	var hasWeapons bool
	var hasDangerousAnimals bool
	var hostAsIndividual = true
	var organizationName = ""
	notes, err := server.store.ListThingToNoteChecked(ctx, option.ID)
	if err != nil {
		log.Printf("Error at PublishHandleHostQuestion in ListThingToNoteChecked: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
	} else {
		for _, note := range notes {
			switch note.Tag {
			case "dangerous_animal":
				hasDangerousAnimals = note.Checked
			case "weapon_on_property":
				hasWeapons = note.Checked
			case "cameras_audio_devices":
				hasSecurityCamera = note.Checked
			}

		}
	}
	hostQuestion, err := server.store.GetOptionQuestion(ctx, option.ID)
	if err != nil {
		log.Printf("Error at PublishHandleHostQuestion in GetOptionQuestion: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, userID)
	} else {
		hostAsIndividual = hostQuestion.HostAsIndividual
		organizationName = hostQuestion.OrganizationName
	}

	res = CreateOptionQuestionParams{
		OptionID:            tools.UuidToString(option.ID),
		UserOptionID:        tools.UuidToString(option.OptionUserID),
		CurrentServerView:   completeOption.CurrentState,
		PreviousServerView:  completeOption.PreviousState,
		OptionType:          option.OptionType,
		Currency:            option.Currency,
		HostAsIndividual:    hostAsIndividual,
		HasWeapons:          hasWeapons,
		HasSecurityCamera:   hasSecurityCamera,
		HasDangerousAnimals: hasDangerousAnimals,
		OrganizationName:    organizationName,
		MainOptionType:      option.MainOptionType,
	}
	return
}

func PublishHandlePhoto(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionInfoPhotoParams) {
	optionPhoto, err := server.store.GetOptionInfoPhoto(ctx, option.ID)
	if err != nil {
		res = CreateOptionInfoPhotoParams{
			OptionID:           tools.UuidToString(option.ID),
			CoverImage:         "",
			Photo:              []string{""},
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
		return
	}
	res = CreateOptionInfoPhotoParams{
		OptionID:           tools.UuidToString(option.ID),
		CoverImage:         optionPhoto.CoverImage,
		Photo:              optionPhoto.Photo,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	return
}
func PublishHandleHighlight(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionInfoHighlight) {
	optionInfoDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		res = CreateOptionInfoHighlight{
			OptionID:           tools.UuidToString(option.ID),
			Highlight:          []string{""},
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
		return
	}
	res = CreateOptionInfoHighlight{
		OptionID:           tools.UuidToString(option.ID),
		Highlight:          optionInfoDetail.OptionHighlight,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	return
}

func ReversePublishToPhoto(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateOptionInfoPhotoParams) {
	photoData, err := server.store.GetOptionInfoPhoto(ctx, option.ID)

	if err != nil {
		res = CreateOptionInfoPhotoParams{
			OptionID:           tools.UuidToString(option.ID),
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
			Photo:              []string{},
			CoverImage:         "",
		}
	} else {
		if len(photoData.Photo) == 0 {
			photoData.Photo = []string{""}
		}
		res = CreateOptionInfoPhotoParams{
			OptionID:           tools.UuidToString(option.ID),
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
			Photo:              photoData.Photo,
			CoverImage:         photoData.CoverImage,
		}
	}
	return
}

func EventLocationHandleSubCategory(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateEventSubCategoryParams) {
	eventInfo, err := server.store.GetEventInfo(ctx, option.ID)
	if err != nil {
		res = CreateEventSubCategoryParams{
			OptionID:           tools.UuidToString(option.ID),
			SubCategoryType:    "none",
			UserOptionID:       tools.UuidToString(option.OptionUserID),
			CurrentServerView:  completeOption.CurrentState,
			PreviousServerView: completeOption.PreviousState,
			MainOptionType:     option.MainOptionType,
			OptionType:         option.OptionType,
			Currency:           option.Currency,
		}
		return
	}
	res = CreateEventSubCategoryParams{
		OptionID:           tools.UuidToString(option.ID),
		SubCategoryType:    eventInfo.SubCategoryType,
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		OptionType:         option.OptionType,
		Currency:           option.Currency,
	}
	return
}

func LocationHandleShortletParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res ReverseShortletInfoParams) {
	var space = []string{}
	var spaceShared = false
	var spaceType string
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("Error at LocationHandleShortletParams in GetShortlet error: %v, option.ID: %v, user: %v", err.Error(), option.ID, userID)
		spaceType = ""
	} else {
		spaceType = shortlet.SpaceType
	}
	guestAndSpace, err := server.store.GetShortletGuestWelcomedAndShared(ctx, option.ID)
	if err != nil {
		log.Printf("Error at LocationHandleShortletParams in GetShortletGuestWelcomed error: %v, option.ID: %v, user: %v", err.Error(), option.ID, userID)
	} else {
		spaceShared = guestAndSpace.AnySpaceShared
		for i := 0; i < int(guestAndSpace.GuestWelcomed); i++ {
			space = append(space, "guest")
		}
	}
	spaceData, err := server.store.GetSpaceAreaType(ctx, option.ID)
	if err != nil {
		log.Printf("Error at LocationHandleShortletParams in GetSpaceAreaType error: %v, option.ID: %v, user: %v", err.Error(), option.ID, userID)
	} else {
		print("space data data", spaceData)
		space = append(space, spaceData...)
	}
	log.Println("Space data", space)
	res = ReverseShortletInfoParams{
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		Space:              space,
		SpaceType:          spaceType,
		AnySpaceShared:     spaceShared,
		Currency:           option.Currency,
		OptionType:         option.OptionType,
		OptionID:           tools.UuidToString(option.ID),
	}
	return
}

func DescriptionHandleEventParams(server *Server, ctx *gin.Context, option db.OptionsInfo, userID uuid.UUID, completeOption db.CompleteOptionInfo) (res CreateEventSubCategoryParams) {
	var eventSubType string
	var eventType string
	eventInfo, err := server.store.GetEventInfo(ctx, option.ID)
	if err != nil {
		log.Printf("Error at DescriptionHandleEventParams in GetEventInfo error: %v, option.ID: %v, user: %v", err.Error(), option.ID, userID)
		eventSubType = ""
		eventType = ""
	} else {
		eventSubType = eventInfo.SubCategoryType
		eventType = eventInfo.EventType
	}

	res = CreateEventSubCategoryParams{
		UserOptionID:       tools.UuidToString(option.OptionUserID),
		CurrentServerView:  completeOption.CurrentState,
		PreviousServerView: completeOption.PreviousState,
		MainOptionType:     option.MainOptionType,
		SubCategoryType:    eventSubType,
		EventType:          eventType,
		Currency:           option.Currency,
		OptionType:         option.OptionType,
		OptionID:           tools.UuidToString(option.ID),
	}
	return
}
