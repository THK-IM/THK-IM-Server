package handler

import (
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"github.com/THK-IM/THK-IM-Server/pkg/service/websocket"
	"strconv"
	"time"
)

func RegisterMsgPushHandlers(ctx *app.Context) {
	server := ctx.WebsocketServer()
	server.SetUidGetter(func(token string, pf string) (int64, error) {
		if ctx.Config().Mode == "debug" {
			if uId, err := strconv.Atoi(token); err != nil {
				return 0, err
			} else {
				return int64(uId), nil
			}
		} else {
			req := rpc.GetUserIdByTokenReq{Token: token, Platform: pf}
			if res, err := ctx.RpcUserApi().GetUserIdByToken(req); err != nil {
				return 0, err
			} else {
				return res.UserId, nil
			}
		}
	})
	server.SetOnClientConnected(func(client websocket.Client) {
		ctx.Logger().Infof("OnClientConnected: %v", client.Info())
		{
			// 下发服务器时间
			now := time.Now().UnixMilli()
			serverTimeBody, err := event.BuildSignalBody(event.SignalSyncTime, strconv.Itoa(int(now)))
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
			err = client.WriteMessage(serverTimeBody)
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
		}
		{
			// 下发连接id
			connIdBody, err := event.BuildSignalBody(event.SignalConnId, fmt.Sprintf("%d", client.Info().Id))
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
			err = client.WriteMessage(connIdBody)
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
		}
		// 发送用户上线事件
		{
			if userOnlineEvent, err := event.BuildUserOnlineEvent(ctx.NodeId(), true,
				client.Info().UId, client.Info().Id, client.Info().FirstOnLineTime, client.Info().Platform); err != nil {
				ctx.Logger().Error("UserOnlineEvent Build err:", err)
			} else {
				if err = ctx.ServerEventPublisher().Pub(fmt.Sprintf("uid-%d", client.Info().UId), userOnlineEvent); err != nil {
					ctx.Logger().Error("UserOnlineEvent Pub err:", err)
				}
			}
		}
		// rpc通知api服务用户上线
		{
			sendUserOnlineStatus(ctx, client, true)
		}
	})

	server.SetOnClientClosed(func(client websocket.Client) {
		ctx.Logger().Infof("OnClientClosed: %v", client.Info())
		sendUserOnlineStatus(ctx, client, false)
	})

	server.SetOnClientMsgReceived(func(msg string, client websocket.Client) {
		signal := &event.SignalBody{}
		if err := json.Unmarshal([]byte(msg), signal); err != nil {
			ctx.Logger().Errorf("json Unmarshal err: %s, msg: %s", err.Error(), msg)
		} else {
			err = onWsClientMsgReceived(ctx, client, signal.Type, signal.Body)
		}
	})

	ctx.MsgPusherSubscriber().Sub(func(m map[string]interface{}) error {
		return onMqPushMsgReceived(m, server, ctx)
	})

	ctx.ServerEventSubscriber().Sub(func(m map[string]interface{}) error {
		return onMqServerEventReceived(m, server, ctx)
	})
}

func onMqPushMsgReceived(m map[string]interface{}, server websocket.Server, ctx *app.Context) error {
	ctx.Logger().Info("onMqPushMsgReceived", m)
	eventType, okType := m[event.PushEventTypeKey].(string)
	uIdsStr, okId := m[event.PushEventReceiversKey].(string)
	body, okBody := m[event.PushEventBodyKey].(string)
	if !okType || !okId || !okBody {
		return errorx.ErrMessageFormat
	}
	iType, eType := strconv.Atoi(eventType)
	if eType != nil {
		return errorx.ErrMessageFormat
	}
	uIds := make([]int64, 0)
	if err := json.Unmarshal([]byte(uIdsStr), &uIds); err != nil {
		return errorx.ErrMessageFormat
	}
	if content, err := event.BuildSignalBody(iType, body); err != nil {
		return err
	} else {
		return server.SendMessageToUsers(uIds, content)
	}
}

func onWsClientMsgReceived(ctx *app.Context, client websocket.Client, ty int, body *string) error {
	if ty == event.SignalHeatBeat {
		return onWsHeatBeatMsgReceived(ctx, client, body)
	}
	return nil
}

func onWsHeatBeatMsgReceived(ctx *app.Context, client websocket.Client, body *string) error {
	// 心跳
	heatBody, err := event.BuildSignalBody(event.SignalHeatBeat, "pong")
	if err != nil {
		return err
	}
	ctx.Logger().Info(client.Info())
	sendUserOnlineStatus(ctx, client, true)
	return client.WriteMessage(heatBody)
}

func sendUserOnlineStatus(ctx *app.Context, client websocket.Client, online bool) {
	now := time.Now().UnixMilli()
	client.SetLastOnlineTime(now)
	req := dto.PostUserOnlineReq{
		NodeId:    ctx.NodeId(),
		ConnId:    client.Info().Id,
		Online:    online,
		UId:       client.Info().UId,
		Platform:  client.Info().Platform,
		Timestamp: time.Now().UnixMilli(),
	}
	if err := ctx.RpcMsgApi().PostUserOnlineStatus(req); err != nil {
		ctx.Logger().Errorf("sendUserOnlineStatus, err: %v", err)
	}
}

func onMqServerEventReceived(m map[string]interface{}, server websocket.Server, ctx *app.Context) error {
	tp, okType := m[event.ServerEventTypeKey].(string)
	receivers, okReceivers := m[event.ServerEventReceiversKey].(string)
	body, okBody := m[event.ServerEventBodyKey].(string)
	if !okType || !okReceivers || !okBody {
		return errorx.ErrMessageFormat
	}
	uIds := make([]int64, 0)
	if err := json.Unmarshal([]byte(receivers), &uIds); err != nil {
		return errorx.ErrMessageFormat
	}
	if tp == event.ServerEventUserOnline {
		onlineBody := event.ParserOnlineBody(body)
		if onlineBody != nil {
			if e := server.OnUserOnLineEvent(*onlineBody); e != nil {
				ctx.Logger().Error("OnUserOnLineEvent, err:", e, " onlineBody: ", onlineBody)
			}
		} else {
			ctx.Logger().Error("ServerEventUserOnline, onlineBody is nil, body is: ", body)
		}
	}
	return nil
}
