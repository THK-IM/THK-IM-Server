package logic

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
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
	id, errId := l.appCtx.ObjectModel().Insert(req.SId, engineMinio, uploadKey)
	if errId != nil {
		return nil, err
	}

	id, errId = l.appCtx.SessionObjectModel().Insert(id, req.SId, req.UId, req.ClientId)
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
	var (
		object *model.Object
		err    error
	)
	if req.UId == 0 {
		// 后端模式ip鉴权情况下UId为0
		object, err = l.appCtx.ObjectModel().FindObject(req.Id)
	} else {
		// 鉴权
		userSessionTableName := l.appCtx.UserSessionModel().GenUserSessionTableName(req.UId)
		object, err = l.appCtx.ObjectModel().FindObjectByUId(req.Id, req.UId, userSessionTableName)
	}

	if err != nil || object.Id == 0 {
		return nil, errorx.ErrParamsError
	}
	return l.appCtx.ObjectStorage().GetDownloadUrl(object.Key)
}
