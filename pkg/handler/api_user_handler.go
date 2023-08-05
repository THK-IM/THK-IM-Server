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
	return func(ctx *gin.Context) {
		var req dto.PostUserOnlineReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		l := logic.NewUserLogic(ctx, appCtx)
		if err := l.UpdateUserOnlineStatus(&req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func getUsersOnlineStatus(appCtx *app.Context) gin.HandlerFunc {
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

		l := logic.NewUserLogic(ctx, appCtx)
		if res, err := l.GetUsersOnlineStatus(uIds); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}
