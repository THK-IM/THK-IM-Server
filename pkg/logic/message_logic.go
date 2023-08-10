package logic

import (
	"encoding/json"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/THK-IM/THK-IM-Server/pkg/event"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type MessageLogic struct {
	ctx    *gin.Context
	appCtx *app.Context
}

func NewMessageLogic(ctx *gin.Context, appCtx *app.Context) MessageLogic {
	return MessageLogic{
		ctx:    ctx,
		appCtx: appCtx,
	}
}

func (l *MessageLogic) convSessionMessage2Message(sessionMsg *model.SessionMessage) *dto.Message {
	vo := dto.Message{
		ClientId:  sessionMsg.ClientId,
		FUid:      sessionMsg.FromUserId,
		SessionId: sessionMsg.SessionId,
		MsgId:     sessionMsg.MsgId,
		CTime:     sessionMsg.CreateTime,
		Body:      sessionMsg.MsgContent,
		AtUsers:   sessionMsg.AtUsers,
		Type:      sessionMsg.MsgType,
		RMsgId:    sessionMsg.ReplyMsgId,
	}
	return &vo
}

func (l *MessageLogic) convUserMessage2Message(userMsg *model.UserMessage) *dto.Message {
	msg := dto.Message{
		ClientId:  userMsg.ClientId,
		SessionId: userMsg.SessionId,
		Type:      userMsg.MsgType,
		MsgId:     userMsg.MsgId,
		FUid:      userMsg.FromUserId,
		CTime:     userMsg.CreateTime,
		RMsgId:    userMsg.ReplyMsgId,
		Body:      userMsg.MsgContent,
		Status:    &userMsg.Status,
		AtUsers:   userMsg.AtUsers,
	}
	return &msg
}

func (l *MessageLogic) GetUserMessages(req dto.GetMessageReq) (*dto.GetMessageRes, error) {
	userMessages, err := l.appCtx.UserMessageModel().GetUserMessages(req.UId, req.CTime, req.Offset, req.Count)
	if err != nil {
		return nil, err
	}
	messages := make([]*dto.Message, 0)
	for _, userMessage := range userMessages {
		userMessageV3 := l.convUserMessage2Message(userMessage)
		messages = append(messages, userMessageV3)
	}
	return &dto.GetMessageRes{Data: userMessages}, nil
}

func (l *MessageLogic) GetSessionMessages(req dto.GetSessionMessageReq) (*dto.GetMessageRes, error) {
	msgIds := make([]int64, 0)
	if req.MsgIds != "" {
		strIds := strings.Split(req.MsgIds, ",")
		for _, str := range strIds {
			if id, err := strconv.ParseInt(str, 10, 64); err != nil {
				return nil, err
			} else {
				msgIds = append(msgIds, id)
			}
		}
	}
	sessionMessages, err := l.appCtx.SessionMessageModel().GetSessionMessages(req.SessionId, req.CTime, req.Offset, req.Count, msgIds)
	if err != nil {
		return nil, err
	}
	messages := make([]*dto.Message, 0)
	for _, sessionMessage := range sessionMessages {
		message := l.convSessionMessage2Message(sessionMessage)
		messages = append(messages, message)
	}
	return &dto.GetMessageRes{Data: messages}, nil
}

func (l *MessageLogic) DelSessionMessage(req *dto.DelSessionMessageReq) error {
	err := l.appCtx.SessionMessageModel().DelMessages(req.SessionId, req.MessageIds, req.TimeFrom, req.TimeTo)
	return err
}

func (l *MessageLogic) SendMessage(req dto.SendMessageReq) (*dto.SendMessageRes, error) {
	session, e1 := l.appCtx.SessionModel().FindSession(req.SessionId, nil)
	if e1 != nil {
		l.appCtx.Logger().Error(e1)
		return nil, errorx.ErrInvalidSession
	}
	// req.FUid为0是系统消息, 不需要校验是否能对session发送消息
	if req.FUid > 0 {
		if session.Type != model.SingleSessionType {
			if session.Status&model.MutedBitInSessionStatus > 0 {
				return nil, errorx.ErrCannotSendMessage
			}
			userSession, e2 := l.appCtx.UserSessionModel().GetUserSession(req.FUid, req.SessionId)
			if e2 != nil {
				l.appCtx.Logger().Error(e2)
				return nil, errorx.ErrInvalidSession
			}
			if userSession.Status&model.MutedBitInUserSessionStatus > 0 {
				return nil, errorx.ErrCannotSendMessage
			}
		}
	}
	receivers := l.appCtx.SessionUserModel().FindUIdsInSessionWithoutStatus(req.SessionId, model.RejectBitInUserSessionStatus, req.Receivers)
	if receivers == nil || len(receivers) == 0 {
		return nil, errorx.ErrOtherRejectMessage
	}

	// 根据clientId和fromUserId查询是否已经发送过消息
	sessionMessage, errSession := l.appCtx.SessionMessageModel().FindMessageByClientId(req.SessionId, req.ClientId, req.FUid)
	if errSession != nil && errSession != gorm.ErrRecordNotFound {
		l.appCtx.Logger().Error(errSession, req)
		return nil, errSession
	}
	// 如果已经发送过，直接取数据库里的数据库, 没有发送过则插入数据库
	if sessionMessage == nil || sessionMessage.SessionId == 0 {
		// 插入数据库发送消息
		msgId := int64(l.appCtx.SnowflakeNode().Generate())
		sessionMessage, errSession = l.appCtx.SessionMessageModel().InsertMessage(req.ClientId, req.FUid, req.SessionId, msgId, req.Body, req.Type, req.AtUsers, req.RMsgId)
		if errSession != nil || sessionMessage == nil {
			l.appCtx.Logger().Error(errSession, req)
			return nil, errSession
		}
	}
	if onlineUIds, offlineUIds, err := l.publishSendMessageEvents(sessionMessage, session.Type, receivers); err != nil {
		return nil, errorx.ErrMessageDeliveryFailed
	} else {
		return &dto.SendMessageRes{
			MsgId:      sessionMessage.MsgId,
			CreateTime: sessionMessage.CreateTime,
			OnlineIds:  onlineUIds,
			OfflineIds: offlineUIds,
		}, nil
	}

}

