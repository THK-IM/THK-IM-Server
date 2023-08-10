package errorx

var ErrParamsError = NewErrorX(4000000, "Params Error")

var ErrUserNotOnLine = NewErrorX(5000000, "User not on line")

var ErrSessionType = NewErrorX(5001000, "Session type error")
var ErrSessionInvalid = NewErrorX(5001001, "Invalid session")
var ErrSessionMuted = NewErrorX(5001002, "Session muted")
var ErrUserMuted = NewErrorX(5001003, "User muted")
var ErrUserReject = NewErrorX(5001004, "user reject your message")

var ErrGroupMemberCountBeyond = NewErrorX(5002001, "group member count beyond")
var ErrGroupAlreadyDeleted = NewErrorX(5002002, "group has been deleted")

var ErrMessageFormat = NewErrorX(5004000, "Message format error")
var ErrMessageDeliveryFailed = NewErrorX(5004002, "Message delivery failed")
