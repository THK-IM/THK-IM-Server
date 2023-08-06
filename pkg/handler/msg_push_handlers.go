package handler

import (
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"github.com/THK-IM/THK-IM-Server/pkg/websocket"
	"strconv"
	"time"
)

func RegisterMsgPushHandlers(ctx *app.Context) {
	server := ctx.WebsocketServer()
	server.SetUidGetter(func(token string, pf string) (int64, error) {
		req := rpc.GetUserIdByTokenReq{Token: token, Platform: pf}
		if res, err := ctx.RpcUserApi().GetUserIdByToken(req); err != nil {
			return 0, err
		} else {
			return res.UserId, nil
		}
	})
	server.SetOnClientConnected(func(client websocket.Client) {
		ctx.Logger().Infof("OnClientConnected: %v", client.Info())
		{
			// 下发服务器时间
			now := time.Now().UnixMilli()
			serverTimeBody := &event.PushBody{
				Type:    event.PushCommonEventType,
				SubType: 3,
				Body:    strconv.Itoa(int(now)),
			}
			serverTimeBodyBytes, err := json.Marshal(serverTimeBody)
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
			err = client.WriteMessage(string(serverTimeBodyBytes))
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
		}
		{
			// 下发连接id
			connIdBody := &event.PushBody{
				Type:    event.PushCommonEventType,
				SubType: 4,
				Body:    fmt.Sprintf("%d", client.Info().Id),
			}
			connIdBodyBytes, err := json.Marshal(connIdBody)
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
			err = client.WriteMessage(string(connIdBodyBytes))
			if err != nil {
				ctx.Logger().Errorf("OnClientConnected: %s", err.Error())
			}
		}
		{
			sendUserOnlineStatus(ctx, client, true)
		}
	})

	server.SetOnClientClosed(func(client websocket.Client) {
		ctx.Logger().Infof("OnClientClosed: %v", client.Info())
		sendUserOnlineStatus(ctx, client, false)
	})

	server.SetOnClientMsgReceived(func(msg string, client websocket.Client) {
		wsBody := &event.PushBody{}
		if err := json.Unmarshal([]byte(msg), wsBody); err != nil {
			ctx.Logger().Errorf("json Unmarshal err: %s, msg: %s", err.Error(), msg)
		} else {
			err = onWsClientMsgReceived(ctx, client, wsBody.Type, wsBody.SubType, wsBody.Body)
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
	ctx.Logger().Info(m)
	eventType, okType := m[event.PushEventTypeKey].(string)
	eventSubtype, okSubType := m[event.PushEventSubTypeKey].(string)
	uIdsStr, okId := m[event.PushEventReceiversKey].(string)
	body, okBody := m[event.PushEventBodyKey].(string)
	if !okType || !okSubType || !okId || !okBody {
		return errorx.ErrMsgBodyContent
	}
	iType, eType := strconv.Atoi(eventType)
	if eType != nil {
		return errorx.ErrMsgBodyContent
	}
	iSubtype, eSubType := strconv.Atoi(eventSubtype)
	if eSubType != nil {
		return errorx.ErrMsgBodyContent
	}
	uIds := make([]int64, 0)
	if err := json.Unmarshal([]byte(uIdsStr), &uIds); err != nil {
		return errorx.ErrMsgBodyContent
	}
	if content, err := event.BuildPushBody(iType, iSubtype, body); err != nil {
		return err
	} else {
		return server.SendMessageToUsers(uIds, content)
	}
}

func onWsClientMsgReceived(ctx *app.Context, client websocket.Client, ty, subType int, body string) error {
	if ty == event.PushCommonEventType {
		if subType == 1 {
			return onWsHeatBeatMsgReceived(ctx, client, body)
		} else {
			return nil
		}
	}
	return nil
}

func onWsHeatBeatMsgReceived(ctx *app.Context, client websocket.Client, body string) error {
	// 心跳
	heatBody := &event.PushBody{
		Type:    event.PushCommonEventType,
		SubType: 2,
		Body:    "pong",
	}
	heatBeat, err := json.Marshal(heatBody)
	if err != nil {
		return err
	}
	ctx.Logger().Info(client.Info())
	sendUserOnlineStatus(ctx, client, true)
	return client.WriteMessage(string(heatBeat))
}

func sendUserOnlineStatus(ctx *app.Context, client websocket.Client, online bool) {
	now := time.Now().UnixMilli()
	client.SetLastOnlineTime(now)
	req := dto.PostUserOnlineReq{
		NodeId: ctx.NodeId(),
		ConnId: client.Info().Id,
		Online: online,
		UId:    client.Info().UId,
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
		return errorx.ErrMsgBodyContent
	}
	uIds := make([]int64, 0)
	if err := json.Unmarshal([]byte(receivers), &uIds); err != nil {
		return errorx.ErrMsgBodyContent
	}
	if tp == event.ServerEventUserOnline {
		onlineBody := event.ParserOnlineBody(body)
		if onlineBody != nil {
			if e := server.KickOffClient(onlineBody.UserId, onlineBody.ConnId, onlineBody.Platform); e != nil {
				ctx.Logger().Error("KickOffClient, err:", e, " onlineBody: ", onlineBody)
			}
		} else {
			ctx.Logger().Error("ServerEventUserOnline, onlineBody is nil, body is: ", body)
		}
	}
	return nil
}
