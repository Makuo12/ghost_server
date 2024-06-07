package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListRequestNotify(ctx *gin.Context) {
	var req ListRequestNotifyParams
	var onLastIndex = false
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at  ListRequestNotify in ShouldBindJSON: %v, offset: %v \n", err.Error(), req.Offset)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	contactID, err := tools.StringToUuid(req.ContactID)
	if err != nil {
		log.Printf("Error at  ListRequestNotify in GetMessageContactCount err: %v, contactID: %v\n", err, req.ContactID)
		err = fmt.Errorf("this contact does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	count, err := server.store.CountRequestNotifyID(ctx, db.CountRequestNotifyIDParams{
		SenderID:   contactID,
		ReceiverID: user.UserID,
		Type:       "user_request",
		Type_2:     "user_request",
		Cancelled:  false,
		Approved:   false,
	})
	if err != nil {
		log.Printf("Error at ListRequestNotify in .CountRequestNotifyID err: %v, user: %v\n", err, user.ID)
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	if count <= int64(req.Offset) || count == 0 {
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	requests, err := server.store.ListRequestNotifyID(ctx, db.ListRequestNotifyIDParams{
		SenderID:   contactID,
		ReceiverID: user.UserID,
		Type:       "user_request",
		Type_2:     "user_request",
		Cancelled:  false,
		Approved:   false,
		Limit:      30,
		Offset:     int32(req.Offset),
	})
	if err != nil {
		log.Printf("Error at  ListRequestNotify in ListRequestNotify err: %v, user: %v\n", err, user.ID)
		ctx.JSON(http.StatusNoContent, "none")
		return
	}
	var text string
	var res ListRequestNotifyRes
	var resData []RequestNotifyItem
	for _, r := range requests {
		optionUserID, err := tools.StringToUuid(r.ItemID)
		if err != nil {
			log.Printf("Error at  ListRequestNotify in tools.StringToUuid err: %v, user: %v, itemID: %v\n", err, user.ID, r.ItemID)
			continue
		}
		option, err := server.store.GetRequestNotifyItem(ctx, optionUserID)
		if err != nil {
			log.Printf("Error at  ListRequestNotify in GetRequestNotifyItem err: %v, user: %v, itemID: %v\n", err, user.ID, r.ItemID)
			continue
		}
		switch r.Type {
		case "user_request":
			text = "Can you host " + r.FirstName
		case "host_change_dates":
			text = r.FirstName + " changed dates"
		}
		data := RequestNotifyItem{
			MID:            tools.UuidToString(r.MID),
			MsgID:          tools.UuidToString(r.MsgID),
			HostNameOption: option.HostNameOption,
			Text:           text,
			StartDate:      r.StartDate,
			EndDate:        r.EndDate,
			MainOptionType: option.MainOptionType,
			Category:       option.Category,
			SpecialType:    r.Type,
			Reference:      r.Reference,
		}
		resData = append(resData, data)
	}
	if count <= int64(req.Offset+len(requests)) {
		onLastIndex = true
	}
	res = ListRequestNotifyRes{
		List:              resData,
		Offset:            req.Offset + len(requests),
		OnLastIndex:       onLastIndex,
		SelectedContactID: req.ContactID,
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) GetOUserRequestNotifyDetail(ctx *gin.Context) {
	var req RequestNotifyDetailParams
	dollarToNaira := server.config.DollarToNaira
	dollarToCAD := server.config.DollarToCAD
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at GetOUserRequestNotifyDetail in ShouldBindJSON: %v, reference: %v \n", err.Error(), req.Reference)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	contactID, err := tools.StringToUuid(req.ContactID)
	if err != nil {
		log.Printf("Error at GetOUserRequestNotifyDetail in GetMessageContactCount err: %v, contactID: %v\n", err, req.ContactID)
		err = fmt.Errorf("this contact does not exist")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	detail, err := server.store.GetChargeOptionReferenceByMsg(ctx, db.GetChargeOptionReferenceByMsgParams{
		Reference:  req.Reference,
		SenderID:   contactID,
		ReceiverID: user.UserID,
	})
	if err != nil {
		log.Printf("Error atGetOUserRequestNotifyDetail in GetChargeOptionReferenceByMsg err: %v, user: %v\n", err, user.ID)
		err = fmt.Errorf("could not get the details for this request, it must have expired")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var priceString string
	fee := tools.ConvertFloatToString(tools.ConvertStringToFloat(tools.IntToMoneyString(detail.TotalFee)) - tools.ConvertStringToFloat(tools.IntToMoneyString(detail.ServiceFee)))
	price, err := tools.ConvertPrice(fee, detail.Currency, req.Currency, dollarToNaira, dollarToCAD, user.ID)
	if err != nil {
		log.Printf("Error atGetOUserRequestNotifyDetail in tools.ConvertPrice err: %v, user: %v\n", err, user.ID)
		priceString = "0.00"
	} else {
		priceString = tools.ConvertFloatToString(price)
	}
	res := OUserRequestNotifyDetailRes{
		StartDate:        tools.ConvertDateOnlyToString(detail.StartDate),
		EndDate:          tools.ConvertDateOnlyToString(detail.EndDate),
		CoverImage:       detail.CoverImage,
		Guests:           detail.Guests,
		Price:            priceString,
		EmailVerified:    !tools.ServerStringEmpty(detail.Email),
		PhoneVerified:    !tools.ServerStringEmpty(detail.PhoneNumber),
		IdentityVerified: detail.IsVerified,
		Text:             fmt.Sprintf("Can you host %v?", detail.FirstName),
	}
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) MsgRequestResponse(ctx *gin.Context) {
	var req MsgRequestResponseParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at MsgRequestResponse in ShouldBindJSON: %v, MsgID: %v \n", err.Error(), req.MsgID)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msgID, err := tools.StringToUuid(req.MsgID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	msg, err := server.store.GetMessageByMsgID(ctx, db.GetMessageByMsgIDParams{
		MsgID:      msgID,
		SenderID:   user.UserID,
		ReceiverID: user.UserID,
	})
	if err != nil {
		log.Printf("Error at MsgRequestResponse in GetMessageByMsgID: %v, MsgID: %v \n", err.Error(), req.MsgID)
		err = fmt.Errorf("this message cannot be accessed by you")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	switch msg.Type {
	case "user_request":
		res, err := HandleOptionUserRequest(req, msg, user, server, ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		} else {
			ctx.JSON(http.StatusOK, res)
		}
	}
}
