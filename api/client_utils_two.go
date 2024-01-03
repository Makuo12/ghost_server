package api

import (
	"bytes"
	"encoding/json"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
)

//func HandleHostReserve(ctx *connection, payload []byte) (data []byte, err error) {
//	var reserve *ListReservationDetailParams = &ListReservationDetailParams{}
//	err = json.Unmarshal(payload, reserve)
//	if err != nil {
//		log.Printf("error decoding HandleHostReserve response: %v, user: %v", err, ctx.username)
//		if e, ok := err.(*json.SyntaxError); ok {
//			log.Printf("syntax error at byte offset %d\n", e.Offset)
//		}
//		log.Printf("HandleHostReserve response: %q", payload)
//		return
//	}
//	switch reserve.MainOption {
//	case "options":
//		mainHostRes, hasHostData, err := HandleReserveHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  HandleHostReserve in HandleReserveHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasHostData = false
//		}
//		coHostRes, hasCoData, err := HandleReserveCoHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  HandleHostReserve in HandleReserveCoHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasCoData = false
//		}
//		var list []ReserveHostItem
//		if hasCoData && hasHostData {
//			// Has co host data and main host data
//			list = ConcatSlicesReserveItem(mainHostRes, coHostRes)

//		} else if hasHostData && !hasCoData {
//			// Has main host data but no co host data
//			list = mainHostRes
//		} else if hasCoData && !hasHostData {
//			// Has co host data but not main host data
//			list = coHostRes
//		}
//		reserveIDs := tools.HandleListReq(reserve.ReferenceIDs)
//		list = HandleHostReserveOptionSelected(reserveIDs, list)
//		if len(list) > 0 {
//			res := ListReservationDetailRes{
//				List:      list,
//				Selection: reserve.Selection,
//			}
//			resBytes := new(bytes.Buffer)
//			json.NewEncoder(resBytes).Encode(res)
//			data = resBytes.Bytes()
//		}
//	case "events":
//		mainHostRes, hasHostData, err := HandleReserveEventHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  ListReservationDetail in HandleReserveEventHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasHostData = false
//		}
//		coHostRes, hasCoData, err := HandleReserveEventCoHost(reserve.Selection, ctx.server, ctx.ctx, ctx.user)
//		if err != nil {
//			log.Printf("Error at  ListReservationDetail in HandleReserveEventCoHost err: %v, user: %v\n", err, ctx.user.ID)
//			hasCoData = false
//		}
//		var list []DateHostItem
//		if hasCoData && hasHostData {
//			// Has co host data and main host data
//			list = ConcatSlicesDateItem(mainHostRes, coHostRes)

//		} else if hasHostData && !hasCoData {
//			// Has main host data but no co host data
//			list = mainHostRes
//		} else if hasCoData && !hasHostData {
//			// Has co host data but not main host data
//			list = coHostRes
//		}
//		if len(list) > 0 {
//			res := ReserveEventHostItem{
//				List:      list,
//				Selection: reserve.Selection,
//			}
//			resBytes := new(bytes.Buffer)
//			json.NewEncoder(resBytes).Encode(res)
//			data = resBytes.Bytes()
//		}
//	}
//	return
//}

func HandleHostReserveOptionSelected(referenceIDs []string, list []ReserveHostItem) (res []ReserveHostItem) {
	for _, item := range list {
		exist := false
		for _, id := range referenceIDs {
			if item.ReferenceID == id {
				exist = true
			}
		}
		if !exist {
			res = append(res, item)
		}

	}
	return
}

func HandleNotificationListen(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	//var resData []MessageContactItem
	var currentTime *CurrentTime = &CurrentTime{}
	err = json.Unmarshal(payload, currentTime)
	if err != nil {
		log.Printf("error decoding HandleNotificationListen response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	current, err := tools.ConvertStringToTime(currentTime.Time)
	if err != nil {
		log.Printf("error HandleNotificationListen decoding ConvertStringToTime response: %v, user: %v", err, ctx.userID)
		return
	}
	// We first check the redis storage
	notifications, err := ctx.server.store.ListNotificationByTime(ctx.ctx, db.ListNotificationByTimeParams{
		UserID:    ctx.userID,
		CreatedAt: current,
	})
	if err != nil {
		log.Printf("error at HandleNotificationListen at ListNotificationByTime err:%v, user: %v \n", err, ctx.userID)
		return
	}
	var resData []NotificationItem
	for _, n := range notifications {
		dataNotification := NotificationItem{
			ID:        tools.UuidToString(n.ID),
			Type:      n.Type,
			Header:    n.Header,
			Message:   n.Message,
			Handled:   n.Handled,
			CreatedAt: tools.ConvertTimeToString(n.CreatedAt),
		}
		resData = append(resData, dataNotification)
	}
	if len(resData) > 0 {
		res := ListNotificationListenRes{
			List:        resData,
			CurrentTime: tools.ConvertTimeToString(current),
		}
		//log.Println("res", res)
		resBytes := new(bytes.Buffer)

		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		hasData = true
		data = resBytes.Bytes()
		return
	}
	return
}
