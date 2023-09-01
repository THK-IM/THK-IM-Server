package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"time"
)

func (l *MessageLogic) AckUserMessages(req dto.AckUserMessagesReq) error {
	return l.appCtx.UserMessageModel().AckUserMessages(req.UId, req.SId, req.MsgIds)
}

func (l *MessageLogic) ReadUserMessages(req dto.ReadUserMessageReq) error {
	// 对消息发件人发送已读消息
	for _, msgId := range req.MsgIds {
		if userMessage, err := l.appCtx.UserMessageModel().FindUserMessage(req.UId, req.SId, msgId); err == nil {
			if userMessage.MsgType < 0 || userMessage.Status&model.MsgStatusRead == 1 { // 小于0的类型消息为状态操作消息或者已经是已读了，不需要发送已读
				continue
			}
			sendMessageReq := dto.SendMessageReq{
				CId:       l.genClientId(),
				SId:       req.SId,
				Type:      model.MsgTypeRead,
				FUid:      req.UId,
				CTime:     time.Now().UnixMilli(),
				RMsgId:    &userMessage.MsgId,
				Receivers: []int64{userMessage.FromUserId, req.UId}, // 发送给对方和自己
			}
			if _, err = l.SendMessage(sendMessageReq); err != nil {
				l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SId, msgId, err)
			}
		} else {
			l.appCtx.Logger().Errorf("ReadUserMessages err:%d, %d, %d, %v", req.UId, req.SId, msgId, err)
		}
	}
	return nil
}

func (l *MessageLogic) RevokeUserMessage(req dto.RevokeUserMessageReq) error {
	if sessionMessage, err := l.appCtx.SessionMessageModel().FindSessionMessage(req.SId, req.MsgId, req.UId); err == nil {
		if sessionMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能发送撤回
			return errorx.ErrMessageTypeNotSupport
		}
		if sessionMessage.Deleted == 1 { // 已经撤回了
			return nil
		}
		// 删除session的消息
		affectedRow, errRevoke := l.appCtx.SessionMessageModel().RevokeSessionMessage(sessionMessage.SessionId, sessionMessage.MsgId, sessionMessage.FromUserId)
		if errRevoke != nil {
			return errRevoke
		}
		l.appCtx.Logger().Infof("RevokeUserMessage, affectedRow: %d, req: %v", affectedRow, req)
		if affectedRow == 0 {
			return nil
		}
		sendMessageReq := dto.SendMessageReq{
			CId:    l.genClientId(),
			SId:    req.SId,
			Type:   model.MsgTypeRevoke,
			FUid:   req.UId,
			CTime:  time.Now().UnixMilli(),
			RMsgId: &req.MsgId,
		} // 发送给session下的所有人
		if _, err = l.SendMessage(sendMessageReq); err != nil {
			l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SId, req.MsgId, err)
		}
	} else {
		l.appCtx.Logger().Errorf("RevokeUserMessage err:%d, %d, %d, %v", req.UId, req.SId, req.MsgId, err)
	}
	return nil
}

func (l *MessageLogic) ReeditUserMessage(req dto.ReeditUserMessageReq) error {
	if sessionMessage, err := l.appCtx.SessionMessageModel().FindSessionMessage(req.UId, req.SId, req.MsgId); err == nil {
		if sessionMessage.MsgType < 0 { // 小于0的类型消息为状态操作消息，不能发送撤回
			return errorx.ErrMessageTypeNotSupport
		}
		// 更新session的消息内容
		_, err = l.appCtx.SessionMessageModel().UpdateSessionMessageContent(sessionMessage.SessionId, sessionMessage.MsgId, sessionMessage.FromUserId, req.Content, model.MsgStatusReedit)
		if err != nil {
			return err
		}
		sendMessageReq := dto.SendMessageReq{
			CId:    l.genClientId(),
			SId:    req.SId,
			Type:   model.MsgTypeReedit,
			FUid:   req.UId,
			CTime:  time.Now().UnixMilli(),
			Body:   req.Content,
			RMsgId: &req.MsgId,
		} // 发送给session下的所有人
		if _, err = l.SendMessage(sendMessageReq); err != nil {
			l.appCtx.Logger().Errorf("ReeditUserMessage err:%d, %d, %d, %v", req.UId, req.SId, req.MsgId, err)
		}
	} else {
		l.appCtx.Logger().Errorf("ReeditUserMessage err:%d, %d, %d, %v", req.UId, req.SId, req.MsgId, err)
	}
	return nil
}

func (l *MessageLogic) genClientId() int64 {
	return l.appCtx.SnowflakeNode().Generate().Int64()
}
