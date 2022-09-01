package api

import (
	"database/sql"
	"net/http"
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	arg := db.CreateAccountParams{
		Owner: request.Owner,
		Currency: request.Currency,
		Balance: 0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountReq struct {
	ID int64	`uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountReq
	err := ctx.ShouldBindUri(&request)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	account, err := server.store.GetAccount(ctx, request.ID)

	if err != nil {

		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return;
		}

		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusOK, account)

}

type getAccountsReq struct {
	PageID int32	`form:"pageId" binding:"required,min=1"`
	PageSize int32	`form:"pageSize" binding:"required,min=3,max=10"`
}

func (server *Server) getAccounts(ctx *gin.Context){
	var request getAccountsReq
	err := ctx.ShouldBindQuery(&request)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}	

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit: request.PageSize,
		Offset: (request.PageID -1) * request.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusOK, accounts)

}

type updateAccountReq struct {
	 ID int64	`json:"id" binding:"required,min=1"`
	 BALANCE int64 `json:"balance" binding:"required,min=1"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var request updateAccountReq 
	var err error

	err = ctx.ShouldBindJSON(&request)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	account, err := server.store.UpdateAccount(ctx, db.UpdateAccountParams{
		ID: request.ID,
		Balance: request.BALANCE,
	})
	
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account) 
}


type deleteReq struct {
	 ID int64	`uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var request deleteReq
	var err error

	err = ctx.ShouldBindUri(&request)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	err = server.store.DeleteAccount(ctx, request.ID)
	
	if err != nil{
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, "deleted") 
}