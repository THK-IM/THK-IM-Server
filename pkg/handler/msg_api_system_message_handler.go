package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
)

func pushSystemMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Type == event.PushCommonEventType || req.Type == event.PushUserEventType ||
			req.Type == event.PushFriendEventType || req.Type == event.PushGroupEventType ||
			req.Type == event.PushOtherEventType {
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

func sendSystemMessage(appCtx *app.Context) gin.HandlerFunc {
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
