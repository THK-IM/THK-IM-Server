package logic

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"gorm.io/gorm"
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
	lockKey := fmt.Sprintf(sessionCreateLockKey, l.appCtx.Config().Name, req.UId, req.EntityId)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return nil, errorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()

	if req.Type == model.SingleSessionType {
		if len(req.Members) > 0 {
			return nil, errorx.ErrParamsError
		}
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.UId, req.EntityId, req.Type, true)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// 如果session已经存在
		if userSession.UserId > 0 {
			// 如果单聊被删除，需要恢复
			if userSession.Deleted == 1 {
				session, errSession := l.appCtx.SessionModel().FindSession(userSession.SessionId)
				if errSession != nil {
					return nil, errSession
				}
				userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(session,
					[]int64{userSession.EntityId}, []int64{userSession.UserId}, []int{model.SessionOwner}, 2)
				if errUserSessions != nil {
					return nil, errUserSessions
				}
				userSession = userSessions[0]
			}
			return &dto.CreateSessionRes{
				SId:      userSession.SessionId,
				Type:     userSession.Type,
				EntityId: userSession.EntityId,
				Name:     userSession.Name,
				Remark:   userSession.Remark,
				Role:     userSession.Role,
				Mute:     userSession.Mute,
				CTime:    userSession.CreateTime,
				MTime:    userSession.UpdateTime,
				Status:   userSession.Status,
				Top:      userSession.Top,
				IsNew:    false,
			}, nil
		}
	} else if req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType {
		userSession, err := l.appCtx.UserSessionModel().FindUserSessionByEntityId(req.UId, req.EntityId, req.Type, true)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// 如果session已经存在
		if userSession.UserId > 0 {
			// 群删除后不能恢复
			if userSession.Deleted == 1 {
				return nil, errorx.ErrGroupAlreadyDeleted
			} else {
				return nil, errorx.ErrGroupAlreadyCreated
			}
		}
	} else {
		return nil, errorx.ErrSessionType
	}
	return l.createNewSession(req)
}

func (l *SessionLogic) createNewSession(req dto.CreateSessionReq) (*dto.CreateSessionRes, error) {
	session, err := l.appCtx.SessionModel().CreateEmptySession(req.Type, req.ExtData, req.Name, req.Remark)
	if err != nil {
		return nil, err
	}
	var userSession *model.UserSession
	if req.Type == model.SingleSessionType {
		entityIds := []int64{req.EntityId, req.UId}
		uIds := []int64{req.UId, req.EntityId}
		roles := []int{model.SessionOwner, model.SessionOwner}
		if userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(session, entityIds, uIds, roles, 2); err != nil {
			return nil, errUserSessions
		} else {
			userSession = userSessions[0]
		}
	} else if req.Type == model.GroupSessionType || req.Type == model.SuperGroupSessionType {
		if req.EntityId <= 0 {
			err = errorx.ErrParamsError
			return nil, err
		}
		entityIds := make([]int64, 0)
		roles := make([]int, 0)
		// 插入自己的角色和entity_id
		entityIds = append(entityIds, req.EntityId)
		roles = append(roles, model.SessionOwner)
		// 插入群成员的角色和entity_id
		for range req.Members {
			entityIds = append(entityIds, req.EntityId)
			roles = append(roles, model.SessionMember)
		}
		maxMember := l.appCtx.Config().IM.MaxSuperGroupMember
		if req.Type == model.GroupSessionType {
			maxMember = l.appCtx.Config().IM.MaxGroupMember
		}
		if userSessions, errUserSessions := l.appCtx.SessionUserModel().AddUser(session, entityIds, req.Members, roles, maxMember); err != nil {
			return nil, errUserSessions
		} else {
			userSession = userSessions[0]
		}
	} else {
		err = errorx.ErrParamsError
		return nil, err
	}

	res := &dto.CreateSessionRes{
		SId:      userSession.SessionId,
		Type:     userSession.Type,
		EntityId: userSession.EntityId,
		Name:     userSession.Name,
		Remark:   userSession.Remark,
		Role:     userSession.Role,
		Mute:     userSession.Mute,
		Status:   userSession.Status,
		Top:      userSession.Top,
		CTime:    userSession.CreateTime,
		MTime:    userSession.UpdateTime,
		IsNew:    true,
	}
	return res, nil
}

func (l *SessionLogic) UpdateSession(req dto.UpdateSessionReq) error {
	lockKey := fmt.Sprintf(sessionUpdateLockKey, l.appCtx.Config().Name, req.Id)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return errorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()
	err := l.appCtx.SessionModel().UpdateSession(req.Id, req.Name, req.Remark, req.Mute, req.ExtData)
	if err != nil {
		return err
	}
	sessionUsers, errSessionUsers := l.appCtx.SessionUserModel().FindAllSessionUsers(req.Id)
	if errSessionUsers != nil {
		return errSessionUsers
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
	return l.appCtx.UserSessionModel().UpdateUserSession(uIds, req.Id, req.Name, req.Remark, mute, req.ExtData, nil, nil, nil, nil)
}

func (l *SessionLogic) UpdateUserSession(req dto.UpdateUserSessionReq) (err error) {
	lockKey := fmt.Sprintf(userSessionUpdateLockKey, l.appCtx.Config().Name, req.UId, req.SId)
	locker := l.appCtx.NewLocker(lockKey, 1000, 1000)
	success, lockErr := locker.Lock()
	if lockErr != nil || !success {
		return errorx.ErrServerBusy
	}
	defer func() {
		if success, lockErr = locker.Release(); lockErr != nil {
			l.appCtx.Logger().Errorf("release locker success: %t, error: %s", success, lockErr.Error())
		}
	}()
	err = l.appCtx.UserSessionModel().UpdateUserSession([]int64{req.UId}, req.SId, nil, nil, nil, nil, req.Top, req.Status, nil, req.ParentId)
	if err == nil {
		err = l.appCtx.SessionUserModel().UpdateUser(req.SId, []int64{req.UId}, nil, req.Status, nil)
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
		ExtData:  userSession.ExtData,
		CTime:    userSession.CreateTime,
		MTime:    userSession.UpdateTime,
	}
}
