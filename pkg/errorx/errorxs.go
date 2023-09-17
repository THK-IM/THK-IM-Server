package errorx

var (
	ErrParamsError            = NewErrorX(4000000, "Params Error")
	ErrPermission             = NewErrorX(4001001, "Permission denied")
	ErrUserNotOnLine          = NewErrorX(4000002, "User not on line")
	ErrSessionInvalid         = NewErrorX(4000003, "Invalid session")
	ErrGroupMemberCountBeyond = NewErrorX(4000004, "group member count beyond")
	ErrGroupAlreadyDeleted    = NewErrorX(4000005, "group has been deleted")
	ErrSessionType            = NewErrorX(4000006, "Session type error")
	ErrMessageFormat          = NewErrorX(4000007, "Message format error")
	ErrMessageTypeNotSupport  = NewErrorX(4000008, "Message type not support")
	ErrSessionMessageInvalid  = NewErrorX(4000009, "Invalid session or message")
	ErrSessionMuted           = NewErrorX(4001001, "Session muted")
	ErrUserMuted              = NewErrorX(4001002, "User muted")
	ErrUserReject             = NewErrorX(4001003, "user reject your message")

	ErrServerUnknown         = NewErrorX(5000000, "Server unknown err")
	ErrMessageDeliveryFailed = NewErrorX(5004001, "Message delivery failed")
)
