package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleReserveEvent(user db.User, server *Server, ctx *gin.Context, tickets []TicketItem, optionUserID string, userCurrency string) (reference string, err error) {
	_, eventData, err := ReserveEventCalculate(user, server, ctx, tickets, optionUserID, userCurrency)
	if err != nil {
		return
	}
	reference, err = HandleEventReserveRedis(user, eventData)
	return
}

func ReserveEventCalculate(user db.User, server *Server, ctx *gin.Context, tickets []TicketItem, optionUserID string, userCurrency string) (optionID uuid.UUID, eventData EventDateReserveDB, err error) {
	// Server Currency setup
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	var servicePercent float64
	if userCurrency != utils.NGN {
		servicePercent = tools.ConvertStringToFloat(server.config.IntServiceEventUserPercent)
	} else {
		servicePercent = tools.ConvertStringToFloat(server.config.LocServiceEventUserPercent)
	}
	// First let us group the tickets to make sure we only one occurrence
	optionID, err = tools.StringToUuid(optionUserID)
	if err != nil {
		log.Printf("Error at ReserveEventCalculate in StringToUuid %v for user: %v. optionID: %v\n", err.Error(), user.ID, optionUserID)
		err = fmt.Errorf("the event does not exist")
		return
	}

	ticketData := make(map[string][]TicketItem)

	for _, t := range tickets {
		ticketID := fmt.Sprintf("%v&%v", t.TicketID, t.StartDate)
		ticketData[ticketID] = append(ticketData[ticketID], t)
	}

	// We want to check if this tickets are available for the user
	for ticketID, ticketItem := range ticketData {
		id, startDate := strings.Split(ticketID, "&")[0], strings.Split(ticketID, "&")[1]
		err = TicketQuantityAvailable(server, ctx, len(ticketItem), id, ticketItem[0].EventDateID, startDate)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	ticketMainData := make(map[string]GetTicketByIDAndOptionIDRow)
	// Now we are going to build the data
	for ticketDateID, ticket := range ticketData {
		id, startDate := strings.Split(ticketDateID, "&")[0], strings.Split(ticketDateID, "&")[1]
		ticketID, errTicket := tools.StringToUuid(id)
		if errTicket != nil {
			log.Printf("Error at ReserveEventCalculate in StringToUuid %v for user: %v. ticket: %v\n", errTicket.Error(), user.ID, ticketID)
			err = fmt.Errorf("one of the tickets selected does not exist")
			return
		}
		eventDateID, errEventDate := tools.StringToUuid(ticket[0].EventDateID)
		if errEventDate != nil {
			log.Printf("Error at ReserveEventCalculate date in StringToUuid %v for user: %v. ticket: %v\n", errEventDate.Error(), user.ID, ticketID)
			err = fmt.Errorf("event date selected does not exist")
			return
		}
		data, errOption := server.store.GetTicketByIDAndOptionID(ctx, db.GetTicketByIDAndOptionIDParams{
			OptionUserID: optionID,
			TicketID:     ticketID,
			EventDateID:  eventDateID,
		})
		if errOption != nil {
			log.Printf("Error at ReserveEventCalculate in GetTicketByIDAndOptionID %v for user: %v. ticket: %v\n", errOption.Error(), user.ID, ticketID)
			err = fmt.Errorf("ticket is unavailable at this moment")
			return
		}

		setData := GetTicketByIDAndOptionIDRow{
			StartDate:    startDate,
			TicketID:     data.TicketID,
			Level:        data.Level,
			Price:        data.Price,
			TicketType:   data.TicketType,
			TicketName:   data.TicketName,
			EventDateID:  data.EventDateID,
			OptionUserID: data.OptionUserID,
			Currency:     data.Currency,
			StartTime:    data.StartTime,
			EndTime:      data.EndTime,
			TimeZone:     data.TimeZone,
			PayType:      data.PayType,
			AbsorbFees:   data.AbsorbFees,
		}
		log.Println("data option user id", data.OptionUserID)
		ticketMainData[ticketDateID] = setData
	}
	eventData, err = HandleReserveEventData(user, ticketData, ticketMainData, dollarToNaira, dollarToCAD, servicePercent, userCurrency)
	if err != nil {
		return
	}
	return
}

func GetEventDateData(user db.User, eventDateID string, ticketIDs []string, ticketData map[string][]TicketItem, ticketMainData map[string]GetTicketByIDAndOptionIDRow, userCurrency string, dollarToNaira string, dollarToCAD string, servicePercent float64) (eventDate DateReserveItemDB, err error) {
	var tickets []TicketReserveItemDB
	var totalPrice float64
	var totalServicePrice float64
	var totalDateAbsorbFee float64
	var data GetTicketByIDAndOptionIDRow
	var startDate string
	var endDate string
	for _, ticketID := range ticketIDs {
		data = ticketMainData[ticketID]
		startDate = ticketData[ticketID][0].StartDate
		endDate = ticketData[ticketID][0].EndDate
		var serviceFee float64
		var absorbFee float64
		price, errPrice := tools.ConvertPrice(tools.IntToMoneyString(data.Price), data.Currency, userCurrency, dollarToNaira, dollarToCAD, data.TicketID)
		if errPrice != nil {
			err = errPrice
			log.Printf("Error at ReserveEventCalculate in GetTicketByIDAndOptionID %v for user: %v. ticket: %v\n", err.Error(), user.ID, ticketID)
			return
		}
		if data.AbsorbFees {
			absorbFee = price * (servicePercent / 100)
		} else {
			serviceFee = price * (servicePercent / 100)
		}

		res := TicketReserveItemDB{
			ID:         tools.UuidToString(data.TicketID),
			Grade:      data.Level,
			Price:      price,
			ServiceFee: serviceFee,
			Type:       data.TicketType,
			PayType:    data.PayType,
			AbsorbFees: absorbFee,
			GroupPrice: price * float64(len(ticketData[ticketID])),
		}
		for i := 0; i < len(ticketData[ticketID]); i++ {
			tickets = append(tickets, res)
		}
		totalPrice += price * float64(len(ticketData[ticketID]))
		totalServicePrice += serviceFee * float64(len(ticketData[ticketID]))
		totalDateAbsorbFee += absorbFee * float64(len(ticketData[ticketID]))
	}
	eventDate = DateReserveItemDB{
		ID:                  tools.UuidToString(data.EventDateID),
		StartDate:           startDate,
		EndDate:             endDate,
		StartTime:           data.StartTime,
		EndTime:             data.EndTime,
		TotalDateFee:        totalPrice,
		TotalDateServiceFee: totalServicePrice,
		TotalDateAbsorbFee:  totalDateAbsorbFee,
		TimeZone:            data.TimeZone,
		Tickets:             tickets,
	}
	return
}

func HandleReserveEventData(user db.User, ticketData map[string][]TicketItem, ticketMainData map[string]GetTicketByIDAndOptionIDRow, dollarToNaira string, dollarToCAD string, servicePercent float64, userCurrency string) (eventData EventDateReserveDB, err error) {
	// First let us ground events in the same date time together
	eventDateTimes := make(map[string][]string)
	var eventID string
	for ticketDateID, data := range ticketMainData {
		eventID = tools.UuidToString(data.OptionUserID)
		dateTimeID := fmt.Sprintf("%v&%v", data.EventDateID, data.StartDate)

		// So that ticket with same start date and event_date_id are grouped together
		eventDateTimes[dateTimeID] = append(eventDateTimes[dateTimeID], ticketDateID)
	}
	var dateTimes []DateReserveItemDB
	var totalFee float64 = 0
	var totalServiceFee float64 = 0
	var totalAbsorbFee float64 = 0
	for dateTime, ticketIDs := range eventDateTimes {
		// Lets first start by arranging ticket data
		dateTimeID := strings.Split(dateTime, "&")[0]
		eventDate, errDate := GetEventDateData(user, dateTimeID, ticketIDs, ticketData, ticketMainData, userCurrency, dollarToNaira, dollarToCAD, servicePercent)
		if errDate != nil {
			err = errDate
			log.Printf("Error at HandleReserveEventData in GetEventDateData %v for user: %v. dateTimeID: %v\n", err.Error(), user.ID, dateTimeID)
			return
		}
		log.Println("eventDate data ", eventDate)
		totalFee += eventDate.TotalDateFee
		totalServiceFee += eventDate.TotalDateServiceFee
		totalAbsorbFee += eventDate.TotalDateAbsorbFee
		dateTimes = append(dateTimes, eventDate)
	}
	totalFee += totalServiceFee

	eventData = EventDateReserveDB{
		ID:              eventID,
		DateTimes:       dateTimes,
		Currency:        userCurrency,
		TotalFee:        totalFee,
		TotalServiceFee: totalServiceFee,
		TotalAbsorbFee:  totalAbsorbFee,
	}
	log.Println("main eventData data ", eventData.ID)
	return
}

// We want to save a recept in the database
// We also want to store a snap shot of what the event looks like
func HandleEventReserveComplete(server *Server, ctx *gin.Context, reserveData EventDateReserve, paystackReference string, user db.User, msg string, reference string) (err error) {
	// First we store the receipt
	optionUserID, chargeID, err := HandleEventReserveReceipt(server, ctx, reserveData, paystackReference, user, true, reference, "HandleEventReserveComplete")
	if err != nil {
		log.Printf("Error at HandleEventReserveComplete in HandleEventReserveReceipt: %v optionID: %v referenceID: %v, paystackReference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		return
	}
	// We want to store a snap shot of the option
	err = HandleEventSnapShot(server, ctx, optionUserID, reference, paystackReference, chargeID)
	if err != nil {
		log.Printf("Error at HandleEventReserveComplete in HandleEventSnapShot: %v optionID: %v referenceID: %v, paystackReference: %v\n", err.Error(), optionUserID, reference, paystackReference)
	}
	err = nil
	return
}

func HandleEventReserveReceipt(server *Server, ctx *gin.Context, reserveData EventDateReserve, paystackReference string, user db.User, isComplete bool, reference string, functionName string) (optionUserID uuid.UUID, chargeID uuid.UUID, err error) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	optionUserID, err = tools.StringToUuid(reserveData.ID)
	if err != nil {
		log.Printf("Error at HandleEventReserveReceipt in StringToUuid: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.ID, reference, paystackReference, functionName)
		err = fmt.Errorf("error 651 occur, pls contact us")
		return
	}
	// First we start with create event reference
	eventData, err := server.store.CreateChargeEventReference(ctx, db.CreateChargeEventReferenceParams{
		UserID:           user.UserID,
		OptionUserID:     optionUserID,
		TotalFee:         tools.MoneyStringToInt(reserveData.TotalFee),
		ServiceFee:       tools.MoneyStringToInt(reserveData.TotalServiceFee),
		TotalAbsorbFee:   tools.MoneyStringToInt(reserveData.TotalAbsorbFee),
		DateBooked:       time.Now().Add(time.Hour),
		Currency:         reserveData.Currency,
		Reference:        reference,
		PaymentReference: paystackReference,
		IsComplete:       isComplete,
	})

	if err != nil {
		log.Printf("Error at HandleEventReserveReceipt inCreateChargeEventReference: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.ID, reference, paystackReference, functionName)
		err = fmt.Errorf("error 652 occur, pls contact us")
		return
	}
	chargeID = eventData.ID

	for _, date := range reserveData.DateTimes {
		dateID, err := tools.StringToUuid(date.ID)
		if err != nil {
			log.Printf("Error at HandleEventReserveReceipt in StringToUuid: %v dateID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), date.ID, reference, paystackReference, functionName)
			continue
		}
		startDate, err := tools.ConvertDateOnlyStringToDate(date.StartDate)
		if err != nil {
			log.Printf("Error at HandleEventReserveReceipt in ConvertDateOnlyStringToDate: %v dateID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), date.ID, reference, paystackReference, functionName)
			startDate = time.Now()
		}
		endDate, err := tools.ConvertDateOnlyStringToDate(date.EndDate)
		if err != nil {
			log.Printf("Error at HandleEventReserveReceipt in ConvertDateOnlyStringToDate: %v dateID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), date.ID, reference, paystackReference, functionName)
			endDate = time.Now()
		}
		dateData, err := server.store.CreateChargeDateReference(ctx, db.CreateChargeDateReferenceParams{
			ChargeEventID:       eventData.ID,
			EventDateID:         dateID,
			StartDate:           startDate,
			TotalDateServiceFee: tools.MoneyStringToInt(date.TotalDateServiceFee),
			TotalDateAbsorbFee:  tools.MoneyStringToInt(date.TotalDateAbsorbFee),
			EndDate:             endDate,
			StartTime:           date.StartTime,
			DateBooked:          time.Now().Add(time.Hour),
			EndTime:             date.EndTime,
			TotalDateFee:        tools.MoneyStringToInt(date.TotalDateFee),
		})
		if err != nil {
			log.Printf("Error at HandleEventReserveReceipt in CreateChargeDateReference: %v optionID: %v dateID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), date.ID, reference, paystackReference, functionName)
		}
		// We handle tickets
		for _, ticket := range date.Tickets {
			ticketID, err := tools.StringToUuid(ticket.ID)
			if err != nil {
				log.Printf("Error at HandleEventReserveReceipt in StringToUuid: %v ticketID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), ticket.ID, reference, paystackReference, functionName)
				continue
			}
			chargeTicket, err := server.store.CreateChargeTicketReference(ctx, db.CreateChargeTicketReferenceParams{
				ChargeDateID: dateData.ID,
				TicketID:     ticketID,
				Grade:        ticket.Grade,
				ServiceFee:   tools.MoneyStringToInt(ticket.ServiceFee),
				AbsorbFee:    tools.MoneyStringToInt(ticket.AbsorbFee),
				Price:        tools.MoneyStringToInt(ticket.Price),
				DateBooked:   time.Now().Add(time.Hour),
				Type:         ticket.PayType,
				TicketType:   ticket.Type,
				GroupPrice:   tools.MoneyStringToInt(ticket.GroupPrice),
			})
			if err != nil {
				log.Printf("Error at HandleEventReserveReceipt inCreateChargeTicketReference: %v ticketID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), ticket.ID, reference, paystackReference, functionName)
				continue
			}
			//if totalFee is less than 1 we want to make payoutComplete true so we do have to payout the host
			var payoutComplete bool
			var payoutAmount float64
			var servicePercent float64
			var serviceFee float64
			amount := tools.ConvertStringToFloat(tools.IntToMoneyString(chargeTicket.Price)) - tools.ConvertStringToFloat(tools.IntToMoneyString(chargeTicket.ServiceFee))
			amount = amount - tools.ConvertStringToFloat(tools.IntToMoneyString(chargeTicket.AbsorbFee))
			amount, err = tools.ConvertPrice(tools.ConvertFloatToString(amount), eventData.Currency, utils.PayoutCurrency, dollarToNaira, dollarToCAD, user.ID)
			if err != nil {
				log.Printf("Error at HandleEventReserveReceipt for tools.ConvertPrice: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.ID, reference, paystackReference, functionName)
				continue
			}
			switch utils.PayoutCurrency {
			case utils.NGN:
				// server.config.LocServiceEventHostPercent is at 0.0
				servicePercent = tools.ConvertStringToFloat(server.config.LocServiceEventHostPercent)
			default:
				servicePercent = tools.ConvertStringToFloat(server.config.IntServiceEventHostPercent)
			}
			if tools.ConvertStringToFloat(ticket.Price) < 0.1 {
				payoutComplete = true
			} else {
				serviceFee = (servicePercent / 100) * amount
				payoutAmount = amount - serviceFee
			}
			err = server.store.CreateMainPayout(ctx, db.CreateMainPayoutParams{
				ChargeID:   chargeTicket.ID,
				Type:       constants.CHARGE_TICKET_REFERENCE,
				IsComplete: payoutComplete,
				Amount:     tools.MoneyFloatToInt(payoutAmount),
				ServiceFee: tools.MoneyFloatToInt(serviceFee),
				Currency:   utils.PayoutCurrency,
			})
			if err != nil {
				log.Printf("Error at HandleEventReserveReceipt for ticket CreateMainPayout: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.ID, reference, paystackReference, functionName)
			}
		}
	}
	return
}

func HandleEventSnapShot(server *Server, ctx *gin.Context, optionUserID uuid.UUID, reference string, paystackReference string, chargeID uuid.UUID) (err error) {
	var eventDetailString string
	var eventInfoString string
	var eventDateTimeString string
	var eventLocationString string
	option, err := server.store.GetEventInfoByUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in GetEventInfoByUserID: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		err = fmt.Errorf("error 500 occur, pls contact us")
		return
	}
	eventInfoPolicy, err := server.store.GetEventInfoAnyPolicy(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in GetEventInfo: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
	}
	eventInfo := db.EventInfo{
		OptionID:        eventInfoPolicy.OptionID,
		SubCategoryType: eventInfoPolicy.SubCategoryType,
		EventType:       eventInfoPolicy.EventType,
		CreatedAt:       eventInfoPolicy.CreatedAt,
		UpdatedAt:       eventInfoPolicy.UpdatedAt,
	}
	eventDetails := []db.EventDateDetail{}
	eventLocations := []db.EventDateLocation{}
	eventDateTimes, err := server.store.ListEventDateTime(ctx, db.ListEventDateTimeParams{
		EventInfoID: option.ID,
	})
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in ListEventDateTime: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		eventDateTimes = []db.EventDateTime{}
	}

	for _, date := range eventDateTimes {
		// Get details
		detail, err := server.store.GetEventDateDetail(ctx, date.ID)
		if err != nil {
			log.Printf("Error at HandleEventSnapShot in GetEventDateDetail: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		} else {
			eventDetails = append(eventDetails, detail)
		}
		// Get Location
		location, err := server.store.GetEventDateLocation(ctx, date.ID)
		if err != nil {
			log.Printf("Error at HandleEventSnapShot in GetEventDateLocation: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		} else {
			eventLocations = append(eventLocations, location)
		}
	}

	eventLocationString, err = StructToStringEventLocation(eventLocations)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in StructToStringEventLocation: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		eventLocationString = "none"
	}

	eventInfoString, err = StructToStringEventInfo(eventInfo)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in StructToStringEventInfo: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		eventInfoString = "none"
	}

	eventDetailString, err = StructToStringEventDateDetail(eventDetails)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in StructToStringEventDateDetail: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		eventDetailString = "none"
	}
	eventDateTimeString, err = StructToStringEventDateTime(eventDateTimes)
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in StructToStringEventDateTime: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		eventDateTimeString = "none"
	}

	err = server.store.CreateEventReferenceInfo(ctx, db.CreateEventReferenceInfoParams{
		EventChargeID:     chargeID,
		EventDateLocation: eventLocationString,
		EventInfo:         eventInfoString,
		EventDateTimes:    eventDateTimeString,
		CancelPolicyOne:   eventInfoPolicy.TypeOne,
		CancelPolicyTwo:   eventInfoPolicy.TypeTwo,
		EventDateDetails:  eventDetailString,
		HostAsIndividual:  option.HostAsIndividual,
		OrganizationName:  option.OrganizationName,
	})
	if err != nil {
		log.Printf("Error at HandleEventSnapShot in CreateEventReferenceInfo: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
	}
	return

}

// This makes sure the selected tickets to be booked are actually available
func TicketQuantityAvailable(server *Server, ctx *gin.Context, quantity int, ticketIDString string, eventDateIDString string, startDateString string) (err error) {
	eventDateID, err := tools.StringToUuid(eventDateIDString)
	if err != nil {
		log.Printf("Error at TicketQuantityAvailable in eventDateID StringToUuid: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateIDString, eventDateIDString)
		return
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(startDateString)
	if err != nil {
		log.Printf("Error at TicketQuantityAvailable in ConvertDateOnlyStringToDate: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateIDString, ticketIDString)
		return
	}
	ticketID, err := tools.StringToUuid(ticketIDString)
	if err != nil {
		log.Printf("Error at TicketQuantityAvailable in StringToUuid: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateIDString, ticketIDString)
		return
	}
	eventDate, err := server.store.GetEventDateTime(ctx, eventDateID)
	if err != nil {
		log.Printf("Error at TicketQuantityAvailable in GetEventDateTime: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateIDString, ticketIDString)
		return
	}
	// We want to get the actual ticket
	ticket, err := server.store.GetEventDateTicket(ctx, db.GetEventDateTicketParams{
		ID:              ticketID,
		EventDateTimeID: eventDateID,
	})
	if err != nil {
		log.Printf("Error at TicketQuantityAvailable in GetEventDateTicket: %v eventDateIDString: %v ticketIDString: %v\n", err.Error(), eventDateIDString, ticketIDString)
		return
	}
	//if ticket.Type == "free" {
	//	err = nil
	//	return
	//}
	if eventDate.Type == "single" {
		ticketCount, errCharge := server.store.CountChargeTicketReference(ctx, db.CountChargeTicketReferenceParams{
			TicketID:   ticketID,
			Cancelled:  false,
			IsComplete: true,
		})
		if errCharge != nil {
			if err == db.ErrorRecordNotFound {
				err = nil
			} else {
				log.Printf("Error at TicketQuantityAvailable in ListChargeTicketReference %v eventDateIDString: %v ticketIDString: %v\n", errCharge.Error(), eventDateIDString, ticketIDString)
				err = nil
				return
			}
		}
		if (int(ticketCount) + quantity) > int(ticket.Capacity) {
			err = fmt.Errorf("there are only about %v left and the quantity you request exceed the capacity", ticketCount)
			return
		} else {
			err = nil
		}

	} else {
		ticketCount, errCharge := server.store.CountChargeTicketReferenceByStartDate(ctx, db.CountChargeTicketReferenceByStartDateParams{
			Date:       startDate,
			TicketID:   ticketID,
			Cancelled:  false,
			IsComplete: true,
		})
		if errCharge != nil {
			if err == db.ErrorRecordNotFound {
				err = nil
			} else {
				log.Printf("Error at TicketQuantityAvailable in CountChargeTicketReference %v eventDateIDString: %v ticketIDString: %v\n", errCharge.Error(), eventDateIDString, ticketIDString)
				err = nil
				return
			}
		}
		if (int(ticketCount) + quantity) > int(ticket.Capacity) {
			err = fmt.Errorf("there are only about %v left and the quantity you request exceed the capacity", ticketCount)
			return
		} else {
			err = nil
		}
	}
	return
}
