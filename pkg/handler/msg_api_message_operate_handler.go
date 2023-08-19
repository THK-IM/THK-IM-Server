package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
)

func ackUserMessages(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.AckUserMessagesReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.AckUserMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func readUserMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ReadUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.ReadUserMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func revokeUserMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.RevokeUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.RevokeUserMessage(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func reeditUserMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ReeditUserMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.ReeditUserMessage(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}
