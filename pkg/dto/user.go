package dto

type PostUserOnlineReq struct {
	NodeId int64 `json:"node_id" binding:"required"`
	ConnId int64 `json:"conn_id" binding:"required"`
	Online bool  `json:"online"`
	UId    int64 `json:"u_id" binding:"required"`
}

type GetUsersOnlineStatusReq struct {
	UIds []int64 `json:"u_ids" form:"u_ids"`
}

type UserOnlineStatus struct {
	UId            int64 `json:"u_id"`
	Online         bool  `json:"online"`
	LastOnlineTime int64 `json:"last_online_time"`
}

type GetUsersOnlineStatusRes struct {
	UsersOnlineStatus []*UserOnlineStatus `json:"data"`
}
