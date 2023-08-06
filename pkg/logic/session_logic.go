package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/gin-gonic/gin"
)

type SessionLogic struct {
	ctx    *gin.Context
	appCtx *app.Context
}

func NewSessionLogic(ctx *gin.Context, appCtx *app.Context) SessionLogic {
	return SessionLogic{
		ctx:    ctx,
		appCtx: appCtx,
	}
}

func (l *SessionLogic) CreateSession(req dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	if req.Type == model.SingleSessionType {
		if len(req.Members) != 2 {
			return nil, errorx.ErrParamsError
		}
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.Members[0], req.Members[1], req.Type)
		if err == nil && userSession.UserId > 0 {
			return &dto.CreateSessionRes{
				SessionId: userSession.SessionId,
				Type:      userSession.Type,
				EntityId:  userSession.EntityId,
				CTime:     userSession.CreateTime,
				MTime:     userSession.UpdateTime,
				Status:    userSession.Status,
				Top:       userSession.Top,
			}, nil
		}
	} else {
		if len(req.Members) < 1 || req.EntityId == 0 {
			return nil, errorx.ErrParamsError
		}
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.Members[0], req.Members[1], req.Type)
		if err == nil && userSession.UserId > 0 {
			if len(req.Members) > 1 {
				session := &model.Session{
					Id:   userSession.SessionId,
					Type: userSession.Type,
				}
				members := req.Members[1:]
				if err = l.appCtx.SessionUserModel().AddUser(session, req.EntityId, members); err != nil {
					return nil, err
				}
			}
			return &dto.CreateSessionRes{
				SessionId: userSession.SessionId,
				Type:      userSession.Type,
				EntityId:  userSession.EntityId,
				CTime:     userSession.CreateTime,
				MTime:     userSession.UpdateTime,
				Status:    userSession.Status,
				Top:       userSession.Top,
			}, nil
		}
	}
	return l.createNewSession(req)
}

func (l *SessionLogic) createNewSession(req dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	session, err := l.appCtx.SessionModel().CreateEmptySession(req.Type, nil)
	if err != nil {
		return nil, err
	}
	entityId := req.EntityId
	if req.Type == model.SingleSessionType {
		entityId = req.Members[1]
		entityIds := []int64{req.Members[1], req.Members[0]}
		for i, member := range req.Members {
			if err = l.appCtx.SessionUserModel().AddUser(session, entityIds[i], []int64{member}); err != nil {
				return nil, err
			}
		}
	} else if req.Type == model.GroupSessionType {
		if req.EntityId <= 0 {
			return nil, errorx.ErrParamsError
		}
		if err = l.appCtx.SessionUserModel().AddUser(session, req.EntityId, req.Members); err != nil {
			return nil, err
		}
	} else {
		return nil, errorx.ErrParamsError
	}

	return &dto.CreateSessionRes{
		SessionId: session.Id,
		Type:      req.Type,
		EntityId:  entityId,
		CTime:     session.CreateTime,
		MTime:     session.UpdateTime,
	}, nil
}

func (l *SessionLogic) UpdateSession(req dto.UpdateSessionReq) error {
	return l.appCtx.SessionModel().UpdateSession(req.Id, req.Status, req.Name, req.Remark)
}

func (l *SessionLogic) UpdateUserSession(req dto.UpdateUserSessionReq) error {
	db := l.appCtx.Database()
	var err error
	tx := db.Begin()
	if err = l.appCtx.UserSessionModel().UpdateUserSession(req.UId, req.SId, req.Top, req.Status, tx); err == nil {
		err = l.appCtx.SessionUserModel().UpdateUser(req.SId, req.UId, *req.Status, tx)
	}
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func (l *SessionLogic) GetUserSessions(req dto.GetUserSessionsReq) (*dto.GetUserSessionsRes, error) {
	userSessions, err := l.appCtx.UserSessionModel().GetUserSessions(req.UId, req.MTime, req.Offset, req.Count)
	if err != nil {
		return nil, err
	}
	dtoUserSessions := make([]*dto.UserSession, 0)
	for _, userSession := range userSessions {
		userSession := l.convUserSession(userSession)
		dtoUserSessions = append(dtoUserSessions, userSession)
	}
	return &dto.GetUserSessionsRes{Data: dtoUserSessions}, nil
}

func (l *SessionLogic) GetUserSession(uid, sid int64) (*dto.UserSession, error) {
	userSession, err := l.appCtx.UserSessionModel().GetUserSession(uid, sid)
	if err != nil {
		return nil, err
	}
	dtoUserSession := l.convUserSession(userSession)
	return dtoUserSession, nil
}

func (l *SessionLogic) convUserSession(userSession *model.UserSession) *dto.UserSession {
	return &dto.UserSession{
		SessionId: userSession.SessionId,
		Type:      userSession.Type,
		Top:       userSession.Top,
		Status:    userSession.Status,
		EntityId:  userSession.EntityId,
		CTime:     userSession.CreateTime,
		MTime:     userSession.UpdateTime,
	}
}
