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

	// PushCommonEvent 通用推送事件类型
	PushCommonEvent = 0
	PushUserEvent   = 1
	PushFriendEvent = 2
	PushGroupEvent  = 3
	PushMsgEvent    = 4
	PushOtherEvent  = 5
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
