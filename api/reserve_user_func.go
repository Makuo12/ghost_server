package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"

	"github.com/gin-gonic/gin"
)

// Current
func HandleCurrentListReserveUserOptionItem(server *Server, ctx *gin.Context, user db.User, req ListReserveUserItemParams) (res ListReserveUserItemRes, hasData bool, err error) {
	count, err := server.store.CountChargeOptionReferenceCurrent(ctx, db.CountChargeOptionReferenceCurrentParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  HandleCurrentListReserveUserOptionItem in GetChargeOptionReferenceCurrent err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		hasData = false
		return
	}
	reserves, err := server.store.ListChargeOptionReferenceCurrent(ctx, db.ListChargeOptionReferenceCurrentParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		Limit:      10,
		Offset:     int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error atHandleCurrentListReserveUserOptionItem in GetChargeOptionReferenceCurrent err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	var resData []ReserveUserItem

	for _, r := range reserves {
		timeString, timeType := HandleReserveUserTime(r.ArriveAfter, r.ArriveBefore)
		roomID, err := SingleContextRoom(ctx, server, user.UserID, r.UserID, "HandleMessageListen")
		if err != nil {
			continue
		}
		data := ReserveUserItem{
			ID:               tools.UuidToString(r.ID),
			MainOption:       r.MainOptionType,
			HostNameOption:   r.HostNameOption,
			HostUserID:       tools.UuidToString(r.UserID),
			StartDate:        tools.ConvertDateOnlyToString(r.StartDate),
			EndDate:          tools.ConvertDateOnlyToString(r.EndDate),
			HostName:         r.FirstName,
			ProfilePhoto:     r.Photo,
			OptionCoverImage: r.CoverImage,
			OptionPhotos:     r.Photo_2,
			StartTime:        timeString,
			Timezone:         r.TimeZone,
			EndTime:          r.LeaveBefore,
			CheckInMethod:    r.CheckInMethod,
			Type:             timeType,
			Grade:            "none",
			OptionType:       r.TypeOfShortlet,
			SpaceType:        r.SpaceType,
			State:            r.State,
			Country:          r.Country,
			Street:           r.Street,
			City:             r.City,
			ReviewStatus:     r.ReviewStage,
			RoomID:           tools.UuidToString(roomID),
		}
		resData = append(resData, data)
	}
	onLastIndex := false
	hasData = true
	if count <= int64(req.Offset+len(reserves)) {
		onLastIndex = true
	}
	res = ListReserveUserItemRes{
		List:        resData,
		MainOption:  req.MainOption,
		Offset:      req.Offset + len(reserves),
		OnLastIndex: onLastIndex,
		UserID:      tools.UuidToString(user.UserID),
	}
	return
}

func HandleCurrentListReserveUserEventItem(server *Server, ctx *gin.Context, user db.User, req ListReserveUserItemParams) (res ListReserveUserItemRes, hasData bool, err error) {
	log.Printf("reserves 126 %v \n", user.UserID)
	count, err := server.store.CountChargeTicketReferenceCurrent(ctx, db.CountChargeTicketReferenceCurrentParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
	})
	log.Printf("reserves 12 %v \n", count)
	if err != nil {
		log.Printf("reserves 123 %v \n", count)
		log.Printf("Error at HandleCurrentListReserveUserEventItem in CountChargeTicketReferenceCurrent err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		log.Printf("reserves 124 %v \n", count)
		hasData = false
		return
	}
	reserves, err := server.store.ListChargeTicketReferenceCurrent(ctx, db.ListChargeTicketReferenceCurrentParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		Limit:      10,
		Offset:     int32(req.Offset),
	})
	log.Printf("reserves 12 %v \n")
	if err != nil {
		log.Printf("Error at HandleCurrentListReserveUserEventItem in ListChargeTicketReferenceCurrent err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	var resData []ReserveUserItem

	for _, r := range reserves {
		data := ReserveUserItem{
			ID:               tools.UuidToString(r.ChargeID),
			MainOption:       r.MainOptionType,
			HostNameOption:   r.HostNameOption,
			HostUserID:       tools.UuidToString(r.UserID),
			StartDate:        tools.ConvertDateOnlyToString(r.StartDate),
			EndDate:          tools.ConvertDateOnlyToString(r.EndDate),
			HostName:         r.FirstName,
			ProfilePhoto:     r.Photo,
			OptionCoverImage: r.CoverImage,
			OptionPhotos:     r.Photo_2,
			StartTime:        r.StartTime,
			EndTime:          r.EndTime,
			Timezone:         r.TimeZone,
			CheckInMethod:    r.CheckInMethod,
			Type:             r.TicketType,
			Grade:            r.Grade,
			OptionType:       r.EventType,
			SpaceType:        "none",
			State:            HandleSqlNullString(r.State),
			Country:          HandleSqlNullString(r.Country),
			Street:           HandleSqlNullString(r.Street),
			City:             HandleSqlNullString(r.City),
			ReviewStatus:     r.ReviewStage,
			RoomID:           "none",
		}
		resData = append(resData, data)
	}
	onLastIndex := false
	hasData = true
	if count <= int64(req.Offset+len(reserves)) {
		onLastIndex = true
	}
	res = ListReserveUserItemRes{
		List:        resData,
		Offset:      req.Offset + len(reserves),
		OnLastIndex: onLastIndex,
		MainOption:  req.MainOption,
		UserID:      tools.UuidToString(user.UserID),
	}
	return
}

// End Current

// Visited
func HandleVisitedListReserveUserOptionItem(server *Server, ctx *gin.Context, user db.User, req ListReserveUserItemParams) (res ListReserveUserItemRes, hasData bool, err error) {
	count, err := server.store.CountChargeOptionReferenceVisited(ctx, db.CountChargeOptionReferenceVisitedParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  HandleVisitedListReserveUserOptionItem in GetChargeOptionReferenceVisited err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		hasData = false
		return
	}
	reserves, err := server.store.ListChargeOptionReferenceVisited(ctx, db.ListChargeOptionReferenceVisitedParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		Limit:      10,
		Offset:     int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error atHandleVisitedListReserveUserOptionItem in GetChargeOptionReferenceVisited err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	var resData []ReserveUserItem

	for _, r := range reserves {
		roomID, err := SingleContextRoom(ctx, server, user.UserID, r.UserID, "HandleMessageListen")
		if err != nil {
			continue
		}
		timeString, timeType := HandleReserveUserTime(r.ArriveAfter, r.ArriveBefore)
		data := ReserveUserItem{
			ID:               tools.UuidToString(r.ID),
			MainOption:       r.MainOptionType,
			HostNameOption:   r.HostNameOption,
			HostUserID:       tools.UuidToString(r.UserID),
			StartDate:        tools.ConvertDateOnlyToString(r.StartDate),
			EndDate:          tools.ConvertDateOnlyToString(r.EndDate),
			HostName:         r.FirstName,
			ProfilePhoto:     r.Photo,
			OptionCoverImage: r.CoverImage,
			OptionPhotos:     r.Photo_2,
			StartTime:        timeString,
			Timezone:         r.TimeZone,
			EndTime:          r.LeaveBefore,
			CheckInMethod:    r.CheckInMethod,
			Type:             timeType,
			Grade:            "none",
			OptionType:       r.TypeOfShortlet,
			SpaceType:        r.SpaceType,
			State:            r.State,
			Country:          r.Country,
			Street:           r.Street,
			City:             r.City,
			ReviewStatus:     "none",
			RoomID:           tools.UuidToString(roomID),
		}
		resData = append(resData, data)
	}
	onLastIndex := false
	hasData = true
	if count <= int64(req.Offset+len(reserves)) {
		onLastIndex = true
	}
	res = ListReserveUserItemRes{
		List:        resData,
		MainOption:  req.MainOption,
		Offset:      req.Offset + len(reserves),
		OnLastIndex: onLastIndex,
		UserID:      tools.UuidToString(user.UserID),
	}
	return
}

func HandleVisitedListReserveUserEventItem(server *Server, ctx *gin.Context, user db.User, req ListReserveUserItemParams) (res ListReserveUserItemRes, hasData bool, err error) {
	log.Printf("reserves 126 %v \n", user.UserID)
	count, err := server.store.CountChargeTicketReferenceVisited(ctx, db.CountChargeTicketReferenceVisitedParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
	})
	log.Printf("reserves 12 %v \n", count)
	if err != nil {
		log.Printf("reserves 123 %v \n", count)
		log.Printf("Error at HandleVisitedListReserveUserEventItem in CountChargeTicketReferenceVisited err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		log.Printf("reserves 124 %v \n", count)
		hasData = false
		return
	}
	reserves, err := server.store.ListChargeTicketReferenceVisited(ctx, db.ListChargeTicketReferenceVisitedParams{
		UserID:     user.UserID,
		Cancelled:  false,
		IsComplete: true,
		Limit:      10,
		Offset:     int32(req.Offset),
	})
	log.Printf("reserves 12 %v \n")
	if err != nil {
		log.Printf("Error at HandleVisitedListReserveUserEventItem in ListChargeTicketReferenceVisited err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	log.Printf("reserves %v \n", reserves)
	var resData []ReserveUserItem

	for _, r := range reserves {
		data := ReserveUserItem{
			ID:               tools.UuidToString(r.ChargeID),
			MainOption:       r.MainOptionType,
			HostNameOption:   r.HostNameOption,
			HostUserID:       tools.UuidToString(r.UserID),
			StartDate:        tools.ConvertDateOnlyToString(r.StartDate),
			EndDate:          tools.ConvertDateOnlyToString(r.EndDate),
			HostName:         r.FirstName,
			ProfilePhoto:     r.Photo,
			OptionCoverImage: r.CoverImage,
			OptionPhotos:     r.Photo_2,
			StartTime:        r.StartTime,
			EndTime:          r.EndTime,
			Timezone:         r.TimeZone,
			CheckInMethod:    r.CheckInMethod,
			Type:             r.TicketType,
			Grade:            r.Grade,
			OptionType:       r.EventType,
			SpaceType:        "none",
			State:            HandleSqlNullString(r.State),
			Country:          HandleSqlNullString(r.Country),
			Street:           HandleSqlNullString(r.Street),
			City:             HandleSqlNullString(r.City),
			ReviewStatus:     "none",
			RoomID:           "none",
		}
		resData = append(resData, data)
	}
	onLastIndex := false
	hasData = true
	if count <= int64(req.Offset+len(reserves)) {
		onLastIndex = true
	}
	res = ListReserveUserItemRes{
		List:        resData,
		Offset:      req.Offset + len(reserves),
		OnLastIndex: onLastIndex,
		MainOption:  req.MainOption,
		UserID:      tools.UuidToString(user.UserID),
	}
	return
}

// End Visited

func HandleReserveUserTime(arriveAfter string, arriveBefore string) (timeString string, timeType string) {
	if !tools.ServerStringEmpty(arriveAfter) {
		timeString = arriveAfter
		timeType = "arrive_after"
	} else if !tools.ServerStringEmpty(arriveBefore) {
		timeString = arriveBefore
		timeType = "arrive_before"
	} else {
		timeString = "Not set"
		timeType = "not_set"
	}
	return
}

func HandleGetRUOptionDirection(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveUserDirectionRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	log.Println("ticketID dir ", req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionDirection in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	d, err := server.store.GetChargeOptionReferenceDirection(ctx, db.GetChargeOptionReferenceDirectionParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUOptionDirection in GetChargeOptionReferenceDirection err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	var text string
	if HandleSqlNullString(d.Type) == "direction" {
		text = HandleSqlNullString(d.Info)
	}
	hasData = true
	res = ReserveUserDirectionRes{
		Street:    d.Street,
		City:      d.City,
		State:     d.State,
		Country:   d.Country,
		Postcode:  d.Postcode,
		Lat:       tools.ConvertFloatToLocationString(d.Geolocation.P.Y, 9),
		Lng:       tools.ConvertFloatToLocationString(d.Geolocation.P.X, 9),
		Direction: text,
	}
	return
}

func HandleGetRUEventDirection(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveUserDirectionRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	log.Println("ticketID d ", req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUEventDirection in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	d, err := server.store.GetChargeTicketReferenceDirection(ctx, db.GetChargeTicketReferenceDirectionParams{
		ID:         id,
		Type:       "direction",
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUEventDirection in GetChargeTicketReferenceDirection err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = ReserveUserDirectionRes{
		Street:    d.Street,
		City:      d.City,
		State:     d.State,
		Country:   d.Country,
		Postcode:  d.Postcode,
		Lat:       tools.ConvertFloatToLocationString(d.Geolocation.P.Y, 9),
		Lng:       tools.ConvertFloatToLocationString(d.Geolocation.P.X, 9),
		Direction: HandleSqlNullString(d.Info),
	}
	return
}

func HandleGetRUOptionCheckInStep(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res RUListCheckInStepRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	log.Println("ticketID c ", req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionCheckInStep in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	steps, err := server.store.GetChargeOptionReferenceCheckInStep(ctx, db.GetChargeOptionReferenceCheckInStepParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("Error at HandleGetRUOptionCheckInStep in GetChargeOptionReferenceCheckInStep err: %v, user: %v\n", err, user.ID)
		}
		hasData = false
		return
	}
	hasData = true
	var resData []RUCheckInStepRes
	for _, s := range steps {
		data := RUCheckInStepRes{
			ID:    tools.UuidToString(s.ID),
			Des:   s.Des,
			Photo: s.Photo,
		}
		resData = append(resData, data)
	}
	res = RUListCheckInStepRes{
		List: resData,
	}
	return
}

func HandleGetRUEventCheckInStep(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res RUListCheckInStepRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUEventCheckInStep in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	steps, err := server.store.GetChargeTicketReferenceCheckInStep(ctx, db.GetChargeTicketReferenceCheckInStepParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil || len(steps) == 0 {
		if err != nil {
			log.Printf("Error at HandleGetRUEventCheckInStep in GetChargeTicketReferenceCheckInStep err: %v, user: %v\n", err, user.ID)
		}
		hasData = false
		return
	}
	hasData = true
	var resData []RUCheckInStepRes
	for _, s := range steps {
		data := RUCheckInStepRes{
			ID:    tools.UuidToString(s.ID),
			Des:   s.Des,
			Photo: s.Photo,
		}
		resData = append(resData, data)
	}
	res = RUListCheckInStepRes{
		List: resData,
	}
	return
}

func HandleGetRUOptionHelp(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res RUHelpManualRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	info, err := server.store.GetChargeOptionReferenceHelp(ctx, db.GetChargeOptionReferenceHelpParams{
		ID:         id,
		Type:       "help_manual",
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in GetChargeOptionReferenceHelp err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = RUHelpManualRes{
		Help: info,
	}
	return
}

func HandleGetRUEventHelp(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res RUHelpManualRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUEventHelp in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	info, err := server.store.GetChargeTicketReferenceHelp(ctx, db.GetChargeTicketReferenceHelpParams{
		ID:         id,
		Type:       "help_manual",
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUEventHelp in GetChargeTicketReferenceHelp err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = RUHelpManualRes{
		Help: info,
	}
	return
}

func HandleGetRUEventReceipt(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveEventReceiptRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUEventReceipt in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	r, err := server.store.GetChargeTicketReferenceReceipt(ctx, db.GetChargeTicketReferenceReceiptParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUEventReceipt in GetChargeTicketReferenceReceipt err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = ReserveEventReceiptRes{
		Grade:          r.Grade,
		Price:          tools.IntToMoneyString(r.Price),
		Currency:       r.Currency,
		Type:           r.Type,
		TicketType:     r.TicketType,
		HostNameOption: r.HostNameOption,
	}
	return
}

func HandleGetRUOptionReceipt(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveOptionReceiptRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionReceipt in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	r, err := server.store.GetChargeOptionReferenceReceipt(ctx, db.GetChargeOptionReferenceReceiptParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUOptionReceipt in GetChargeOptionReferenceReceipt err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = ReserveOptionReceiptRes{
		HostNameOption:  r.HostNameOption,
		Discount:        r.Discount,
		MainPrice:       tools.IntToMoneyString(r.MainPrice),
		ServiceFee:      tools.IntToMoneyString(r.ServiceFee),
		TotalFee:        tools.IntToMoneyString(r.TotalFee),
		DatePrice:       r.DatePrice,
		Currency:        r.Currency,
		GuestFee:        tools.IntToMoneyString(r.GuestFee),
		PetFee:          tools.IntToMoneyString(r.PetFee),
		CleanFee:        tools.IntToMoneyString(r.CleanFee),
		NightlyPetFee:   tools.IntToMoneyString(r.NightlyPetFee),
		NightlyGuestFee: tools.IntToMoneyString(r.NightlyGuestFee),
	}
	return
}

func HandleGetRUOptionWifi(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveUserWifiRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	wifi, err := server.store.GetChargeOptionReferenceWifi(ctx, db.GetChargeOptionReferenceWifiParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in GetChargeOptionReferenceWifi err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = ReserveUserWifiRes{
		NetworkName: wifi.NetworkName,
		Password:    wifi.Password,
	}
	return
}

func HandleGetRUEventWifi(server *Server, ctx *gin.Context, user db.User, req ReserveUserInfoParams) (res ReserveUserWifiRes, hasData bool, err error) {
	id, err := tools.StringToUuid(req.ID)
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in StringToUuid err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	wifi, err := server.store.GetChargeTicketReferenceWifi(ctx, db.GetChargeTicketReferenceWifiParams{
		ID:         id,
		UserID:     user.UserID,
		IsComplete: true,
		Cancelled:  false,
	})
	if err != nil {
		log.Printf("Error at HandleGetRUOptionHelp in GetChargeTicketReferenceWifi err: %v, user: %v\n", err, user.ID)
		hasData = false
		return
	}
	hasData = true
	res = ReserveUserWifiRes{
		NetworkName: wifi.NetworkName,
		Password:    wifi.Password,
	}
	return
}
