package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) ListMessageContact(ctx *gin.Context) {
	var req ListMessageContactParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListMessageContact in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.GetMessageContactCount(ctx, user.UserID)
	if err != nil {
		log.Printf("Error at  ListMessageContact in GetMessageContactCount err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	contacts, err := server.store.ListMessageContact(ctx, db.ListMessageContactParams{
		SenderID: user.UserID,
		Limit:    10,
		Offset:   int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListMessageContact in ListMessageContact err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var res ListMessageContactRes
	var resData []MessageContactItem
	for _, contact := range contacts {
		roomID, err := SingleContextRoom(ctx, server, user.UserID, contact.ConnectedUserID, "ListMessage")
		if err != nil {
			continue
		}
		data := MessageContactItem{
			MsgID:                      tools.UuidToString(contact.MessageID),
			ConnectedUserID:            tools.UuidToString(contact.ConnectedUserID),
			FirstName:                  contact.FirstName,
			MainImage:                  contact.Image,
			LastMessage:                contact.LastMessage,
			LastMessageTime:            tools.ConvertTimeToString(contact.LastMessageTime),
			UnreadMessageCount:         int(contact.UnreadMessageCount),
			UnreadUserRequestCount:     int(contact.UnreadUserRequestCount),
			UnreadUserCancelCount:      int(contact.UnreadUserCancelCount),
			UnreadHostCancelCount:      int(contact.UnreadHostCancelCount),
			UnreadHostChangeDatesCount: int(contact.UnreadHostChangeDatesCount),
			RoomID:                     tools.UuidToString(roomID),
		}
		resData = append(resData, data)
	}
	var timeString string = tools.ConvertTimeToString(time.Now().Add(time.Hour))
	if len(resData) > 0 {
		timeString = resData[0].LastMessageTime
	}
	if count <= int64(req.Offset+len(contacts)) {
		onLastIndex = true
	}
	res = ListMessageContactRes{
		List:        resData,
		Offset:      req.Offset + len(contacts),
		OnLastIndex: onLastIndex,
		UserID:      tools.UuidToString(user.UserID),
		Time:        timeString,
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListMessage(ctx *gin.Context) {
	var req ListMessageParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListMessage in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	contactID, err := tools.StringToUuid(req.ContactID)
	if err != nil {
		log.Printf("Error at  ListMessage in GetMessageContactCount err: %v, contactID: %v\n", err, req.ContactID)
		err = fmt.Errorf("this contact does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	emptyData := MessageMainItem{MessageItem{"none", "none", "none", "none", "none", true, "none", "none", "none", "none"}, true, MessageItem{"none", "none", "none", "none", "none", true, "none", "none", "none", "none"}, true}
	roomID, err := SingleContextRoom(ctx, server, user.UserID, contactID, "ListMessage")
	if err != nil {
		err = fmt.Errorf("could not setup your room, please try again")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	log.Println("room id: ", roomID)
	emptyMsg := ListMessageRes{
		List:              []MessageMainItem{emptyData},
		Offset:            0,
		OnLastIndex:       false,
		SelectedContactID: req.ContactID,
		RoomID:            tools.UuidToString(roomID),
		IsEmpty:           true,
	}
	count, err := server.store.GetMessageCount(ctx, db.GetMessageCountParams{
		SenderID:     user.UserID,
		ReceiverID:   contactID,
		SenderID_2:   contactID,
		ReceiverID_2: user.UserID,
	})
	if err != nil {
		log.Printf("Error at ListMessage in GetMessageCount err: %v, user: %v\n", err, user.ID)
		ctx.JSON(http.StatusOK, emptyMsg)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		ctx.JSON(http.StatusOK, emptyMsg)
		return
	}
	messages, err := server.store.ListMessage(ctx, db.ListMessageParams{
		SenderID:     user.UserID,
		ReceiverID:   contactID,
		SenderID_2:   contactID,
		ReceiverID_2: user.UserID,
		Limit:        8,
		Offset:       int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListMessage in ListMessage err: %v, user: %v\n", err, user.ID)
		ctx.JSON(http.StatusOK, emptyMsg)
		return
	}
	var res ListMessageRes
	var resData []MessageMainItem
	if count > int64(req.Offset) {
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
				Image:      msg.MainImage,
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
					Image:      HandleSqlNullString(msg.ParentMainImage),
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
					Image:      "none",
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
	}
	if count <= int64(req.Offset+len(messages)) {
		onLastIndex = true
	}
	res = ListMessageRes{
		List:              resData,
		Offset:            req.Offset + len(messages),
		OnLastIndex:       onLastIndex,
		SelectedContactID: req.ContactID,
		RoomID:            tools.UuidToString(roomID),
		IsEmpty:           false,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateMessage(ctx *gin.Context) {
	var req CreateMessageParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateMessage in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.SenderID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	senderID, err := tools.StringToUuid(req.SenderID)
	if err != nil {
		log.Printf("error decoding CreateMessage tools.StringToUuid(req.SenderID) response: %v, user: %v", err, user.UserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	receiverID, err := tools.StringToUuid(req.ReceiverID)
	if err != nil {
		log.Printf("error decoding CreateMessage tools.StringToUuid(req.ReceiverID response: %v, user: %v", err, user.UserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msgTime := time.Now().Add(time.Hour)
	msgD, err := server.store.CreateMessage(ctx, db.CreateMessageParams{
		MsgID:      uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    req.Message,
		Type:       req.Type,
		MainImage:  req.Image,
		Read:       false,
		ParentID:   req.ParentID,
		Reference:  req.Reference,
		CreatedAt:  msgTime,
		UpdatedAt:  msgTime,
	})
	if err != nil {
		log.Printf("error decoding CreateMessage server.store.CreateMessage response: %v, user: %v", err, user.UserID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We a notification to the receiver
	HandleUserIdMessageApn(ctx, server, receiverID, req.Message, user.FirstName)
	res := CreateMessageRes{
		MsgID:     tools.UuidToString(msgD.MsgID),
		CreatedAt: tools.ConvertTimeToString(msgTime),
	}
	ctx.JSON(http.StatusOK, res)
}
