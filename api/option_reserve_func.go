package api

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/makuo12/ghost_server/constants"
	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
	"github.com/makuo12/ghost_server/utils"
	"github.com/makuo12/ghost_server/val"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func ReserveOptionCalculate(user db.User, server *Server, ctx *gin.Context, startDate string, endDate string, guests []string, optionUserID string, userCurrency string) (canInstantBook bool, datePriceFloat []DatePriceFloat, totalDatePrice float64, discountType string, discount float64, cleanFee float64, extraGuestFee float64, petFee float64, petStayFee float64, extraGuestStayFee float64, totalPrice float64, serviceFee float64, requireRequest bool, requestType string, optionID uuid.UUID, err error) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	var servicePercent float64
	if userCurrency != utils.NGN {
		servicePercent = tools.ConvertStringToFloat(server.config.IntServiceOptionUserPercent)
	} else {
		servicePercent = tools.ConvertStringToFloat(server.config.LocServiceOptionUserPercent)
	}
	optionID, err = tools.StringToUuid(optionUserID)
	if err != nil {
		log.Printf("Error at ReserveOptionCalculate in StringToUuid %v for user: %v. optionID: %v\n", err.Error(), user.ID, optionUserID)
		err = fmt.Errorf("the listing does not exist")
		return
	}
	option, err := server.store.GetOptionInfoCustomer(ctx, db.GetOptionInfoCustomerParams{
		OptionUserID:    optionID,
		IsComplete:      true,
		IsActive:        true, // Option is active
		IsActive_2:      true, // Host is active
		OptionStatusOne: "list",
		OptionStatusTwo: "staged",
	})
	if err != nil {
		log.Printf("Error at ReserveOptionCalculate in GetOptionInfoCustomer %v for user: %v. optionID: %v\n", err.Error(), user.ID, optionUserID)
		err = fmt.Errorf("could not perform this request")
		return
	}
	// List of dates
	userDateTime, err := tools.GenerateDateListString(startDate, endDate)
	if err != nil {
		log.Printf("Error at ReserveOptionList in .ListOptionDiscount: %v optionID: %v\n", err.Error(), optionID)
		return
	}
	// Day count
	dayCount, err := tools.FindDateDifference(startDate, endDate)
	if err != nil {
		log.Printf("Error at ReserveAddCharge in FindDateDifference: %v optionID: %v\n", err.Error(), optionID)
		return
	}
	charges, discounts, dateTime := ReserveOptionList(server, ctx, option.ID)
	// Get dateTimeString
	dateTimeString := OptionDateTimeString(dateTime)
	// Check if the dates are available
	confirm, err := ReserveDatesAvailable(startDate, endDate, userDateTime, dateTimeString, option.ID)
	if err != nil || !confirm {
		err = tools.HandleConfirmError(err, confirm, "your dates are not confirmed please try selecting dates that are available")
		return
	}
	free, err := HandleReserveAvailable(ctx, server, option.OptionUserID, option.PreparationTime, option.AutoBlockDates, userDateTime)
	if err != nil || !free {
		err = tools.HandleConfirmError(err, confirm, "a date you selected is currently booked")
		return
	}
	confirm, err = ReserveAvailableSetting(startDate, endDate, option.ID, option.AdvanceNotice, option.PreparationTime, option.AvailabilityWindow)
	if err != nil || !confirm {
		err = tools.HandleConfirmError(err, confirm, "your dates are not confirmed please try selecting dates that are available")
		return
	}
	confirm, requireRequest, requestType, err = ReserveTripLength(startDate, endDate, option.ID, int(option.MinStayDay), int(option.MaxStayNight), option.ManualApproveRequestPassMax, option.AllowReservationRequest, dayCount)
	if err != nil || !confirm {
		err = tools.HandleConfirmError(err, confirm, "your dates are not confirmed please try selecting dates that are available")
		return
	}

	canInstantBook, err = ReserveBookMethods(option.ID, ctx, server, user, option.InstantBook, option.IdentityVerified, option.ProfilePhoto)
	log.Println("canInstantBook ", canInstantBook)
	if err != nil {
		return
	}
	datePriceFloat, totalDatePrice, err = ReserveDateTimePrice(option.ID, startDate, endDate, userDateTime, tools.IntToMoneyString(option.WeekendPrice), tools.IntToMoneyString(option.Price), dateTimeString, option.Currency, userCurrency, dollarToNaira, dollarToCAD)
	if err != nil {
		return
	}
	discountType, discountInt, err := ReserveDiscount(option.ID, discounts, startDate, endDate, dayCount)
	if err != nil {
		return
	}
	cleanFee, petFee, extraGuestFee, petStayFee, extraGuestStayFee, err = ReserveAddCharge(option.ID, user, charges, int(option.GuestWelcomed), guests, startDate, endDate, dayCount, option.Currency, userCurrency, dollarToNaira, dollarToCAD)
	totalPrice, serviceFee, discount = HandleEndPrice(totalDatePrice, discountInt, cleanFee, petFee, extraGuestFee, servicePercent)
	return
}