func (l *MessageLogic) publishSendMessageEvents(sessionMsg *model.SessionMessage, sessionType int, receivers []int64) ([]int64, []int64, error) {
	userMsg := &dto.Message{
		ClientId:  sessionMsg.ClientId,
		MsgId:     sessionMsg.MsgId,
		SessionId: sessionMsg.SessionId,
		FUid:      sessionMsg.FromUserId,
		AtUsers:   sessionMsg.AtUsers,
		Type:      sessionMsg.MsgType,
		Body:      sessionMsg.MsgContent,
		RMsgId:    sessionMsg.ReplyMsgId,
		CTime:     sessionMsg.CreateTime,
	}
	msgJson, err := json.Marshal(userMsg)
	if err != nil {
		l.appCtx.Logger().Error(err)
		return nil, nil, err
	}
	msgJsonStr := string(msgJson)

	onlineUIds, offlineUIds, errPubPush := l.pubPushMessageEvent(event.PushMsgEventType, 0, msgJsonStr, receivers)
	if errPubPush != nil {
		l.appCtx.Logger().Error("pubPushMessageEvent, err:", errPubPush)
		return nil, nil, errPubPush
	}
	if sessionType != model.SuperGroupSessionType {
		errPubSave := l.pubSaveMsgEvent(msgJsonStr, receivers)
		if errPubSave != nil {
			l.appCtx.Logger().Error("pubSaveMsgEvent, err:", errPubSave)
			return nil, nil, errPubPush
		}
	}
	return onlineUIds, offlineUIds, nil
}

// PushMessage 业务消息推送
func (l *MessageLogic) PushMessage(req dto.PushMessageReq) (*dto.PushMessageRes, error) {
	// 在线推送
	onlineUIds, offlineUIds, err := l.pubPushMessageEvent(req.Type, req.SubType, req.Body, req.UIds)
	if err == nil {
		rsp := &dto.PushMessageRes{}
		rsp.OnlineUIds = onlineUIds
		rsp.OfflineUIds = offlineUIds
		return rsp, err
	} else {
		return nil, err
	}
}

func (l *MessageLogic) pubSaveMsgEvent(msgBody string, receivers []int64) error {
	if receiversStr, errJson := json.Marshal(receivers); errJson != nil {
		return errJson
	} else {
		m := make(map[string]interface{})
		m[event.SaveMsgEventKey] = msgBody
		m[event.SaveMsgUsersKey] = receiversStr
		return l.appCtx.MsgSaverPublisher().Pub("", m)
	}
}

// 发布推送消息
func (l *MessageLogic) pubPushMessageEvent(t, subType int, body string, uIds []int64) ([]int64, []int64, error) {
	onlineTimeout := l.appCtx.Config().OnlineTimeout
	offlineTime := time.Now().UnixMilli() - onlineTimeout*int64(time.Second)
	offlineUsers := make([]int64, 0)
	onlineUsers, err := l.appCtx.UserOnlineStatusModel().GetOnlineUsers(uIds, offlineTime)
	if err != nil {
		// 如果查询报错 默认全部用户为离线
		l.appCtx.Logger().Error("get userOnlineStatus error:", err)
		onlineUsers = make([]int64, 0)
	}
	onlineUserMap := make(map[int64]bool, 0)
	for _, uid := range onlineUsers {
		onlineUserMap[uid] = true
	}
	for _, uid := range uIds {
		online, ok := onlineUserMap[uid]
		if !ok && online {
			offlineUsers = append(offlineUsers, uid)
		}
	}
	if receiverStr, errJson := json.Marshal(onlineUsers); errJson != nil {
		return nil, nil, errJson
	} else {
		m := make(map[string]interface{})
		m[event.PushEventTypeKey] = t
		m[event.PushEventSubTypeKey] = subType
		m[event.PushEventBodyKey] = body
		m[event.PushEventReceiversKey] = string(receiverStr)

		err = l.appCtx.MsgPusherPublisher().Pub("", m)
		return onlineUsers, offlineUsers, err
	}
}

func (l *MessageLogic) DeleteUserMessage(req *dto.DeleteMessageReq) error {
	return l.appCtx.UserMessageModel().DeleteMessages(req.UId, req.SessionId, req.MessageIds, req.TimeFrom, req.TimeTo)
}
