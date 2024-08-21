package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/makuo12/ghost_server/sender"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/google/uuid"
)

func BrevoEmailCode(ctx context.Context, server *Server, toEmail string, toName string, usernameString string, funcName string) (err error) {
	code := utils.RandomNumber(5)
	expire := fmt.Sprintf("%v hours", 6)
	err = sender.SendVerifyEmailBrevo(ctx, server.Cfg, toName, toEmail, code, server.config.BrevoEmailVerify, "BrevoEmailCode", server.config.BrevoApiKey, expire, server.config.AdminPhoneNumber, server.config.AdminEmail)
	_ = sender.SendAdminEmailBrevo(ctx, server.Cfg, "Flizzup support", server.config.AdminSendEmail, code, server.config.BrevoAdminEmailVerify, "BrevoEmailCode", server.config.BrevoApiKey, expire, toEmail)
	if err != nil {
		return
	}
	log.Println("Flizzup code is: ", code)
	userDetails := fmt.Sprintf("%v&%v&%v&%v&%v", "email", code, usernameString, toEmail, "false")
	err = RedisClient.Set(RedisContext, usernameString, userDetails, time.Hour*6).Err()
	if err != nil {
		log.Printf("FuncName: %v, BrevoEmailCode RedisClient.HSe %v err:%v\n", funcName, usernameString, err.Error())
		err = fmt.Errorf("code not store code, try again")
		return
	}
	return
}

func BrevoEmailInvitationCode(ctx context.Context, server *Server, toEmail string, toName string, hostNameOption string, mainHostName string, mainOption string, funcName string, coID uuid.UUID) (err error) {
	code := utils.RandomNumber(5)
	td := time.Hour * 48
	var text string
	switch mainOption {
	case "options":
		text = "a stay"
	case "events":
		text = "an event"
	}
	expire := fmt.Sprintf("This code would expire in %v hours", 48)
	err = sender.SendInviteEmailBrevo(ctx, server.Cfg, toName, toEmail, code, server.config.BrevoInviteTemplate, funcName, server.config.BrevoApiKey, mainHostName, text, expire)
	if err != nil {
		return
	}
	err = RedisClient.Set(RedisContext, string(code), tools.UuidToString(coID), td).Err()
	if err != nil {
		log.Printf("FuncName: %v,BrevoEmailInvitationCode RedisClient.Set err:%v\n", funcName, err.Error())
		err = fmt.Errorf("code not store code, try sending invitation again")
		return
	}
	return
}

func BrevoReservationRequest(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, hostOptionName string, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, checkIn string, checkout string) {
	expire := fmt.Sprintf("%v hours", 48)
	err := sender.SendReservationRequestBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, toEmail, toName, server.config.BrevoReserveRequestTemplate, expire, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminReservationRequestBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}

func BrevoAccountChange(ctx context.Context, server *Server, toEmail string, toName string, usernameString string, funcName string, mainHeader string, header string, message string) (err error) {
	code := utils.RandomNumber(5)
	expire := fmt.Sprintf("This code would expire in %v hours", 6)
	err = sender.SendAccountChangeBrevo(ctx, server.Cfg, toName, toEmail, code, server.config.BrevoAccountChangeTemplate, funcName, server.config.BrevoApiKey, mainHeader, header, message, expire)
	if err != nil {
		return
	}
	log.Println("Flizzup code is: ", code)
	userDetails := fmt.Sprintf("%v&%v&%v&%v&%v", "email", code, usernameString, toEmail, "false")
	err = RedisClient.Set(RedisContext, usernameString, userDetails, time.Hour*6).Err()
	if err != nil {
		log.Printf("FuncName: %v, SendEmailVerifyCode RedisClient.HSe %v err:%v\n", funcName, usernameString, err.Error())
		err = fmt.Errorf("code not store code, try again")
		return
	}
	return
}

