package handler

import (
	"encoding/json"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
)

func RegisterSaveMsgHandlers(appCtx *app.Context) {
	appCtx.MsgSaverSubscriber().Sub(func(m map[string]interface{}) error {
		return onMqSaveMsgEventReceived(m, appCtx)
	})
}

func onMqSaveMsgEventReceived(m map[string]interface{}, appCtx *app.Context) error {
	msgJsonStr, okMsg := m[event.SaveMsgEventKey].(string)
	receiversStr, okReceiver := m[event.SaveMsgUsersKey].(string)
	if okMsg && okReceiver {
		message := &dto.Message{}
		err := json.Unmarshal([]byte(msgJsonStr), message)
		if err != nil {
			return errorx.ErrMessageFormat
		}
		receivers := make([]int64, 0)
		err = json.Unmarshal([]byte(receiversStr), &receivers)
		for _, r := range receivers {
			status := 0
			if r == message.FUid {
				status = model.MsgStatusAcked | model.MsgStatusRead
			}
			userMessage := &model.UserMessage{
				MsgId:      message.MsgId,
				ClientId:   message.CId,
				UserId:     r,
				SessionId:  message.SId,
				FromUserId: message.FUid,
				AtUsers:    message.AtUsers,
				MsgType:    message.Type,
				MsgContent: message.Body,
				ReplyMsgId: message.RMsgId,
				Status:     status,
				CreateTime: message.CTime,
				UpdateTime: message.CTime,
				Deleted:    0,
			}
			err = appCtx.UserMessageModel().InsertUserMessage(userMessage)
			if err != nil {
				return errorx.ErrMessageFormat
			}
			if message.Type == model.MsgTypeRead && message.RMsgId != nil && r == message.FUid { // 发送已读的人将自己的消息标记为已读
				err = appCtx.UserMessageModel().UpdateUserMessage(r, message.SId, []int64{*message.RMsgId}, model.MsgStatusRead, nil)
				if err != nil {
					return err
				}
			}
			if message.Type == model.MsgTypeRevoke && message.RMsgId != nil {
				err = appCtx.UserMessageModel().UpdateUserMessage(r, message.SId, []int64{*message.RMsgId}, model.MsgStatusRevoke, nil)
				if err != nil {
					return err
				}
			}
			if message.Type == model.MsgTypeReedit && message.RMsgId != nil {
				err = appCtx.UserMessageModel().UpdateUserMessage(r, message.SId, []int64{*message.RMsgId}, model.MsgStatusReedit, &message.Body)
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else {
		return errorx.ErrMessageFormat
	}
}
