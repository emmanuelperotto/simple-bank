package api

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	db "simplebank/db/sqlc"
)

//accountHandler handles all HTTP requests in Accounts domain.
type (
	accountHandler struct {
		store db.Store
	}
	createAccountRequest struct {
		Owner    string `json:"owner" binding:"required"`
		Currency string `json:"currency" binding:"required,oneof=USD EUR"`
	}
	getAccountRequest struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	listAccountsRequest struct {
		PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
		PageID   int32 `form:"page_id" binding:"required,min=1"`
	}
)

//newAccountHandler builds accountHandler struct
func newAccountHandler(store db.Store) accountHandler {
	return accountHandler{
		store: store,
	}
}

func (h accountHandler) post(ctx *gin.Context) {
	var req createAccountRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := h.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (h accountHandler) get(ctx *gin.Context) {
	var req getAccountRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := h.store.GetAccount(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("account not found")))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (h accountHandler) list(ctx *gin.Context) {
	var req listAccountsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := h.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("unknown error")))
		return
	}

	ctx.JSON(http.StatusOK, account)
}
