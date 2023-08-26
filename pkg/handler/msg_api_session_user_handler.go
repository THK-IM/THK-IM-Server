package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
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

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 64, 10)
		if errSessionId != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				dto.ResponseForbidden(ctx)
				return
			}
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if e := l.AddUser(sessionId, req); e != nil {
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
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

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 64, 10)
		if errSessionId != nil {
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				dto.ResponseForbidden(ctx)
				return
			}
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if e := l.DelUser(sessionId, req); e != nil {
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
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

		sessionId, errSessionId := strconv.ParseInt(ctx.Param("id"), 64, 10)
		if errSessionId != nil {
			dto.ResponseBadRequest(ctx)
			return
		}
		req.SId = sessionId

		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 { // 检查角色权限
			if hasPermission := checkPermission(appCtx, requestUid, sessionId, req.UIds); !hasPermission {
				dto.ResponseForbidden(ctx)
				return
			}
		}

		l := logic.NewSessionLogic(ctx, appCtx)
		if e := l.UpdateSessionUser(req); e != nil {
			dto.ResponseInternalServerError(ctx, e)
		} else {
			dto.ResponseSuccess(ctx, nil)
		}
	}
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
