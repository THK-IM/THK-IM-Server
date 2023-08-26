package handler

import (
	"errors"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/gin-gonic/gin"
)

func RegisterApiHandlers(ctx *app.Context) {
	httpEngine := ctx.HttpEngine()
	userAuth := userTokenAuth(ctx)
	ipAuth := whiteIpAuth(ctx)
	var authMiddleware gin.HandlerFunc
	if ctx.Config().DeployMode == conf.DeployExposed {
		authMiddleware = userAuth
	} else if ctx.Config().DeployMode == conf.DeployBackend {
		authMiddleware = ipAuth
	} else {
		panic(errors.New("check deployMode conf"))
	}
	sessionGroup := httpEngine.Group("/session")
	sessionGroup.Use(authMiddleware)
	{
		sessionGroup.POST("", createSession(ctx))                      // 创建/获取session
		sessionGroup.PUT("/:id", updateSession(ctx))                   // 修改session相关信息
		sessionGroup.POST("/:id/user", addSessionUser(ctx))            // 会话增员
		sessionGroup.DELETE("/:id/user", deleteSessionUser(ctx))       // 会话减员
		sessionGroup.PUT("/:id/user", updateSessionUser(ctx))          // 会话成员修改
		sessionGroup.GET("/:id/message", getSessionMessages(ctx))      // 获取session下的消息列表
		sessionGroup.DELETE("/:id/message", deleteSessionMessage(ctx)) // 删除session下的消息列表
	}

	userSessionGroup := httpEngine.Group("/user_session")
	userSessionGroup.Use(authMiddleware)
	{
		userSessionGroup.GET("/:uid/:sid", getUserSession(ctx)) // 用户获取自己的session
		userSessionGroup.GET("/latest", getUserSessions(ctx))   // 用户获取自己的session列表
		userSessionGroup.PUT("", updateUserSession(ctx))        // 用户修改自己的session
	}

	messageGroup := httpEngine.Group("/message")
	messageGroup.Use(authMiddleware)
	{
		messageGroup.POST("", sendMessage(ctx))                 // 发送消息
		messageGroup.DELETE("", deleteUserMessage(ctx))         // 删除消息
		messageGroup.POST("/ack", ackUserMessages(ctx))         // 用户消息设置ack(已接收) 不支持超级群
		messageGroup.POST("/read", readUserMessage(ctx))        // 用户消息设置已读 不支持超级群
		messageGroup.POST("/revoke", revokeUserMessage(ctx))    // 用户消息撤回
		messageGroup.POST("/reedit", reeditUserMessage(ctx))    // 更新用户消息
		messageGroup.GET("/latest", getUserLatestMessages(ctx)) // 获取最近消息
	}

	systemGroup := httpEngine.Group("/system")
	systemGroup.Use(ipAuth)
	{
		systemGroup.POST("/user/online", updateUserOnlineStatus(ctx)) // 更新用户在线状态
		systemGroup.GET("/user/online", getUsersOnlineStatus(ctx))    // 获取用户上线状态
		systemGroup.POST("/user/kickoff", kickOffUser(ctx))           // 踢下线用户
		systemGroup.POST("/message", sendSystemMessage(ctx))          // 发送会话中的系统消息
		systemGroup.POST("/message/push", pushSystemMessage(ctx))     // 推送消息(用户消息/好友消息/群组消息/自定义消息)
	}
}
