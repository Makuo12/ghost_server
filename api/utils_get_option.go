package api

import (
	"fmt"
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// This would return the data only if the option is not completed
func HandleGetOptionIncomplete(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (user db.User, option db.OptionsInfo, err error) {
	user, err = HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	option, err = server.store.GetOptionInfo(ctx, db.GetOptionInfoParams{
		ID:         reqID,
		HostID:     user.ID,
		IsComplete: false,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetOptionIncomplete in GetOptionInfoID %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("error occurred while performing your request")
		return
	}
	return
}

// This would return the data only if the option is completed (co-host post)
func HandleGetCompleteOptionPost(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionPost in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.Post {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host scan code)
func HandleGetCompleteOptionScanCode(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionScanCode in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.ScanCode {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host edit option info)
func HandleGetCompleteOptionEditOptionInfo(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionEditOptionInfo in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.EditOptionInfo {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host edit event date times)
func HandleGetCompleteOptionEditEventDateTimes(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionEditEventDateTimes in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.EditEventDatesTimes || !optionMain.EditOptionInfo {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host edit co host)
func HandleGetCompleteOptionEditCOHosts(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionEditCoHosts in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.EditCoHosts {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host calender)
func HandleGetCompleteOptionCalender(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionCalender in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.Calender {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

// This would return the data only if the option is completed (co-host edit reservations)
func HandleGetCompleteOptionReservation(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionReservation in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.Reservations {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

func HandleCoHostUpdateMsg(ctx *gin.Context, server *Server, userCoHost db.User, user db.User, option db.OptionsInfo, funcOne string, typeOne string, typeTwo string) {
	optionDetail, err := server.store.GetOptionInfoDetail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at HandleCoHostUpdateMsg at %v in user is cost host at GetOptionInfoDetail optionID: %v, optionID: %v, userID: %v \n", funcOne, err.Error(), option.ID, user.ID)
	}
	msg := "an update made to " + optionDetail.HostNameOption + "by " + userCoHost.Email + ". This change was made to " + typeOne + " specifically the " + typeTwo
	_, err = server.store.CreateOptionMessage(ctx, db.CreateOptionMessageParams{
		OptionID: option.ID,
		UserID:   user.ID,
		Message:  msg,
		Type:     "co_host",
	})
	if err != nil {
		log.Printf("There an error at HandleCoHostUpdateMsg at %v in user is cost host at CreateOptionMessage optionID: %v, optionID: %v, userID: %v \n", funcOne, err.Error(), option.ID, user.ID)
	}
}

func HandleGetCompleteOptionInsights(reqID uuid.UUID, ctx *gin.Context, server *Server, isForComplete bool) (userHost db.User, optionMain db.GetOptionInfoMainRow, option db.OptionsInfo, isCoHost bool, userCoHost db.User, err error) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		return
	}
	optionMain, err = server.store.GetOptionInfoMain(ctx, db.GetOptionInfoMainParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       reqID,
		OptionCoHostID: reqID,
	})
	if err != nil {
		log.Printf("Error occurred at HandleGetCompleteOptionInsights in store.GetOptionInfoMain %v, optionID: %v, userID: %v \n", err.Error(), reqID, user.ID)
		err = fmt.Errorf("you do not have access to this resource")
		return
	}
	if !optionMain.IsActive {
		err = fmt.Errorf("your request is forbidden as your option is deactivated. Contact our support team to know how to activate your option")
		return
	}
	// optionMain.HostType == "co_host" this means user is co-host
	if optionMain.HostType == "co_host" {
		if !optionMain.Insight {
			err = fmt.Errorf("you do not have access to this resource")
			return
		}
		userCoHost = user
		isCoHost = true
	}
	userHost = GetUserFromOptionMain(optionMain)
	option = GetOptionInfoFromOptionMain(optionMain)
	return
}

func GetUserFromOptionMain(optionMain db.GetOptionInfoMainRow) db.User {
	data := db.User{
		ID:                optionMain.UID,
		UserID:            optionMain.UserID,
		FirebaseID:        optionMain.FirebaseID,
		HashedPassword:    optionMain.HashedPassword,
		FirebasePassword:  optionMain.FirebasePassword,
		Email:             optionMain.Email,
		PhoneNumber:       optionMain.PhoneNumber,
		FirstName:         optionMain.FirstName,
		Username:          optionMain.Username,
		LastName:          optionMain.LastName,
		DateOfBirth:       optionMain.DateOfBirth,
		DialCode:          optionMain.DialCode,
		DialCountry:       optionMain.DialCountry,
		CurrentOptionID:   optionMain.CurrentOptionID,
		Currency:          optionMain.UCurrency,
		DefaultCard:       optionMain.DefaultCard,
		DefaultPayoutCard: optionMain.DefaultPayoutCard,
		DefaultAccountID:  optionMain.DefaultAccountID,
		IsActive:          optionMain.UIsActive,
		Image:             optionMain.HostImage,
		PasswordChangedAt: optionMain.UPasswordChangedAt,
		CreatedAt:         optionMain.UCreatedAt,
		UpdatedAt:         optionMain.UUpdatedAt,
		DeepLinkID:        optionMain.UserDeepLinkID,
	}
	return data
}

func GetOptionInfoFromOptionMain(optionMain db.GetOptionInfoMainRow) db.OptionsInfo {
	data := db.OptionsInfo{
		ID:             optionMain.ID,
		CoHostID:       optionMain.CoHostID,
		OptionUserID:   optionMain.OptionUserID,
		HostID:         optionMain.HostID,
		PrimaryUserID:  optionMain.PrimaryUserID,
		IsActive:       optionMain.IsActive,
		IsComplete:     optionMain.IsComplete,
		IsVerified:     optionMain.IsVerified,
		Category:       optionMain.Category,
		CategoryTwo:    optionMain.CategoryTwo,
		CategoryThree:  optionMain.CategoryThree,
		CategoryFour:   optionMain.CategoryFour,
		IsTopSeller:    optionMain.IsTopSeller,
		TimeZone:       optionMain.TimeZone,
		Currency:       optionMain.Currency,
		OptionImg:      optionMain.OptionImg,
		OptionType:     optionMain.OptionType,
		MainOptionType: optionMain.MainOptionType,
		CreatedAt:      optionMain.CreatedAt,
		Completed:      optionMain.Completed,
		UpdatedAt:      optionMain.UpdatedAt,
		DeepLinkID:     optionMain.DeepLinkID,
	}
	return data
}
