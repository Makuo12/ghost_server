package api

import (
	db "flex_server/db/sqlc"
	"flex_server/tools"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (server *Server) ListEventPayout(ctx *gin.Context) {
	var req ListEventPayoutParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  EventDateParams in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	eventDates, err := server.store.ListEventDateTimeHost(ctx, db.ListEventDateTimeHostParams{
		CoUserID: tools.UuidToString(user.UserID),
		HostID:   user.ID,
	})
	if err != nil {
		log.Printf("Error at ListEventPayout in ListEventDateTimeHost err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	var resData []PayoutEventItem
	for _, ed := range eventDates {
		var startDate time.Time
		var endDate time.Time
		switch ed.Type {
		case "single":
			startDate = ed.StartDate
			endDate = ed.EndDate
		case "recurring":
			startDate, err = tools.ConvertDateOnlyStringToDate(ed.EventDate)
			if err != nil {
				log.Printf("Error at ListEventPayout in ListChargeTicketReferencePayout err: %v, user: %v, eventID: %v\n", err, user.ID, ed.EventDateTimeID)
				continue
			}
			endDate = startDate
		}
		log.Println("payout start time, ", startDate)
		log.Println("payout end time, ", endDate)
		payout, err := server.store.ListChargeTicketReferencePayout(ctx, db.ListChargeTicketReferencePayoutParams{
			Date:            startDate,
			Cancelled:       false,
			EventDateID:     ed.EventDateTimeID,
			PaymentComplete: true,
			PayoutComplete:  req.IsComplete,
		})
		if err != nil || len(payout) == 0 {
			if err != nil {
				log.Printf("Error at ListEventPayout in ListChargeTicketReferencePayout err: %v, user: %v, eventID: %v\n",
					err, user.ID, ed.EventDateTimeID)
			}
			log.Println("payout start time continue, ", startDate)
			log.Println("payout end time continue, ", endDate)
			continue
		}

		amount, date := GetDateAndAmount(payout)
		data := PayoutEventItem{
			ID:              tools.UuidToString(uuid.New()),
			DatePaid:        tools.ConvertDateOnlyToString(date),
			Amount:          tools.ConvertFloatToString(amount),
			Currency:        payout[0].Currency,
			StartDate:       tools.ConvertDateOnlyToString(startDate),
			EndDate:         tools.ConvertDateOnlyToString(endDate),
			HostOptionName:  ed.HostNameOption,
			AccountNumber:   payout[0].AccountNumber,
			EventDateType:   ed.Type,
			EventDateTimeID: tools.UuidToString(ed.EventDateTimeID),
		}
		log.Println("payout start time data, ", data)
		resData = append(resData, data)
	}
	log.Println("payout start time resData, ", resData)
	resDataOffset := PayoutEventItemOffset(resData, req.Offset, 10)
	log.Println("payout start time resDataOffset, ", resDataOffset)
	res := ListPayoutEventRes{
		List:       resDataOffset,
		IsComplete: req.IsComplete,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListOptionPayout(ctx *gin.Context) {
	var req ListOptionPayoutParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListOptionPayout in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountOptionMainPayout(ctx, db.CountOptionMainPayoutParams{
		PayoutComplete:        req.IsComplete,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
		HostID:                user.ID,
	})
	if err != nil {
		log.Printf("Error at  ListOptionPayout in .CountOptionMainPayout err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	payouts, err := server.store.ListOptionMainPayout(ctx, db.ListOptionMainPayoutParams{
		PayoutComplete:        req.IsComplete,
		ChargePaymentComplete: true,
		ChargeCancelled:       false,
		HostID:                user.ID,
		Limit:                 40,
		Offset:                int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListOptionPayout in .ListOptionMainPayout err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []PayoutOptionItem
	for _, p := range payouts {
		data := PayoutOptionItem{
			ID:             tools.UuidToString(uuid.New()),
			DatePaid:       tools.ConvertDateOnlyToString(p.TimePaid),
			Amount:         tools.IntToMoneyString(p.Amount),
			GuestName:      p.GuestName,
			HostOptionName: p.HostNameOption,
			StartDate:      tools.ConvertDateOnlyToString(p.StartDate),
			EndDate:        tools.ConvertDateOnlyToString(p.EndDate),
			AccountNumber:  tools.AccountNumberToFour(p.AccountNumber),
			Currency:       p.Currency,
		}
		resData = append(resData, data)
	}
	res := ListPayoutOptionRes{
		List:       resData,
		IsComplete: req.IsComplete,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListOptionPayment(ctx *gin.Context) {
	var req ListOptionPaymentParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListOptionPayment in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountOptionPaymentByUserID(ctx, db.CountOptionPaymentByUserIDParams{
		UserID:     user.UserID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  ListOptionPayment in .CountOptionPaymentByUserID err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	payments, err := server.store.ListOptionPaymentByUserID(ctx, db.ListOptionPaymentByUserIDParams{
		UserID:     user.UserID,
		IsComplete: true,
		Limit:      40,
		Offset:     int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListOptionPayment in ListOptionPaymentByUserID err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []PaymentOptionItem
	for _, p := range payments {
		data := PaymentOptionItem{
			ID:             tools.UuidToString(uuid.New()),
			DatePaid:       tools.ConvertDateOnlyToString(p.DateBooked),
			Amount:         tools.IntToMoneyString(p.TotalFee),
			HostName:       p.FirstName,
			HostOptionName: p.HostNameOption,
			StartDate:      tools.ConvertDateOnlyToString(p.StartDate),
			EndDate:        tools.ConvertDateOnlyToString(p.EndDate),
			Cancelled:      p.Cancelled,
			Currency:       p.Currency,
		}
		resData = append(resData, data)
	}
	res := ListPaymentOptionRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListTicketPayment(ctx *gin.Context) {
	var req ListTicketPaymentParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListTicketPayment in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountTicketPaymentUser(ctx, db.CountTicketPaymentUserParams{
		UserID:     user.UserID,
		IsComplete: true,
	})
	if err != nil {
		log.Printf("Error at  ListTicketPayment in .CountTicketPaymentUser err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	payments, err := server.store.ListTicketPaymentUser(ctx, db.ListTicketPaymentUserParams{
		UserID:     user.UserID,
		IsComplete: true,
		Limit:      40,
		Offset:     int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListTicketPayment in ListTicketPaymentUser err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []PaymentTicketItem
	for _, p := range payments {
		data := PaymentTicketItem{
			ID:             tools.UuidToString(uuid.New()),
			DatePaid:       tools.ConvertDateOnlyToString(p.DateBooked),
			Amount:         tools.IntToMoneyString(p.Price),
			HostName:       p.FirstName,
			HostOptionName: p.HostNameOption,
			StartDate:      tools.ConvertDateOnlyToString(p.StartDate),
			EndDate:        tools.ConvertDateOnlyToString(p.EndDate),
			Cancelled:      p.Cancelled,
			Currency:       p.Currency,
		}
		resData = append(resData, data)
	}
	res := ListPaymentTicketRes{
		List: resData,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListRefund(ctx *gin.Context) {
	var req ListRefundParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListRefund in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountRefund(ctx, db.CountRefundParams{
		RefundComplete: req.IsComplete,
		UID:              user.UserID,
	})
	if err != nil {
		log.Printf("Error at  ListRefund in .CountRefund err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	refunds, err := server.store.ListRefund(ctx, db.ListRefundParams{
		PayoutIsComplete: req.IsComplete,
		UID:              user.UserID,
		Limit:            40,
		Offset:           int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListRefund in ListRefund err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []RefundItem
	for _, p := range refunds {
		var startDate time.Time
		var endDate time.Time
		var cancelled bool
		var currency string
		switch p.Type {
		case "charge_option_reference":
			startDate = HandleSqlNullTime(p.OptionStartDate)
			endDate = HandleSqlNullTime(p.OptionEndDate)
			cancelled = HandleSqlNullBool(p.OptionCancelled)
			currency = HandleSqlNullString(p.OptionCurrency)
		case "charge_ticket_reference":
			startDate = HandleSqlNullTime(p.EventStartDate)
			endDate = HandleSqlNullTime(p.EventEndDate)
			cancelled = HandleSqlNullBool(p.TicketCancelled)
			currency = HandleSqlNullString(p.EventCurrency)
		}
		data := RefundItem{
			ID:             tools.UuidToString(uuid.New()),
			DatePaid:       tools.ConvertDateOnlyToString(p.TimePaid),
			Amount:         tools.IntToMoneyString(p.Amount),
			HostName:       HandleSqlNullString(p.HostName),
			HostOptionName: HandleSqlNullString(p.HostNameOption),
			StartDate:      tools.ConvertDateOnlyToString(startDate),
			EndDate:        tools.ConvertDateOnlyToString(endDate),
			Cancelled:      cancelled,
			Currency:       currency,
		}
		resData = append(resData, data)
	}
	res := ListRefundRes{
		List:       resData,
		IsComplete: req.IsComplete,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListRefundPayout(ctx *gin.Context) {
	var req ListRefundParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListRefundPayout in ShouldBindJSON: %v, Offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountRefundPayout(ctx, db.CountRefundPayoutParams{
		PayoutIsComplete: req.IsComplete,
		UID:              user.ID,
	})
	if err != nil {
		log.Printf("Error at  ListRefundPayout in .CountRefundPayout err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
		return
	}
	refundPayouts, err := server.store.ListRefundPayout(ctx, db.ListRefundPayoutParams{
		PayoutIsComplete: req.IsComplete,
		UID:              user.ID,
		Limit:            40,
		Offset:           int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListRefundPayout in ListRefundPayout err: %v, user: %v\n", err, user.ID)
		res := "none"
		ctx.JSON(http.StatusNoContent, res)
	}
	var resData []RefundPayoutItem
	for _, p := range refundPayouts {
		var startDate time.Time
		var endDate time.Time
		var cancelled bool
		var currency string
		switch p.Type {
		case "charge_option_reference":
			startDate = HandleSqlNullTime(p.OptionStartDate)
			endDate = HandleSqlNullTime(p.OptionEndDate)
			cancelled = HandleSqlNullBool(p.OptionCancelled)
			currency = HandleSqlNullString(p.OptionCurrency)
		case "charge_ticket_reference":
			startDate = HandleSqlNullTime(p.EventStartDate)
			endDate = HandleSqlNullTime(p.EventEndDate)
			cancelled = HandleSqlNullBool(p.TicketCancelled)
			currency = HandleSqlNullString(p.EventCurrency)
		}
		data := RefundPayoutItem{
			ID:             tools.UuidToString(uuid.New()),
			DatePaid:       tools.ConvertDateOnlyToString(p.TimePaid),
			Amount:         tools.IntToMoneyString(p.Amount),
			GuestName:      HandleSqlNullString(p.GuestName),
			HostOptionName: HandleSqlNullString(p.HostNameOption),
			StartDate:      tools.ConvertDateOnlyToString(startDate),
			EndDate:        tools.ConvertDateOnlyToString(endDate),
			Cancelled:      cancelled,
			Currency:       currency,
		}
		resData = append(resData, data)
	}
	res := ListRefundPayoutRes{
		List:       resData,
		IsComplete: req.IsComplete,
	}
	ctx.JSON(http.StatusOK, res)
}
