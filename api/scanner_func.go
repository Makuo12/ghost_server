package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/google/uuid"
)

// HandleCodeExist checks if this user chargeID already has a code with it, if so then it deletes it
func HandleCodeExist(userID uuid.UUID, chargeID string, funcName string) (string, error) {
	userDetails := fmt.Sprintf("%v&%v", userID, chargeID)
	result, err := RedisClient.Exists(RedisContext, userDetails).Result()
	if err != nil {
		log.Printf("Error at FuncName: %v HandleCodeExist in RedisClient.Exists err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("something went wrong")
		return "", err
	}
	if result == 1 {
		// If it exist we remove it from redis
		key, err := RedisClient.Get(RedisContext, userDetails).Result()
		if err != nil {
			log.Printf("Error at FuncName: %v HandleCodeExist in RedisClient.Get err: %v, chargeID: %v\n", funcName, err, chargeID)
			err = errors.New("something went wrong")
			return "", err
		}
		err = RedisClient.Del(RedisContext, key, userDetails).Err()
		if err != nil {
			log.Printf("Error at FuncName: %v HandleCodeExist in RedisClient.Del err: %v, chargeID: %v\n", funcName, err, chargeID)
			err = errors.New("something went wrong")
			return "", err
		}
	}
	return userDetails, nil
}

func GetCodeDetails(code string, funcName string) (guestID uuid.UUID, chargeID uuid.UUID, optionUserID uuid.UUID, firstName string, err error) {
	result, err := RedisClient.Exists(RedisContext, code).Result()
	if err != nil {
		log.Printf("Error at FuncName: %v GetCodeDetails in RedisClient.Exists err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("something went wrong")
		return
	}
	if result == 0 {
		err = errors.New("this code does not exist, tell the guest to regenerate the code has it might have expired")
		return
	}
	// If it exist we want to get the code
	details, err := RedisClient.Get(RedisContext, code).Result()
	if err != nil {
		log.Printf("Error at FuncName: %v GetCodeDetails in RedisClient.Get err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("something went wrong")
		return
	}
	// We know the code exist
	split := strings.Split(details, "&")
	if len(split) != 4 {
		err = errors.New("details in code does not match, tell the guest to regenerate the code has it might have expired")
		return
	}
	guestID, err = tools.StringToUuid(split[0])
	if err != nil {
		log.Printf("Error userID at FuncName: %v GetCodeDetails in tools.StringToUuid err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("guest could not the decoded")
		return
	}
	chargeID, err = tools.StringToUuid(split[1])
	if err != nil {
		log.Printf("Error chargeID at FuncName: %v GetCodeDetails in tools.StringToUuid err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("reference could not the decoded")
		return
	}
	optionUserID, err = tools.StringToUuid(split[2])
	if err != nil {
		log.Printf("Error optionUserID at FuncName: %v GetCodeDetails in tools.StringToUuid err: %v, chargeID: %v\n", funcName, err, chargeID)
		err = errors.New("reference could not the decoded")
		return
	}
	firstName = split[3]
	return
}

// GetChargeForScanned Is to know if code has been scanned
// chargeID uuid.UUID, optionUserID uuid.UUID, scanned bool, grade string, scannedBy uuid.UUID, scannedTime time, scannedByName string, ScannedUserImage, err error
func GetChargeForScanned(ctx context.Context, server *Server, req GetChargeCodeParams, user db.User, funcName string) (uuid.UUID, uuid.UUID, bool, string, uuid.UUID, time.Time, string, string, string, error) {
	reqChargeID, err := tools.StringToUuid(req.ID)
	if err != nil {
		err = errors.New("not found, try again later")
		return uuid.UUID{}, uuid.UUID{}, false, "none", uuid.UUID{}, time.Time{}, "none", "none", "none", err
	}

	switch req.MainOption {
	case "options":
		charge, err := server.store.GetScannedChargeOption(ctx, db.GetScannedChargeOptionParams{
			ChargeID:         reqChargeID,
			UserID:           user.UserID,
			PaymentCompleted: true,
			Cancelled:        false,
		})
		if err != nil {
			log.Printf("Error at FuncName: %v GetChargeForScanned in .server.store.GetScannedChargeOption err: %v, user: %v\n", funcName, err, user.ID)
			err = errors.New("could not find your booking, try again later. Or try contacting us")
			return uuid.UUID{}, uuid.UUID{}, false, "none", uuid.UUID{}, time.Time{}, "none", "none", "none", err
		} else {
			return charge.ChargeID, charge.OptionUserID, HandleSqlNullBool(charge.Scanned), "none", HandleSqlNullUUID(charge.ScannedBy), HandleSqlNullTimestamp(charge.ScannedTime), HandleSqlNullString(charge.ScannedByName), HandleSqlNullString(charge.ScannedUserImage), "none", nil
		}
	case "events":
		charge, err := server.store.GetScannedChargeTicket(ctx, db.GetScannedChargeTicketParams{
			ChargeID:         reqChargeID,
			UserID:           user.UserID,
			PaymentCompleted: true,
			Cancelled:        false,
		})
		if err != nil {
			log.Printf("Error at FuncName: %v GetChargeForScanned in .GetScannedChargeTicket err: %v, user: %v\n", funcName, err, user.ID)
			if err == db.ErrorRecordNotFound {
				err = nil
			} else {
				err = errors.New("could not find your booking, try again later. Or try contacting us")
			}
			return uuid.UUID{}, uuid.UUID{}, false, "none", uuid.UUID{}, time.Time{}, "none", "none", "none", err
		} else {
			return charge.ChargeID, charge.OptionUserID, HandleSqlNullBool(charge.Scanned), charge.Grade, HandleSqlNullUUID(charge.ScannedBy), HandleSqlNullTimestamp(charge.ScannedTime), HandleSqlNullString(charge.ScannedByName), HandleSqlNullString(charge.ScannedUserImage), charge.TicketType, nil
		}
	default:
		err = errors.New("unsure whether this is an event or listing, try sending the right info")
		return uuid.UUID{}, uuid.UUID{}, false, "none", uuid.UUID{}, time.Time{}, "none", "none", "none", err
	}
}

func HandleCodeEncrypt(server *Server, key string, user db.User, funcName string) (code string, err error) {
	hashKey := []byte(server.config.TokenSymmetricKey)
	// Hash the code using same format with password
	code, err = tools.Encrypt(hashKey, key)
	if err != nil {
		log.Printf("Error at FuncName: %v HandleCodeEncrypt in tools.Encrypt err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("could not generate your code")
	}
	return
}

func HandleCodeDecrypt(server *Server, encrypted string, id string, funcName string) (code string, err error) {
	hashKey := []byte(server.config.TokenSymmetricKey)
	// Hash the code using same format with password
	code, err = tools.Decrypt(hashKey, encrypted)
	if err != nil {
		log.Printf("Error at FuncName: %v HandleCodeDecrypt in tools.Decrypt err: %v, user: %v\n", funcName, err, id)
		err = errors.New("could not validate this code")
	}
	return
}

func CreateAndValidateChargeCodeTicket(ctx context.Context, server *Server, req ValidateChargeCodeParams, funcName string, user db.User, guestID uuid.UUID, chargeID uuid.UUID, chargeOptionUserID uuid.UUID, optionID uuid.UUID) (chargeMain db.GetScannedChargeTicketByHostRow, used bool, err error) {
	eventDateTimeID, err := tools.StringToUuid(req.EventDateTimeID)
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeTicket in tools.StringToUuid err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("this event date cannot be found")
		return
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(req.StartDate)
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeTicket in tools.ConvertDateOnlyStringToDate err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("this event date cannot be found")
		return
	}
	chargeMain, err = server.store.GetScannedChargeTicketByHost(ctx, db.GetScannedChargeTicketByHostParams{
		MainHostID:        user.ID,
		CoUserID:          tools.UuidToString(user.UserID),
		OptionID:          optionID,
		OptionCoHostID:    optionID,
		ChargeTicketID:    chargeID,
		ChargeStartDate:   startDate,
		ChargeEventDateID: eventDateTimeID,
		GuestUserID:       guestID,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeTicket in store.GetScannedChargeTicketByHost err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("your not allowed to verify this code")
		return
	}

	if !chargeMain.ScanCode {
		err = errors.New("your not allowed to verify this code")
		return
	}
	if chargeMain.OptionUserID != chargeOptionUserID {
		err = errors.New("this code isn't properly registered with this event")
		return
	}
	if chargeMain.TicketScanned {
		used = true
		return
	}
	_, err = server.store.CreateScannedCharge(ctx, db.CreateScannedChargeParams{
		ChargeID:    chargeID,
		Scanned:     true,
		ScannedBy:   user.UserID,
		ScannedTime: time.Now().Add(time.Hour),
		ChargeType:  constants.CHARGE_TICKET_REFERENCE,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeTicket in CreateScannedCharge err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("this code could not be verified")
		return
	}
	return
}

func CreateAndValidateChargeCodeOption(ctx context.Context, server *Server, req ValidateChargeCodeParams, funcName string, user db.User, guestID uuid.UUID, chargeID uuid.UUID, chargeOptionUserID uuid.UUID, optionID uuid.UUID) (chargeMain db.GetScannedChargeOptionByHostRow, used bool, err error) {
	chargeMain, err = server.store.GetScannedChargeOptionByHost(ctx, db.GetScannedChargeOptionByHostParams{
		MainHostID:     user.ID,
		CoUserID:       tools.UuidToString(user.UserID),
		OptionID:       optionID,
		OptionCoHostID: optionID,
		ChargeOptionID: chargeID,
		GuestUserID:    guestID,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeOption in store.GetScannedChargeOptionByHost err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("your not allowed to verify this code")
		return
	}

	if !chargeMain.ScanCode {
		err = errors.New("your not allowed to verify this code")
		return
	}
	if chargeMain.OptionUserID != chargeOptionUserID {
		err = errors.New("this code isn't properly registered with this event")
		return
	}
	if chargeMain.StayScanned {
		used = true
		return
	}
	_, err = server.store.CreateScannedCharge(ctx, db.CreateScannedChargeParams{
		ChargeID:    chargeID,
		Scanned:     true,
		ScannedBy:   user.UserID,
		ScannedTime: time.Now().Add(time.Hour),
		ChargeType:  constants.CHARGE_OPTION_REFERENCE,
	})
	if err != nil {
		log.Printf("Error at FuncName: %v CreateAndValidateChargeCodeOption in CreateScannedCharge err: %v, user: %v\n", funcName, err, user.ID)
		err = errors.New("this code could not be verified")
		return
	}
	return
}
