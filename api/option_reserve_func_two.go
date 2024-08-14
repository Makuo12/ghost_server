package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/payment"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

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

func HandleFinalOptionReserveDetail(server *Server, ctx context.Context, reference string, user db.User, msg string, payMethodReference string) (res FinalOptionReserveRequestDetailRes, hasResData bool, totalFee string, refRes string, reserveData ExperienceReserveOModel, fromCharge bool, chargeData db.ChargeOptionReference, err error) {
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
				err = fmt.Errorf("your reservation has already been payed for")
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
		log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveRedisData: %v, payMethodReference: %v, userID: %v \n", err.Error(), payMethodReference, user.ID)
		err = fmt.Errorf("reference has expired, please try again")
		return
	}
	// We want to handle request reservations
	if !reserveData.CanInstantBook || reserveData.RequireRequest {
		err = HandleOptionReserveRequest(server, ctx, payMethodReference, reserveData, user, msg)
		if err != nil {
			log.Printf("Error at FinalOptionReserveDetail in HandleOptionReserveRequest: %v, payMethodReference: %v, userID: %v \n", err.Error(), payMethodReference, user.ID)
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

// request_notifies
// it has 3 types of status
// normal
// request_payment
// sent_request
func HandleOptionUserRequest(req MsgRequestResponseParams, msg db.Message, user db.User, server *Server, ctx *gin.Context) (res MessageContactItem, err error) {
	status := "normal"
	if user.UserID != msg.ReceiverID {
		err = fmt.Errorf("only the host is allowed to response to this request")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if req.Approved {
		status = constants.REQUEST_PAYMENT
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
		Status: pgtype.Text{
			String: status,
			Valid:  true,
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
	// We want to check if the dates are still available
	option, err := server.store.GetOptionInfoCustomerWithRef(ctx, db.GetOptionInfoCustomerWithRefParams{
		Reference:       msg.Reference,
		IsComplete:      true,
		IsActive:        true, // Option is active
		IsActive_2:      true, // Host is active
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		log.Printf("Error at HandleOptionUserRequest in GetOptionInfoCustomerWithRef: %v, MsgID: %v \n", err.Error(), req.MsgID)
		err = fmt.Errorf("the dates the guest selected are no more available, please try again")
		return
	}
	available, err := HandleChargeDatesAvailableWithRef(ctx, server, option, tools.ConvertDateOnlyToString(option.StartDate), tools.ConvertDateOnlyToString(option.EndDate), "FinalOptionReserveDetail")
	if err != nil || !available {
		if err != nil {
			log.Printf("Error at FinalOptionReserveDetail in HandleChargeDatesAvailableWithRef: %v, option.ID: %v, userID: %v \n", err.Error(), option.ID, user.ID)
		}
		err = fmt.Errorf("the dates the guest selected are no more available, please try again")
		return
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
			log.Printf("Error at HandleOptionUserRequest in UpdateChargeOptionReferenceByRef: %v, MsgID: %v \n", errUpdate.Error(), req.MsgID)
			err = fmt.Errorf("your response was recorded, but could not update this request notification")
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		guest, err := server.store.GetUserByUserID(ctx, charge.UserID)
		if err != nil {
			log.Printf("Error at HandleOptionUserRequest in GetUserByUserID: %v, MsgID: %v \n", err.Error(), req.MsgID)
		} else {
			header := "Reservation approved"
			msg := fmt.Sprintf("Hey %s,\n your reservation request has just been approved by %s", tools.CapitalizeFirstCharacter(guest.FirstName), tools.CapitalizeFirstCharacter(user.FirstName))
			// We send a notification to the guest to notify them
			CreateTypeNotification(ctx, server, charge.ID, guest.UserID, "USER_REQUEST_APPROVE", msg, false, header)
			// We send an email
			checkIn := tools.HandleReadableDate(option.StartDate, tools.DateDMMYyyy)
			checkout := tools.HandleReadableDate(option.EndDate, tools.DateDMMYyyy)
			BrevoReservationRequestApproved(ctx, server, guest.Email, tools.CapitalizeFirstCharacter(guest.FirstName), header, msg, "HandleOptionUserRequest", charge.ID, user.Email, tools.CapitalizeFirstCharacter(user.FirstName), tools.CapitalizeFirstCharacter(user.LastName), tools.UuidToString(charge.ID), tools.UuidToString(user.UserID), guest.Email, tools.CapitalizeFirstCharacter(guest.FirstName), tools.CapitalizeFirstCharacter(guest.LastName), tools.UuidToString(guest.UserID), option.HostNameOption, checkIn, checkout)
		}

	} else if notify.Cancelled {
		charge, errGet := server.store.GetChargeOptionReferenceDetailByRef(ctx, msg.Reference)
		if errGet != nil {
			log.Printf("error at HandleOptionUserRequest timeData at RedisClient.SAdd err:%v, user: %v , id: %v \n", errGet.Error(), msg.ReceiverID, msg.ID)
		} else {
			header := fmt.Sprintf("Reservation for %v disapproved", charge.HostNameOption)
			msgString := fmt.Sprintf("Hey %v,\n%v could not approved your booking for %v. No worry's remember you can try again later or find other stays that match the experience you want to have", charge.UserFirstName, user.FirstName, charge.HostNameOption)
			CreateTypeNotification(ctx, server, msg.ID, msg.SenderID, constants.HOST_DISAPPROVED_BOOKING, msgString, false, header)
			guest, err := server.store.GetUserByUserID(ctx, charge.GuestID)
			if err != nil {
				log.Printf("Error at HandleOptionUserRequest in GetUserByUserID: %v, MsgID: %v \n", err.Error(), req.MsgID)
			} else {
				header := "Reservation approved"
				msg := fmt.Sprintf("Hey %s,\n your reservation request has just been approved by %s", guest.FirstName, user.FirstName)
				// We send a notification to the guest to notify them
				CreateTypeNotification(ctx, server, charge.ChargeID, user.UserID, "USER_REQUEST_APPROVE", msg, false, header)
				// We send an email
				checkIn := tools.HandleReadableDate(option.StartDate, tools.DateDMMYyyy)
				checkout := tools.HandleReadableDate(option.EndDate, tools.DateDMMYyyy)
				BrevoReservationRequestApproved(ctx, server, guest.Email, guest.FirstName, header, msg, "HandleOptionUserRequest", charge.ChargeID, user.Email, user.FirstName, user.LastName, tools.UuidToString(charge.ChargeID), tools.UuidToString(user.UserID), guest.Email, guest.FirstName, guest.LastName, tools.UuidToString(guest.UserID), option.HostNameOption, checkIn, checkout)
			}
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
			MainImage:                  "none",
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
		MainImage:                  contact.Image,
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

func HandleChargeDatesAvailableWithRef(ctx context.Context, server *Server, option db.GetOptionInfoCustomerWithRefRow, startDate string, endDate string, funcName string) (available bool, err error) {
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
// The job that handles when a reservation request has been approved
func HandleURApproved(ctx context.Context, server *Server, mid uuid.UUID, reference string, chargeID uuid.UUID, guestFirstName string, hostFirstName string, guestUserID uuid.UUID, hostUserID uuid.UUID) {
	log.Println("starting payment")
	successUrl := server.config.PaymentSuccessUrl
	failureUrl := server.config.PaymentFailUrl
	_, err := server.store.UpdateRequestNotify(ctx, db.UpdateRequestNotifyParams{
		Status: pgtype.Text{
			String: constants.SENT_REQUEST,
			Valid:  true,
		},
		MID: mid,
	})
	if err != nil {
		log.Printf("1 Error at HandleURApproved guestUserID in tools.StringToUuid: %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		msgString := fmt.Sprintf("%v has approved your booking request. However, your payment was unable to go through. Please try to complete your payment", guestFirstName)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, "Payment unsuccessful")
		return
	}
	charge, err := server.store.GetChargeOptionReferenceDetailByRef(ctx, reference)
	if err != nil {
		log.Printf("2 Error at HandleURApproved guestUserID in tools.StringToUuid: %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		msgString := fmt.Sprintf("%v has approved your booking request. However, your payment was unable to go through. Please try to complete your payment", guestFirstName)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, "Payment unsuccessful")
		return
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
	msgString := fmt.Sprintf("%v has approved your booking request for %v. However, the dates are no more available so we could not proceed further with the payments.", guestFirstName, charge.HostNameOption)
	if err != nil {
		log.Printf("3 Error at HandleURApproved guestUserID in store.GetOptionInfoCustomer %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, false)
		return
	}
	available, err := HandleChargeDatesAvailable(ctx, server, option, tools.ConvertDateOnlyToString(charge.StartDate), tools.ConvertDateOnlyToString(charge.EndDate), "HandleURApproved")
	if err != nil || !available {
		// We handle when dates are not available
		if err != nil {
			log.Printf("4 Error at HandleURApproved guestUserID in HandleChargeDatesAvailable %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		}
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, false)
		return
	}
	header = fmt.Sprintf("Payment for %v unsuccessful", charge.HostNameOption)
	msgString = fmt.Sprintf("%v has approved your booking request for %v. However, your payment was unable to go through. Please try to complete your payment", guestFirstName, charge.HostNameOption)
	cardID, err := tools.StringToUuid(charge.UserDefaultCard)
	if err != nil {
		log.Printf("5 Error at HandleURApproved guestUserID in StringToUuid(sender.DefaultCard): %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
		return
	}
	card, err := server.store.GetCard(ctx, db.GetCardParams{
		ID:     cardID,
		UserID: charge.GuestID,
	})
	if err != nil {
		log.Printf("6 Error at HandleURApproved guestUserID in store.GetCard(: %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
		return
	}
	guest, err := server.store.GetUser(ctx, charge.GuestID)
	if err != nil {
		log.Printf("7 Error at HandleURApproved guestUserID in GetUser: %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
		return
	}
	if charge.Currency != card.Currency {
		msgStringData := "selected currency doesn't match card currency"
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgStringData, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
		return
	}
	// Next we want to make the payment
	resData, _, resChallenged, err := payment.HandlePaystackChargeAuthorization(ctx, successUrl, failureUrl, server.config.PaystackSecretLiveKey, card, tools.IntToMoneyString(charge.TotalFee))
	if err != nil {
		log.Printf("8 Error at HandleURApproved guestUserID in HandlePaystackChargeAuthorization: %v, MsgID: %v guestUserID: %v \n", err.Error(), mid, guestUserID)
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
		return
	}
	if !resChallenged {
		if resData.Data.Status != "success" {
			msgString := fmt.Sprintf("%s has approved your booking request for %s. However, the dates are no more available so we could not proceed further with the payments. The issue was because of %s", guestFirstName, charge.HostNameOption, resData.Data.GatewayResponse)
			CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
			HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
			log.Println("pay 8: ", resData.Data.GatewayResponse)
			return
		} else {
			// Payment was successful
			// We want to save a recept in the database
			// We also want to store a snap shot of what the shortlet looks like
			// Creating snapshot
			HandleURComplete(ctx, server, charge, option)
			dateString := tools.HandleReadableDates(charge.StartDate, charge.EndDate, tools.DateDMM)
			msg := fmt.Sprintf("Hey %v,\nI would like to thank you for letting me stay in your place from %v.", charge.HostFirstName, dateString)
			err := HandleOptionReserveComplete(server, ctx, ExperienceReserveOModel{}, charge.Reference, resData.Data.Reference, guest, msg, true)
			log.Println("payment good")
			
			if err != nil {
				header = fmt.Sprintf("Payment for %v successful, but couldn't generate receipt", charge.HostNameOption)
				log.Printf("Error at HandleURApproved in HandleOptionReserveComplete: %v, cardID: %v, hostUserID: %v \n", err.Error(), cardID, hostUserID)
				msg = fmt.Sprintf("Hey %v,\nYour payment was successful for %v, however we could not generate a receipt for you. This error was from our servers. Please contact us immediately so we can take care of this.", charge.UserFirstName, charge.HostNameOption)
				CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_PAYMENT_SUCCESSFUL_NO_RECEIPT, msg, false, header)
				HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
				return
			}
		}
	} else {
		CreateTypeNotification(ctx, server, mid, guestUserID, constants.OPTION_BOOKING_PAYMENT_UNSUCCESSFUL, msgString, false, header)
		HandleURApprovedFailed(ctx, server, charge, header, msgString, true)
	}
}

func HandleURApprovedFailed(ctx context.Context, server *Server, charge db.GetChargeOptionReferenceDetailByRefRow, header string, message string, isPayment bool) {
	checkIn := tools.HandleReadableDate(charge.StartDate, tools.DateDMMYyyy)
	checkout := tools.HandleReadableDate(charge.EndDate, tools.DateDMMYyyy)
	if isPayment {

		BrevoPaymentFailed(ctx, server, charge.GuestEmail, charge.UserFirstName, header, message, "HandleURApproved", charge.ChargeID, charge.HostEmail, charge.HostFirstName, charge.HostLastName, tools.UuidToString(charge.ChargeID), tools.UuidToString(charge.HostUserID), charge.GuestEmail, charge.UserFirstName, charge.UserLastName, tools.UuidToString(charge.GuestUserID), charge.HostNameOption, checkIn, checkout)
	} else {
		BrevoPaymentFailed(ctx, server, charge.GuestEmail, charge.UserFirstName, header, message, "HandleURApproved", charge.ChargeID, charge.HostEmail, charge.HostFirstName, charge.HostLastName, tools.UuidToString(charge.ChargeID), tools.UuidToString(charge.HostUserID), charge.GuestEmail, charge.UserFirstName, charge.UserLastName, tools.UuidToString(charge.GuestUserID), charge.HostNameOption, checkIn, checkout)
	}

}

func HandleURComplete(ctx context.Context, server *Server, charge db.GetChargeOptionReferenceDetailByRefRow, option db.GetOptionInfoCustomerRow) {
	host, err := server.store.GetOptionInfoUserIDByUserID(ctx, charge.OptionUserID)
		if err != nil {
			log.Printf("Error at HandleOptionReserveRequest in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), charge.OptionUserID, charge.Reference, charge.Reference)
			err = nil
		} else {
			header := "Payment and Reservation confirmed"
			msg := "Thank you for using Flizzup"
			checkIn := tools.HandleReadableDate(charge.StartDate, tools.DateDMMYyyy)
			checkout := tools.HandleReadableDate(charge.EndDate, tools.DateDMMYyyy)
			BrevoOptionPaymentSuccess(ctx, server, header, msg, "FinalOptionReserveDetail", charge.ChargeID, host.Email, host.FirstName, host.LastName, tools.UuidToString(charge.ChargeID), tools.UuidToString(host.UserID), option.Email, option.FirstName, option.LastName, tools.UuidToString(option.UserID), host.HostNameOption, checkIn, checkout)
			// Notification for Guest
			msg = "Payment received, and reservation confirmed! Check your email for further information about your booking, including scheduling an inspection and our 100% refund policy if the property does not match the app's description."
			header = fmt.Sprintf("Hey %v, booking confirmed", option.FirstName)
			CreateTypeNotification(ctx, server, charge.ChargeID, option.UserID, constants.OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
			// Notification for Host
			msg = fmt.Sprintf("You have a new booking at %v! Check your email for further details.", host.HostNameOption)
			header = fmt.Sprintf("Hey %v", host.FirstName)
			CreateTypeNotification(ctx, server, charge.ChargeID, host.UserID, constants.HOST_OPTION_PAYMENT_SUCCESSFUL, msg, false, header)
		}
}

func HandleReserveCard(ctx context.Context, server *Server, user db.User, funcName string) (string, payment.CardDetailResponse, bool) {

	cardDetail := payment.CardDetailResponse{CardID: "none", CardLast4: "none", CardType: "none", ExpMonth: "none", ExpYear: "none", Currency: "none"}
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

func HandleConvertCardToResponse(card db.Card) payment.CardDetailResponse {
	res := payment.CardDetailResponse{
		CardID:    tools.UuidToString(card.ID),
		CardLast4: card.Last4,
		CardType:  utils.MatchCardType(card.CardType),
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
