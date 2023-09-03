package dto

type CreateSessionReq struct {
	Type     int     `json:"type" binding:"required"`
	EntityId *int64  `json:"entity_id"`                  // type 为group或supergroup时存在
	Members  []int64 `json:"members" binding:"required"` // 数组0位置为创建人
}

type CreateSessionRes struct {
	SId      int64  `json:"s_id"`
	EntityId int64  `json:"entity_id"`
	Type     int    `json:"type"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	Mute     int    `json:"mute"`
	Role     int    `json:"role"`
	CTime    int64  `json:"c_time"`
	MTime    int64  `json:"m_time"`
	Top      int64  `json:"top"`
	Status   int    `json:"status"`
	Success  bool   `json:"success"` // 如果之前已经创建，false
}

type UpdateSessionReq struct {
	Id     int64   `json:"id" uri:"id"`
	Mute   *int    `json:"mute"`
	Name   *string `json:"name"`
	Remark *string `json:"remark"`
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
	SId    int64  `json:"s_id" form:"id"`
	CTime  int64  `json:"c_time" form:"c_time"`
	Offset int    `json:"offset" form:"offset"`
	Count  int    `json:"count" form:"count"`
	MsgIds string `json:"msg_ids" form:"msg_ids"`
}

type DelSessionMessageReq struct {
	SId      int64   `json:"s_id"`
	MsgIds   []int64 `json:"msg_ids"`
	TimeFrom int64   `json:"time_from"`
	TimeTo   int64   `json:"time_to"`
}

type UserSession struct {
	SId      int64  `json:"s_id"`
	Name     string `json:"name"`
	Remark   string `json:"remark"`
	Type     int    `json:"type"`
	Status   int    `json:"status"`
	Role     int    `json:"role"`
	Mute     int    `json:"mute"`
	Top      int64  `json:"top"`
	EntityId int64  `json:"entity_id"`
	CTime    int64  `json:"c_time"`
	MTime    int64  `json:"m_time"`
}

type SessionUser struct {
	SId    int64 `json:"s_id"`
	UId    int64 `json:"u_id"`
	Type   int   `json:"type"`
	Mute   int   `json:"mute"`
	Role   int   `json:"role"`
	Status int   `json:"status"`
	CTime  int64 `json:"c_time"`
	MTime  int64 `json:"m_time"`
}

type GetUserSessionsRes struct {
	Data []*UserSession `json:"data"`
}

type GetSessionUserReq struct {
	SId   int64 `json:"s_id" form:"s_id"`
	Role  *int  `json:"role" form:"role"`
	MTime int64 `json:"m_time" form:"m_time"`
	Count int   `json:"count" form:"count"`
}

type GetSessionUserRes struct {
	Data []*SessionUser `json:"data"`
}

type SessionAddUserReq struct {
	EntityId int64   `json:"entity_id" binding:"required"`
	UIds     []int64 `json:"u_ids" binding:"required"`
	Role     int     `json:"role" binding:"required"`
}

type SessionDelUserReq struct {
	UIds []int64 `json:"u_ids" binding:"required"`
}

type SessionUserUpdateReq struct {
	SId  int64   `json:"s_id" binding:"required"`
	UIds []int64 `json:"u_ids" binding:"required"`
	Role *int    `json:"role"`
	Mute *int    `json:"mute"`
}
