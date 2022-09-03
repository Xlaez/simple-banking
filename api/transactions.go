package api

import (
	"errors"
	"fmt"
	"net/http"
	db "simple-bank/db/sqlc"
	"simple-bank/token"

	"github.com/gin-gonic/gin"
)

type createTransactionReq struct{
	FromAccountID int64 `json:"fromAccountId" binding:"required,min=0"`
	ToAccountID int64	`json:"toAccountId" binding:"required,gt=0"`
	Amount int64		`json:"amount" binding:"required,min=0"`
	Currency string     `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransaction(ctx *gin.Context){
	var request createTransactionReq

	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	fromAcc, valid := server.validateAccount(ctx, request.FromAccountID, request.Currency)
	
	if !valid {
		return;
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if fromAcc.Owner != authPayload.Username{
		err := errors.New("from account doesn't belong to this user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return;
	}

	_,valid = server.validateAccount(ctx, request.ToAccountID, request.Currency)

	if !valid {
		return;
	}

	arg := db.TransferTxParams{
		FromAccountID: request.FromAccountID,
		ToAccountID: request.ToAccountID,
		Amount: request.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusCreated, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountId int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false;
	}

	if account.Currency != currency{
		err := fmt.Errorf("account [%v] currency mismatch %v vs %v", account.ID, account.Currency, currency)

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false;
	}

	return account,true;
}


// type getTransactionReq struct{
// 	ID int64 `uri:"id" binding:"required,min=0"`
// }

// func (server *Server) getTransaction(ctx *gin.Context){
// 	var request getTransactionReq

// 	err := ctx.ShouldBindUri(&request)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	transaction, err := server.store.GetTransfer(ctx, request.ID)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	ctx.JSON(http.StatusCreated, transaction)
// }


// type listTransfersReq struct {
// 	PageId  int32 `form:"pageId"`
// 	PageSize int32 `form:"pageSize"`
// }

// func (server *Server) getTransactions (ctx *gin.Context) {
// 	var request listTransfersReq 

// 	err := ctx.ShouldBindQuery(&request)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	transactions, err := server.store.ListTransfers(ctx, db.ListTransfersParams{
// 		Limit: request.PageId,
// 		Offset: (request.PageId -1) * request.PageSize,
// 	})

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	ctx.JSON(http.StatusCreated, transactions)
// }


// type updateTransfersReq struct {
// 	ID     int64 `json:"id"`
// 	Amount int64 `json:"amount"`
// }

// func (server *Server) updateTransactions (ctx *gin.Context) {
// 	var request updateTransfersReq 

// 	err := ctx.ShouldBindJSON(&request)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	transactions, err := server.store.UpdateTransfer(ctx, db.UpdateTransferParams{
// 		ID: request.ID,
// 		Amount: request.Amount,
// 	})

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	ctx.JSON(http.StatusCreated, transactions)
// }


// type deleteTransfersReq struct {
// 	ID     int64 `uri:"id"`
// }

// func (server *Server) deleteTransactions (ctx *gin.Context) {
// 	var request deleteTransfersReq
// 	var err error
	
// 	err = ctx.ShouldBindUri(&request)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	err = server.store.DeleteTransfer(ctx, request.ID)

// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return;
// 	}

// 	ctx.JSON(http.StatusCreated, "deleted")
// }