// This would handle all factors affecting option reservation that is list format
func ReserveOptionList(server *Server, ctx context.Context, optionID uuid.UUID) (charges []db.ListOptionAddChargeRow, discounts []db.ListOptionDiscountRow, dateTimes []db.OptionDateTime) {
	// Option Additional Charges
	charges, err := server.store.ListOptionAddCharge(ctx, optionID)
	if err != nil {
		// We don't need to send an error because host might have no additional charges
		log.Printf("Error at ReserveOptionList in ListOptionAddCharge: %v optionID: %v\n", err.Error(), optionID)
		charges = []db.ListOptionAddChargeRow{}
	}
	// Discounts
	discounts, err = server.store.ListOptionDiscount(ctx, optionID)
	if err != nil {
		// We don't need to send an error because host might have no discounts
		log.Printf("Error at ReserveOptionList in .ListOptionDiscount: %v optionID: %v\n", err.Error(), optionID)
		discounts = []db.ListOptionDiscountRow{}
	}
	dateTimes, err = server.store.ListOptionDateTimeMore(ctx, optionID)
	if err != nil {
		// We don't need to send an error because host might have no special days
		log.Printf("Error at ReserveOptionList in .ListOptionDiscount: %v optionID: %v\n", err.Error(), optionID)
		dateTimes = []db.OptionDateTime{}
	}
	return
}

func ReserveDatesAvailable(startDate string, endDate string, userDateTime []string, dateTimeString []OptionDateTime, optionID uuid.UUID) (confirm bool, err error) {
	for _, userDate := range userDateTime {
		for _, date := range dateTimeString {
			if userDate == date.Date {
				// Check if the date is available
				if !date.Available {
					err = fmt.Errorf("%v is unavailable", date.Date)
					confirm = false
					return
				}
			}
		}
	}
	// We also need to check if any of the dates are book
	// We set confirm to true because means all user dates are available
	confirm = true
	err = nil

	return
}

func HandleReserveAvailable(ctx context.Context, server *Server, optionUserID uuid.UUID, prepareTime string, autoBlock bool, userDateTime []string) (free bool, err error) {
	charges, err := server.store.ListChargeOptionReferenceByOptionUserID(ctx, db.ListChargeOptionReferenceByOptionUserIDParams{
		OptionUserID: optionUserID,
		Cancelled:    false,
		IsComplete:   true,
	})
	if err != nil || len(charges) == 0 {
		if err != nil {
			log.Printf("Error at HandleExPrepareTime in ListChargeOptionReferenceByOptionUserID err: %v, user: %v, optionUserID: %v\n", err, ctx, optionUserID)
		}
		err = nil
		free = true
		return
	}
	// First we handle dates that are booked
	if autoBlock {
		bookDates := GetExDateTime(prepareTime, tools.UuidToString(optionUserID), charges, tools.UuidToString(optionUserID))
		for _, userDate := range userDateTime {
			for _, book := range bookDates {
				if userDate == book.Date {
					err = fmt.Errorf("%v is booked", userDate)
					free = false
					return
				}

			}
		}
	}
	err = nil
	free = true
	return
}

