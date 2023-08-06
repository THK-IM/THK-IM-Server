package event

import "encoding/json"

const (
	// PushEventTypeKey 推送事件类型Key
	PushEventTypeKey = "push_type_key"
	// PushEventSubTypeKey 推送事件子类型Key
	PushEventSubTypeKey = "push_sub_type_key"
	// PushEventReceiversKey 推送事件子类型Key
	PushEventReceiversKey = "push_receivers_key"
	// PushEventBodyKey 推送事件Body Key
	PushEventBodyKey = "push_body_key"

	// PushCommonEventType 通用推送事件类型
	PushCommonEventType = 0
	PushUserEventType   = 1
	PushFriendEventType = 2
	PushGroupEventType  = 3
	PushMsgEventType    = 4
	PushOtherEventType  = 5

	CommonEventSubtypePing     = 1
	CommonEventSubtypePong     = 2
	CommonEventSubtypeSyncTime = 3

	UserEventSubtypeKickOff    = 1
	UserEventSubtypeOnline     = 2
	UserEventSubtypeInfoUpdate = 3
)

type (
	PushBody struct {
		Type    int    `json:"type"`
		SubType int    `json:"sub_type"`
		Body    string `json:"body"`
	}
)

func BuildPushBody(t int, subtype int, body string) (string, error) {
	pushBody := &PushBody{
		Type:    t,
		SubType: subtype,
		Body:    body,
	}
	content, err := json.Marshal(pushBody)
	return string(content), err
}
