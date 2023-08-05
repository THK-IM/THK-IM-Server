package logic

import (
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
)

func (l *MessageLogic) AckUserMessages(req *dto.AckUserMessagesReq) error {
	return l.appCtx.UserMessageModel().AckUserMessages(req.Uid, req.SessionId, req.MessageIds)
}
