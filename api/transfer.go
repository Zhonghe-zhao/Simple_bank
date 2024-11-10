package api

import (
	"errors"
	"fmt"
	"net/http"
	db "project/simplebank/db/sqlc"
	"project/simplebank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	var req transferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// 获取并处理 FromAccount
	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account 不属于当前账户")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// 获取并处理 ToAccount
	toAccount, valid := server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        req.Amount,
	}
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// 检查id是否存在和货币是否匹配
func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		//第一种情况账户没有找到
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		//其他情况 服务器内部错误
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
