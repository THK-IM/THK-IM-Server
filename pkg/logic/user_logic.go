package logic

import (
	"encoding/json"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
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
	onlineTime := req.Timestamp
	if !req.Online {
		onlineTime = 0
	}
	err := l.appCtx.UserOnlineStatusModel().UpdateUserOnlineStatus(req.UId, onlineTime, req.ConnId, req.Platform)
	go func() {
		onlineReq := rpc.PostUserOnlineReq{
			UserId:    req.UId,
			IsOnline:  req.Online,
			Timestamp: req.Timestamp,
			ConnId:    req.ConnId,
			Platform:  req.Platform,
		}
		if l.appCtx.RpcUserApi() != nil {
			if e := l.appCtx.RpcUserApi().PostUserOnlineStatus(onlineReq); e != nil {
				l.appCtx.Logger().Errorf("UpdateUserOnlineStatus, RpcUserApi, call err: %s", e.Error())
			}
		}
	}()
	return err
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64) (*dto.GetUsersOnlineStatusRes, error) {
	usersOnlineStatus, err := l.appCtx.UserOnlineStatusModel().GetUsersOnlineStatus(uIds)
	if err != nil {
		return nil, err
	} else {
		dtoUsersOnlineStatus := make([]*dto.UserOnlineStatus, 0)
		for _, user := range usersOnlineStatus {
			dtoUserOnlineStatus := &dto.UserOnlineStatus{
				UId:            user.UserId,
				Platform:       user.Platform,
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
