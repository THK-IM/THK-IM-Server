package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"time"
)

func (l *MessageLogic) AckUserMessages(req dto.AckUserMessagesReq) error {
	return l.appCtx.UserMessageModel().AckUserMessages(req.UId, req.SessionId, req.MessageIds)
}

func (l *MessageLogic) ReadUserMessages(req dto.ReadUserMessageReq) error {
	// 对消息发件人发送已读消息
	for _, msgId := range req.MessageIds {
		if userMessage, err := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SessionId, msgId); err == nil {
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
				Receivers: []int64{userMessage.FromUserId, req.UId}, // 发送给对方和自己
			}
			if _, err = l.SendMessage(sendMessageReq); err != nil {
				l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SessionId, msgId, err)
			}
		} else {
			l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SessionId, msgId, err)
		}
	}
	return nil
}

func (l *MessageLogic) RevokeUserMessage(req dto.RevokeUserMessageReq) error {
	if sessionMessage, err := l.appCtx.SessionMessageModel().FindSessionMessage(req.SessionId, req.MessageId, req.UId); err == nil {
		if sessionMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能发送撤回
			return errorx.ErrMessageTypeNotSupportOperated
		}
		// 删除session的消息
		_, err = l.appCtx.SessionMessageModel().UpdateSessionMessageStatus(sessionMessage.SessionId, sessionMessage.MsgId, sessionMessage.FromUserId, model.MsgStatusRevoke)
		if err != nil {
			return err
		}
		sendMessageReq := dto.SendMessageReq{
			ClientId:  l.genClientId(),
			SessionId: req.SessionId,
			Type:      model.MsgTypeRevoke,
			FUid:      req.UId,
			CTime:     time.Now().UnixMilli(),
			RMsgId:    &req.MessageId,
		} // 发送给session下的所有人
		if _, err = l.SendMessage(sendMessageReq); err != nil {
			l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
		}
	} else {
		l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
	}
	return nil
}

func (l *MessageLogic) ReeditUserMessage(req dto.ReeditUserMessageReq) error {
	if sessionMessage, err := l.appCtx.SessionMessageModel().FindSessionMessage(req.UId, req.SessionId, req.MessageId); err == nil {
		if sessionMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能发送撤回
			return errorx.ErrMessageTypeNotSupportOperated
		}
		// 更新session的消息内容
		_, err = l.appCtx.SessionMessageModel().UpdateSessionMessageContent(sessionMessage.SessionId, sessionMessage.MsgId, sessionMessage.FromUserId, req.Content, model.MsgStatusReedit)
		if err != nil {
			return err
		}
		sendMessageReq := dto.SendMessageReq{
			ClientId:  l.genClientId(),
			SessionId: req.SessionId,
			Type:      model.MsgTypeReedit,
			FUid:      req.UId,
			CTime:     time.Now().UnixMilli(),
			Body:      req.Content,
			RMsgId:    &req.MessageId,
		} // 发送给session下的所有人
		if _, err = l.SendMessage(sendMessageReq); err != nil {
			l.appCtx.Logger().Errorf("ReeditUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
		}
	} else {
		l.appCtx.Logger().Errorf("ReeditUserMessage err:%d, %d, %d, %v", req.UId, req.SessionId, req.MessageId, err)
	}
	return nil
}

func (l *MessageLogic) genClientId() int64 {
	return l.appCtx.SnowflakeNode().Generate().Int64()
}
