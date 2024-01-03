package api

// HandleListMsgRedisWebsocket this list messages for websocket
//func HandleListMsgRedisWebsocket(messageIDs []string, ctx *connection, current time.Time, userID uuid.UUID) (resData []MessageContactItem) {
//	for _, id := range messageIDs {
//		msg, err := RedisClient.HGetAll(RedisContext, id).Result()
//		if tools.ConvertStringToBool(msg[constants.READ]) {
//			continue
//		}
//		if err != nil {
//			log.Printf("error at HandleListMsgRedis at HGetAll err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, id)
//			continue
//		} else {
//			var connectedUserID uuid.UUID
//			var messageCount int
//			// First we check if the time is greater than the current time
//			msgTime, err := tools.ConvertStringToTime(msg[constants.CREATED_AT])
//			if err != nil {
//				log.Printf("error at HandleListMsgRedis at ConvertStringToTime err:%v, user: %v , id: %v \n", err.Error(), ctx.userID, id)
//				continue
//			}
//			//log.Println("current time ", current)
//			if msgTime.Before(current) {
//				continue
//			}
//			senderID, err := tools.StringToUuid(msg[constants.SENDER_ID])
//			if err != nil {
//				log.Printf("error at HandleListMsgRedisDB at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), userID, id)
//				continue
//			}
//			receiverID, err := tools.StringToUuid(msg[constants.RECEIVER_ID])
//			if err != nil {
//				log.Printf("error at HandleListMsgRedisDB at StringToUuid err:%v, user: %v , id: %v \n", err.Error(), userID, id)
//				continue
//			}
//			// Lets us know if the current user is the sender or receiver
//			// If senderID == userID then we know it is the user that sent the message
//			if senderID == userID {
//				connectedUserID = receiverID
//			}else {
//				connectedUserID = senderID
//				messageCount = 1
//			}
//			connectedUser, err := ctx.server.store.GetUserByUserID(ctx.ctx, connectedUserID)
//			if err != nil {
//				log.Printf("error at HandleListMsgRedisDB at GetUser err:%v, user: %v , id: %v \n", err.Error(), userID, id)
//				continue
//			}
//			userRequestCount, userCancelCount, hostCancelCount, hostChangeDatesCount := MessageType(msg[constants.TYPE])
//			contact := MessageContactItem{
//				ConnectedUserID:            tools.UuidToString(connectedUserID),
//				FirstName:                  connectedUser.FirstName,
//				Photo:                      connectedUser.Photo,
//				LastMessage:                msg[constants.MESSAGE],
//				LastMessageTime:            msg[constants.CREATED_AT],
//				UnreadMessageCount:         messageCount,
//				UnreadUserRequestCount:     userRequestCount,
//				UnreadUserCancelCount:      userCancelCount,
//				UnreadHostCancelCount:      hostCancelCount,
//				UnreadHostChangeDatesCount: hostChangeDatesCount,
//			}
//			resData = append(resData, contact)
//		}
//	}
//	return

//}