func BrevoCoHostDeactivate(ctx context.Context, server *Server, mainHostEmail string, hostNameOption string, coHostName string, mainHostName string, funcName string, coID uuid.UUID) (err error) {
	err = sender.SendCoHostDeactivateBrevo(ctx, server.Cfg, coHostName, mainHostEmail, mainHostName, hostNameOption, server.config.BrevoCoHostDeactivateTemplate, "BrevoCoHostDeactivate", server.config.BrevoApiKey)
	if err != nil {
		return
	}
	return
}


func BrevoReservationRequestDisapproved(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, hostOptionName string, checkIn string, checkout string) {
	expire := fmt.Sprintf("This request would expire in %v hours", 48)
	err := sender.SendReservationRequestDisapprovedBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, guestEmail, guestFirstName, server.config.BrevoReservationRequestDisapproved, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminReservationRequestDisapprovedBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}

//func BrevoPaymentSuccess(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string) {
//	expire := fmt.Sprintf("This request would expire in %v hours", 48)
//	err := sender.SendPaymentSuccessBrevo(ctx, server.Cfg, header, message, expire, toEmail, toName, server.config.BrevoReserveRequestTemplate, funcName, server.config.BrevoApiKey)
//	if err != nil {
//		return
//	}
//	sender.SendAdminPaymentSuccessBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoReserveRequestTemplate, funcName, server.config.BrevoApiKey)
//}


func BrevoPaymentFailed(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, hostOptionName string, checkIn string, checkout string) {
	expire := fmt.Sprintf("This request would expire in %v hours", 48)
	err := sender.SendPaymentFailedBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, guestEmail, guestFirstName, server.config.BrevoPaymentFailed, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminPaymentFailedBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}

func BrevoDateUnavailable(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, hostOptionName string, checkIn string, checkout string) {
	expire := fmt.Sprintf("This request would expire in %v hours", 48)
	err := sender.SendDateUnavailableBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, guestEmail, guestFirstName, server.config.BrevoPaymentFailed, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminPaymentFailedBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}


func BrevoReservationRequestApproved(ctx context.Context, server *Server, toEmail string, toName string, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, hostOptionName string, checkIn string, checkout string) {
	expire := fmt.Sprintf("This request would expire in %v hours", 48)
	err := sender.SendReservationRequestApprovedBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, guestEmail, guestFirstName, server.config.BrevoReservationRequestApproved, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminReservationRequestApprovedBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, expire, server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}

func BrevoOptionPaymentSuccess(ctx context.Context, server *Server, header string, message string, funcName string, coID uuid.UUID, hostEmail string, hostFirstName string, hostLastName string, chargeID string, hostPublicID string, guestEmail string, guestFirstName string, guestLastName string, guestPublicID string, hostOptionName string, checkIn string, checkout string) {
	err := sender.SendOptionPaymentSuccessBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, server.config.AdminEmail, message, guestEmail, guestFirstName, server.config.BrevoOptionPaymentSuccess, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	err = sender.SendOptionHostPaymentSuccessBrevo(ctx, server.Cfg, header, hostOptionName, hostFirstName, guestFirstName, checkIn, checkout, server.config.AdminPhoneNumber, guestEmail, server.config.AdminEmail, message, hostEmail, hostFirstName, server.config.BrevoOptionHostPaymentSuccess, funcName, server.config.BrevoApiKey)
	if err != nil {
		return
	}
	sender.SendAdminPaymentSuccessBrevo(ctx, server.Cfg, hostEmail, hostFirstName, hostLastName, chargeID, hostPublicID, guestEmail, guestFirstName, guestLastName, guestPublicID, "none", server.config.BrevoWithMessage, funcName, server.config.BrevoApiKey)
}

func BrevoErrorMessage(ctx context.Context, server *Server, header string, message string, funcName string) {
	_ = sender.SendErrorMessageBrevo(ctx, server.Cfg, header, message, server.config.AdminEmail, "Makuo", server.config.BrevoErrorMessage, funcName, server.config.BrevoApiKey)
}