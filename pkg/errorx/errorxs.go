package errorx

var (
	ErrParamsError            = NewErrorX(4000000, "Params Error")
	ErrUserNotOnLine          = NewErrorX(4000001, "User not on line")
	ErrSessionInvalid         = NewErrorX(4000002, "Invalid session")
	ErrGroupMemberCountBeyond = NewErrorX(4000003, "group member count beyond")
	ErrGroupAlreadyDeleted    = NewErrorX(4000004, "group has been deleted")
	ErrSessionType            = NewErrorX(4000005, "Session type error")
	ErrMessageFormat          = NewErrorX(4000006, "Message format error")
	ErrMessageTypeNotSupport  = NewErrorX(4000007, "Message type not support")
	ErrSessionMuted           = NewErrorX(4001001, "Session muted")
	ErrUserMuted              = NewErrorX(4001002, "User muted")
	ErrUserReject             = NewErrorX(4001003, "user reject your message")
	ErrMessageDeliveryFailed  = NewErrorX(5004001, "Message delivery failed")
)
