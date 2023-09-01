package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func updateUserOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewUserLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.PostUserOnlineReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		if err := l.UpdateUserOnlineStatus(&req); err != nil {
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
		if err := ctx.ShouldBindQuery(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		strUIds := strings.Split(req.UIds, "#")
		if len(strUIds) == 0 {
			dto.ResponseBadRequest(ctx)
			return
		}
		uIds := make([]int64, len(strUIds))
		for _, strUid := range strUIds {
			if uId, e := strconv.ParseInt(strUid, 10, 64); e != nil {
				dto.ResponseBadRequest(ctx)
				return
			} else {
				uIds = append(uIds, uId)
			}
		}

		if res, err := l.GetUsersOnlineStatus(uIds); err != nil {
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
			dto.ResponseBadRequest(ctx)
			return
		}
		if err := l.KickUser(&req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}

	}
}
