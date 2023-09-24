package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
)

func updateUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PostUserOnlineReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warnf("req: %+v, err: %s", req, err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if err := l.UpdateUserOnlineStatus(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func getUsersOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUsersOnlineStatusReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if res, err := l.GetUsersOnlineStatus(req.UIds); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func kickOffUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.KickUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		if err := l.KickUser(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}

	}
}
