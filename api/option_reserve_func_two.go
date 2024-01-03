package api

import (
	"context"
	"flex_server/constants"
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleReserveOption(user db.User, server *Server, ctx *gin.Context, startDate string, endDate string, guests []string, optionUserID string, userCurrency string) (reference string, err error) {
	canInstantBook, datePriceFloat, totalDatePrice, discountType, discount, cleanFee, extraGuestFee, petFee, petStayFee, extraGuestStayFee, totalPrice, serviceFee, requireRequest, requestType, optionUserUUID, err := ReserveOptionCalculate(user, server, ctx, startDate, endDate, guests, optionUserID, userCurrency)
	if err != nil {
		return
	}
	reference, err = HandleOptionReserveRedis(user, canInstantBook, datePriceFloat, totalDatePrice, discountType, discount, cleanFee, extraGuestFee, petFee, petStayFee, extraGuestStayFee, totalPrice, serviceFee, requireRequest, requestType, userCurrency, optionUserUUID, startDate, endDate, guests)
	return
}

func HandleFinalOptionReserveDetail(server *Server, ctx context.Context, reference string, user db.User, cardID string, msg string) (res FinalOptionReserveRequestDetailRes, hasResData bool, totalFee string, refRes string, reserveData ExperienceReserveOModel, fromCharge bool, chargeData db.ChargeOptionReference, err error) {
	charge, err := server.store.GetChargeOptionReferenceByRef(ctx, db.GetChargeOptionReferenceByRefParams{
		Reference: reference,
		UserID:    user.UserID,
	})
	if err != nil {
		// If there is an error we know it is not a reservation request
		log.Printf("Error at tools.StringToUuid: %v, userID: %v \n", err.Error(), user.ID)
		err = nil

	} else {
		if charge.RequestApproved || charge.CanInstantBook {
			if charge.IsComplete {
				err = fmt.Errorf("Your reservation has already been payed for")
			}
			hasResData = false
			totalFee = tools.IntToMoneyString(charge.TotalFee)
			refRes = charge.Reference
			chargeData = charge
			fromCharge = true
			return
		}

		res = FinalOptionReserveRequestDetailRes{
			Reference:   reference,
			Message:     "Your reservation request has not yet been approved",
			RequestSent: true,
		}
		hasResData = true
		return
	}

	userID := tools.UuidToString(user.ID)
	reserveData, err = HandleOptionReserveRedisData(userID, reference)
	if err != nil {
		log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveRedisData: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
		err = fmt.Errorf("reference has expired, please try again")
		return
	}
	// We want to handle request reservations
	if !reserveData.CanInstantBook || reserveData.RequireRequest {
		err = HandleOptionReserveRequest(server, ctx, cardID, reserveData, user, msg)
		if err != nil {
			log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveRequest: %v, cardID: %v, userID: %v \n", err.Error(), cardID, user.ID)
			err = fmt.Errorf("your reservation request was unsuccessful")
			return
		}
		res = FinalOptionReserveRequestDetailRes{
			Reference:   reference,
			Message:     "Your reservation request has been sent",
			RequestSent: true,
		}
		hasResData = true
		return
	}
	totalFee = reserveData.TotalFee

	refRes = reserveData.Reference
	return

}

func HandleOptionUserRequest(req MsgRequestResponseParams, msg db.Message, user db.User, server *Server, ctx *gin.Context) (res MessageContactItem, err error) {
	if user.UserID != msg.ReceiverID {
		err = fmt.Errorf("only the host is allowed to response to this request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// This means that it is the host response
	notify, err := server.store.UpdateRequestNotify(ctx, db.UpdateRequestNotifyParams{
		Approved: pgtype.Bool{
			Bool:  req.Approved,
			Valid: true,
		},
		Cancelled: pgtype.Bool{
			Bool:  !req.Approved,
			Valid: true,
		},

		MID: msg.ID,
	})
	if err != nil {
		log.Printf("Error at HandleOptionUserRequest in UpdateRequestNotify: %v, MsgID: %v \n", err.Error(), req.MsgID)
		err = fmt.Errorf("your response was recorded, but could not update this request notification")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We want to set the message field read to true
	_, err = server.store.UpdateMessageReadByID(ctx, db.UpdateMessageReadByIDParams{
		Read:       true,
		MsgID:      msg.MsgID,
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
	})
	if err != nil {
		log.Printf("Error at HandleOptionUserRequest in UpdateMessageReady: %v, MsgID: %v \n", err.Error(), req.MsgID)
	}
	if notify.Approved {
		charge, errUpdate := server.store.UpdateChargeOptionReferenceByRef(ctx, db.UpdateChargeOptionReferenceByRefParams{
			RequestApproved: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			RequireRequest: pgtype.Bool{
				Bool:  false,
				Valid: true,
			},

			DateBooked: pgtype.Timestamptz{
				Time:  time.Now().Add(time.Hour),
				Valid: true,
			},
			CanInstantBook: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			Reference: msg.Reference,
			UserID:    msg.SenderID,
		})
		if errUpdate != nil {
			log.Printf("Error at HandleOptionUserRequest in UpdateChargeOptionReferenceByRef: %v, MsgID: %v \n", err.Error(), req.MsgID)
			err = fmt.Errorf("your response was recorded, but could not update this request notification")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// We want to store in redis to make payment
		timeData := []string{
			constants.MID,
			tools.UuidToString(msg.ID),
			constants.SENDER_ID,
			tools.UuidToString(msg.SenderID),
			constants.RECEIVER_ID,
			tools.UuidToString(msg.ReceiverID),
			constants.FIRSTNAME,
			user.FirstName,
			constants.REFERENCE,
			charge.Reference,
			constants.TIME,
			tools.ConvertTimeToString(time.Now().Add(time.Hour)),
		}
		id := tools.UuidToString(uuid.New())
		err = RedisClient.HSet(RedisContext, id, timeData).Err()
		if err != nil {
			log.Printf("error at HandleOptionUserRequest timeData at RedisClient.HSet err:%v, user: %v , id: %v \n", err.Error(), msg.ReceiverID, msg.ID)
		} else {
			err = RedisClient.SAdd(RedisContext, constants.USER_REQUEST_APPROVE, id).Err()
			if err != nil {
				log.Printf("error at HandleOptionUserRequest timeData at RedisClient.SAdd err:%v, user: %v , id: %v \n", err.Error(), msg.ReceiverID, msg.ID)
			}
		}

	} else if notify.Cancelled {
		charge, errGet := server.store.GetChargeOptionReferenceDetailByRef(ctx, msg.Reference)
		if errGet != nil {
			log.Printf("error at HandleOptionUserRequest timeData at RedisClient.SAdd err:%v, user: %v , id: %v \n", err.Error(), msg.ReceiverID, msg.ID)
		} else {
			header := fmt.Sprintf("Reservation for %v disapproved", charge.HostNameOption)
			msgString := fmt.Sprintf("Hey %v,\n%v could not approved your booking for %v. No worry's remember you can try again later or find other stays that match the experience you want to have", charge.UserFirstName, user.FirstName, charge.HostNameOption)
			CreateTypeNotification(ctx, server, msg.ID, msg.SenderID, constants.HOST_DISAPPROVED_BOOKING, msgString, false, header)
		}

	}
	err = nil
	contact, err := server.store.GetMessageContact(ctx, db.GetMessageContactParams{
		SenderID: msg.ReceiverID,
		MsgID:    msg.MsgID,
	})
	if err != nil {

		log.Printf("Error at  GetMessageContact in GetMessageContact err: %v, user: %v\n", err, user.ID)
		err = nil
		res = MessageContactItem{
			MsgID:                      "none",
			ConnectedUserID:            "none",
			FirstName:                  "none",
			Photo:                      "none",
			LastMessage:                "none",
			LastMessageTime:            "none",
			UnreadMessageCount:         0,
			UnreadUserRequestCount:     0,
			UnreadUserCancelCount:      0,
			UnreadHostCancelCount:      0,
			UnreadHostChangeDatesCount: 0,
		}
		log.Println("res empty")
		return
	}
	res = MessageContactItem{
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
	log.Println("res contact ", res)
	return
}

func HandleChargeDatesAvailable(ctx context.Context, server *Server, option db.GetOptionInfoCustomerRow, startDate string, endDate string, funcName string) (available bool, err error) {
	// List of dates
	userDateTime, err := tools.GenerateDateListString(startDate, endDate)
	if err != nil {
		log.Printf("Error at FuncName: %v at ReserveOptionList in .ListOptionDiscount: %v option.ID: %v\n", funcName, err.Error(), option.ID)
		return
	}
	dateTimes, err := server.store.ListAllOptionDateTime(ctx, option.ID)
	if err != nil {
		// We don't need to send an error because host might have no special days
		log.Printf("Error at ReserveOptionList in .ListOptionDiscount: %v optionID: %v\n", err.Error(), option.ID)
		dateTimes = []db.OptionDateTime{}
	}
	dateTimeString := OptionDateTimeString(dateTimes)
	confirm, err := ReserveDatesAvailable(startDate, endDate, userDateTime, dateTimeString, option.ID)
	if err != nil || !confirm {
		err = tools.HandleConfirmError(err, confirm, "Your reservation was unsuccessful because the dates have now been booked")
		return
	}
	free, err := HandleReserveAvailable(ctx, server, option.OptionUserID, option.PreparationTime, option.AutoBlockDates, userDateTime)
	if err != nil || !free {
		err = tools.HandleConfirmError(err, confirm, "Your reservation was unsuccessful because the dates have now been booked")
		return
	}
	if free && confirm {
		available = true
	}
	return
}

// UserRequest -> UR
// Firstname is the host name
// mid is the request_notify mid
func HandleURApproved(ctx context.Context, server *Server, mid uuid.UUID, senderID string, receiverID string, firstName string, reference string) {
	senderUUID, err := tools.StringToUuid(senderID)
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in tools.StringToUuid: %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		return
	}
	charge, err := server.store.GetChargeOptionReferenceDetailByRef(ctx, reference)
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in tools.StringToUuid: %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		msgString := fmt.Sprintf("%v has approved your booking request. However, your payment was unable to go through. Please try to complete your payment", firstName)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, "Payment unsuccessful")
	}
	// We want to check if the dates are still available
	option, err := server.store.GetOptionInfoCustomer(ctx, db.GetOptionInfoCustomerParams{
		OptionUserID:    charge.OptionUserID,
		IsComplete:      true,
		IsActive:        true, // Option is active
		IsActive_2:      true, // Host is active
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	header := fmt.Sprintf("Selected dates for %v are unavailable", charge.HostNameOption)
	msgString := fmt.Sprintf("%v has approved your booking request for %v. However, the dates are no more available so we could not proceed further with the payments.", firstName, charge.HostNameOption)
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in store.GetOptionInfoCustomer %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_DATES_UNAVAILABLE, msgString, false, header)
		return
	}
	available, err := HandleChargeDatesAvailable(ctx, server, option, tools.ConvertDateOnlyToString(charge.StartDate), tools.ConvertDateOnlyToString(charge.EndDate), "HandleURApproved")
	if err != nil || !available {
		if err != nil {
			log.Printf("Error at HandleURApproved senderID in HandleChargeDatesAvailable %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		}
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_DATES_UNAVAILABLE, msgString, false, header)
		return
	}

	header = fmt.Sprintf("Payment for %v unsuccessful", charge.HostNameOption)
	msgString = fmt.Sprintf("%v has approved your booking request for %v. However, your payment was unable to go through. Please try to complete your payment", firstName, charge.HostNameOption)
	cardID, err := tools.StringToUuid(charge.UserDefaultCard)
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in StringToUuid(sender.DefaultCard): %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		return
	}
	card, err := server.store.GetCard(ctx, db.GetCardParams{
		ID:     cardID,
		UserID: charge.GuestID,
	})
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in store.GetCard(: %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		return
	}
	guest, err := server.store.GetUser(ctx, charge.GuestID)
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in GetUser: %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		return
	}
	if charge.Currency != card.Currency {
		msgStringData := "selected currency doesn't match card currency"
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgStringData, false, header)
		return
	}
	// Next we want to make the payment
	resData, _, resChallenged, err := HandlePaystackChargeAuthorization(server, ctx, card, tools.IntToMoneyString(charge.TotalFee))
	if err != nil {
		log.Printf("Error at HandleURApproved senderID in HandlePaystackChargeAuthorization: %v, MsgID: %v senderID: %v \n", err.Error(), mid, senderID)
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		return
	}
	if !resChallenged {
		if resData.Data.Status != "success" {
			CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
			return
		} else {
			// Payment was successful
			// We want to save a recept in the database
			// We also want to store a snap shot of what the shortlet looks like
			// Creating snapshot
			dateString := tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM)
			msg := fmt.Sprintf("Hey %v,\nI would like to thank you for letting me stay in your place from %v.", charge.HostFirstName, dateString)
			err := HandleOptionReserveComplete(server, ctx, ExperienceReserveOModel{}, charge.Reference, resData.Data.Reference, guest, msg, true)
			if err != nil {
				header = fmt.Sprintf("Payment for %v successful, but couldn't generate receipt", charge.HostNameOption)
				log.Printf("Error at HandleURApproved in HandleOptionReserveComplete: %v, cardID: %v, receiverID: %v \n", err.Error(), cardID, receiverID)
				msg = fmt.Sprintf("Hey %v,\nYour payment was successful for %v, however we could not generate a receipt for you. This error was from our servers. Please contact us immediately so we can take care of this.", charge.UserFirstName, charge.HostNameOption)
				CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_PAYMENT_SUCCESSFUL_NO_RECEIPT, msg, false, header)
				return
			}
		}
	} else {
		CreateTypeNotification(ctx, server, mid, senderUUID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
	}
}

