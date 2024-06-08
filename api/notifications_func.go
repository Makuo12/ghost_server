package api

import (
	"context"
	"fmt"
	"log"
	"strings"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"
)

func HandleOptionChargeToReserve(server *Server, charge db.ChargeOptionReference, userCurrency string) (ExperienceReserveOModel, error) {
	// First lets handle discount
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	discountSplit := strings.Split(charge.Discount, "&")
	discount := ReDiscount{"none", "none"}
	if len(discountSplit) == 2 {
		price, err := tools.ConvertPrice(discountSplit[0], charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
		if err != nil {
			err = fmt.Errorf("could not convert to your currency correctly")
			return ExperienceReserveOModel{}, err
		}
		discount = ReDiscount{
			Price: tools.ConvertFloatToString(price),
			Type:  discountSplit[1],
		}
	}
	var datePrice []DatePrice
	if len(charge.DatePrice) == 0 {
		err := fmt.Errorf("prices of dates cannot be found")
		return ExperienceReserveOModel{}, err
	}
	for _, d := range charge.DatePrice {
		split := strings.Split(d, "&")
		if len(split) == 3 {
			price, err := tools.ConvertPrice(split[0], charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
			if err != nil {
				err = fmt.Errorf("could not convert to your currency correctly")
				return ExperienceReserveOModel{}, err
			}
			groupPrice, err := tools.ConvertPrice(split[2], charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
			if err != nil {
				err = fmt.Errorf("could not convert to your currency correctly")
				return ExperienceReserveOModel{}, err
			}
			data := DatePrice{
				Price:      tools.ConvertFloatToString(price),
				Date:       split[1],
				GroupPrice: tools.ConvertFloatToString(groupPrice),
			}
			datePrice = append(datePrice, data)
		}
	}
	mainFee, serviceFee, totalFee, guestFee, petFee, cleanFee, nightPetFee, nightGuestFee, err := GetOptionPrice(charge, userCurrency, dollarToNaira, dollarToCAD, "HandleOptionChargeToReserve")
	if err != nil {
		return ExperienceReserveOModel{}, err
	}
	res := ExperienceReserveOModel{
		Discount:        discount,
		MainPrice:       tools.ConvertFloatToString(mainFee),
		ServiceFee:      tools.ConvertFloatToString(serviceFee),
		TotalFee:        tools.ConvertFloatToString(totalFee),
		DatePrice:       datePrice,
		Currency:        userCurrency,
		Guests:          charge.Guests,
		GuestFee:        tools.ConvertFloatToString(guestFee),
		PetFee:          tools.ConvertFloatToString(petFee),
		CleaningFee:     tools.ConvertFloatToString(cleanFee),
		NightlyPetFee:   tools.ConvertFloatToString(nightPetFee),
		NightlyGuestFee: tools.ConvertFloatToString(nightGuestFee),
		CanInstantBook:  charge.CanInstantBook,
		RequireRequest:  charge.RequireRequest,
		RequestType:     charge.RequestType,
		Reference:       charge.Reference,
		OptionUserID:    tools.UuidToString(charge.OptionUserID),
		StartDate:       tools.ConvertDateOnlyToString(charge.StartDate),
		EndDate:         tools.ConvertDateOnlyToString(charge.EndDate),
	}
	return res, nil
}

func UpdateChargeOptionReferencePrice(ctx context.Context, server *Server, reserveData ExperienceReserveOModel, user db.User, functionName string) (err error) {
	var datePrice []string
	for _, d := range reserveData.DatePrice {
		data := d.Price + "&" + d.Date + "&" + d.GroupPrice
		datePrice = append(datePrice, data)
	}
	discount := reserveData.Discount.Price + "&" + reserveData.Discount.Type
	_, err = server.store.UpdateChargeOptionReferencePriceByRef(ctx, db.UpdateChargeOptionReferencePriceByRefParams{
		UserID:          user.UserID,
		Discount:        discount,
		MainPrice:       tools.MoneyStringToInt(reserveData.MainPrice),
		ServiceFee:      tools.MoneyStringToInt(reserveData.ServiceFee),
		TotalFee:        tools.MoneyStringToInt(reserveData.TotalFee),
		DatePrice:       datePrice,
		Currency:        reserveData.Currency,
		GuestFee:        tools.MoneyStringToInt(reserveData.GuestFee),
		PetFee:          tools.MoneyStringToInt(reserveData.PetFee),
		CleanFee:        tools.MoneyStringToInt(reserveData.CleaningFee),
		NightlyPetFee:   tools.MoneyStringToInt(reserveData.NightlyPetFee),
		NightlyGuestFee: tools.MoneyStringToInt(reserveData.NightlyGuestFee),
		Reference:       reserveData.Reference,
	})
	if err != nil {
		log.Printf("Error at endDate UpdateChargeOptionReferencePrice in UpdateChargeOptionReferencePriceByRef: %v optionID: %v referenceID: %v, paystack_reference: %v, functionName: %v\n", err.Error(), reserveData.OptionUserID, reserveData.Reference, reserveData.Reference, functionName)
		err = fmt.Errorf("error 453 occur, pls contact us")
		return
	}
	return
}

func GetOptionPrice(charge db.ChargeOptionReference, userCurrency string, dollarToNaira string, dollarToCAD string, funcName string) (mainFee float64, serviceFee float64, totalFee float64, guestFee float64, petFee float64, cleanFee float64, nightPetFee float64, nightGuestFee float64, err error) {
	mainFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.MainPrice), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification tools.ConvertPrice(charge.MainPrice err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	serviceFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.ServiceFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification tools.ConvertPrice(charge.ServiceFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	totalFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.TotalFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.TotalFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	guestFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.GuestFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.GuestFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	petFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.PetFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.PetFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	cleanFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.CleanFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.CleanFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	nightPetFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.NightlyPetFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.NightlyPetFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	nightGuestFee, err = tools.ConvertPrice(tools.IntToMoneyString(charge.NightlyGuestFee), charge.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName %v GetOptionPrice in CountNotification charge.NightlyGuestFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("could not convert to your currency correctly")
		return
	}
	return
}

func HandleChargeToOptionData(ctx context.Context, server *Server, charge db.ChargeOptionReference, funcName string, userCurrency string) (ExperienceOptionData, error) {
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	data, err := server.store.GetOptionExperienceByOptionUserID(ctx, db.GetOptionExperienceByOptionUserIDParams{
		OptionUserID: charge.OptionUserID,
		IsComplete:   true,
		IsActive:     true,
		IsActive_2:   true,
	})
	if err != nil {
		log.Printf("Error at FuncName %v HandleChargeToOptionData in GetOptionExperienceByOptionUserID charge.NightlyGuestFee err: %v, user: %v\n", funcName, err, charge.ID)
		err = fmt.Errorf("listing is currently unavailable")
		return ExperienceOptionData{}, err
	}
	basePrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.Price), data.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName: %v basePrice HandleChargeToOptionData in ConvertPrice err: %v, user: %v\n", funcName, err, charge.ID)
		basePrice = 0.0

	}
	weekendPrice, err := tools.ConvertPrice(tools.IntToMoneyString(data.WeekendPrice), data.Currency, userCurrency, dollarToNaira, dollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at FuncName: %v weekendPrice HandleChargeToOptionData in ConvertPrice err: %v, user: %v\n", funcName, err, charge.ID)
		weekendPrice = 0.0
	}
	addDateFound, startDateBook, endDateBook, addPrice := HandleOptionRedisExAddPrice(ctx, server, data.ID, data.OptionUserID, data.PreparationTime, data.AvailabilityWindow, data.AdvanceNotice, data.Price, data.WeekendPrice)
	addedPrice, err := tools.ConvertPrice(tools.IntToMoneyString(addPrice), data.Currency, data.Currency, server.config.DollarToNaira, server.config.DollarToCAD, charge.ID)
	if err != nil {
		log.Printf("Error at addedPrice GetDeepLinkExperience in ConvertPrice err: %v, user: %v\n", err, charge.ID)
		addedPrice = 0.0
	}
	res := ExperienceOptionData{
		UserOptionID:       tools.UuidToString(data.OptionUserID),
		Name:               data.HostNameOption,
		IsVerified:         data.IsVerified,
		CoverImage:         data.CoverImage,
		HostAsIndividual:   data.HostAsIndividual,
		BasePrice:          tools.ConvertFloatToString(basePrice),
		WeekendPrice:       tools.ConvertFloatToString(weekendPrice),
		Photos:             data.Photo,
		TypeOfShortlet:     data.TypeOfShortlet,
		State:              data.State,
		Country:            data.Country,
		ProfilePhoto:       data.Photo_2,
		HostName:           data.FirstName,
		HostJoined:         tools.ConvertDateOnlyToString(data.CreatedAt),
		HostVerified:       data.IsVerified_2,
		Category:           data.Category,
		AddedPrice:         tools.ConvertFloatToString(addedPrice),
		AddPriceFound:      addDateFound,
		StartDate:          tools.ConvertDateOnlyToString(startDateBook),
		EndDate:            tools.ConvertDateOnlyToString(endDateBook),
		PublicPhotos:       data.OptionPublicPhoto,
		PublicCoverImage:   data.PublicCoverImage,
		PublicProfilePhoto: data.HostPublicPhoto,
	}
	return res, nil
}
