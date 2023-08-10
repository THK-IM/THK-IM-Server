package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
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
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.Members[0], req.Members[1], req.Type, true)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if userSession.UserId > 0 {
			if userSession.Deleted == 1 {
				now := time.Now().UnixMilli()
				if err = l.appCtx.UserSessionModel().RecoverUserSession(userSession.UserId, userSession.SessionId, now); err != nil {
					return nil, err
				} else {
					userSession.CreateTime = now
					userSession.UpdateTime = now
					userSession.Status = 0
					userSession.Top = 0
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
	} else if req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType {
		if len(req.Members) < 1 || req.EntityId == nil {
			return nil, errorx.ErrParamsError
		}
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.Members[0], req.Members[1], req.Type, true)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if userSession.UserId > 0 {
			if userSession.Deleted == 1 {
				return nil, errorx.ErrGroupAlreadyDeleted
			}
			if len(req.Members) > 1 {
				session := &model.Session{
					Id:   userSession.SessionId,
					Type: userSession.Type,
				}
				members := req.Members[1:]
				maxCount := 0
				if userSession.Type == model.GroupSessionType {
					maxCount = l.appCtx.Config().IM.MaxGroupMember
				} else if userSession.Type == model.SuperGroupSessionType {
					maxCount = l.appCtx.Config().IM.MaxSuperGroupMember
				}
				if err = l.appCtx.SessionUserModel().AddUser(session, *req.EntityId, members, maxCount); err != nil {
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
	} else {
		return nil, errorx.ErrSessionType
	}
	return l.createNewSession(req)
}

func (l *SessionLogic) createNewSession(req dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	session, err := l.appCtx.SessionModel().CreateEmptySession(req.Type, nil)
	if err != nil {
		return nil, err
	}
	entityId := int64(0)
	if req.Type == model.SingleSessionType {
		entityId = req.Members[1]
		entityIds := []int64{req.Members[1], req.Members[0]}
		for i, member := range req.Members {
			if err = l.appCtx.SessionUserModel().AddUser(session, entityIds[i], []int64{member}, 2); err != nil {
				return nil, err
			}
		}
	} else if req.Type == model.GroupSessionType {
		if req.EntityId == nil {
			return nil, errorx.ErrParamsError
		}
		entityId = *req.EntityId
		if err = l.appCtx.SessionUserModel().AddUser(session, entityId, req.Members, l.appCtx.Config().IM.MaxGroupMember); err != nil {
			return nil, err
		}
	} else if req.Type == model.SuperGroupSessionType {
		if req.EntityId == nil {
			return nil, errorx.ErrParamsError
		}
		entityId = *req.EntityId
		if err = l.appCtx.SessionUserModel().AddUser(session, entityId, req.Members, l.appCtx.Config().IM.MaxSuperGroupMember); err != nil {
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
		dtoUserSession := l.convUserSession(userSession)
		dtoUserSessions = append(dtoUserSessions, dtoUserSession)
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
