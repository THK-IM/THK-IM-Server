package dto

type AckUserMessagesReq struct {
	UId        int64   `json:"u_id"`
	SessionId  int64   `json:"session_id" binding:"required"`
	MessageIds []int64 `json:"msg_ids" binding:"required"`
}

type ReadUserMessageReq struct {
	UId        int64   `json:"u_id"`
	SessionId  int64   `json:"session_id" binding:"required"`
	MessageIds []int64 `json:"msg_ids" binding:"required"`
}

type RevokeUserMessageReq struct {
	UId       int64 `json:"u_id"`
	SessionId int64 `json:"session_id" binding:"required"`
	MessageId int64 `json:"msg_id" binding:"required"`
}

type ReeditUserMessageReq struct {
	UId       int64  `json:"u_id"`
	SessionId int64  `json:"session_id" binding:"required"`
	MessageId int64  `json:"msg_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}
