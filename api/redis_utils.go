package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/google/uuid"
)

// This stores the date in redis
func HandleOptionReserveRedis(user db.User, canInstantBook bool, datePriceFloat []DatePriceFloat, totalDatePrice float64, discountType string, discount float64, cleanFee float64, extraGuestFee float64, petFee float64, petStayFee float64, extraGuestStayFee float64, totalPrice float64, serviceFee float64, requireRequest bool, requestType string, userCurrency string, optionUserID uuid.UUID, startDate string, endDate string, guests []string) (reference string, err error) {
	reference = tools.UuidToString(uuid.New())
	userID := tools.UuidToString(user.ID)
	//mainReference would hold reference to main redis data
	mainReference := userID + "&" + reference
	dateReference := tools.UuidToString(uuid.New())
	dateReferenceList := []string{}
	guestString := ""
	// DatePriceFloat
	for _, date := range datePriceFloat {
		// First we generate a reference we would use as the id
		referenceNew := tools.UuidToString(uuid.New())
		data := []string{
			constants.DATE_PRICE_FLOAT,
			tools.ConvertFloatToString(date.Price),
			constants.DATE_PRICE_DATE,
			date.Date,
			constants.DATE_PRICE_GROUP_PRICE,
			tools.ConvertFloatToString(date.GroupPrice),
		}
		errData := RedisClient.HSet(RedisContext, referenceNew, data).Err()

		if errData != nil {
			err = errData
			log.Printf("HandleOptionReserveRedis HSet for mainReference %v err:%v\n", mainReference, err.Error())
		}
		dateReferenceList = append(dateReferenceList, referenceNew)
	}
	err = RedisClient.SAdd(RedisContext, dateReference, dateReferenceList).Err()
	if err != nil {
		log.Printf("HandleOptionReserveRedis SAdd for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	for i := 0; i < len(guests); i++ {
		if i == (len(guests) - 1) {
			guestString += guests[i]
		} else {
			guestString += guests[i] + "&"
		}
	}
	mainData := []string{
		constants.START_DATE,
		startDate,
		constants.END_DATE,
		endDate,
		constants.OPTION_USER_ID,
		tools.UuidToString(optionUserID),
		constants.CAN_INSTANT_BOOK,
		tools.ConvertBoolToString(canInstantBook),
		constants.DATE_PRICE_FLOAT,
		dateReference,
		constants.TOTAL_DATE_PRICE,
		tools.ConvertFloatToString(totalDatePrice),
		constants.GUESTS,
		guestString,
		constants.DISCOUNT,
		tools.ConvertFloatToString(discount),
		constants.DISCOUNT_TYPE,
		discountType,
		constants.CLEAN_FEE,
		tools.ConvertFloatToString(cleanFee),
		constants.EXTRA_GUEST_FEE,
		tools.ConvertFloatToString(extraGuestFee),
		constants.PET_FEE,
		tools.ConvertFloatToString(petFee),
		constants.PET_STAY_FEE,
		tools.ConvertFloatToString(petStayFee),
		constants.EXTRA_GUEST_STAY_FEE,
		tools.ConvertFloatToString(extraGuestStayFee),
		constants.TOTAL_PRICE,
		tools.ConvertFloatToString(totalPrice),
		constants.SERVICE_FEE,
		tools.ConvertFloatToString(serviceFee),
		constants.REQUIRE_REQUEST,
		tools.ConvertBoolToString(requireRequest),
		constants.REQUEST_TYPE,
		requestType,
		constants.USER_CURRENCY,
		userCurrency,
	}
	err = RedisClient.HSet(RedisContext, mainReference, mainData).Err()
	if err != nil {
		log.Printf("HandleOptionReserveRedis HSet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	reserveTime := tools.TimeReserveUser()
	timeData := []string{
		constants.MAIN_REFERENCE,
		mainReference,
		constants.TIME,
		tools.ConvertTimeToString(reserveTime),
	}
	id := tools.UuidToString(uuid.New())
	err = RedisClient.HSet(RedisContext, id, timeData).Err()
	if err != nil {
		log.Printf("HandleOptionReserveRedis RedisClient.HSe %v err:%v\n", mainReference, err.Error())
		return
	}
	err = RedisClient.SAdd(RedisContext, constants.REMOVE_OPTION_RESERVE_USER, id).Err()
	if err != nil {
		log.Printf("HandleOptionReserveRedis SAdd(RedisContext, constants.REMOVE_OPTION_RESERVE_USER, id).Err() %v err:%v\n", mainReference, err.Error())
		return
	}
	return
}

func HandleEventReserveRedis(user db.User, eventData EventDateReserveDB) (reference string, err error) {
	reference = tools.UuidToString(uuid.New())
	userID := tools.UuidToString(user.ID)
	//mainReference would hold reference to main redis data
	mainReference := userID + "&" + reference
	dateRedisID := tools.UuidToString(uuid.New())
	dateReferences := []string{}
	for _, date := range eventData.DateTimes {
		// First we start by storing the ticket reserve item in DB
		ticketRedisID := tools.UuidToString(uuid.New())
		ticketReferences := []string{}
		for _, ticket := range date.Tickets {
			ticketID := fmt.Sprintf("%v&%v", ticket.ID, uuid.New())
			data := []string{
				constants.ID,
				ticket.ID,
				constants.GRADE,
				ticket.Grade,
				constants.PAY_TYPE,
				ticket.PayType,
				constants.PRICE,
				tools.ConvertFloatToString(ticket.Price),
				constants.TYPE,
				ticket.Type,
				constants.TICKET_SERVICE_FEE,
				tools.ConvertFloatToString(ticket.ServiceFee),
				constants.TICKET_ABSORB_FEE,
				tools.ConvertFloatToString(ticket.AbsorbFees),
				constants.EVENT_DATE_GROUP_PRICE,
				tools.ConvertFloatToString(ticket.GroupPrice),
			}
			errTicket := RedisClient.HSet(RedisContext, ticketID, data).Err()
			if errTicket != nil {
				err = errTicket
				log.Printf("HandleEventReserveRedis ticket HSet for mainReference %v err:%v\n", mainReference, err.Error())
				return
			}
			ticketReferences = append(ticketReferences, ticketID)
		}
		//if err != nil {
		//	return
		//}
		// Store all ticketReferences in the database
		errTicketReference := RedisClient.SAdd(RedisContext, ticketRedisID, ticketReferences).Err()

		if errTicketReference != nil {
			err = errTicketReference
			log.Printf("HandleEventReserveRedis ticket SAdd for mainReference %v err:%v\n", mainReference, err.Error())
			return
		}

		// Next we setup date
		dateID := fmt.Sprintf("%v&%v", date.ID, uuid.New())
		data := []string{
			constants.ID,
			date.ID,
			constants.START_DATE,
			date.StartDate,
			constants.END_DATE,
			date.EndDate,
			constants.START_TIME,
			date.StartTime,
			constants.END_TIME,
			date.EndTime,
			constants.TIME_ZONE,
			date.TimeZone,
			constants.TOTAL_DATE_FEE,
			tools.ConvertFloatToString(date.TotalDateFee),
			constants.TOTAL_DATE_ABSORB_FEE,
			tools.ConvertFloatToString(date.TotalDateAbsorbFee),
			constants.TOTAL_DATE_SERVICE_FEE,
			tools.ConvertFloatToString(date.TotalDateServiceFee),
			constants.TICKETS,
			ticketRedisID,
		}
		errDate := RedisClient.HSet(RedisContext, dateID, data).Err()
		if errDate != nil {
			err = errDate
			log.Printf("HandleEventReserveRedis date HSet for mainReference %v err:%v\n", mainReference, err.Error())
			return
		}
		dateReferences = append(dateReferences, dateID)
	}
	//if err != nil {
	//	return
	//}
	// We save the dateReferences
	err = RedisClient.SAdd(RedisContext, dateRedisID, dateReferences).Err()
	if err != nil {
		log.Printf("HandleEventReserveRedis eventData SAdd for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	// Lastly we setup eventData
	eventDataID := eventData.ID
	data := []string{
		constants.ID,
		eventDataID,
		constants.CURRENCY,
		eventData.Currency,
		constants.TOTAL_FEE,
		tools.ConvertFloatToString(eventData.TotalFee),
		constants.SERVICE_FEE,
		tools.ConvertFloatToString(eventData.TotalServiceFee),
		constants.TOTAL_ABSORB_FEE,
		tools.ConvertFloatToString(eventData.TotalAbsorbFee),
		constants.DATE_TIMES,
		dateRedisID,
	}
	err = RedisClient.HSet(RedisContext, mainReference, data).Err()
	if err != nil {
		log.Printf("HandleEventReserveRedis eventData HSet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	reserveTime := tools.TimeReserveUser()
	timeData := []string{
		constants.MAIN_REFERENCE,
		mainReference,
		constants.TIME,
		tools.ConvertTimeToString(reserveTime),
	}
	id := tools.UuidToString(uuid.New())
	err = RedisClient.HSet(RedisContext, id, timeData).Err()
	if err != nil {
		log.Printf("HandleEventReserveRedis RedisClient.HSe %v err:%v\n", mainReference, err.Error())
		return
	}
	err = RedisClient.SAdd(RedisContext, constants.REMOVE_EVENT_RESERVE_USER, id).Err()
	if err != nil {
		log.Printf("HandleEventReserveRedis SAdd(RedisContext, constants.REMOVE_EVENT_RESERVE_USER, id).Err() %v err:%v\n", mainReference, err.Error())
		return
	}
	return
}

// This gets the data using EventDateReserve format
func HandleEventReserveRedisData(uID, reference string) (reserveData EventDateReserve, err error) {
	mainReference := uID + "&" + reference
	eventData, err := RedisClient.HGetAll(RedisContext, mainReference).Result()
	if err != nil {
		log.Printf("HandleEventReserveRedisData HGet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	dateReferences, err := RedisClient.SMembers(RedisContext, eventData[constants.DATE_TIMES]).Result()
	if err != nil {
		log.Printf("HandleEventReserveRedisData SMembers for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	dates := []DateReserveItem{}
	for _, dateRef := range dateReferences {
		// We want to handle tickets
		dateData, errDate := RedisClient.HGetAll(RedisContext, dateRef).Result()
		if errDate != nil {
			err = errDate
			log.Printf("HandleEventReserveRedisData HGet errDate for DateRef %v err:%v\n", dateRef, err.Error())
			return
		}
		ticketReferences, errTicketRef := RedisClient.SMembers(RedisContext, dateData[constants.TICKETS]).Result()
		if errTicketRef != nil {
			err = errTicketRef
			log.Printf("HandleEventReserveRedisData SMembers errTicketRef for ticketRedisID %v err:%v\n", dateData[constants.TICKETS], err.Error())
			return
		}
		tickets := []TicketReserveItem{}
		for _, ticketRef := range ticketReferences {
			// We setup TicketReserveItem
			ticket, errTicketItem := RedisClient.HGetAll(RedisContext, ticketRef).Result()
			if errTicketItem != nil {
				err = errTicketItem
				log.Printf("HandleEventReserveRedisData HGetAll ticketRef for ticketRef %v err:%v\n", ticketRef, err.Error())
			}
			data := TicketReserveItem{
				ID:         ticket[constants.ID],
				Grade:      ticket[constants.GRADE],
				Price:      ticket[constants.PRICE],
				Type:       ticket[constants.TYPE],
				PayType:    ticket[constants.PAY_TYPE],
				ServiceFee: ticket[constants.TICKET_SERVICE_FEE],
				AbsorbFee:  ticket[constants.TICKET_ABSORB_FEE],
				GroupPrice: ticket[constants.EVENT_DATE_GROUP_PRICE],
			}
			tickets = append(tickets, data)

		}
		if err != nil {
			return
		}

		// We want to setup DateReserveItem
		data := DateReserveItem{
			ID:                  dateData[constants.ID],
			StartDate:           dateData[constants.START_DATE],
			EndDate:             dateData[constants.END_DATE],
			StartTime:           dateData[constants.START_TIME],
			EndTime:             dateData[constants.END_TIME],
			TimeZone:            dateData[constants.TIME_ZONE],
			TotalDateServiceFee: dateData[constants.TOTAL_DATE_SERVICE_FEE],
			TotalDateFee:        dateData[constants.TOTAL_DATE_FEE],
			TotalDateAbsorbFee:  dateData[constants.TOTAL_DATE_ABSORB_FEE],
			Tickets:             tickets,
		}

		dates = append(dates, data)

	}
	if err != nil {
		return
	}
	reserveData = EventDateReserve{
		ID:              eventData[constants.ID],
		DateTimes:       dates,
		Currency:        eventData[constants.CURRENCY],
		TotalFee:        eventData[constants.TOTAL_FEE],
		TotalServiceFee: eventData[constants.SERVICE_FEE],
		TotalAbsorbFee:  eventData[constants.TOTAL_ABSORB_FEE],
	}
	return
}

// This gets the data using ExperienceReserveOModel format
func HandleOptionReserveRedisData(userID string, reference string) (reserveData ExperienceReserveOModel, err error) {
	var datePrice []DatePrice
	// We start with datePrice
	mainReference := userID + "&" + reference
	// We date the date price float reference
	dateReference, err := RedisClient.HGet(RedisContext, mainReference, constants.DATE_PRICE_FLOAT).Result()
	if err != nil {
		log.Printf("HandleOptionReserveRedisData HGet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	// We get all the references for the dateReference
	dateFloatReferences, err := RedisClient.SMembers(RedisContext, dateReference).Result()
	if err != nil {
		log.Printf("HandleOptionReserveRedisData SMembers for dateReference %v err:%v\n", dateReference, err.Error())
	}
	for _, dateRef := range dateFloatReferences {
		dateRefData, err := RedisClient.HGetAll(RedisContext, dateRef).Result()
		if err != nil {
			log.Printf("HandleOptionReserveRedisData HGetAll for dateRef %v err:%v\n", dateRef, err.Error())
			break
		}
		data := DatePrice{
			Price:      dateRefData[constants.DATE_PRICE_FLOAT],
			Date:       dateRefData[constants.DATE_PRICE_DATE],
			GroupPrice: dateRefData[constants.DATE_PRICE_GROUP_PRICE],
		}
		datePrice = append(datePrice, data)
	}
	// This err checks if there was any error in the loop
	if err != nil {
		return
	}
	mainData, err := RedisClient.HGetAll(RedisContext, mainReference).Result()
	if err != nil {
		return
	}
	discount := ReDiscount{
		Price: mainData[constants.DISCOUNT],
		Type:  mainData[constants.DISCOUNT_TYPE],
	}
	guests := strings.Split(mainData[constants.GUESTS], "&")
	reserveData = ExperienceReserveOModel{
		Discount:        discount,
		MainPrice:       mainData[constants.TOTAL_DATE_PRICE],
		ServiceFee:      mainData[constants.SERVICE_FEE],
		TotalFee:        mainData[constants.TOTAL_PRICE],
		DatePrice:       datePrice,
		Currency:        mainData[constants.USER_CURRENCY],
		Guests:          guests,
		GuestFee:        mainData[constants.EXTRA_GUEST_FEE],
		PetFee:          mainData[constants.PET_FEE],
		CleaningFee:     mainData[constants.CLEAN_FEE],
		NightlyPetFee:   mainData[constants.PET_STAY_FEE],
		NightlyGuestFee: mainData[constants.EXTRA_GUEST_STAY_FEE],
		CanInstantBook:  tools.ConvertStringToBool(mainData[constants.CAN_INSTANT_BOOK]),
		RequireRequest:  tools.ConvertStringToBool(mainData[constants.REQUIRE_REQUEST]),
		RequestType:     mainData[constants.REQUEST_TYPE],
		Reference:       reference,
		OptionUserID:    mainData[constants.OPTION_USER_ID],
		StartDate:       mainData[constants.START_DATE],
		EndDate:         mainData[constants.END_DATE],
	}
	return
}

func TimerRemoveEventReserveUser(mainReference string) (err error) {
	dateRedisID, err := RedisClient.HGet(RedisContext, mainReference, constants.DATE_TIMES).Result()
	if err != nil {
		log.Printf("TimeRemoveEventReserveUser HGet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	dateReferences, err := RedisClient.SMembers(RedisContext, dateRedisID).Result()
	if err != nil {
		log.Printf("TimeRemoveEventReserveUser SMembers for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	for _, dateRef := range dateReferences {
		// We want to handle tickets
		ticketRedisID, errDateTicket := RedisClient.HGet(RedisContext, dateRef, constants.TICKETS).Result()
		if errDateTicket != nil {
			err = errDateTicket
			log.Printf("TimeRemoveEventReserveUser HGet errDateTicket for DateRef %v err:%v\n", dateRef, err.Error())
			return
		}
		ticketReferences, errTicketRef := RedisClient.SMembers(RedisContext, ticketRedisID).Result()
		if errTicketRef != nil {
			err = errTicketRef
			log.Printf("TimeRemoveEventReserveUser SMembers errTicketRef for ticketRedisID %v err:%v\n", ticketRedisID, err.Error())
			return
		}
		for _, ticketRef := range ticketReferences {
			// We want to remove each ticket
			err = RedisClient.Del(RedisContext, ticketRef).Err()
			if err != nil {
				log.Printf("TimeRemoveEventReserveUser Del for ticketRef %v err:%v\n", ticketRef, err.Error())
			}

		}
		// After we are done deleting tickets
		// We want to delete the references to the tickets
		err = RedisClient.Del(RedisContext, ticketRedisID).Err()
		if err != nil {
			log.Printf("TimeRemoveEventReserveUser Del for ticketRedisID %v err:%v\n", ticketRedisID, err.Error())
		}
		// We want to delete each date data
		err = RedisClient.Del(RedisContext, dateRef).Err()
		if err != nil {
			log.Printf("TimeRemoveEventReserveUser Del for dateRef %v err:%v\n", dateRef, err.Error())
		}

	}

	// We want to delete all date references
	err = RedisClient.Del(RedisContext, dateRedisID).Err()
	if err != nil {
		log.Printf("TimeRemoveEventReserveUser Del for dateRedisID %v err:%v\n", dateRedisID, err.Error())
		return
	}
	// We want to delete main data
	err = RedisClient.Del(RedisContext, mainReference).Err()
	if err != nil {
		log.Printf("TimeRemoveEventReserveUser Del for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	err = nil
	return
}

func TimerRemoveOptionReserveUser(mainReference string) (err error) {
	// We date the date price float reference
	dateReference, err := RedisClient.HGet(RedisContext, mainReference, constants.DATE_PRICE_FLOAT).Result()
	if err != nil {
		log.Printf("TimerRemoveOptionReserveUser HGet for mainReference %v err:%v\n", mainReference, err.Error())
		return
	}

	// We get all the references for the dateReference
	dateFloatReferences, err := RedisClient.SMembers(RedisContext, dateReference).Result()
	if err != nil {
		log.Printf("TimerRemoveOptionReserveUser SMembers for dateReference %v err:%v\n", dateReference, err.Error())
		return
	}
	for _, dateRef := range dateFloatReferences {
		// We delete all for the HSet for each dateRef
		err = RedisClient.Del(RedisContext, dateRef).Err()
		if err != nil {
			log.Printf("TimerRemoveOptionReserveUser for dateReF Del dateRef %v err:%v\n", dateRef, err.Error())
		}
	}
	// We then delete the date references using dateReference
	err = RedisClient.Del(RedisContext, dateReference).Err()
	if err != nil {
		log.Printf("TimerRemoveOptionReserveUser for Del mainReference %v err:%v\n", mainReference, err.Error())
		return
	}

	// We then delete the main data using mainReference
	err = RedisClient.Del(RedisContext, mainReference).Err()
	if err != nil {
		log.Printf("TimerRemoveOptionReserveUser for Del mainReference %v err:%v\n", mainReference, err.Error())
		return
	}
	err = nil
	return
}