func HandleReserveCard(ctx context.Context, server *Server, user db.User, funcName string) (string, CardDetailResponse, bool) {

	cardDetail := CardDetailResponse{"none", "none", "none", "none", "none", "none"}
	if tools.ServerStringEmpty(user.DefaultCard) {
		// If default card is empty we get any card to make the payment
		card, exist := HandleGetAnyCard(ctx, server, user, funcName)
		if exist {
			return tools.UuidToString(card.ID), HandleConvertCardToResponse(card), true
		}
		return "none", cardDetail, false
	}
	cardID, err := tools.StringToUuid(user.DefaultCard)
	if err != nil {
		// If there is an error we get any card
		log.Printf("error at funcName: %v HandleReserveCard at StringToUuid for card: %v, userID: %v\n", funcName, err.Error(), user.ID)
		card, exist := HandleGetAnyCard(ctx, server, user, funcName)
		if exist {
			return tools.UuidToString(card.ID), HandleConvertCardToResponse(card), true
		}
		return "none", cardDetail, false
	}
	defaultCard, err := server.store.GetCard(ctx, db.GetCardParams{
		UserID: user.ID,
		ID:     cardID,
	})
	if err != nil {
		// If there is an error we get any card
		log.Printf("error at funcName: %v HandleReserveCard at server.store.GetCard for card: %v, userID: %v\n", funcName, err.Error(), user.ID)
		card, exist := HandleGetAnyCard(ctx, server, user, funcName)
		if exist {
			return tools.UuidToString(card.ID), HandleConvertCardToResponse(card), true
		}
		return "none", cardDetail, false
	}
	return tools.UuidToString(defaultCard.ID), HandleConvertCardToResponse(defaultCard), true
}

func HandleConvertCardToResponse(card db.Card) CardDetailResponse {
	res := CardDetailResponse{
		CardID:    tools.UuidToString(card.ID),
		CardLast4: card.Last4,
		CardType:  card.CardType,
		ExpMonth:  card.ExpMonth,
		ExpYear:   card.ExpYear,
		Currency:  card.Currency,
	}
	return res
}

func HandleGetAnyCard(ctx context.Context, server *Server, user db.User, funcName string) (db.Card, bool) {
	card, err := server.store.GetCardAny(ctx, user.ID)
	if err != nil {
		log.Printf("error at funcName: %v HandleGetAnyCard at StringToUuid for card: %v, userID: %v\n", funcName, err.Error(), user.ID)
		return db.Card{}, false
	}
	return card, true
}