func ReserveAvailableSetting(startDate string, endDate string, optionID uuid.UUID, advanceNotice string, prepareTime string, availableWindow string) (confirm bool, err error) {
	// First we start with advance notice
	var advanceNoticeDate string
	if advanceNotice == "same day" {
		advanceNoticeDate = tools.ConvertTimeToStringDateOnly(time.Now())
	} else {
		day, err := tools.ExtractNumberFromString(advanceNotice)
		if err != nil {
			day = 0
		}
		advanceNoticeDate = tools.CurrentDatePlusDaysToString(day)
	}
	startDateTime, err := tools.ConvertDateOnlyStringToDate(startDate)
	if err != nil {
		log.Printf("Error at start time ReserveAvailableSetting in .ConvertDateOnlyStringToDate: %v optionID: %v\n", err.Error(), optionID)
		confirm = false
		return
	}
	advanceNoticeDateTime, err := tools.ConvertDateOnlyStringToDate(advanceNoticeDate)
	if err != nil {
		log.Printf("Error at advance date ReserveAvailableSetting in .ConvertDateOnlyStringToDate: %v optionID: %v\n", err.Error(), optionID)
		confirm = false
		return
	}
	if !startDateTime.After(advanceNoticeDateTime) {
		confirm = false
		err = fmt.Errorf("your dates don't not aline with the host available dates")
		return
	}

	// Next we move to available window
	// This deals with months not days
	var availableWindowDate string
	if availableWindow == "all future dates" {
		availableWindowDate = tools.CurrentDatePlusMonthsToString(36)
	} else if availableWindow == "dates unavailable by default" {
		availableWindowDate = tools.CurrentDatePlusMonthsToString(24)
	} else {
		month, err := tools.ExtractNumberFromString(availableWindow)
		if err != nil {
			month = 0
		}
		availableWindowDate = tools.CurrentDatePlusMonthsToString(month)
	}

	endDateTime, err := tools.ConvertDateOnlyStringToDate(endDate)
	if err != nil {
		log.Printf("Error at start time ReserveAvailableSetting in .ConvertDateOnlyStringToDate: %v optionID: %v\n", err.Error(), optionID)
		confirm = false
		return
	}
	availableWindowDateTime, err := tools.ConvertDateOnlyStringToDate(availableWindowDate)
	if err != nil {
		log.Printf("Error at advance date ReserveAvailableSetting in .ConvertDateOnlyStringToDate: %v optionID: %v\n", err.Error(), optionID)
		confirm = false
		return
	}
	if !endDateTime.Before(availableWindowDateTime) {
		confirm = false
		err = fmt.Errorf("your dates don't not aline with the host available dates")
		return
	}

	// Last we do prepareTime
	// This is be done later

	// We set confirm to true because means it pass available settings
	confirm = true
	return
}

func ReserveTripLength(startDate string, endDate string, optionID uuid.UUID, minStay int, maxStay int, manualApprove bool, allowRequest bool, dayCount int) (confirm bool, requireRequest bool, requestType string, err error) {
	// minStay
	requestType = "normal_stay"
	if dayCount < minStay {
		if allowRequest {
			requireRequest = true
			requestType = "min_stay"
		} else {
			err = fmt.Errorf("trip length is less than host minimum stay")
			confirm = false
			return
		}
	}
	if dayCount > maxStay {
		if manualApprove {
			requireRequest = true
			requestType = "max_stay"
		} else {
			err = fmt.Errorf("trip length is greater than host maximum stay")
			confirm = false
			return
		}
	}

	// We set confirm to true because means it pass available settings
	confirm = true
	return
}

func ReserveBookMethods(optionID uuid.UUID, ctx *gin.Context, server *Server, user db.User, instantBook bool, identityVerified bool, profilePhoto bool) (canInstantBook bool, err error) {
	if identityVerified {
		// This means host requires that user identity to be verified
		identity, errIdent := server.store.GetIdentity(ctx, user.ID)
		if errIdent != nil {
			errIdent = fmt.Errorf("to book this stay your identity must be verified")
			err = errIdent
			log.Printf("Error at advance date ReserveBookMethods in GetIdentity: %v optionID: %v\n", err.Error(), optionID)

		}
		if !identity.IsVerified {
			// We send an error because the user isn't verified
			err = fmt.Errorf("to book this stay your identity must be verified")
		}
	}
	if profilePhoto {
		// This means host requires user to have a profile photo
		if tools.ServerStringEmpty(user.Image) {
			err = fmt.Errorf("to book this stay you need to have a profile photo")
		}
	}
	if err != nil {
		return
	}
	canInstantBook = instantBook
	return
}

