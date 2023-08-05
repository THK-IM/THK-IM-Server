package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
	"strconv"
)

func addSessionUser(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SessionAddUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		if sid, err := strconv.Atoi(ctx.Param("id")); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		} else {
			l := logic.NewSessionLogic(ctx, appCtx)
			if e := l.AddUser(int64(sid), req); e != nil {
				dto.ResponseInternalServerError(ctx, e)
			} else {
				dto.ResponseSuccess(ctx, nil)
			}
		}
	}
}

func deleteSessionUser(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SessionDelUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		if sid, err := strconv.Atoi(ctx.Param("id")); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		} else {
			l := logic.NewSessionLogic(ctx, appCtx)
			if e := l.DelMember(int64(sid), req); e != nil {
				dto.ResponseInternalServerError(ctx, e)
			} else {
				dto.ResponseSuccess(ctx, nil)
			}
		}
	}
}

func updateSessionUser(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SessionUserUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		if uId, err := strconv.Atoi(ctx.Param("uid")); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		} else {
			req.UId = int64(uId)
		}

		if sid, err := strconv.Atoi(ctx.Param("id")); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		} else {
			req.SId = int64(sid)
			l := logic.NewSessionLogic(ctx, appCtx)
			if e := l.UpdateSessionUser(req); e != nil {
				dto.ResponseInternalServerError(ctx, e)
			} else {
				dto.ResponseSuccess(ctx, nil)
			}
		}

	}
}
