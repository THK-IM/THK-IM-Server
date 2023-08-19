package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func createSession(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CreateSessionReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Type == model.SingleSessionType && (req.EntityId != nil || len(req.Members) != 2) {
			dto.ResponseBadRequest(ctx)
			return
		} else if (req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType) && req.EntityId == nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		for _, member := range req.Members {
			if member <= 0 {
				dto.ResponseBadRequest(ctx)
				return
			}
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if resp, err := l.CreateSession(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func updateSession(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UpdateSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		if id, err := strconv.Atoi(ctx.Param("id")); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		} else {
			req.Id = int64(id)
		}
		l := logic.NewSessionLogic(ctx, appCtx)
		if err := l.UpdateSession(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateUserSession(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UpdateUserSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if err := l.UpdateUserSession(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func getUserSessions(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.GetUserSessionsReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if resp, err := l.GetUserSessions(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func getUserSession(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, e1 := strconv.ParseInt(uid, 10, 64)
		iSid, e2 := strconv.ParseInt(sid, 10, 64)
		if e1 != nil || e2 != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if res, err := l.GetUserSession(iUid, iSid); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func getSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, e1 := strconv.ParseInt(sessionId, 10, 64)
		if e1 != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		var req dto.GetSessionMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SessionId = iSessionId
		l := logic.NewMessageLogic(ctx, appCtx)
		if res, err := l.GetSessionMessages(req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func deleteSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, e1 := strconv.ParseInt(sessionId, 10, 64)
		if e1 != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		var req dto.DelSessionMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SessionId = iSessionId
		l := logic.NewMessageLogic(ctx, appCtx)
		if err := l.DelSessionMessage(&req); err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}
