package api

import (
	//"database/sql"

	db "flex_server/db/sqlc"
	"flex_server/tools"
	"flex_server/val"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// U stands for update
func (server *Server) CreateOptionCoHost(ctx *gin.Context) {
	var req CreateOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  CreateOptionCOHostParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var email = strings.TrimSpace(strings.ToLower(req.Email))
	if user.Email == email {
		err = fmt.Errorf("you cannot create a co host using your own email")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHostEmails, err := server.store.ListOptionCOHostEmail(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at CreateOptionCoHost at ListOptionCOHostEmail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		coHostEmails = []string{""}
	}
	if tools.ContainsString(coHostEmails, email) {
		log.Printf("There an error at CreateOptionCoHost at ListOptionCOHostEmail: %v, optionID: %v, userID: %v \n", "err.Error()", option.ID, user.ID)
		err := fmt.Errorf("this email account already exist as a co host for this option")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var coHost db.CreateOptionCOHostRow
	switch option.MainOptionType {
	case "options":
		coHost, err = server.store.CreateOptionCOHost(ctx, db.CreateOptionCOHostParams{
			OptionID:            option.ID,
			Email:               email,
			Reservations:        req.Reservations,
			Post:                req.Post,
			ScanCode:            req.ScanCode,
			Calender:            req.Calender,
			EditOptionInfo:      req.EditOptionInfo,
			EditEventDatesTimes: false,
			EditCoHosts:         req.EditCoHosts,
		})
		if err != nil {
			log.Printf("There an error at CreateOptionCoHost options at CreateOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("was not able to add this account to one of your co-host")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	case "events":
		coHost, err = server.store.CreateOptionCOHost(ctx, db.CreateOptionCOHostParams{
			OptionID:            option.ID,
			Email:               email,
			Reservations:        req.Reservations,
			Post:                req.Post,
			ScanCode:            req.ScanCode,
			Calender:            req.Calender,
			EditOptionInfo:      req.EditOptionInfo,
			EditEventDatesTimes: req.EditEventDateTimes,
			EditCoHosts:         req.EditCoHosts,
		})
		if err != nil {
			log.Printf("There an error at CreateOptionCoHost events at CreateOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("was not able to add this account to one of your co-host")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
	}
	var mainOption string
	switch option.MainOptionType {
	case "options":
		mainOption = "Stays"
	case "events":
		mainOption = "Events"
	}
	err = SendEmailInvitationCode(server, coHost.Email, "Co-host", option.MainOptionType, user.FirstName, mainOption, "CreateOptionCoHost", coHost.ID)
	// we want to send an email to the cost host
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := OptionCOHostItem{
		ID:            tools.UuidToString(coHost.ID),
		Email:         coHost.Email,
		IsPrimaryHost: false,
		IsMainHost:    false,
		Accepted:      coHost.Accepted,
		FirstName:     "",
		ProfilePhoto:  "",
		Date:          "",
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "CreateOptionCoHost", "co-host", "create co-host")
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UpdateOptionCoHost(ctx *gin.Context) {
	var req UpdateOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UpdateOptionCOHostParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	coHostID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, id: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var email = strings.TrimSpace(strings.ToLower(req.Email))
	if user.Email == email {
		err = fmt.Errorf("you cannot create a co host using your own email")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var coHost db.UpdateOptionCOHostRow
	switch option.MainOptionType {
	case "options":
		coHost, err = server.store.UpdateOptionCOHost(ctx, db.UpdateOptionCOHostParams{
			OptionID:            option.ID,
			ID:                  coHostID,
			Reservations:        req.Reservations,
			Post:                req.Post,
			ScanCode:            req.ScanCode,
			Calender:            req.Calender,
			EditOptionInfo:      req.EditOptionInfo,
			EditEventDatesTimes: false,
			EditCoHosts:         req.EditCoHosts,
			Insights:            req.Insights,
		})
		if err != nil {
			log.Printf("There an error at UpdateOptionCoHost options at UpdateOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("was not able to add this account to one of your co-host")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return

		}
	case "events":
		coHost, err = server.store.UpdateOptionCOHost(ctx, db.UpdateOptionCOHostParams{
			OptionID:            option.ID,
			ID:                  coHostID,
			Reservations:        req.Reservations,
			Post:                req.Post,
			ScanCode:            req.ScanCode,
			Calender:            req.Calender,
			EditOptionInfo:      req.EditOptionInfo,
			EditEventDatesTimes: req.EditEventDateTimes,
			EditCoHosts:         req.EditCoHosts,
			Insights:            req.Insights,
		})
		if err != nil {
			log.Printf("There an error at UpdateOptionCoHost events at UpdateOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err = fmt.Errorf("was not able to add this account to one of your co-host")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return

		}
	}
	var firstName string
	var profilePhoto string
	var isPrimaryHost bool
	date := tools.ConvertTimeToStringDateOnly(coHost.CreatedAt)
	coUserID, err := tools.StringToUuid(coHost.CoUserID)
	if err != nil {
		log.Printf("There an error at UpdateOptionCoHost at tools.StringToUuid: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
	} else {
		u, err := server.store.GetUserByUserID(ctx, coUserID)
		if err != nil {
			log.Printf("There an error at UpdateOptionCoHost at GetUserWithEmail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		} else {
			firstName = u.FirstName
			profilePhoto = u.Photo
		}
	}
	if option.PrimaryUserID == coUserID {
		isPrimaryHost = true
	}
	res := OptionCOHostItem{
		ID:            tools.UuidToString(coHost.ID),
		Email:         coHost.Email,
		IsMainHost:    false,
		Accepted:      coHost.Accepted,
		FirstName:     firstName,
		ProfilePhoto:  profilePhoto,
		Date:          date,
		IsPrimaryHost: isPrimaryHost,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionCoHost", "co-host", "update co-host")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListOptionCoHost(ctx *gin.Context) {
	var req ListOptionCoHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListOptionCoHost in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var resData []OptionCOHostItem
	// First we add the main host
	if req.Offset == 0 {
		var isPrimaryHost bool
		if option.PrimaryUserID == user.UserID {
			// user.UserID is always gonna be the main user that created the item
			isPrimaryHost = true
		}
		data := OptionCOHostItem{
			ID:            tools.UuidToString(uuid.New()),
			Email:         user.Email,
			IsPrimaryHost: isPrimaryHost,
			IsMainHost:    true,
			Accepted:      true,
			FirstName:     user.FirstName,
			ProfilePhoto:  user.Photo,
			Date:          tools.ConvertTimeToStringDateOnly(user.CreatedAt),
		}
		resData = append(resData, data)
	}
	count, err := server.store.CountOptionCOHost(ctx, option.ID)
	if err != nil {
		log.Printf("Error at  ListOptionCoHost in store.CountOptionCoHostByCoHost err: %v, user: %v\n", err, user.ID)
		if len(resData) == 0 {
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
			return
		} else {
			err = nil
		}

	}
	if (count <= int64(req.Offset) || count == 0) && len(resData) == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count > 0 {
		coHosts, err := server.store.ListOptionCOHost(ctx, db.ListOptionCOHostParams{
			Limit:    15,
			Offset:   int32(req.Offset),
			OptionID: option.ID,
		})
		if err != nil {
			log.Printf("Error at  ListOptionCoHost in store.ListOptionCOHost err: %v, user: %v\n", err, user.ID)
			res := "none"
			ctx.JSON(http.StatusNoContent, res)
		}
		for _, co := range coHosts {
			var isPrimaryHost bool
			var userFound bool = true
			var data OptionCOHostItem
			var u db.User
			date := tools.ConvertTimeToStringDateOnly(co.CreatedAt)
			userID, err := tools.StringToUuid(co.CoUserID)
			if err != nil {
				log.Printf("There an error at ListOptionCoHost at tools.StringToUuid: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
				userFound = false
			}
			if tools.UuidToString(option.PrimaryUserID) == co.CoUserID {
				isPrimaryHost = true
			}
			if userFound {
				u, err = server.store.GetUserByUserID(ctx, userID)
				if err != nil {
					log.Printf("There an error at ListOptionCoHost at GetUserWithEmail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
					userFound = false
				} else {
					userFound = true
				}
			}
			if userFound {
				data = OptionCOHostItem{
					ID:            tools.UuidToString(co.ID),
					Email:         co.Email,
					IsPrimaryHost: isPrimaryHost,
					IsMainHost:    false,
					Accepted:      co.Accepted,
					FirstName:     u.FirstName,
					ProfilePhoto:  u.Photo,
					Date:          date,
				}
			} else {
				data = OptionCOHostItem{
					ID:            tools.UuidToString(co.ID),
					Email:         co.Email,
					IsPrimaryHost: isPrimaryHost,
					IsMainHost:    false,
					Accepted:      co.Accepted,
					FirstName:     "",
					ProfilePhoto:  "",
					Date:          date,
				}
			}
			resData = append(resData, data)
		}
	}
	res := ListOptionCOHostRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionCoHost(ctx *gin.Context) {
	var req GetOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  GetOptionCOHostParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	coHostID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, id: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	co, err := server.store.GetOptionCOHost(ctx, db.GetOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at GetOptionCOHost events at GetOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	res := GetOptionCOHostRes{
		Reservations:       co.Reservations,
		Post:               co.Post,
		ScanCode:           co.ScanCode,
		Calender:           co.Calender,
		EditOptionInfo:     co.EditOptionInfo,
		EditEventDateTimes: co.EditEventDatesTimes,
		EditCoHosts:        co.EditCoHosts,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ResendInviteCoHost(ctx *gin.Context) {
	var req GetOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ResendInviteCoHost in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	coHostID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, id: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHost, err := server.store.GetOptionCOHost(ctx, db.GetOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at ResendInviteCoHost events at GetOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if coHost.Accepted {
		err = fmt.Errorf("%v has already accepted the invitation", coHost.CoHostEmail)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var mainOption string
	switch coHost.MainOptionType {
	case "options":
		mainOption = "Stays"
	case "events":
		mainOption = "Events"
	}
	err = SendEmailInvitationCode(server, coHost.CoHostEmail, "Co-host", coHost.HostNameOption, coHost.MainHostName, mainOption, "ResendInviteCoHost", coHost.CoID)
	// we want to send an email to the cost host
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ResendInviteRes{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CancelInviteCoHost(ctx *gin.Context) {
	var req GetOptionCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  CancelInviteCoHost in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	coHostID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, id: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHost, err := server.store.GetOptionCOHost(ctx, db.GetOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at CancelInviteCoHost events at GetOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	// If the user has already accepted the invite we send an error because this route is only when the user hasn't accepted the invite
	if coHost.Accepted {
		err = fmt.Errorf("cannot cancel invite because request was already accepted. Instead try removing account with Remove co-host")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// Because the user hasn't accepted the invite we can easily remove the user account
	err = server.store.RemoveOptionCOHost(ctx, db.RemoveOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at CancelInviteCoHost events at RemoveOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := RemoveItemRes{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveCoHost(ctx *gin.Context) {
	var req RemoveCOHostParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  RemoveCoHost in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
		err = fmt.Errorf("an error occurred while getting up this amenity detail, try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println(req)
	coHostID, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at tools.StringToUuid: %v, id: %v \n", err.Error(), req.ID)
		err = fmt.Errorf("error occurred while processing your request")
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
	user, _, option, _, _, err := HandleGetCompleteOptionEditCOHosts(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	coHost, err := server.store.GetOptionCOHost(ctx, db.GetOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCoHost events at GetOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return

	}
	// If the user has not accepted the invite we send an error because this route is only when the user has accepted the invite
	if !coHost.Accepted {
		err = fmt.Errorf("cannot remove co-host because request is not already accepted")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.RemoveOptionCOHost(ctx, db.RemoveOptionCOHostParams{
		OptionID: option.ID,
		ID:       coHostID,
	})
	if err != nil {
		log.Printf("There an error at RemoveCoHost events at RemoveOptionCOHost: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get this co-host account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// we want to send a message to the user telling them they were removed
	res := RemoveItemRes{
		Success: true,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetCancelPolicy(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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

	cancelPolicy, err := server.store.GetCancelPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetCancelPolicy at GetCancelPolicy: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("could not get your cancel policy")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateCancelPolicyRes{
		TypeOne:       cancelPolicy.TypeOne,
		TypeTwo:       cancelPolicy.TypeTwo,
		RequestRefund: cancelPolicy.RequestARefund,
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UpdateCancelPolicy(ctx *gin.Context) {
	var req UpdateCancelPolicyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateCancelPolicyRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	if val.ValidateRegStayCancel(req.TypeOne) {
		err := fmt.Errorf("unable to update your cancel policy, your standard cancel policy cannot be empty")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We are not yet handling cancel policy for long stays
	cancelPolicy, err := server.store.UpdateCancelPolicy(ctx, db.UpdateCancelPolicyParams{
		TypeOne:        req.TypeOne,
		TypeTwo:        "none",
		RequestARefund: req.RequestRefund,
		OptionID:       option.ID,
	})
	if err != nil {
		log.Printf("There an error at UpdateCancelPolicy at UpdateCheckInOutDetail: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("unable to update your cancel policy")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UpdateCancelPolicyRes{
		TypeOne:       cancelPolicy.TypeOne,
		TypeTwo:       cancelPolicy.TypeTwo,
		RequestRefund: cancelPolicy.RequestARefund,
	}

	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UCheckInOutDetail", "option cancel policy", "update cancel policy")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOptionInfoStatus(ctx *gin.Context) {
	var req GetOptionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOptionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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

	optionStatus, err := server.store.GetOptionInfoStatus(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at GetOptionInfoStatus at GetOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err = fmt.Errorf("unable to get the status")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := UOptionInfoStatusRes{
		Status:          tools.HandleOptionStatus(optionStatus.OptionStatus),
		StatusReason:    optionStatus.StatusReason,
		SnoozeStartDate: tools.ConvertDateOnlyToString(optionStatus.SnoozeStartDate),
		SnoozeEndDate:   tools.ConvertDateOnlyToString(optionStatus.SnoozeEndDate),
		UnlistReason:    optionStatus.UnlistReason,
		UnlistDes:       optionStatus.UnlistDes,
	}
	ctx.JSON(http.StatusOK, res)
}

// U stands for update
func (server *Server) UpdateOptionInfoStatus(ctx *gin.Context) {
	var req UOptionInfoStatusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  UOptionInfoStatusRes in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	user, _, option, isCoHost, userCoHost, err := HandleGetCompleteOptionEditOptionInfo(requestID, ctx, server, true)
	if err != nil {
		err = fmt.Errorf("you cannot access this resource")
		ctx.JSON(http.StatusForbidden, errorResponse(err))
		return
	}
	var optionStatus db.UpdateOptionInfoStatusRow
	switch req.Status {
	case "unlist":
		startDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
		if err != nil {
			log.Printf("There an error start at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("start date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
		if err != nil {
			log.Printf("There an error end at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// We want to sent a reminder to remove the option info
		// We would delete it manually after 6 months
		optionStatus, err = server.store.UpdateOptionInfoStatus(ctx, db.UpdateOptionInfoStatusParams{
			Status:          req.Status,
			StatusReason:    req.StatusReason,
			SnoozeStartDate: startDate,
			SnoozeEndDate:   endDate,
			UnlistReason:    req.UnlistReason,
			UnlistDes:       req.UnlistDes,
			OptionID:        option.ID,
		})
		if err != nil {
			log.Printf("There an error unlist at UpdateOptionInfoStatus at UpdateOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("unable to unlist your option status, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		//// After we snooze we want to change option info to isActive to false
		//_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		//	IsActive: pgtype.Bool{
		//		Bool:  false,
		//		Valid: true,
		//	},
		//	ID: option.ID,
		//})

		//if err != nil {
		//	log.Printf("There an error snooze at UpdateOptionInfoStatus at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		//}
	case "list":
		startDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
		if err != nil {
			log.Printf("There an error start at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("start date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(tools.FakeDate)
		if err != nil {
			log.Printf("There an error end at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		optionStatus, err = server.store.UpdateOptionInfoStatus(ctx, db.UpdateOptionInfoStatusParams{
			Status:          req.Status,
			StatusReason:    req.StatusReason,
			SnoozeStartDate: startDate,
			SnoozeEndDate:   endDate,
			UnlistReason:    "none",
			UnlistDes:       "none",
			OptionID:        option.ID,
		})
		if err != nil {
			log.Printf("There an error at list UpdateOptionInfoStatus at UpdateOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("unable to list your option status, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

	case "snooze":
		startDate, err := tools.ConvertDateOnlyStringToDate(req.SnoozeStartDate)
		// We need to also make sure that the time
		if err != nil {
			log.Printf("There an error start at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("start date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(req.SnoozeEndDate)
		if err != nil {
			log.Printf("There an error end at UpdateOptionInfoStatus at ConvertDateOnlyStringToDate: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("end date is not in the right date format")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		err = tools.SnoozeDatesGood(startDate, endDate)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		// we put snooze on staged because it is not yet snooze and we also don't want to leave it on list
		optionStatus, err = server.store.UpdateOptionInfoStatus(ctx, db.UpdateOptionInfoStatusParams{
			Status:          "staged",
			StatusReason:    req.StatusReason,
			SnoozeStartDate: startDate,
			SnoozeEndDate:   endDate,
			UnlistReason:    "none",
			UnlistDes:       "none",
			OptionID:        option.ID,
		})

		if err != nil {
			log.Printf("There an error snooze at UpdateOptionInfoStatus at UpdateOptionInfoStatus: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("unable to snooze your option status, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		//// After we snooze we want to change option info to isActive to false
		//_, err = server.store.UpdateOptionInfo(ctx, db.UpdateOptionInfoParams{
		//	IsActive: pgtype.Bool{
		//		Bool:  false,
		//		Valid: true,
		//	},
		//	ID: option.ID,
		//})

		//if err != nil {
		//	log.Printf("There an error snooze at UpdateOptionInfoStatus at UpdateOptionInfo: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)

		//}
	}
	var optionType string
	if option.MainOptionType == "options" {
		optionType = "listing status"
	} else {
		optionType = "event status"
	}
	res := UOptionInfoStatusRes{
		Status:          tools.HandleOptionStatus(optionStatus.OptionStatus),
		StatusReason:    optionStatus.StatusReason,
		SnoozeStartDate: tools.ConvertDateOnlyToString(optionStatus.SnoozeStartDate),
		SnoozeEndDate:   tools.ConvertDateOnlyToString(optionStatus.SnoozeEndDate),
		UnlistReason:    optionStatus.UnlistReason,
		UnlistDes:       optionStatus.UnlistDes,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionInfoStatus", optionType, "update status")
	}
	ctx.JSON(http.StatusOK, res)
}

//

func (server *Server) UpdateOptionQuestion(ctx *gin.Context) {
	var req UpdateOptionQuestionParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at UpdateOptionQuestionParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	var res UpdateOptionQuestionRes
	if req.HostAsIndividual {
		question, err := server.store.UpdateOptionQuestion(ctx, db.UpdateOptionQuestionParams{
			OptionID: option.ID,
			HostAsIndividual: pgtype.Bool{
				Bool:  req.HostAsIndividual,
				Valid: true,
			},
		})
		if err != nil {
			log.Printf("There an error host as individual at UpdateOptionQuestion at .UpdateOptionQuestion: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("unable to snooze your option status, try again")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res = UpdateOptionQuestionRes{
			HostAsIndividual:  question.HostAsIndividual,
			OrganizationName:  question.OrganizationName,
			OrganizationEmail: question.OrganizationEmail,
		}

	} else {
		question, err := server.store.UpdateOptionQuestion(ctx, db.UpdateOptionQuestionParams{
			OptionID: option.ID,
			HostAsIndividual: pgtype.Bool{
				Bool:  req.HostAsIndividual,
				Valid: true,
			},
			OrganizationName: pgtype.Text{
				String: req.OrganizationName,
				Valid:  true,
			},
			OrganizationEmail: pgtype.Text{
				String: req.OrganizationEmail,
				Valid:  true,
			},
		})
		if err != nil {
			log.Printf("There an error host as professional at UpdateOptionQuestion at .UpdateOptionQuestion: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
			err := fmt.Errorf("request was unsuccessful")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res = UpdateOptionQuestionRes{
			HostAsIndividual:  question.HostAsIndividual,
			OrganizationName:  question.OrganizationName,
			OrganizationEmail: question.OrganizationEmail,
		}
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "UpdateOptionQuestion", "host question", "update host question")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) AddLegalRepresent(ctx *gin.Context) {
	var req AddLegalRepresentParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at AddLegalRepresentParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	legals, err := server.store.GetOptionQuestionLegal(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at UpdateOptionQuestion at .GetOptionQuestionLegal: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("request was unsuccessful")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(req.Name) == 0 || tools.ContainsString(legals, req.Name) {
		err := fmt.Errorf("try entering another legal representative this one might already exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	legals = append(legals, req.Name)
	legals, err = server.store.UpdateOptionQuestionLegal(ctx, legals)
	if err != nil {
		log.Printf("There an error host at UpdateOptionQuestion at .UpdateOptionQuestionLegal: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("request was unsuccessful")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := LegalRepresentRes{
		LegalRepresents: legals,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "AddLegalRepresent", "adding a legal representative", "add legal representative")
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) RemoveLegalRepresent(ctx *gin.Context) {
	var req RemoveLegalRepresentParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveLegalRepresentParams in ShouldBindJSON: %v, optionID: %v \n", err.Error(), req.OptionID)
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
	legals, err := server.store.GetOptionQuestionLegal(ctx, option.ID)
	if err != nil {
		log.Printf("There an error at UpdateOptionQuestion at .GetOptionQuestionLegal: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("request was unsuccessful")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if len(req.Name) == 0 || !tools.ContainsString(legals, req.Name) {
		err := fmt.Errorf("try entering a legal representative that already exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	legals = tools.RemoveString(legals, req.Name)
	if len(legals) == 0 {
		legals = []string{"none"}
	}
	legals, err = server.store.UpdateOptionQuestionLegal(ctx, legals)
	if err != nil {
		log.Printf("There an error host at UpdateOptionQuestion at .UpdateOptionQuestionLegal: %v, optionID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		err := fmt.Errorf("request was unsuccessful")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := LegalRepresentRes{
		LegalRepresents: legals,
	}
	if isCoHost {
		HandleCoHostUpdateMsg(ctx, server, userCoHost, user, option, "RemoveLegalRepresent", "Removing a legal representative", "Remove legal representative")
	}
	ctx.JSON(http.StatusOK, res)
}
