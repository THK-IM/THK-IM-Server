package handler

import "github.com/THK-IM/THK-IM-Server/pkg/app"

func RegisterApiHandlers(ctx *app.Context) {
	httpEngine := ctx.HttpEngine()
	sessionGroup := httpEngine.Group("/session")
	{
		sessionGroup.POST("", createSession(ctx))                      // 创建/获取session
		sessionGroup.PUT("/:id", updateSession(ctx))                   // 修改session相关信息
		sessionGroup.POST("/:id/user", addSessionUser(ctx))            // 会话增员
		sessionGroup.DELETE("/:id/user", deleteSessionUser(ctx))       // 会话减员
		sessionGroup.PUT("/:id/user/:uid", updateSessionUser(ctx))     // 会话成员修改
		sessionGroup.GET("/:id/message", getSessionMessage(ctx))       // 获取session下的消息列表
		sessionGroup.DELETE("/:id/message", deleteSessionMessage(ctx)) // 删除session下的消息列表
	}

	messageGroup := httpEngine.Group("/message")
	{
		messageGroup.POST("", sendMessage(ctx))                 // 发送消息
		messageGroup.DELETE("", deleteUserMessage(ctx))         // 删除消息
		messageGroup.POST("/ack", ackUserMessages(ctx))         // 用户消息设置ack(已接受)
		messageGroup.POST("/read", readUserMessage(ctx))        // 用户消息设置已读
		messageGroup.POST("/revoke", revokeUserMessage(ctx))    // 用户消息撤回
		messageGroup.POST("/reedit", reeditUserMessage(ctx))    // 更新用户消息
		messageGroup.POST("/push", pushMessage(ctx))            // 推送消息(用户消息/好友消息/群组消息/自定义消息)
		messageGroup.GET("/latest", getUserLatestMessages(ctx)) // 获取最近消息
	}

	userSessionGroup := httpEngine.Group("/user_session")
	{
		userSessionGroup.GET("/:uid/:sid", getUserSession(ctx)) // 用户获取自己的session
		userSessionGroup.GET("/latest", getUserSessions(ctx))   // 用户获取自己的session列表
		userSessionGroup.PUT("", updateUserSession(ctx))        // 用户修改自己的session
	}

	userGroup := httpEngine.Group("/user")
	{
		userGroup.POST("/online", updateUserOnlineStatus(ctx)) // 更新用户在线状态
		userGroup.GET("/online", getUsersOnlineStatus(ctx))    // 获取用户上线状态
		userGroup.POST("/kickoff", kickOffUser(ctx))           // 踢下线用户
	}
}
