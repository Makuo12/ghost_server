package api

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/makuo12/ghost_server/db/sqlc"
	"github.com/makuo12/ghost_server/tools"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// Account Number and USSD

func (server *Server) ListBank(ctx *gin.Context) {
	var req ListBankParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at ListBank in ShouldBindJSON: %v, Country: %v \n", err.Error(), req.Country)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res, err := HandleListBank(server, ctx, req.Country)
	if err != nil {
		log.Printf("Error at ListBank in HandleListBank: %v, country: %v, userID: %v \n", err.Error(), req.Country, user.ID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	log.Printf("ListBank sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) ListAccountNumber(ctx *gin.Context) {
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	var res ListAccountNumberRes
	numbers, err := server.store.ListAccountNumber(ctx, user.ID)
	if err != nil || len(numbers) == 0 {
		if err != nil {
			log.Printf("Error at ListAccountNumber in ListAccountNumber: %v, userID: %v \n", err.Error(), user.ID)
		}
		data := []AccountNumberItem{{"none", "none", "none", "none", "none"}}
		res = ListAccountNumberRes{
			List:             data,
			IsEmpty:          true,
			DefaultAccountID: "none",
		}
	} else {
		var resData []AccountNumberItem
		for _, n := range numbers {
			data := AccountNumberItem{
				AccountNumber: n.AccountNumber,
				ID:            tools.UuidToString(n.ID),
				BankName:      n.BankName,
				Currency:      n.Currency,
				AccountName:   n.AccountName,
			}
			resData = append(resData, data)
		}
		res = ListAccountNumberRes{
			List:             resData,
			IsEmpty:          false,
			DefaultAccountID: user.DefaultAccountID,
		}
	}
	log.Printf("ListAccountNumber sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)
}

func (server *Server) CreateAccountNumber(ctx *gin.Context) {
	var req CreateAccountNumberParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at CreateAccountNumber in ShouldBindJSON: %v, Country: %v \n", err.Error(), req.Country)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Next want want to verify account number
	accountData, err := VerifyAccountNumber(server, ctx, req.AccountNumber, req.Code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	recipientData, err := CreateTransferRecipient(server, ctx, accountData.Data.AccountNumber, req.Code, accountData.Data.AccountName, req.Currency)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if !recipientData.Data.Active || !recipientData.Status {
		err = fmt.Errorf("%v is either not active or cannot be used as a recipient, try again", req.AccountNumber)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, err := server.store.CreateAccountNumber(ctx, db.CreateAccountNumberParams{
		UserID:        user.ID,
		AccountNumber: accountData.Data.AccountNumber,
		AccountName:   accountData.Data.AccountName,
		BankCode:      req.Code,
		Country:       req.Country,
		Currency:      req.Currency,
		BankName:      req.BankName,
		RecipientCode: recipientData.Data.RecipientCode,
		Type:          "nuban",
		BankID:        int32(accountData.Data.BankID),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	// We make it the default
	_, err = server.store.UpdateUser(ctx, db.UpdateUserParams{
		DefaultAccountID: pgtype.Text{
			String: tools.UuidToString(account.ID),
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("Error at CreateAccountNumber in UpdateUser: %v, Country: %v \n", err.Error(), req.Country)
	}

	res := AccountNumberItem{
		AccountNumber: account.AccountNumber,
		ID:            tools.UuidToString(account.ID),
		BankName:      account.BankName,
		Currency:      account.Currency,
		AccountName:   account.AccountName,
	}
	log.Printf("CreateAccountNumberRes sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) SetDefaultAccountNumber(ctx *gin.Context) {
	var req SetDefaultAccountNumberParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at SetDefaultAccountNumber in ShouldBindJSON: %v, AccountID: %v \n", err.Error(), req.AccountID)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	updateUser, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		DefaultAccountID: pgtype.Text{
			String: req.AccountID,
			Valid:  true,
		},
		ID: user.ID,
	})
	if err != nil {
		log.Printf("Error at SetDefaultAccountNumber in UpdateUser: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("could not set this account number as your default account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := SetDefaultAccountNumberRes{
		Success:   true,
		AccountID: updateUser.DefaultAccountID,
	}
	log.Printf("CreateAccountNumberRes sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)

}

func (server *Server) RemoveAccountNumber(ctx *gin.Context) {
	var req RemoveAccountNumberParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("Error at RemoveAccountNumber in ShouldBindJSON: %v, AccountID: %v \n", err.Error(), req.AccountID)
		err = fmt.Errorf("select one of the provided options")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	accountID, err := tools.StringToUuid(req.AccountID)

	if err != nil {
		log.Printf("Error at RemoveAccountNumber in UpdateUser: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("could not remove this account number")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.RemoveAccountNumber(ctx, db.RemoveAccountNumberParams{
		UserID: user.ID,
		ID:     accountID,
	})
	if err != nil {
		log.Printf("Error at RemoveAccountNumber in RemoveAccountNumber: %v, userID: %v \n", err.Error(), user.ID)
		err = fmt.Errorf("could not remove this account number")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := RemoveAccountNumberRes{
		Success:   true,
		AccountID: req.AccountID,
	}
	log.Printf("RemoveAccountNumberRes sent successfully (%v) id: %v\n", user.Email, user.ID)
	ctx.JSON(http.StatusOK, res)

}


func (server *Server) ListUSSD(ctx *gin.Context) {
	_, err := HandleGetUser(ctx, server)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	res := ListUSSDRes {
		List: tools.USSDNames,
	}
	ctx.JSON(http.StatusOK, res)
}