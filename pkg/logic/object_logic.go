package logic

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"time"
)

const (
	engineMinio = "minio"
)

type ObjectLogic struct {
	appCtx *app.Context
}

func NewObjectLogic(appCtx *app.Context) ObjectLogic {
	return ObjectLogic{
		appCtx: appCtx,
	}
}

func (l *ObjectLogic) GetUploadParams(req dto.GetUploadParamsReq) (*dto.GetUploadParamsRes, error) {
	// 鉴权
	su, errSu := l.appCtx.SessionUserModel().FindSessionUser(req.SId, req.UId)
	if errSu != nil || su.UserId <= 0 {
		return nil, errorx.ErrPermission
	}
	now := time.Now().UnixMilli()
	uploadKey := fmt.Sprintf("%d/%d/%d-%s", req.SId, req.UId, now, req.FileName)
	uploadUrl, uploadMethod, params, err := l.appCtx.ObjectStorage().GetUploadParams(uploadKey)
	if err != nil {
		return nil, err
	}
	id, errInsert := l.appCtx.ObjectModel().Insert(req.SId, engineMinio, uploadKey)
	if errInsert != nil {
		return nil, errInsert
	}
	return &dto.GetUploadParamsRes{
		Id:     id,
		Url:    uploadUrl,
		Method: uploadMethod,
		Params: params,
	}, nil
}

func (l *ObjectLogic) GetObjectByKey(req dto.GetDownloadUrlReq) (*string, error) {
	object, err := l.appCtx.ObjectModel().FindOne(req.Id)
	if err != nil {
		return nil, errorx.ErrParamsError
	}
	// 鉴权
	if object.SId > 0 && req.UId > 0 {
		su, errSu := l.appCtx.SessionUserModel().FindSessionUser(object.SId, req.UId)
		if errSu != nil || su.UserId <= 0 {
			return nil, errorx.ErrPermission
		}
	}
	return l.appCtx.ObjectStorage().GetDownloadUrl(object.Key)
}
