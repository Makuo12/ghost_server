package api

import (
	"bytes"
	"encoding/json"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"time"

	"log"
	"strings"

	"github.com/google/uuid"
	//"github.com/go-redis/redis/v8"
)

func HandleSearchTextRes(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var search *SearchText = &SearchText{}
	err = json.Unmarshal(payload, search)
	if err != nil {
		log.Printf("error decoding HandleSearchTextRes response: %v, user: %v", err, ctx.username)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	searchName := "%" + search.Text + "%"
	optionDataList, err := ctx.server.store.ListOIDSearchByNameNoPhoto(ctx.ctx, db.ListOIDSearchByNameNoPhotoParams{
		HostNameOption: strings.ToLower(searchName),
		HostID:         ctx.username,
		CoUserID:       tools.UuidToString(ctx.userID),
		IsActive:       true,
	})
	if err != nil {
		log.Printf("error at HandleSearchTextRes at ListOIDSearchByName err:%v, user: %v \n", err, ctx.username)
		return
	}
	var resData []UHMOptionSelectionRes
	for _, item := range optionDataList {
		// i want to check if it has a photo
		var coverPhoto string
		var isActive bool
		var isCoHost bool
		if item.OptionStatus == "unlist" || item.OptionStatus == "snooze" {
			isActive = false
		} else {
			isActive = true
		}
		switch item.HostType {
		case "co_host":
			isCoHost = true
		case "main_host":
			isCoHost = false
		}

		d := UHMOptionSelectionRes{
			HostNameOption: item.HostNameOption,
			CoverImage:     coverPhoto,
			OptionID:       tools.UuidToString(item.ID),
			HasName:        true,
			MainOptionType: item.MainOptionType,
			IsComplete:     item.IsComplete,
			IsActive:       isActive,
			IsCoHost:       isCoHost,
		}
		resData = append(resData, d)
	}
	if len(resData) > 0 {
		res := SearchTextRes{
			List: resData,
		}
		resBytes := new(bytes.Buffer)

		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		data = resBytes.Bytes()
		hasData = true
		return
	}
	return
}

// Used for CalenderView [CalenderOptionItem]
func HandleSearchTextCalRes(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var search *SearchText = &SearchText{}
	err = json.Unmarshal(payload, search)
	if err != nil {
		log.Printf("error decoding HandleSearchTextRes response: %v, user: %v", err, ctx.username)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		log.Printf("HandleSearchTextRes response: %q", payload)
		return
	}
	log.Println("data formatted", search)
	log.Println("data formatted 1", search.Text)
	searchName := "%" + search.Text + "%"
	optionDataList, err := ctx.server.store.ListOIDSearchByNameCal(ctx.ctx, db.ListOIDSearchByNameCalParams{
		HostNameOption: strings.ToLower(searchName),
		ID:             ctx.username,
		IsActive:       true,
		IsComplete:     true,
	})
	if err != nil {
		log.Printf("error at HandleSearchTextCalRes at ListOIDSearchByNameCal err:%v, user: %v \n", err, ctx.username)
		return
	}
	var resData []CalenderOptionItem
	for _, item := range optionDataList {
		d := CalenderOptionItem{
			HostNameOption: item.HostNameOption,
			OptionID:       tools.UuidToString(item.ID),
			MainOptionType: item.MainOptionType,
		}
		resData = append(resData, d)
	}
	if len(resData) > 0 {
		res := SearchTextCalRes{
			List: resData,
		}
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

// EDT event_date_times
func HandleEDTSearchTextRes(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var search *SearchTextEDT = &SearchTextEDT{}
	err = json.Unmarshal(payload, search)
	if err != nil {
		log.Printf("error decoding HandleEDTSearchTextRes response: %v, user: %v", err, ctx.username)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		log.Printf("HandleEDTSearchTextRes response: %q", payload)
		return
	}
	eventID, err := tools.StringToUuid(search.EventInfoID)
	if err != nil {
		log.Printf("HandleEDTSearchTextRes response err at tools.StringToUuid: %q", payload)
		return
	}
	searchName := "%" + search.Text + "%"
	eventDateTimes, err := ctx.server.store.ListEDTSearchByName(ctx.ctx, db.ListEDTSearchByNameParams{
		Name:        searchName,
		ID:          ctx.username,
		EventInfoID: eventID,
	})

	if err != nil {
		log.Printf("error at HandleEDTSearchTextRes at ListEDTSearchByName err:%v, user: %v \n", err, ctx.username)
		return
	}
	var resData []EventDateItem
	for _, item := range eventDateTimes {
		var date EventDateItem
		detailIsEmpty := false
		var tickets int
		eventDateDetail, err := ctx.server.store.GetEventDateDetail(ctx.ctx, item.ID)
		if err != nil {
			log.Printf("Error at  ListEventDateItems in ListEventDateTime err: %v, user: %v\n", err, ctx.username)
			detailIsEmpty = true
		} //: IF End
		ticketCount, err := ctx.server.store.GetEventDateTicketCount(ctx.ctx, item.ID)
		// We check to error
		if err != nil {
			tickets = 0
			log.Printf("Error at  ListEventDateItems in ListEventDateTime err: %v, user: %v\n", err, ctx.username)
		} else {
			tickets = int(ticketCount)
		}
		log.Println("Setting up event date")
		if detailIsEmpty {
			log.Println("Setting up event empty event date")
			date = EventDateItem{
				ID:               tools.UuidToString(item.ID),
				Name:             item.Name,
				StartTime:        "",
				EndTime:          "",
				StartDate:        tools.ConvertDateOnlyToString(item.StartDate),
				Status:           item.Status,
				EndDate:          tools.ConvertDateOnlyToString(item.EndDate),
				Tickets:          tickets,
				Note:             "",
				TimeZone:         "",
				NeedBands:        item.NeedBands,
				NeedTickets:      item.NeedTickets,
				AbsorbBandCharge: item.AbsorbBandCharge,
			}
		} else {
			log.Println("Setting up event not empty event date")
			date = EventDateItem{
				ID:               tools.UuidToString(item.ID),
				Name:             item.Name,
				StartTime:        eventDateDetail.StartTime,
				EndTime:          eventDateDetail.EndTime,
				Status:           item.Status,
				StartDate:        tools.ConvertDateOnlyToString(item.StartDate),
				EndDate:          tools.ConvertDateOnlyToString(item.EndDate),
				Tickets:          tickets,
				Note:             item.Note,
				TimeZone:         eventDateDetail.TimeZone,
				NeedBands:        item.NeedBands,
				NeedTickets:      item.NeedTickets,
				AbsorbBandCharge: item.AbsorbBandCharge,
			}
		}
		resData = append(resData, date)
	}
	if len(resData) > 0 {
		res := SearchTextEDTRes{
			List: resData,
		}
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

func HandleSearchEventByName(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var search *SearchText = &SearchText{}
	err = json.Unmarshal(payload, search)
	if err != nil {
		log.Printf("error decoding HandleSearchTextRes response: %v, user: %v", err, ctx.username)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		log.Printf("HandleSearchTextRes response: %q", payload)
		return
	}
	log.Println("data formatted", search)
	log.Println("data formatted 1", search.Text)
	searchName := "%" + search.Text + "%"
	optionDataList, err := ctx.server.store.ListUserSearchEventByName(ctx.ctx, db.ListUserSearchEventByNameParams{
		HostNameOption: strings.ToLower(searchName),
		IsActive:       true,
		IsComplete:     true,
	})
	if err != nil {
		log.Printf("error at HandleSearchTextCalRes at ListOIDSearchByNameCal err:%v, user: %v \n", err, ctx.username)
		return
	}
	var resData []UserEventSearchItem
	for _, item := range optionDataList {
		d := UserEventSearchItem{
			HostNameOption: item.HostNameOption,
			OptionUserID:   tools.UuidToString(item.OptionUserID),
			IsVerified:     item.IsVerified,
			CoverImage:     item.CoverImage,
		}
		resData = append(resData, d)
	}
	if len(resData) > 0 {
		res := SearchEventByNameRes{
			List: resData,
		}
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

//func HandleMapExperienceLocation(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
//	var location *MapExperienceLocationParams = &MapExperienceLocationParams{}
//	err = json.Unmarshal(payload, location)
//	if err != nil {
//		log.Printf("error decoding HandleSearchTextRes response: %v, user: %v", err, ctx.username)
//		if e, ok := err.(*json.SyntaxError); ok {
//			log.Printf("syntax error at byte offset %d\n", e.Offset)
//		}
//		log.Printf("HandleSearchTextRes response: %q", payload)
//		return
//	}
//	log.Println("data formatted", location)
//	log.Println("data formatted 1", location.Lat)
//	lat := tools.ConvertLocationStringToFloat(location.Lat, 9)
//	lng := tools.ConvertLocationStringToFloat(location.Lng, 9)
//	if lat == 0.0 && lng == 0.0 {
//		return
//	}
//	query := &redis.GeoRadiusQuery{
//		Radius: 200,
//		Unit:   "km",
//		Sort:   "ASC",
//	}
//	rads, err := RedisClient.GeoRadius(RedisContext, constants.ALL_EXPERIENCE_LOCATION, lng, lat, query).Result()
//	if err != nil || len(rads) == 0 {
//		if err != nil {
//			log.Printf("error at HandleSearchTextCalRes at RedisClient.GeoRadius err:%v, user: %v \n", err, ctx.username)
//		}
//		return
//	}
//	var resData MapExperienceLocationItem
//	for _, r := range rads {
//		key := strings.Split(r.Name, "&")
//		if len(key) != 2 {
//			continue
//		}
//		log.Println("Keys", key)
//		mainOption := key[0]
//		optionUserID, err := tools.StringToUuid(key[1])
//		if err != nil {
//			log.Printf("error at HandleSearchTextCalRes at tools.StringToUuid(key[1] err:%v, user: %v \n", err, ctx.username)
//			continue
//		}
//		switch mainOption {
//		case constants.EXPERIENCE_OPTION_LOCATION:
//		case constants.EXPERIENCE_OPTION_LOCATION:
//			eventDate, err := ctx.server.store.GetEventDateTimeByUIDMap(ctx.ctx, db.GetEventDateTimeByUIDMapParams{
//				EventDateTimeID: optionUserID,
//				IsActive:        true,
//			})
//			if err != nil {
//				log.Printf("error at HandleSearchTextCalRes at GetEventDateTimeByUIDMap err:%v, user: %v \n", err, ctx.username)
//				continue
//			}
//		}
//	}
//	var resData []UserEventSearchItem
//	for _, item := range optionDataList {
//		d := UserEventSearchItem{
//			HostNameOption: item.HostNameOption,
//			OptionUserID:   tools.UuidToString(item.OptionUserID),
//			IsVerified:     item.IsVerified,
//			CoverImage:     item.CoverImage,
//		}
//		resData = append(resData, d)
//	}
//	if len(resData) > 0 {
//		res := SearchEventByNameRes{
//			List: resData,
//		}
//		resBytes := new(bytes.Buffer)

//		err = json.NewEncoder(resBytes).Encode(res)
//		if err != nil {
//			log.Println("error showing", err)
//		}
//		hasData = true
//		data = resBytes.Bytes()
//		return
//	}
//	return
//}

// The receiver listens for any incoming messages
func HandleMessageListen(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	//var resData []MessageContactItem
	var currentTime *CurrentTime = &CurrentTime{}
	err = json.Unmarshal(payload, currentTime)
	if err != nil {
		log.Printf("error decoding HandleMessageListen response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	log.Println("msg_contact time ", currentTime.Time)
	current, err := tools.ConvertStringToTime(currentTime.Time)
	if err != nil {
		log.Printf("error HandleMessageListen decoding ConvertStringToTime response: %v, user: %v", err, ctx.userID)
		return
	}
	contacts, err := ctx.server.store.ListMessageContactByTime(ctx.ctx, db.ListMessageContactByTimeParams{
		SenderID:  ctx.userID,
		CreatedAt: current,
	})
	if err != nil {
		log.Printf("error HandleMessageListen decoding ListMessageContactByTime response: %v, user: %v", err, ctx.userID)
		return
	}
	var resData []MessageContactItem
	for _, contact := range contacts {
		contactData := MessageContactItem{
			MsgID:                      tools.UuidToString(contact.MessageID),
			ConnectedUserID:            tools.UuidToString(contact.ConnectedUserID),
			FirstName:                  contact.FirstName,
			Photo:                      contact.Photo,
			LastMessage:                contact.LastMessage,
			LastMessageTime:            tools.ConvertTimeToString(contact.LastMessageTime),
			UnreadMessageCount:         int(contact.UnreadMessageCount),
			UnreadUserRequestCount:     int(contact.UnreadUserRequestCount),
			UnreadUserCancelCount:      int(contact.UnreadUserCancelCount),
			UnreadHostCancelCount:      int(contact.UnreadHostCancelCount),
			UnreadHostChangeDatesCount: int(contact.UnreadHostChangeDatesCount),
		}
		resData = append(resData, contactData)
	}
	if len(resData) > 0 {

		res := ListMessageContactListenRes{
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

func HandleUnreadMessageList(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var unreadMsg *UnreadMessageParams = &UnreadMessageParams{}
	err = json.Unmarshal(payload, unreadMsg)
	if err != nil {
		log.Printf("error decoding HandleUnreadMessageList response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	contactID, err := tools.StringToUuid(unreadMsg.SelectedContactID)
	if err != nil {
		log.Printf("error decoding HandleUnreadMessageList contact tools.StringToUuid(id): %v, user: %v, id: %v", err, ctx.userID, unreadMsg.SelectedContactID)
		return
	}
	userID, err := tools.StringToUuid(unreadMsg.UserID)
	if err != nil {
		log.Printf("error decoding HandleUnreadMessageList user tools.StringToUuid(id): %v, user: %v, id: %v", err, ctx.userID, unreadMsg.SelectedContactID)
		return
	}
	// Update in database
	updatedIDs, err := ctx.server.store.UpdateMessageRead(ctx.ctx, db.UpdateMessageReadParams{
		Read:       true,
		SenderID:   contactID,
		ReceiverID: userID,
	})
	if err != nil {
		log.Printf("error decoding HandleUnreadMessageList user server.store.UpdateMessageRead(id): %v, user: %v", err, ctx.userID)
	}
	readIDs := tools.ListUuidToString(updatedIDs)
	if len(readIDs) == 0 {
		readIDs = []string{"none"}
	}
	if len(readIDs) == 0 {
		return
	} else {
		res := UnreadMessageRes{
			List:              readIDs,
			UserID:            unreadMsg.UserID,
			SelectedContactID: unreadMsg.SelectedContactID,
		}
		resBytes := new(bytes.Buffer)
		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		hasData = true
		data = resBytes.Bytes()
		return
	}
}

func HandleGetMessage(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var getMessage *GetMessageParams = &GetMessageParams{}
	err = json.Unmarshal(payload, getMessage)
	if err != nil {
		log.Printf("error decoding HandleGetMessage response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	messages, msgIsEmpty := HandleMainMessage(ctx, getMessage.LatestTime, getMessage.SelectedContactID, getMessage.UserID, "HandleGetMessage")
	if msgIsEmpty {
		return
	} else {
		res := GetMessageRes{
			MsgList:           messages,
			MsgEmpty:          msgIsEmpty,
			UserID:            getMessage.UserID,
			SelectedContactID: getMessage.SelectedContactID,
		}
		resBytes := new(bytes.Buffer)
		err = json.NewEncoder(resBytes).Encode(res)
		if err != nil {
			log.Println("error showing", err)
		}
		hasData = true
		data = resBytes.Bytes()
		return
	}
}

func HandleMainMessage(ctx *connection, msgTimeString string, contactIDString string, userIDString string, funcName string) ([]MessageMainItem, bool) {
	emptyData := []MessageMainItem{{MessageItem{"none", "none", "none", "none", "none", true, "none", "none", "none", "none"}, true, MessageItem{"none", "none", "none", "none", "none", true, "none", "none", "none", "none"}, true}}

	msgTime, err := tools.ConvertStringToTime(msgTimeString)
	if err != nil {
		log.Printf("error msgTime at funcName %v, error at HandleMainMessage tools.StringToUuid error: %v, user: %v, contact: %v", funcName, err, ctx.userID, contactIDString)
		return emptyData, true
	}
	contactID, err := tools.StringToUuid(contactIDString)
	if err != nil {
		log.Printf("error contactID at funcName %v, error at HandleMainMessage tools.StringToUuid error: %v, user: %v, contact: %v", funcName, err, ctx.userID, contactIDString)
		return emptyData, true
	}
	userID, err := tools.StringToUuid(userIDString)
	if err != nil {
		log.Printf("error at userID funcName %v, error at HandleMainMessage tools.StringToUuid error: %v, user: %v, contact: %v", funcName, err, ctx.userID, contactIDString)
		return emptyData, true
	}
	messages, err := ctx.server.store.ListMessageWithTime(ctx.ctx, db.ListMessageWithTimeParams{
		SenderID:     userID,
		ReceiverID:   contactID,
		SenderID_2:   contactID,
		ReceiverID_2: userID,
		CreatedAt:    msgTime,
	})
	if err != nil || len(messages) == 0 {
		if err != nil {
			log.Printf("error at userID funcName %v, error at HandleMainMessage ListMessageWithTime error: %v, user: %v, contact: %v", funcName, err, ctx.userID, contactIDString)
		}
		return emptyData, true
	}
	var resData []MessageMainItem
	for _, msg := range messages {
		var mainMsg MessageItem
		var parentMsg MessageItem
		var parentEmpty bool
		mainMsg = MessageItem{
			ID:         tools.UuidToString(msg.MsgID),
			SenderID:   tools.UuidToString(msg.SenderID),
			ReceiverID: tools.UuidToString(msg.ReceiverID),
			Message:    msg.Message,
			Type:       msg.Type,
			Read:       msg.Read,
			Photo:      msg.Photo,
			ParentID:   msg.ParentID,
			Reference:  msg.Reference,
			CreatedAt:  tools.ConvertTimeToString(msg.CreatedAt),
		}
		if msg.MainParentID.Valid {
			parentMsg = MessageItem{
				ID:         tools.UuidToString(HandleSqlNullUUID(msg.ParentMsgID)),
				SenderID:   tools.UuidToString(HandleSqlNullUUID(msg.ParentSenderID)),
				ReceiverID: tools.UuidToString(HandleSqlNullUUID(msg.ParentReceiverID)),
				Message:    HandleSqlNullString(msg.ParentMessage),
				Type:       HandleSqlNullString(msg.ParentType),
				Read:       HandleSqlNullBool(msg.ParentRead),
				Photo:      HandleSqlNullString(msg.ParentPhoto),
				ParentID:   HandleSqlNullString(msg.ParentParentID),
				Reference:  HandleSqlNullString(msg.ParentReference),
				CreatedAt:  tools.ConvertTimeToString(HandleSqlNullTimestamp(msg.ParentCreatedAt)),
			}
			parentEmpty = false
		} else {
			parentMsg = MessageItem{
				ID:         "none",
				SenderID:   "none",
				ReceiverID: "none",
				Message:    "none",
				Type:       "none",
				Read:       false,
				Photo:      "none",
				ParentID:   "none",
				Reference:  "none",
				CreatedAt:  "none",
			}
			parentEmpty = true
		}
		msgData := MessageMainItem{
			MainMsg:        mainMsg,
			MainMsgEmpty:   false,
			ParentMsg:      parentMsg,
			ParentMsgEmpty: parentEmpty,
		}
		resData = append(resData, msgData)
	}
	return resData, false
}

// The sender creates the message and the message is sent to both the sender and also to the listener (if the listener is in the message room currently)
func HandleMessage(ctx *connection, payload []byte) (data []byte, hasData bool, err error) {
	var msg *CreateMessageParams = &CreateMessageParams{}
	err = json.Unmarshal(payload, msg)
	if err != nil {
		log.Printf("error decoding HandleMessage response: %v, user: %v", err, ctx.userID)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		return
	}
	switch msg.ForUnlist {
	case true:
		contactID, errData := tools.StringToUuid(msg.SelectedContactID)
		if errData != nil {
			log.Printf("error decoding HandleUnreadMessageList contact tools.StringToUuid(id): %v, user: %v, id: %v", err, ctx.userID, msg.SelectedContactID)
			return
		}
		userID, errData := tools.StringToUuid(msg.UserID)
		if errData != nil {
			log.Printf("error decoding HandleUnreadMessageList user tools.StringToUuid(id): %v, user: %v, id: %v", err, ctx.userID, msg.SelectedContactID)
			return
		}
		// Update in database
		updatedIDs, errData := ctx.server.store.UpdateMessageRead(ctx.ctx, db.UpdateMessageReadParams{
			Read:       true,
			SenderID:   contactID,
			ReceiverID: userID,
		})
		if errData != nil {
			log.Printf("error decoding HandleUnreadMessageList user server.store.UpdateMessageRead(id): %v, user: %v", err, ctx.userID)
			return
		}
		readIDs := tools.ListUuidToString(updatedIDs)
		if len(readIDs) != 0 {
			res := UnreadMessageRes{
				List:              readIDs,
				UserID:            msg.UserID,
				SelectedContactID: msg.SelectedContactID,
			}
			resBytes := new(bytes.Buffer)
			err = json.NewEncoder(resBytes).Encode(res)
			if err != nil {
				log.Println("error showing", err)
			}
			hasData = true
			data = resBytes.Bytes()
		}
	case false:
		senderID, errData := tools.StringToUuid(msg.SenderID)
		if errData != nil {
			log.Printf("error decoding HandleMessage tools.StringToUuid(msg.SenderID) response: %v, user: %v", err, ctx.userID)
			return
		}
		receiverID, errData := tools.StringToUuid(msg.ReceiverID)
		if errData != nil {
			log.Printf("error decoding HandleMessage tools.StringToUuid(msg.ReceiverID response: %v, user: %v", err, ctx.userID)
			return
		}
		msgTime := time.Now().Add(time.Hour)
		msgD, errData := ctx.server.store.CreateMessage(ctx.ctx, db.CreateMessageParams{
			MsgID:      uuid.New(),
			SenderID:   senderID,
			ReceiverID: receiverID,
			Message:    msg.Message,
			Type:       msg.Type,
			Photo:      msg.Photo,
			Read:       false,
			ParentID:   msg.ParentID,
			Reference:  msg.Reference,
			CreatedAt:  msgTime,
			UpdatedAt:  msgTime,
		})
		if errData != nil {
			log.Printf("error decoding HandleMessage ctx.server.store.CreateMessage response: %v, user: %v", err, ctx.userID)
			return
		}
		msgRes := MessageItem{
			ID:         tools.UuidToString(msgD.MsgID),
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			Message:    msg.Message,
			Type:       msg.Type,
			Read:       false,
			Photo:      msg.Photo,
			ParentID:   msg.ParentID,
			Reference:  msg.Reference,
			CreatedAt:  tools.ConvertTimeToString(msgTime),
		}
		resBytes := new(bytes.Buffer)
		log.Println("resData handle message", msgRes)
		err = json.NewEncoder(resBytes).Encode(msgRes)
		if errData != nil {
			log.Println("error showing", err)
		}
		hasData = true
		data = resBytes.Bytes()
	}
	return
}
