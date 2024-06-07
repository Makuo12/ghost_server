package api

import (
	"log"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HandleListMsgRedisDB this list all the messages in redis in the MessageContactItem format
func HandleListMsgRedisDB(messageIDs []string, server *Server, ctx *gin.Context, user db.User) (resData []MessageContactItem) {
	for _, id := range messageIDs {
		msg, err := RedisClient.HGetAll(RedisContext, id).Result()
		if err != nil {
			log.Printf("error at HandleListMsgRedisDB at HGetAll err:%v, user: %v , id: %v \n", err.Error(), user.ID, id)
			continue
		} else {
			var connectedUserID uuid.UUID
			var messageCount int
			senderID, err := tools.StringToUuid(msg[constants.SENDER_ID])
			if err != nil {
				log.Printf("error at HandleListMsgRedisDB at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), user.ID, id)
				continue
			}
			receiverID, err := tools.StringToUuid(msg[constants.RECEIVER_ID])
			if err != nil {
				log.Printf("error at HandleListMsgRedisDB at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), user.ID, id)
				continue
			}
			// Lets us know if the current user is the sender or receiver
			// If senderID == user.UserID then we know it is the user that sent the message
			if senderID == user.UserID {
				connectedUserID = receiverID
			} else {
				connectedUserID = senderID
				messageCount = 1
			}
			connectedUser, err := server.store.GetUserByUserID(ctx, connectedUserID)
			if err != nil {
				log.Printf("error at HandleListMsgRedisDB at GetUser err:%v, user: %v , id: %v \n", err.Error(), user.ID, id)
				continue
			}
			userRequestCount, userCancelCount, hostCancelCount, hostChangeDatesCount := MessageType(msg[constants.TYPE])
			contact := MessageContactItem{
				ConnectedUserID:            tools.UuidToString(connectedUserID),
				FirstName:                  connectedUser.FirstName,
				Photo:                      connectedUser.Photo,
				LastMessage:                msg[constants.MESSAGE],
				LastMessageTime:            msg[constants.CREATED_AT],
				UnreadMessageCount:         messageCount,
				UnreadUserRequestCount:     userRequestCount,
				UnreadUserCancelCount:      userCancelCount,
				UnreadHostCancelCount:      hostCancelCount,
				UnreadHostChangeDatesCount: hostChangeDatesCount,
			}
			resData = append(resData, contact)
		}
	}
	return

}

func HandleListMessageRedis(messageIDs []string, contactID string, server *Server, ctx *gin.Context, user db.User) (resData []MessageItem) {
	for _, id := range messageIDs {
		msg, err := RedisClient.HGetAll(RedisContext, id).Result()
		if err != nil {
			log.Printf("error at HandleListMessageRedis at HGetAll err:%v, user: %v , id: %v \n", err.Error(), user.ID, id)
			continue
		} else {
			userID := tools.UuidToString(user.UserID)
			senderID := msg[constants.SENDER_ID]
			receiverID := msg[constants.RECEIVER_ID]
			if (senderID == contactID && receiverID == userID) || (senderID == userID && receiverID == contactID) {
				messageItem := MessageItem{
					ID:         msg[constants.ID],
					SenderID:   senderID,
					ReceiverID: receiverID,
					Message:    msg[constants.MESSAGE],
					Type:       msg[constants.TYPE],
					Read:       tools.ConvertStringToBool(msg[constants.READ]),
					Photo:      msg[constants.PHOTO],
					ParentID:   msg[constants.PARENT_ID],
					Reference:  msg[constants.REFERENCE],
					CreatedAt:  msg[constants.CREATED_AT],
				}
				resData = append(resData, messageItem)
			}
		}
	}
	return
}

func MessageType(msgType string) (userRequestCount int, userCancelCount int, hostCancelCount int, hostChangeDatesCount int) {
	switch msgType {
	case "user_request":
		userRequestCount = 1
	case "user_cancel":
		userCancelCount = 1
	case "host_cancel":
		hostCancelCount = 1
	case "host_change_dates":
		hostChangeDatesCount = 1
	}
	return
}

//func HandleMsgToRedis(msg *CreateMessageParams, ctx *connection) (msgData MessageItem) {
//	msgID := tools.UuidToString(uuid.New())
//	createdAt := tools.ConvertTimeToString(time.Now().UTC())
//	// We want to store it in redis
//	data := []string{
//		constants.ID,
//		msgID,
//		constants.SENDER_ID,
//		msg.SenderID,
//		constants.RECEIVER_ID,
//		msg.ReceiverID,
//		constants.MESSAGE,
//		msg.Message,
//		constants.TYPE,
//		msg.Type,
//		constants.READ,
//		tools.ConvertBoolToString(false),
//		constants.PHOTO,
//		msg.Photo,
//		constants.PARENT_ID,
//		msg.ParentID,
//		constants.REFERENCE,
//		msg.Reference,
//		constants.STORED_IN_DB,
//		tools.ConvertBoolToString(false),
//		constants.CREATED_AT,
//		createdAt,
//	}
//	err := RedisClient.HSet(RedisContext, msgID, data).Err()
//	if err != nil {
//		log.Printf("error at HandleMsgToRedis at HGetAll err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, msgID)
//	}
//	// This is used for reference when we need easy access to the data for the user
//	err = RedisClient.SAdd(RedisContext, msg.ReceiverID, msgID).Err()
//	if err != nil {
//		log.Printf("error at HandleMsgToRedis Receiver at SAdd err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, msgID)
//	}
//	// This is used for reference when we need easy access to the data for the user
//	err = RedisClient.SAdd(RedisContext, msg.SenderID, msgID).Err()
//	if err != nil {
//		log.Printf("error at HandleMsgToRedis Sender at SAdd err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, msgID)
//	}
//	err = RedisClient.SAdd(RedisContext, constants.MESSAGE_RECEIVE, msgID).Err()
//	if err != nil {
//		log.Printf("error at HandleMsgToRedis at SAdd err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, msgID)
//	}
//	msgData = MessageItem{
//		ID:         msgID,
//		SenderID:   msg.SenderID,
//		ReceiverID: msg.ReceiverID,
//		Message:    msg.Message,
//		Type:       msg.Type,
//		Read:       false,
//		Photo:      msg.Photo,
//		ParentID:   msg.ParentID,
//		Reference:  msg.Reference,
//		CreatedAt:  createdAt,
//	}
//	return
//	// We want to create a function that sets a time to the db
//}

// stand means data not from websocket
//func HandleMsgStandToRedis(msg CreateMessageParams, server *Server, ctx context.Context, user db.User, storeInDB bool, msgID string, createdAt time.Time) {
//	//createdAt := tools.ConvertTimeToString(time.Now())
//	// We want to store it in redis
//	data := []string{
//		constants.ID,
//		msgID,
//		constants.SENDER_ID,
//		msg.SenderID,
//		constants.RECEIVER_ID,
//		msg.ReceiverID,
//		constants.MESSAGE,
//		msg.Message,
//		constants.TYPE,
//		msg.Type,
//		constants.READ,
//		tools.ConvertBoolToString(false),
//		constants.PHOTO,
//		msg.Photo,
//		constants.PARENT_ID,
//		msg.ParentID,
//		constants.REFERENCE,
//		msg.Reference,
//		constants.STORED_IN_DB,
//		tools.ConvertBoolToString(storeInDB),
//		constants.CREATED_AT,
//		tools.ConvertTimeToString(createdAt),
//	}
//	// We store data in h-set
//	err := RedisClient.HSet(RedisContext, msgID, data).Err()
//	if err != nil {
//		log.Printf("error at  HandleMsgStandToRedis at HGetAll err:%v, user: %v , id: %v \n", err.Error(), user.ID, msgID)
//	}
//	// We then store reference to the data in s-add with the senderID is the key
//	err = RedisClient.SAdd(RedisContext, msg.ReceiverID, msgID).Err()
//	if err != nil {
//		log.Printf("error at  HandleMsgStandToRedis at SAdd err:%v, user: %v , id: %v \n", err.Error(), user.ID, msgID)
//	}
//	err = RedisClient.SAdd(RedisContext, msg.SenderID, msgID).Err()
//	if err != nil {
//		log.Printf("error at  HandleMsgStandToRedis at SAdd err:%v, user: %v , id: %v \n", err.Error(), user.ID, msgID)
//	}
//	// Then we add that sender id to redis
//	err = RedisClient.SAdd(RedisContext, constants.MESSAGE_RECEIVE, msgID).Err()
//	if err != nil {
//		log.Printf("error at  HandleMsgStandToRedis at SAdd err:%v, user: %v , id: %v \n", err.Error(), user.ID, msgID)
//	}
//}

//func HandleAllRedisMsgToDB(ctx context.Context, server *Server) func() {
//	// First we list all the messages
//	return func() {
//		msgIDs, err := RedisClient.SMembers(RedisContext, constants.MESSAGE_RECEIVE).Result()
//		if err != nil {
//			log.Printf("error at HandleAllRedisMsgToDB at SMembers(RedisContext, constants.MESSAGE_RECEIVE). err:%v\n", err.Error())
//			return
//		}
//		for _, id := range msgIDs {
//			var senderID string
//			var receiverID string
//			msg, err := RedisClient.HGetAll(RedisContext, id).Result()
//			if err != nil {
//				log.Printf("error at HandleRedisMsgToDB at HGetAll err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//				continue
//			} else {
//				senderID = msg[constants.SENDER_ID]
//				receiverID = msg[constants.RECEIVER_ID]
//				if !tools.ConvertStringToBool(msg[constants.STORED_IN_DB]) {
//					// If this msg is store in the database we want to skip
//					mainID, err := tools.StringToUuid(id)
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//					receiverUUID, err := tools.StringToUuid(receiverID)
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB receiver at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//					senderUUID, err := tools.StringToUuid(senderID)
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB at sender StringToUuid err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//					createdAt, err := tools.ConvertStringToTime(msg[constants.CREATED_AT])
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB at ConvertStringToTime err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//					_, err = server.store.CreateMessage(ctx, db.CreateMessageParams{
//						MsgID:      mainID,
//						SenderID:   senderUUID,
//						ReceiverID: receiverUUID,
//						Message:    msg[constants.MESSAGE],
//						Type:       msg[constants.TYPE],
//						Read:       tools.ConvertStringToBool(msg[constants.READ]),
//						Photo:      msg[constants.PHOTO],
//						ParentID:   msg[constants.PARENT_ID],
//						Reference:  msg[constants.REFERENCE],
//						CreatedAt:  createdAt,
//						UpdatedAt:  createdAt,
//					})
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB at CreateMessage err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//					_, err = SingleContextRoom(ctx, server, senderUUID, receiverUUID, "HandleRedisMsgToDB")
//					if err != nil {
//						log.Printf("error at HandleRedisMsgToDB at SingleContextRoom err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//						continue
//					}
//				}

//				// We want to remove it from redis
//				err = RedisClient.Del(RedisContext, id).Err()
//				if err != nil {
//					log.Printf("error at HandleRedisMsgToDB at RedisClient.Del err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//				}
//				// We want to remove it from SMembers sender
//				err = RedisClient.SRem(RedisContext, senderID, id).Err()
//				if err != nil {
//					log.Printf("error at HandleRedisMsgToDB at SRem for message ID err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//				}
//				// We want to remove it from SMembers receiver
//				err = RedisClient.SRem(RedisContext, receiverID, id).Err()
//				if err != nil {
//					log.Printf("error at HandleRedisMsgToDB at SRem for message ID err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//				}

//			}
//			// We remove it from constants.MESSAGE_RECEIVE
//			err = RedisClient.SRem(RedisContext, constants.MESSAGE_RECEIVE, id).Err()
//			if err != nil {
//				log.Printf("error at HandleRedisMsgToDB at SRem for message ID err:%v, user: %v , id: %v \n", err.Error(), ctx, id)
//			}
//		}
//	}
//}
