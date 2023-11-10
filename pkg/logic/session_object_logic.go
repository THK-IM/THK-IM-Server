package logic

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
)

const (
	engineMinio = "minio"
)

type SessionObjectLogic struct {
	appCtx *app.Context
}

func NewSessionObjectLogic(appCtx *app.Context) SessionObjectLogic {
	return SessionObjectLogic{
		appCtx: appCtx,
	}
}

func (l *SessionObjectLogic) GetUploadParams(req dto.GetUploadParamsReq) (*dto.GetUploadParamsRes, error) {
	// 鉴权
	su, errSu := l.appCtx.SessionUserModel().FindSessionUser(req.SId, req.UId)
	if errSu != nil || su.UserId <= 0 {
		return nil, errorx.ErrPermission
	}
	uploadKey := fmt.Sprintf("session-%d/%d/%d-%s", req.SId, req.UId, req.ClientId, req.FName)
	uploadUrl, uploadMethod, params, err := l.appCtx.ObjectStorage().GetUploadParams(uploadKey)
	if err != nil {
		return nil, err
	}
	id, errId := l.appCtx.SessionObjectModel().Insert(req.SId, req.UId, req.ClientId, engineMinio, uploadKey)
	if errId != nil {
		return nil, errId
	} else {
		return &dto.GetUploadParamsRes{
			Id:     id,
			Url:    uploadUrl,
			Method: uploadMethod,
			Params: params,
		}, nil
	}

}

func (l *SessionObjectLogic) GetObjectByKey(req dto.GetDownloadUrlReq) (*string, error) {
	object, err := l.appCtx.SessionObjectModel().FindObject(req.Id, req.SId)
	if err != nil || object.Id == 0 {
		return nil, errorx.ErrParamsError
	}
	// 只有在后端模式ip鉴权情况下UId才会是0
	if req.UId == 0 {
		return l.appCtx.ObjectStorage().GetDownloadUrl(object.Key)
	}
	// 鉴权
	su, errSu := l.appCtx.SessionUserModel().FindSessionUser(object.SId, req.UId)
	if errSu != nil || su.UserId == 0 {
		return nil, errorx.ErrPermission
	}
	return l.appCtx.ObjectStorage().GetDownloadUrl(object.Key)
}
