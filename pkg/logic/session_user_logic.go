package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
)

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
	return l.appCtx.SessionUserModel().AddUser(session, entityIds, req.UIds, roles, maxCount)
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
			tx.Rollback()
		} else {
			tx.Commit()
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
