package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
)

func sendMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		l := logic.NewMessageLogic(ctx, appCtx)
		if rsp, err := l.SendMessage(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, rsp)
		}
	}
}

func pushMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Type == event.PushCommonEvent || req.Type == event.PushUserEvent ||
			req.Type == event.PushFriendEvent || req.Type == event.PushGroupEvent ||
			req.Type == event.PushOtherEvent {
			l := logic.NewMessageLogic(ctx, appCtx)
			if rsp, err := l.PushMessage(req); err != nil {
				dto.ResponseInternalServerError(ctx, err)
			} else {
				dto.ResponseSuccess(ctx, rsp)
			}
		} else {
			dto.ResponseBadRequest(ctx)
			return
		}
	}
}

func getUserLatestMessages(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.GetMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if resp, err := l.GetUserMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func getUserOfflineMessages(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.GetMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if resp, err := l.GetUserMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func deleteUserMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.DeleteMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.DeleteUserMessage(&req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}
