package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
)

func (l *SessionLogic) AddUser(sid int64, req dto.SessionAddUserReq) error {
	session, err := l.appCtx.SessionModel().FindSession(sid, nil)
	if err != nil {
		return err
	}
	return l.appCtx.SessionUserModel().AddUser(session, req.EntityId, req.UIds)
}

func (l *SessionLogic) DelUser(sid int64, req dto.SessionDelUserReq) error {
	session, err := l.appCtx.SessionModel().FindSession(sid, nil)
	if err != nil {
		return err
	}
	return l.appCtx.SessionUserModel().DelUser(session, req.UIds)
}

func (l *SessionLogic) UpdateSessionUser(req dto.SessionUserUpdateReq) error {
	db := l.appCtx.Database()
	var err error
	tx := db.Begin()
	if err = l.appCtx.SessionUserModel().UpdateUser(req.SId, req.UId, *req.Status, tx); err == nil {
		err = l.appCtx.UserSessionModel().UpdateUserSession(req.UId, req.SId, nil, req.Status, tx)
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}
