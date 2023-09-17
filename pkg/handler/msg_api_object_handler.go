package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/logic"
	"github.com/gin-gonic/gin"
	"strconv"
)

func getUploadParams(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		var req dto.GetUploadParamsReq
		if err := ctx.BindQuery(&req); err != nil {
			appCtx.Logger().Warn(err.Error())
			dto.ResponseBadRequest(ctx)
			return
		}
		requestUid := ctx.GetInt64(uidKey)
		if requestUid > 0 && requestUid != req.UId {
			appCtx.Logger().Warn("param uid error")
			dto.ResponseForbidden(ctx)
			return
		}

		res, err := l.GetUploadParams(req)
		if err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			dto.ResponseSuccess(ctx, res)
		}
	}
}

func getObject(appCtx *app.Context) gin.HandlerFunc {
	l := logic.NewObjectLogic(appCtx)
	return func(ctx *gin.Context) {
		id, errParams := strconv.Atoi(ctx.Param("id"))
		if errParams != nil {
			appCtx.Logger().Warn(errParams.Error())
			dto.ResponseBadRequest(ctx)
			return
		}

		requestUid := ctx.GetInt64(uidKey)
		req := dto.GetDownloadUrlReq{
			UId: requestUid,
			Id:  int64(id),
		}

		path, err := l.GetObjectByKey(req)
		if err != nil {
			dto.ResponseInternalServerError(ctx, err)
		} else {
			if path != nil {
				dto.Redirect302(ctx, *path)
			} else {
				dto.ResponseInternalServerError(ctx, errorx.ErrServerUnknown)
			}
		}
	}
}
