package dto

type AckUserMessagesReq struct {
	Uid        int64   `json:"uid"`
	SessionId  int64   `json:"session_id" binding:"required"`
	MessageIds []int64 `json:"msg_ids" binding:"required"`
}

type CreateSessionReq struct {
	Type     int     `json:"type" binding:"required"`
	EntityId int64   `json:"entity_id"`
	Members  []int64 `json:"members" binding:"required"`
}

type CreateSessionRes struct {
	SessionId int64 `json:"session_id"`
	EntityId  int64 `json:"entity_id"`
	Type      int   `json:"type"`
	CTime     int64 `json:"c_time"`
	MTime     int64 `json:"m_time"`
	Top       int64 `json:"top,omitempty"`
	Status    int   `json:"status,omitempty"`
}

type UpdateSessionReq struct {
	Id     int64   `json:"id" uri:"id"`
	Status *int    `json:"status"`
	Name   *string `json:"name"`
	Remark *string `json:"remark"`
}

type UpdateSessionRes struct {
}

type UpdateUserSessionReq struct {
	UId    int64  `json:"u_id"`
	SId    int64  `json:"s_id"`
	Top    *int64 `json:"top"`
	Status *int   `json:"status"`
}

type GetUserSessionsReq struct {
	UId    int64 `json:"u_id" form:"u_id"`
	Offset int   `json:"offset" form:"offset"`
	Count  int   `json:"count" form:"count"`
	MTime  int64 `json:"m_time" form:"m_time"`
}

type GetSessionMessageReq struct {
	SessionId int64  `json:"id" form:"id"`
	CTime     int    `json:"c_time" form:"c_time"`
	Offset    int    `json:"offset" form:"offset"`
	Count     int    `json:"count" form:"count"`
	MsgIds    string `json:"msg_ids" form:"msg_ids"`
}

type DelSessionMessageReq struct {
	SessionId  int64   `json:"session_id"`
	MessageIds []int64 `json:"msg_ids"`
	TimeFrom   int64   `json:"time_from"`
	TimeTo     int64   `json:"time_to"`
}

type UserSession struct {
	SessionId int64 `json:"session_id"`
	Type      int   `json:"type"`
	Status    int   `json:"status"`
	Top       int64 `json:"top"`
	EntityId  int64 `json:"entity_id"`
	CTime     int64 `json:"c_time"`
	MTime     int64 `json:"m_time"`
}

type GetUserSessionsRes struct {
	Data []*UserSession `json:"data"`
}

type SessionAddUserReq struct {
	EntityId int64   `json:"entity_id" binding:"required"`
	UIds     []int64 `json:"u_ids" binding:"required"`
}

type SessionDelUserReq struct {
	UIds []int64 `json:"u_ids" binding:"required"`
}

type SessionUserUpdateReq struct {
	SId    int64 `json:"s_id"`
	UId    int64 `json:"u_id"`
	Status *int  `json:"status" binding:"required"`
}
