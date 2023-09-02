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
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.CreateSessionReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		appCtx.Logger().Info(req)
		if req.Type == model.SingleSessionType && (req.EntityId != nil || len(req.Members) != 2) {
			appCtx.Logger().Warn("param type error")
			dto.ResponseBadRequest(ctx)
			return
		} else if (req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType) && req.EntityId == nil {
			appCtx.Logger().Warn("param type error")
			dto.ResponseBadRequest(ctx)
			return
		}

		for _, member := range req.Members {
			if member <= 0 {
				appCtx.Logger().Warn("param members error")
				dto.ResponseBadRequest(ctx)
				return
			}
		}

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 {
			if requestUid != req.Members[0] {
				appCtx.Logger().Warn("param uid error")
				dto.ResponseBadRequest(ctx)
				return
			}
		}

		if resp, err := l.CreateSession(req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func updateSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.UpdateSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil {
			if *req.Mute != 0 && *req.Mute != 1 {
				appCtx.Logger().Warn("param mute error")
				dto.ResponseBadRequest(ctx)
				return
			}
		}
		if id, err := strconv.Atoi(ctx.Param("id")); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		} else {
			req.Id = int64(id)
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(req.Id, requestUid); err != nil {
				dto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role == model.SessionMember {
					appCtx.Logger().Warn("permission error")
					dto.ResponseForbidden(ctx)
					return
				}
			}
		}

		if err := l.UpdateSession(req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.UpdateUserSessionReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Status != nil && (*req.Status < 0 || *req.Status > 3) {
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}

		if err := l.UpdateUserSession(req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func getUserSessions(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUserSessionsReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}

		if resp, err := l.GetUserSessions(req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func getUserSession(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			uid = ctx.Param("uid")
			sid = ctx.Param("sid")
		)

		iUid, e1 := strconv.ParseInt(uid, 10, 64)
		if e1 != nil {
			appCtx.Logger().Warn(e1)
			dto.ResponseBadRequest(ctx)
			return
		}

		iSid, e2 := strconv.ParseInt(sid, 10, 64)
		if e2 != nil {
			appCtx.Logger().Warn(e2)
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 && requestUid != iUid {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}

		if res, err := l.GetUserSession(iUid, iSid); err != nil {
			appCtx.Logger().Warn(e2)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func getSessionMessages(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSession := strconv.ParseInt(sessionId, 10, 64)
		if errSession != nil {
			appCtx.Logger().Warn(errSession)
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 {
			if _, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().Warn(err)
				dto.ResponseForbidden(ctx)
				return
			}
		}
		var req dto.GetSessionMessageReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SId = iSessionId
		if res, err := l.GetSessionMessages(req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func deleteSessionMessage(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewMessageLogic(appCtx)
	return func(ctx *gin.Context) {
		var (
			sessionId = ctx.Param("id")
		)
		iSessionId, errSessionId := strconv.ParseInt(sessionId, 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Warn(errSessionId)
			dto.ResponseBadRequest(ctx)
			return
		}
		var req dto.DelSessionMessageReq
		if err := ctx.BindJSON(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 {
			if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(iSessionId, requestUid); err != nil {
				appCtx.Logger().Warn(err)
				dto.ResponseForbidden(ctx)
				return
			} else {
				if sessionUser.Role != model.SessionOwner {
					appCtx.Logger().Warn("permission error")
					dto.ResponseForbidden(ctx)
					return
				}
			}
		}
		req.SId = iSessionId
		if err := l.DelSessionMessage(&req); err != nil {
			appCtx.Logger().Warn(err)
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}
