package event

import "encoding/json"

const (
	// PushEventTypeKey 推送事件类型Key
	PushEventTypeKey = "push_type_key"
	// PushEventReceiversKey 推送事件子类型Key
	PushEventReceiversKey = "push_receivers_key"
	// PushEventBodyKey 推送事件Body Key
	PushEventBodyKey = "push_body_key"

	SignalNewMessage  = 0
	SignalPing        = 1
	SignalPong        = 2
	SignalSyncTime    = 3
	SignalConnId      = 4
	SignalKickOffUser = 5
	SignalExtended    = 100
)

type (
	SignalBody struct {
		Type int     `json:"type"`
		Body *string `json:"body"`
	}
)

func BuildSignalType(t int) (string, error) {
	pushBody := &SignalBody{
		Type: t,
	}
	content, err := json.Marshal(pushBody)
	return string(content), err
}

func BuildSignalBody(t int, body string) (string, error) {
	pushBody := &SignalBody{
		Type: t,
		Body: &body,
	}
	content, err := json.Marshal(pushBody)
	return string(content), err
}
