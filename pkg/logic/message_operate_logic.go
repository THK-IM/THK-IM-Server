package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"math/rand"
	"time"
)

func (l *MessageLogic) AckUserMessages(req dto.AckUserMessagesReq) error {
	return l.appCtx.UserMessageModel().AckUserMessages(req.UId, req.SessionId, req.MessageIds)
}

func (l *MessageLogic) ReadUserMessages(req dto.ReadUserMessageReq) error {
	readStatus := model.MsgStatusRead | model.MsgStatusAcked
	if err := l.appCtx.UserMessageModel().UpdateUserMessage(req.UId, req.SessionId, req.MessageIds, readStatus); err != nil {
		l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %v, %v", req.UId, req.SessionId, req.MessageIds, err)
		return err
	} else {
		// 对消息发件人发送已读消息
		for _, msgId := range req.MessageIds {
			if userMessage, e := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SessionId, msgId); e == nil {
				if userMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不需要发送已读
					continue
				}
				sendMessageReq := dto.SendMessageReq{
					ClientId:  l.genClientId(),
					SessionId: req.SessionId,
					Type:      model.MsgTypeRead,
					FUid:      req.UId,
					CTime:     time.Now().UnixMilli(),
					RMsgId:    &userMessage.MsgId,
					Receivers: []int64{userMessage.FromUserId},
				}
				if _, e = l.SendMessage(sendMessageReq); e != nil {
					l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SessionId, msgId, err)
				}
			} else {
				l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SessionId, msgId, err)
			}
		}
		return nil
	}
}

func (l *MessageLogic) RevokeUserMessage(req dto.RevokeUserMessageReq) error {
	revokeStatus := model.MsgStatusRevoke | model.MsgStatusRead | model.MsgStatusAcked
	if err := l.appCtx.UserMessageModel().UpdateUserMessage(req.UId, req.SessionId, []int64{req.MessageId}, revokeStatus); err != nil {
		l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
		return err
	} else {
		// 发送撤回消息
		if userMessage, e := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SessionId, req.MessageId); e == nil {
			if userMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不需要发送已读
				return nil
			}
			sendMessageReq := dto.SendMessageReq{
				ClientId:  l.genClientId(),
				SessionId: req.SessionId,
				Type:      model.MsgTypeRevoke,
				FUid:      req.UId,
				CTime:     time.Now().UnixMilli(),
				RMsgId:    &userMessage.MsgId,
			}
			if _, e = l.SendMessage(sendMessageReq); e != nil {
				l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
			}
		} else {
			l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
		}
		return nil
	}
}

func (l *MessageLogic) genClientId() int64 {
	return time.Now().UnixMilli()*100 + rand.Int63n(100)
}
