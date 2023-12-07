package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
)

func pushExtendedMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PushMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Type > event.SignalExtended {
			if rsp, err := l.PushMessage(req); err != nil {
				appCtx.Logger().Warn(err.Error())
				dto.ResponseInternalServerError(ctx, err)
			} else {
				dto.ResponseSuccess(ctx, rsp)
			}
		} else {
			appCtx.Logger().Warn("param type error")
			dto.ResponseBadRequest(ctx)
			return
		}
	}
}

func sendSystemMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SendMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if rsp, err := l.SendMessage(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, rsp)
		}
	}
}