func ReserveDateTimePrice(optionID uuid.UUID, startDate string, endDate string, userDateTime []string, weekendPriceString string, basePriceString string, dateTimeString []OptionDateTime, hostCurrency string, userCurrency string, dollarToNaira string, dollarToCAD string) (datePriceFloat []DatePriceFloat, totalDatePrice float64, err error) {
	weekendPrice, err := tools.ConvertPrice(weekendPriceString, hostCurrency, userCurrency, dollarToNaira, dollarToCAD, optionID)
	if err != nil {
		return
	}
	basePrice, err := tools.ConvertPrice(basePriceString, hostCurrency, userCurrency, dollarToNaira, dollarToCAD, optionID)
	if err != nil {
		return
	}
	for _, date := range userDateTime {
		// datePriceChanged this tells if the price was change in date == dateString.Date
		var datePriceChanged bool
		for _, dateString := range dateTimeString {
			if date == dateString.Date {
				price, errPrice := tools.ConvertPrice(dateString.Price, hostCurrency, userCurrency, dollarToNaira, dollarToCAD, optionID)
				if errPrice != nil {
					err = errPrice
					break
				}
				// we check if the price that was set if price is greater than 0.01 then it was set
				if price > 0.01 {
					data := DatePriceFloat{
						Price:      price,
						Date:       date,
						GroupPrice: price,
					}
					datePriceFloat = append(datePriceFloat, data)
					datePriceChanged = true
				}
				// We break out of this loop
				break
			}
		}
		if !datePriceChanged {
			var isWeekend bool = false
			var price float64
			if weekendPrice > 0.01 {
				// This means there is weekend price
				// so we check if the date is a weekend
				isWeekend, err = tools.IsWeekend(date)
				if err != nil {
					log.Printf("Error at advance date ReserveDateTimePrice in IsWeekend: %v optionID: %v\n", err.Error(), optionID)
					err = nil
				}
			}
			if isWeekend {
				price = weekendPrice
			} else {
				price = basePrice
			}
			data := DatePriceFloat{
				Price:      price,
				Date:       date,
				GroupPrice: price,
			}
			datePriceFloat = append(datePriceFloat, data)
		}
	}
	if err == nil {
		// Now we want to groupPrice
		datePriceFloat = HandleDateGroupPrice(datePriceFloat)
		// Lastly we get the totalDatePrice
		totalDatePrice = GetTotalDatePrice(datePriceFloat)
	}
	return
}

func ReserveDiscount(optionID uuid.UUID, discounts []db.ListOptionDiscountRow, startDate string, endDate string, dayCount int) (discountType string, discount int, err error) {
	for _, dis := range discounts {
		disNumber := val.GetOptionDiscountNumber(dis.Type)
		if dayCount > disNumber && disNumber != 0 {
			// This would update the discount each time the discount meets the condition
			discount = int(dis.Percent)
			discountType = dis.Type
		}
	}
	return
}

