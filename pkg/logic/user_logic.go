package logic

import (
	"encoding/json"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/gin-gonic/gin"
	"time"
)

type UserLogic struct {
	ctx    *gin.Context
	appCtx *app.Context
}

func NewUserLogic(ctx *gin.Context, appCtx *app.Context) UserLogic {
	return UserLogic{
		ctx:    ctx,
		appCtx: appCtx,
	}
}

func (l *UserLogic) UpdateUserOnlineStatus(req *dto.PostUserOnlineReq) error {
	var isOnline int8 = 0
	if req.Online {
		isOnline = 1
	}
	return l.appCtx.UserOnlineStatusModel().UpdateUserOnlineStatus(req.UserId, isOnline)
}

func (l *UserLogic) GetUsersOnlineStatus(uIds []int64) (*dto.GetUsersOnlineStatusRes, error) {
	usersOnlineStatus, err := l.appCtx.UserOnlineStatusModel().GetUsersOnlineStatus(uIds)
	offlineTime := l.appCtx.Config().OfflineInterval
	if err != nil {
		return nil, err
	} else {
		now := time.Now().UnixMilli()
		dtoUsersOnlineStatus := make([]*dto.UserOnlineStatus, 0)
		for _, user := range usersOnlineStatus {
			online := false
			if user.IsOnline > 0 && (now-user.OnlineTime) < offlineTime*int64(time.Second) {
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
		msg[event.PushEventTypeKey] = event.PushUserEvent
		msg[event.PushEventSubTypeKey] = 1
		msg[event.PushEventReceiversKey] = string(idsStr)
		msg[event.PushEventBodyKey] = "kickOff"
		return l.appCtx.MsgPusherPublisher().Pub("", msg)
	}
}
