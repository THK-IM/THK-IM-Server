package logic

import (
	"encoding/json"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"time"
)

type UserLogic struct {
	appCtx *app.Context
}

func NewUserLogic(appCtx *app.Context) UserLogic {
	return UserLogic{
		appCtx: appCtx,
	}
}

func (l *UserLogic) UpdateUserOnlineStatus(req *dto.PostUserOnlineReq) error {
	var isOnline int8 = 0
	if req.Online {
		isOnline = 1
	}
	err := l.appCtx.UserOnlineStatusModel().UpdateUserOnlineStatus(req.UId, isOnline)
	go func() {
		onlineReq := rpc.PostUserOnlineReq{
			UserId:   req.UId,
			IsOnline: req.Online,
		}
		userApi := l.appCtx.RpcUserApi()
		if userApi != nil {
			if e := userApi.PostUserOnlineStatus(onlineReq); e != nil {
				l.appCtx.Logger().Errorf("UpdateUserOnlineStatus, RpcUserApi, call err: %v", e)
			}
		}

	}()
	return err
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64) (*dto.GetUsersOnlineStatusRes, error) {
	usersOnlineStatus, err := l.appCtx.UserOnlineStatusModel().GetUsersOnlineStatus(uIds)
	onlineTimeout := l.appCtx.Config().IM.OnlineTimeout
	if err != nil {
		return nil, err
	} else {
		now := time.Now().UnixMilli()
		dtoUsersOnlineStatus := make([]*dto.UserOnlineStatus, 0)
		for _, user := range usersOnlineStatus {
			online := false
			if user.IsOnline > 0 && (now-user.OnlineTime) < onlineTimeout*int64(time.Second) {
				online = true
			} else {
				online = false
			}
			dtoUserOnlineStatus := &dto.UserOnlineStatus{
				UId:            user.UserId,
				Online:         online,
				LastOnlineTime: user.OnlineTime,
			}
			dtoUsersOnlineStatus = append(dtoUsersOnlineStatus, dtoUserOnlineStatus)
		}
		return &dto.GetUsersOnlineStatusRes{UsersOnlineStatus: dtoUsersOnlineStatus}, nil
	}
}

func (l *UserLogic) KickUser(req *dto.KickUserReq) error {
	ids := []int64{req.UId}
	if idsStr, err := json.Marshal(ids); err != nil {
		return err
	} else {
		msg := make(map[string]interface{})
		msg[event.PushEventTypeKey] = event.PushUserEventType
		msg[event.PushEventSubTypeKey] = event.UserEventSubtypeKickOff
		msg[event.PushEventReceiversKey] = string(idsStr)
		msg[event.PushEventBodyKey] = "kickOff"
		return l.appCtx.MsgPusherPublisher().Pub("", msg)
	}
}
