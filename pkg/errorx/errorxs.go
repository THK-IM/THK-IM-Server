package errorx

var ErrUserNotOnLine = New("user not on line")
var ErrWrite = New("write msg error")
var ErrConnClosed = New("conn closed")
var ErrUnCatchHandler = New("unCatch handler")
var ErrMsgBodyLen = New("error msg body len")
var ErrMsgBodyContent = New("error msg body content")
var ErrTokenCheckFailed = New("user token error")
var ErrClientInvalid = New("client invalid")

var ErrParamsError = New("Params Error")
var ErrCannotSendMessage = New("Can not send message")
var ErrInvalidSession = New("Invalid session")
var ErrOtherRejectMessage = New("The other party refuses to receive the message")
var ErrMessageDeliveryFailed = New("Message delivery failed")

var ErrMessageFormat = New("Message format error")
var ErrSessionType = NewErrorX(5000000, "invalid session type")
var ErrGroupMemberCount = NewErrorX(5000000, "group member beyond")
var ErrGroupAlreadyDeleted = NewErrorX(5000001, "group has been deleted")