func ReserveAddCharge(optionID uuid.UUID, user db.User, charges []db.ListOptionAddChargeRow, numOfGuests int, guests []string, startDate string, endDate string, dayCount int, hostCurrency string, userCurrency string, dollarToNaira string, dollarToCAD string) (cleanFee float64, petFee float64, extraGuestFee float64, petStayFee float64, extraGuestStayFee float64, err error) {
	guestCount := tools.HandleListCount(guests)
	// First we check to see if it follows the numOfGuest requirement
	// guestCount["children"] + guestCount["adult"] cannot be greater than numOfGuest
	totalGuests := guestCount["children"] + guestCount["adult"]
	if totalGuests > numOfGuests {
		log.Printf("Error at ReserveAddCharge in HandleListCount: optionID: %v\n", optionID)
		err = fmt.Errorf("the number of guests coming has exceeded the host maximum number")
		return
	}
	for _, charge := range charges {
		mainFee, errPrice := tools.ConvertPrice(tools.IntToMoneyString(charge.MainFee), hostCurrency, userCurrency, dollarToNaira, dollarToCAD, optionID)
		if errPrice != nil {
			err = errPrice
			break
		}
		if charge.Type == "cleaning_fee" {
			// Check if it is a short stay
			if dayCount < 3 {
				// DayCount shows it is a short stay
				extraFee, errPriceTwo := tools.ConvertPrice(tools.IntToMoneyString(charge.ExtraFee), hostCurrency, userCurrency, dollarToNaira, dollarToCAD, optionID)
				if errPriceTwo != nil {
					err = errPriceTwo
					break
				}
				if extraFee > 0.01 {
					// This means it has a cleaning fee for short stay
					cleanFee = extraFee
				} else {
					cleanFee = mainFee
				}
			} else {
				// This means it is a long stay
				cleanFee = mainFee
			}
		} else if charge.Type == "pet_fee" {
			// check if there are any pets
			if guestCount["pet"] != 0 {
				// We know there are pets
				petStayFee = mainFee
				petFee = mainFee * float64(dayCount)
			}
		} else if charge.Type == "extra_guest_fee" {
			// Remember this charge comes in if the totalGuest coming passes the stated on in the charge
			if totalGuests > int(charge.NumOfGuest) {
				extraGuestStayFee = mainFee
				extraGuestFee = mainFee * float64(dayCount)
			}
		}
	}
	return
}

func OptionDateTimeString(dateTime []db.OptionDateTime) (dateTimeString []OptionDateTime) {
	for _, date := range dateTime {
		data := OptionDateTime{
			ID:        date.ID,
			OptionID:  date.OptionID,
			Date:      tools.ConvertDateOnlyToString(date.Date),
			Price:     tools.IntToMoneyString(date.Price),
			Note:      date.Note,
			CreatedAt: date.CreatedAt,
			UpdatedAt: date.UpdatedAt,
		}
		dateTimeString = append(dateTimeString, data)
	}
	return
}

func OptionDateTimeStringOUD(dateTime []db.ListAllOptionDateTimeByOUDRow) (dateTimeString []OptionDateTime) {
	for _, date := range dateTime {
		data := OptionDateTime{
			ID:        date.ID,
			OptionID:  date.OptionID,
			Date:      tools.ConvertDateOnlyToString(date.Date),
			Price:     tools.IntToMoneyString(date.Price),
			Note:      date.Note,
			CreatedAt: date.CreatedAt,
			UpdatedAt: date.UpdatedAt,
		}
		dateTimeString = append(dateTimeString, data)
	}
	return
}

// This convert a list of datePriceFloat to a well group date price float
func HandleDateGroupPrice(datePriceFloat []DatePriceFloat) (grouped []DatePriceFloat) {
	var data = make(map[float64]float64)

	for _, datePrice := range datePriceFloat {
		data[datePrice.Price] = data[datePrice.Price] + datePrice.Price
	}
	for _, d := range datePriceFloat {
		newData := DatePriceFloat{
			Price:      d.Price,
			GroupPrice: data[d.Price],
			Date:       d.Date,
		}
		grouped = append(grouped, newData)
	}

	return
}

// This convert a list of datePriceFloat to a well group date price float
func GetTotalDatePrice(datePriceFloat []DatePriceFloat) (totalDatePrice float64) {
	for _, datePrice := range datePriceFloat {
		totalDatePrice += datePrice.Price
	}
	return
}

// This would then show all the prices
func HandleEndPrice(totalDatePrice float64, discountInt int, cleanFee float64, petFee float64, extraGuestFee float64, servicePercent float64) (totalPrice float64, serviceFee float64, discount float64) {
	totalStayPrice := totalDatePrice + petFee + extraGuestFee
	discount = totalStayPrice * float64(discountInt/100)
	totalMainPrice := totalStayPrice - discount
	serviceFee = totalMainPrice * float64(servicePercent/100)
	totalPrice = totalMainPrice + cleanFee + serviceFee
	return
}

