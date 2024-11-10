package api

import (
	"errors"
	"fmt"
	"net/http"
	db "project/simplebank/db/sqlc"
	"project/simplebank/token"

	"github.com/gin-gonic/gin"
)

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof= USD EUR RMB"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username, ///疑问为什么能返回username
		Currency: req.Currency,
		Balance:  0,
	}
	//如果省下以下这段store调用CreateAccount，则在测试的时候会提示缺少调用方法
	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {

		errCode := db.ErrorCode(err)
		if errCode == db.ForeignKeyViolation || errCode == db.UniqueViolation {
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// 得到单个用户
func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		//缺少如果没有查到id应该返回 没有找到的消息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("账户不属于认证的账户")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	//Owner    string `form:"owner"`
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1"`
}

// 得到多组用户
func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//向账户管理api添加授权
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	fmt.Printf("Request params: %+v\n", req)
	fmt.Printf("Limit: %d, Offset: %d\n", arg.Limit, arg.Offset)

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		//缺少如果没有查到id应该返回 没有找到的消息
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fmt.Printf("Limit: %d, Offset: %d\n", arg.Limit, arg.Offset)
	fmt.Printf("Accounts: %+v\n", accounts)
	//	fmt.Printf("Owner: %s\n", arg.Owner)

	if len(accounts) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "没有找到用户"})
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
