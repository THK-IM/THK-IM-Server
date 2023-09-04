package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
)

func (l *SessionLogic) GetUser(req dto.GetSessionUserReq) (*dto.GetSessionUserRes, error) {
	sessionUser, err := l.appCtx.SessionUserModel().FindSessionUsersByMTime(req.SId, req.MTime, req.Role, req.Count)
	if err != nil {
		return nil, err
	}
	dtoSessionUsers := make([]*dto.SessionUser, 0)
	for _, su := range sessionUser {
		dtoSu := l.convSessionUser(su)
		dtoSessionUsers = append(dtoSessionUsers, dtoSu)
	}
	return &dto.GetSessionUserRes{Data: dtoSessionUsers}, nil
}

func (l *SessionLogic) AddUser(sid int64, req dto.SessionAddUserReq) error {
	session, err := l.appCtx.SessionModel().FindSession(sid, nil)
	if err != nil {
		return err
	}
	maxCount := 0
	if session.Type == model.GroupSessionType {
		maxCount = l.appCtx.Config().IM.MaxGroupMember
	} else if session.Type == model.SuperGroupSessionType {
		maxCount = l.appCtx.Config().IM.MaxSuperGroupMember
	} else {
		return errorx.ErrSessionType
	}
	roles := make([]int, 0)
	entityIds := make([]int64, 0)
	for range req.UIds {
		roles = append(roles, req.Role)
		entityIds = append(entityIds, req.EntityId)
	}
	tx := l.appCtx.Database().Begin()
	err = l.appCtx.SessionUserModel().AddUser(session, entityIds, req.UIds, roles, maxCount, tx)
	if err != nil {
		_ = tx.Rollback()
	} else {
		err = tx.Commit().Error
	}
	return err
}

func (l *SessionLogic) DelUser(sid int64, req dto.SessionDelUserReq) error {
	session, err := l.appCtx.SessionModel().FindSession(sid, nil)
	if err != nil {
		return err
	}

	return l.appCtx.SessionUserModel().DelUser(session, req.UIds)
}

func (l *SessionLogic) UpdateSessionUser(req dto.SessionUserUpdateReq) (err error) {
	db := l.appCtx.Database()
	tx := db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()
	var mute *string
	if req.Mute == nil {
		mute = nil
	} else if *req.Mute == 0 {
		sql := "mute & (mute ^ 2)"
		mute = &sql
	} else if *req.Mute == 1 {
		sql := "mute | 2"
		mute = &sql
	} else {
		return errorx.ErrParamsError
	}
	err = l.appCtx.SessionUserModel().UpdateUser(req.SId, req.UIds, req.Role, nil, mute, tx)
	return err
}

func (l *SessionLogic) convSessionUser(sessionUser *model.SessionUser) *dto.SessionUser {
	return &dto.SessionUser{
		SId:    sessionUser.SessionId,
		Type:   sessionUser.Type,
		Role:   sessionUser.Role,
		Mute:   sessionUser.Mute,
		Status: sessionUser.Status,
		CTime:  sessionUser.CreateTime,
		MTime:  sessionUser.UpdateTime,
	}
}