// We want to handle the process of sending a host a reservation request that was made by a guest
func HandleOptionReserveRequest(server *Server, ctx context.Context, payMethodReference string, reserveData ExperienceReserveOModel, user db.User, msg string) (err error) {
	// requestApproved bool, isComplete bool would be set to false because no payment has been made just an awaiting request to be sent
	charge, optionUserID, err := HandleOptionReserveReceipt(server, ctx, reserveData, payMethodReference, user, false, false, "HandleOptionReserveRequest")
	if err != nil {
		log.Printf("Error at HandleOptionReserveRequest in HandleOptionReserveReceipt: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), optionUserID, charge.Reference, payMethodReference)
		err = fmt.Errorf("error 300 occur, pls contact us")
		return
	}
	// Next want want to send a message
	// We do this by storing it in the database
	receiver, err := server.store.GetOptionInfoUserIDByUserID(ctx, optionUserID)
	if err != nil {
		log.Printf("Error at HandleOptionReserveRequest in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), optionUserID, charge.Reference, payMethodReference)
		err = fmt.Errorf("error 300 occurred, pls contact us")
		return
	}
	// We want to create the message
	// We want to store the message in redis
	// First we create a CreateMessageParams struct
	// SenderID is the user because that is the person booking
	// ReceiverID is the option.host_id because that is the owner
	startDate, err := tools.ConvertDateFormat(reserveData.StartDate, tools.DateDMM)
	if err != nil {
		startDate = reserveData.StartDate
	}
	endDate, err := tools.ConvertDateFormat(reserveData.EndDate, tools.DateDMM)
	if err != nil {
		endDate = reserveData.EndDate
	}
	if len(msg) == 0 {

		msg = fmt.Sprintf("Hey, I'd want to know if it's possible for me to stay at %v from %v to %v.", receiver.HostNameOption, startDate, endDate)

	}
	header := fmt.Sprintf("Reservation request for %v, from %v.\nDates from %v to %v", receiver.HostNameOption, user.FirstName, startDate, endDate)
	createdAt := time.Now().UTC()
	msgD, err := server.store.CreateMessage(ctx, db.CreateMessageParams{
		MsgID:      uuid.New(),
		SenderID:   user.UserID,
		ReceiverID: receiver.UserID,
		Message:    msg,
		Type:       "user_request",
		MainImage:  "none",
		ParentID:   "none",
		Reference:  charge.Reference,
		CreatedAt:  createdAt,
		UpdatedAt:  createdAt,
	})
	if err != nil {
		log.Printf("Error at HandleOptionReserveRequest in CreateMessage: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), optionUserID, charge.Reference, payMethodReference)
		err = fmt.Errorf("error 301-e occurred, pls contact us")
		return
	}
	// Send an email notification
	BrevoReservationRequest(ctx, server, receiver.Email, receiver.FirstName, header, msg, "HandleOptionReserveRequest", user.ID, receiver.Email, receiver.FirstName, receiver.LastName, tools.UuidToString(charge.ID), receiver.UserID.String(), user.Email, user.FirstName, user.LastName, user.UserID.String())
	//
	HandleUserIdApn(ctx, server, receiver.UserID, header, msg)
	// When we create a message we want to create a room is this user and the receiver doesn't have a room
	_, err = SingleContextRoom(ctx, server, user.UserID, receiver.UserID, "HandleOptionReserveRequest")
	if err != nil {
		log.Printf("Error at HandleOptionReserveRequest in SingleContextRoom: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), optionUserID, charge.Reference, payMethodReference)
	}

	fee := tools.ConvertStringToFloat(tools.IntToMoneyString(charge.TotalFee)) - tools.ConvertStringToFloat(tools.IntToMoneyString(charge.ServiceFee))
	err = server.store.CreateRequestNotify(ctx, db.CreateRequestNotifyParams{
		MID:       msgD.MID,
		StartDate: reserveData.StartDate,
		EndDate:   reserveData.EndDate,
		HasPrice:  true,
		SamePrice: true,
		Price:     tools.MoneyFloatToInt(fee),
		ItemID:    tools.UuidToString(optionUserID),
	})
	if err != nil {
		log.Printf("Error at HandleOptionReserveRequest in CreateRequestNotify: %v optionID: %v referenceID: %v, pay_method_reference: %v\n", err.Error(), optionUserID, charge.Reference, payMethodReference)
		err = fmt.Errorf("error 304-e occurred, pls contact us")
		return
	}
	return
}

// We want to save a recept in the database
// We also want to store a snap shot of what the shortlet looks like
func HandleOptionReserveComplete(server *Server, ctx context.Context, reserveData ExperienceReserveOModel, referenceCharge string, paystackReference string, user db.User, msg string, fromCharge bool) (err error) {
	// If from charge we just want to update the charge data
	var chargeRef string
	var chargeID uuid.UUID
	var optionUserID uuid.UUID
	if fromCharge {
		log.Println("at charge, ")
		log.Println("at charge, reference ", referenceCharge)
		charge, errCharge := server.store.UpdateChargeOptionReferenceByRef(ctx, db.UpdateChargeOptionReferenceByRefParams{
			Reference: referenceCharge,
			UserID:    user.UserID,
			IsComplete: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
		})
		log.Println("at charge, reference here", charge)
		if errCharge != nil {
			log.Printf("Error at HandleOptionReserveComplete  from Charge in HandleOptionReserveReceipt: %v optionID: %v referenceID: %v, paystackReference: %v\n", errCharge.Error(), charge.OptionUserID, charge.Reference, paystackReference)
			errCharge = fmt.Errorf("an error occurred while storing your payment")
			err = errCharge
			return
		}
		// We want update the notification to handle

		chargeRef = charge.Reference
		chargeID = charge.ID
		optionUserID = charge.OptionUserID
		msgData, err := server.store.GetMessageByRef(ctx, charge.Reference)
		if err != nil {
			log.Printf("Error at HandleOptionReserveComplete  from Charge in .store.GetMessageByRef: %v optionID: %v referenceID: %v, paystackReference: %v\n", err.Error(), charge.OptionUserID, charge.Reference, paystackReference)
		} else {
			err = server.store.UpdateNotificationHandled(ctx, db.UpdateNotificationHandledParams{
				UserID:  user.UserID,
				ItemID:  msgData.ID,
				Handled: true,
			})
			if err != nil {
				log.Printf("Error at HandleOptionReserveComplete  from Charge in .store.UpdateNotificationHandled: %v optionID: %v referenceID: %v, paystackReference: %v\n", err.Error(), charge.OptionUserID, charge.Reference, paystackReference)
			}
		}
		log.Println("at charge, ", chargeRef, chargeID, optionUserID)
	} else {
		// First we store the receipt
		charge, _, errReceipt := HandleOptionReserveReceipt(server, ctx, reserveData, paystackReference, user, true, true, "HandleOptionReserveComplete")
		if errReceipt != nil {
			log.Printf("Error at HandleOptionReserveComplete in HandleOptionReserveReceipt: %v optionID: %v referenceID: %v, paystackReference: %v\n", errReceipt.Error(), optionUserID, charge.Reference, paystackReference)
			err = errReceipt
			return
		}
		chargeRef = charge.Reference
		chargeID = charge.ID
		optionUserID = charge.OptionUserID
		log.Println("not at charge, ", chargeRef, chargeID, optionUserID)
	}
	// We want to store a snap shot of the option
	err = HandleOptionSnapShot(server, ctx, chargeRef, paystackReference, optionUserID, chargeID)
	if err != nil {
		log.Printf("Error at HandleOptionReserveComplete in HandleOptionSnapShot: %v optionID: %v referenceID: %v, paystackReference: %v\n", err.Error(), optionUserID, chargeRef, paystackReference)
	}
	err = nil
	return
}

// This generates and stores the receipt in the database
// requestApproved bool, isComplete bool are important because they show whether the stay was approved and if payment was successful
func HandleOptionReserveReceipt(server *Server, ctx context.Context, reserveData ExperienceReserveOModel, paystackReference string, user db.User, requestApproved bool, isComplete bool, functionName string) (charge db.ChargeOptionReference, optionUserID uuid.UUID, err error) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	// First we store the receipt
	optionUserID, err = tools.StringToUuid(reserveData.OptionUserID)
	if err != nil {
		log.Printf("Error at HandleOptionReserveReceipt in StringToUuid: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName)
		err = fmt.Errorf("error 451 occur, pls contact us")
		return
	}

	var datePrice []string
	for _, d := range reserveData.DatePrice {
		data := d.Price + "&" + d.Date + "&" + d.GroupPrice
		datePrice = append(datePrice, data)
	}
	startDate, err := tools.ConvertDateOnlyStringToDate(reserveData.StartDate)
	if err != nil {
		log.Printf("Error at startDate HandleOptionReserveReceipt in ConvertDateOnlyStringToDate: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName)
		err = fmt.Errorf("error 452 occur, pls contact us")
		return
	}
	endDate, err := tools.ConvertDateOnlyStringToDate(reserveData.EndDate)
	if err != nil {
		log.Printf("Error at endDate HandleOptionReserveReceipt in ConvertDateOnlyStringToDate: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName)
		err = fmt.Errorf("error 452 occur, pls contact us")
		return
	}
	discount := reserveData.Discount.Price + "&" + reserveData.Discount.Type
	charge, err = server.store.CreateChargeOptionReference(ctx, db.CreateChargeOptionReferenceParams{
		UserID:           user.UserID,
		OptionUserID:     optionUserID,
		Discount:         discount,
		MainPrice:        tools.MoneyStringToInt(reserveData.MainPrice),
		ServiceFee:       tools.MoneyStringToInt(reserveData.ServiceFee),
		TotalFee:         tools.MoneyStringToInt(reserveData.TotalFee),
		Guests:           reserveData.Guests,
		DatePrice:        datePrice,
		DateBooked:       time.Now().Add(time.Hour),
		Currency:         reserveData.Currency,
		StartDate:        startDate,
		EndDate:          endDate,
		GuestFee:         tools.MoneyStringToInt(reserveData.GuestFee),
		PetFee:           tools.MoneyStringToInt(reserveData.PetFee),
		CleanFee:         tools.MoneyStringToInt(reserveData.CleaningFee),
		NightlyPetFee:    tools.MoneyStringToInt(reserveData.NightlyPetFee),
		NightlyGuestFee:  tools.MoneyStringToInt(reserveData.NightlyGuestFee),
		CanInstantBook:   reserveData.CanInstantBook,
		RequireRequest:   reserveData.RequireRequest,
		RequestType:      reserveData.RequestType,
		Reference:        reserveData.Reference,
		PaymentReference: paystackReference,
		RequestApproved:  requestApproved,
		IsComplete:       isComplete,
	})
	if err != nil {
		log.Printf("Error at endDate HandleOptionReserveReceipt in CreateChargeOptionReference: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName)
		err = fmt.Errorf("error 453 occur, pls contact us")
		return
	}
	var payoutAmount float64
	var servicePercent float64
	var serviceFee float64
	amount := tools.ConvertStringToFloat(tools.IntToMoneyString(charge.TotalFee)) - tools.ConvertStringToFloat(tools.IntToMoneyString(charge.ServiceFee))
	amount, err = tools.ConvertPrice(tools.ConvertFloatToString(amount), charge.Currency, utils.PayoutCurrency, dollarToNaira, dollarToCAD, user.ID)
	if err != nil {
		log.Printf("Error at endDate HandleOptionReserveReceipt in tools.ConvertPrice: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName)
		err = nil
		return
	}
	switch utils.PayoutCurrency {
	case utils.NGN:
		servicePercent = tools.ConvertStringToFloat(server.config.LocServiceOptionHostPercent)
	default:
		servicePercent = tools.ConvertStringToFloat(server.config.IntServiceOptionHostPercent)
	}
	serviceFee = (servicePercent / 100) * amount
	payoutAmount = amount - serviceFee
	err = server.store.CreateMainPayout(ctx, db.CreateMainPayoutParams{
		ChargeID:   charge.ID,
		Type:       constants.CHARGE_OPTION_REFERENCE,
		IsComplete: false,
		Amount:     tools.MoneyStringToInt(tools.ConvertFloatToString(payoutAmount)),
		ServiceFee: tools.MoneyStringToInt(tools.ConvertFloatToString(serviceFee)),
		Currency:   utils.PayoutCurrency,
	})
	if err != nil {
		log.Printf("Error at endDate HandleOptionReserveReceipt in CreateMainPayout: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v chargeID: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, paystackReference, functionName, charge.ID)
	}
	return
}
