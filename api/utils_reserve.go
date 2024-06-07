package api

import (
	"context"
	"fmt"
	"log"

	db "github.com/makuo12/ghost_server/db/sqlc"

	"github.com/google/uuid"
)

func HandleOptionSnapShot(server *Server, ctx context.Context, reference string, paystackReference string, optionUserID uuid.UUID, chargeID uuid.UUID) (err error) {
	option, err := server.store.GetOptionInfoByUserID(ctx, optionUserID)
	var shortletJson string
	var spaceAreaJson string
	var locationJson string
	if err != nil {
		log.Printf("Error at HandleOptionReserveStore in GetOptionInfoByUserID: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		err = fmt.Errorf("error 300 occur, pls contact us")
		return
	}
	shortlet, err := server.store.GetShortlet(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleOptionReserveStore in GetShortlet: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		shortletJson = "none"
	} else {
		data, err := StructToStringShortlet(shortlet)
		if err != nil {
			log.Printf("Error at HandleOptionReserveStore in StructToStringShortlet: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
			shortletJson = "none"
		} else {
			shortletJson = data
		}
	}
	location, err := server.store.GetLocation(ctx, option.ID)
	if err != nil {
		log.Printf("Error at HandleOptionReserveStore in GetLocation: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		locationJson = "none"
	} else {
		data, err := StructToStringLocation(location)
		if err != nil {
			log.Printf("Error at HandleOptionReserveStore in StructToStringLocation: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
			locationJson = "none"
		} else {
			locationJson = data
		}
	}
	amenities, err := server.store.ListAmenitiesTag(ctx, db.ListAmenitiesTagParams{
		OptionID: option.ID,
		HasAm:    true,
	})
	if err != nil {
		// It is possible the host had no amenities
		log.Printf("Error at HandleOptionReserveStore in ListAmenitiesTag: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		amenities = []string{"none"}
	}
	// Things to note that are checked
	checkNotes, err := server.store.ListThingToNoteTag(ctx, db.ListThingToNoteTagParams{
		OptionID: option.ID,
		Checked:  true,
	})
	if err != nil {
		// It is possible the host had no thing to note
		log.Printf("Error at HandleOptionReserveStore in ListThingToNoteTag: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		checkNotes = []string{"none"}
	}
	// Things to not that are unchecked
	unCheckNotes, err := server.store.ListThingToNoteTag(ctx, db.ListThingToNoteTagParams{
		OptionID: option.ID,
		Checked:  false,
	})
	if err != nil {
		// It is possible the host had no thing to note
		log.Printf("Error at HandleOptionReserveStore in unChecked ListThingToNoteTag: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		unCheckNotes = []string{"none"}
	}
	spaceArea, err := server.store.ListSpaceArea(ctx, option.ID)
	if err != nil {
		// It is possible the host had no space area
		log.Printf("Error at HandleOptionReserveStore in unChecked ListSpaceArea: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		spaceAreaJson = "none"
	} else {
		data, err := StructToStringSpaceArea(spaceArea)
		if err != nil {
			log.Printf("Error at HandleOptionReserveStore in StructToStringSpaceArea: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
			spaceAreaJson = "none"
		} else {
			spaceAreaJson = data
		}
	}
	err = server.store.CreateOptionReferenceInfo(ctx, db.CreateOptionReferenceInfoParams{
		OptionChargeID:   chargeID,
		Amenities:        amenities,
		SpaceArea:        spaceAreaJson,
		TimeZone:         option.TimeZone,
		ArriveBefore:     option.ArriveBefore,
		ArriveAfter:      option.ArriveAfter,
		LeaveBefore:      option.LeaveBefore,
		CancelPolicyOne:  option.TypeOne,
		CancelPolicyTwo:  option.TypeTwo,
		PetsAllowed:      option.PetsAllowed,
		RulesChecked:     checkNotes,
		RulesUnchecked:   unCheckNotes,
		Shortlet:         shortletJson,
		Location:         locationJson,
		HostAsIndividual: option.HostAsIndividual,
		OrganizationName: option.OrganizationName,
	})
	if err != nil {
		log.Printf("Error at HandleOptionReserveStore in CreateOptionReferenceInfo: %v optionID: %v referenceID: %v, paystack_reference: %v\n", err.Error(), optionUserID, reference, paystackReference)
		err = fmt.Errorf("error 304 occur, pls contact us")
		return
	}
	err = nil
	return

}
