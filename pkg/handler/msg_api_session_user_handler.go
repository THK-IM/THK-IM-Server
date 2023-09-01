package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetSessionUserReq
		if err := ctx.ShouldBindQuery(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Count <= 0 {
			appCtx.Logger().Warn("param count error")
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().Warn("param role error")
			dto.ResponseBadRequest(ctx)
			return
		}
		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Warn(errSessionId)
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkReadPermission(appCtx, requestUid, sessionId); !hasPermission {
				appCtx.Logger().Warn("permission error")
				dto.ResponseForbidden(ctx)
				return
			}
		}
		if resp, err := l.GetUser(req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, resp)
		}
	}
}

func addSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionAddUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Role > model.SessionOwner || req.Role < model.SessionMember {
			appCtx.Logger().Warn("param role error")
			dto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Warn(errSessionId.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Warn("permission error")
				dto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.AddUser(sessionId, req); e != nil {
			appCtx.Logger().Warn(e.Error())
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func deleteSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionDelUserReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			appCtx.Logger().Warn(errSessionId.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Warn("permission error")
				dto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.DelUser(sessionId, req); e != nil {
			appCtx.Logger().Warn(e.Error())
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func updateSessionUser(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewSessionLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.SessionUserUpdateReq
		if err := ctx.ShouldBindJSON(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		if req.Mute != nil && *req.Mute != 0 && *req.Mute != 1 {
			appCtx.Logger().Warn("param mute error")
			dto.ResponseBadRequest(ctx)
			return
		}

		if req.Role != nil && (*req.Role > model.SessionOwner || *req.Role < model.SessionMember) {
			appCtx.Logger().Warn("param role error")
			dto.ResponseBadRequest(ctx)
			return
		}

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if errSessionId != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				appCtx.Logger().Warn("permission error")
				dto.ResponseForbidden(ctx)
				return
			}
		}

		if e := l.UpdateSessionUser(req); e != nil {
			appCtx.Logger().Warn(e.Error())
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
}

func checkReadPermission(appCtx *app.Context, uId, sessionId int64) bool {
	if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId); err != nil {
		appCtx.Logger().Error(err)
		return false
	} else {
		if sessionUser.UserId > 0 {
			return true
		}
	}
	return true
}

func checkPermission(appCtx *app.Context, uId, sessionId int64, oprUIds []int64) bool {
	if sessionUser, err := appCtx.SessionUserModel().FindSessionUser(sessionId, uId); err != nil {
		appCtx.Logger().Error(err)
		return false
	} else {
		if sessionUser.Role <= model.SessionAdmin {
			return false
		}
		sessionUsers, errSessionUser := appCtx.SessionUserModel().FindSessionUsers(sessionId, oprUIds)
		if errSessionUser != nil {
			return false
		}
		for _, su := range sessionUsers {
			if su.Role >= sessionUser.Role {
				return false
			}
		}
	}
	return true
}
