package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"gorm.io/gorm"
	"time"
)

type SessionLogic struct {
	appCtx *app.Context
}

func NewSessionLogic(appCtx *app.Context) SessionLogic {
	return SessionLogic{
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
		// 如果session已经存在
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
				SId:      userSession.SessionId,
				Type:     userSession.Type,
				EntityId: userSession.EntityId,
				CTime:    userSession.CreateTime,
				MTime:    userSession.UpdateTime,
				Status:   userSession.Status,
				Top:      userSession.Top,
				Success:  false,
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
		// 如果session已经存在
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
				maxCount := l.appCtx.Config().IM.MaxSuperGroupMember
				if userSession.Type == model.GroupSessionType {
					maxCount = l.appCtx.Config().IM.MaxGroupMember
				}
				entityIds := make([]int64, 0)
				role := make([]int, 0)
				for range members {
					entityIds = append(entityIds, *req.EntityId)
					role = append(role, model.SessionMember)
				}
				if err = l.appCtx.SessionUserModel().AddUser(session, entityIds, members, role, maxCount); err != nil {
					return nil, err
				}
			}
			return &dto.CreateSessionRes{
				SId:      userSession.SessionId,
				Type:     userSession.Type,
				EntityId: userSession.EntityId,
				CTime:    userSession.CreateTime,
				MTime:    userSession.UpdateTime,
				Status:   userSession.Status,
				Top:      userSession.Top,
				Success:  false,
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
		uIds := []int64{req.Members[0], req.Members[1]}
		roles := []int{model.SessionOwner, model.SessionOwner}
		if err = l.appCtx.SessionUserModel().AddUser(session, entityIds, uIds, roles, 2); err != nil {
			return nil, err
		}
	} else if req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType {
		if req.EntityId == nil {
			return nil, errorx.ErrParamsError
		}
		entityIds := make([]int64, 0)
		roles := make([]int, 0)
		for range req.Members {
			entityIds = append(entityIds, *req.EntityId)
			roles = append(roles, model.SessionMember)
		}
		roles[0] = model.SessionOwner
		maxMember := l.appCtx.Config().IM.MaxSuperGroupMember
		if req.Type == model.GroupSessionType {
			maxMember = l.appCtx.Config().IM.MaxGroupMember
		}
		if err = l.appCtx.SessionUserModel().AddUser(session, entityIds, req.Members, roles, maxMember); err != nil {
			return nil, err
		}
	} else {
		return nil, errorx.ErrParamsError
	}

	return &dto.CreateSessionRes{
		SId:      session.Id,
		Type:     req.Type,
		EntityId: entityId,
		CTime:    session.CreateTime,
		MTime:    session.UpdateTime,
		Success:  true,
	}, nil
}

func (l *SessionLogic) UpdateSession(req dto.UpdateSessionReq) (err error) {
	tx := l.appCtx.Database().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	err = l.appCtx.SessionModel().UpdateSession(req.Id, req.Name, req.Remark, req.Mute, tx)
	if err != nil {
		return
	}
	sessionUsers, errSessionUsers := l.appCtx.SessionUserModel().FindAllSessionUsers(req.Id)
	if errSessionUsers != nil {
		err = errSessionUsers
		return
	}
	uIds := make([]int64, 0)
	for _, su := range sessionUsers {
		uIds = append(uIds, su.UserId)
	}
	var mute *string
	if req.Mute == nil {
		mute = nil
	} else if *req.Mute == 0 {
		sql := "mute & (mute ^ 1)"
		mute = &sql
	} else if *req.Mute == 1 {
		sql := "mute | 1"
		mute = &sql
	}
	err = l.appCtx.UserSessionModel().UpdateUserSession(uIds, req.Id, req.Name, req.Remark, mute, nil, nil, nil, tx)
	return
}

func (l *SessionLogic) UpdateUserSession(req dto.UpdateUserSessionReq) (err error) {
	tx := l.appCtx.Database().Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	err = l.appCtx.UserSessionModel().UpdateUserSession([]int64{req.UId}, req.SId, nil, nil, nil, req.Top, req.Status, nil, tx)
	if err == nil {
		err = l.appCtx.SessionUserModel().UpdateUser(req.SId, []int64{req.UId}, nil, req.Status, nil, tx)
	} else {
		l.appCtx.Logger().Error("UpdateUserSession, err", err)
	}
	return
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
		SId:      userSession.SessionId,
		Type:     userSession.Type,
		Name:     userSession.Name,
		Remark:   userSession.Remark,
		Role:     userSession.Role,
		Mute:     userSession.Mute,
		Top:      userSession.Top,
		Status:   userSession.Status,
		EntityId: userSession.EntityId,
		CTime:    userSession.CreateTime,
		MTime:    userSession.UpdateTime,
	}
}
