package dto

type GetMessageReq struct {
	UId    int64 `json:"u_id" form:"u_id"`
	Offset int   `json:"offset" form:"offset"`
	Count  int   `json:"count" form:"count"`
	CTime  int   `json:"c_time" form:"c_time"`
}

type GetMessageRes struct {
	Data interface{} `json:"data"`
}

type DeleteMessageReq struct {
	Uid        int64   `json:"uid"`
	SessionId  int64   `json:"session_id"`
	MessageIds []int64 `json:"msg_ids"`
	TimeFrom   int64   `json:"time_from"`
	TimeTo     int64   `json:"time_to"`
}

type PushMessageReq struct {
	UIds        []int64 `json:"u_ids" binding:"required"`
	Type        int     `json:"type" binding:"required"`
	SubType     int     `json:"sub_type" binding:"required"`
	Body        string  `json:"body" binding:"required"`
	OfflinePush bool    `json:"offline_push"`
}

type PushMessageRes struct {
	OnlineUIds  []int64 `json:"online_ids,omitempty"`
	OfflineUIds []int64 `json:"offline_ids,omitempty"`
}

type Message struct {
	ClientId  int64  `json:"client_id"`
	SessionId int64  `json:"session_id"`
	MsgId     int64  `json:"msg_id"`
	Type      int    `json:"type"`
	FUid      int64  `json:"f_uid"`
	CTime     int64  `json:"c_time"`
	RMsgId    int64  `json:"r_msg_id,omitempty"`
	Body      string `json:"body"`
	Status    *int   `json:"status,omitempty"`
	AtUsers   string `json:"at_users,omitempty"`
}

type SendMessageReq struct {
	ClientId  int64   `json:"client_id" binding:"required"`
	SessionId int64   `json:"session_id" binding:"required"`
	Type      int     `json:"type" binding:"required"`
	FUid      int64   `json:"f_uid"`
	CTime     int64   `json:"c_time" binding:"required"`
	RMsgId    int64   `json:"r_msg_id,omitempty"`
	Body      string  `json:"body" binding:"required"`
	AtUsers   string  `json:"at_users,omitempty"`
	Receivers []int64 `json:"receivers,omitempty"`
}

type SendMessageRes struct {
	MsgId      int64   `json:"msg_id"`
	CreateTime int64   `json:"c_time"`
	OnlineIds  []int64 `json:"online_ids,omitempty"`
	OfflineIds []int64 `json:"offline_ids,omitempty"`
}

type KickUserReq struct {
	UId int64 `json:"u_id" binding:"required"`
}
