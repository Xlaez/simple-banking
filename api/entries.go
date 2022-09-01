package api

import (
	"database/sql"
	"net/http"
	db "simple-bank/db/sqlc"

	"github.com/gin-gonic/gin"
)

type createEntryRequest struct {
	AccountID    int64 `json:"accountId" binding:"required,min=0"`
	Amount 		 int64 `json:"amount" binding:"required,min=0"`
}

func (server *Server) createEntry (ctx *gin.Context) {
	var request createEntryRequest

	err := ctx.ShouldBindJSON(&request)
	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	entry, err := server.store.CreateEntry(ctx, db.CreateEntryParams{
		AccountID: request.AccountID,
		Amount: request.Amount,
	})

	if err != nil {

		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return;
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusCreated, entry)
}

type getEntryReq struct {
	ID int64 `uri:"id" binding:"required,min=0"`
}

func (server *Server) getEntry(ctx *gin.Context) {
	var request getEntryReq

	err := ctx.ShouldBindUri(&request)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	entry, err := server.store.GetEntry(ctx, request.ID)

	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound, errorResponse(err))
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))	
		return;
	}

	ctx.JSON(http.StatusOK, entry)
}

type getEntriesReq struct {
	PageID int32	`form:"pageId" binding:"required,min=1"`
	PageSize int32	`form:"pageSize" binding:"required,min=3,max=10"`
}

func (server *Server) getEntries(ctx *gin.Context) {
	var request getEntriesReq

	err := ctx.ShouldBindQuery(&request)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}
	
	entries, err := server.store.ListEntries(ctx, db.ListEntriesParams{
		Limit: request.PageSize,
		Offset: (request.PageID -1) * request.PageSize,
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))	
		return;
	}

	ctx.JSON(http.StatusOK, entries)
}

type updateEntryReq struct {
	ID     int64 `json:"id" binding:"required,min=1"`
	Amount int64 `json:"amount" binding:"required,min=1"`
}

func (server *Server) updateEntry(ctx *gin.Context) {
	var request updateEntryReq
	
	err:= ctx.ShouldBindJSON(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	entry, err := server.store.UpdateEntry(ctx, db.UpdateEntryParams{
		ID: request.ID,
		Amount: request.Amount,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusOK, entry)
}

type deleteEntryReq struct {
	ID     int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteEntry(ctx *gin.Context) {
	var request deleteEntryReq
	var err error

	err = ctx.ShouldBindUri(&request)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	err = server.store.DeleteEntry(ctx, request.ID) 

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return;
	}

	ctx.JSON(http.StatusOK, "deleted")
}