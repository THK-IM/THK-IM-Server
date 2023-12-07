package handler

import (
	"errors"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/gin-gonic/gin"
)

func RegisterMsgApiHandlers(ctx *app.Context) {
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
	sessionRoute := httpEngine.Group("/session")
	sessionRoute.Use(authMiddleware)
	{
		sessionRoute.POST("", createSession(ctx))                      // 创建/获取session
		sessionRoute.PUT("/:id", updateSession(ctx))                   // 修改session相关信息
		sessionRoute.GET("/:id/user", getSessionUser(ctx))             // 会话成员查询
		sessionRoute.POST("/:id/user", addSessionUser(ctx))            // 会话增员
		sessionRoute.DELETE("/:id/user", deleteSessionUser(ctx))       // 会话减员
		sessionRoute.PUT("/:id/user", updateSessionUser(ctx))          // 会话成员修改
		sessionRoute.GET("/:id/message", getSessionMessages(ctx))      // 获取session下的消息列表
		sessionRoute.DELETE("/:id/message", deleteSessionMessage(ctx)) // 删除session下的消息列表

		// 如果提供内置对象存储服务，则开放接口
		if ctx.ObjectStorage() != nil {
			sessionRoute.GET("/object/upload_params", getObjectUploadParams(ctx)) // 获取对象上传参数
			sessionRoute.GET("/object/download_url", getObjectDownloadUrl(ctx))   // 获取对象,鉴权后重定向到签名后的minio地址
		}
	}

	userSessionRoute := httpEngine.Group("/user_session")
	userSessionRoute.Use(authMiddleware)
	{
		userSessionRoute.GET("/latest", getUserSessions(ctx))         // 用户获取自己最近的session列表
		userSessionRoute.GET("/:uid/:sid", getUserSession(ctx))       // 用户获取自己的session
		userSessionRoute.PUT("", updateUserSession(ctx))              // 用户修改自己的session
		userSessionRoute.DELETE("/:uid/:sid", deleteUserSession(ctx)) // 用户删除自己的session
	}

	messageRoute := httpEngine.Group("/message")
	messageRoute.Use(authMiddleware)
	{
		messageRoute.GET("/latest", getUserLatestMessages(ctx)) // 获取最近消息
		messageRoute.POST("", sendMessage(ctx))                 // 发送消息
		messageRoute.DELETE("", deleteUserMessage(ctx))         // 删除消息
		messageRoute.POST("/ack", ackUserMessages(ctx))         // 用户消息设置ack(已接收) 不支持超级群
		messageRoute.POST("/read", readUserMessage(ctx))        // 用户消息设置已读 不支持超级群
		messageRoute.POST("/revoke", revokeUserMessage(ctx))    // 用户消息撤回
		messageRoute.POST("/reedit", reeditUserMessage(ctx))    // 更新用户消息
		messageRoute.POST("/forward", forwardUserMessage(ctx))  // 转发用户消息
	}

	systemRoute := httpEngine.Group("/system")
	systemRoute.Use(ipAuth)
	{
		systemRoute.POST("/user/online", updateUserOnlineStatus(ctx)) // 更新用户在线状态
		systemRoute.GET("/user/online", getUsersOnlineStatus(ctx))    // 获取用户上线状态
		systemRoute.POST("/user/kickoff", kickOffUser(ctx))           // 踢下线用户
		systemRoute.POST("/message/send", sendSystemMessage(ctx))     // 发送会话中的系统消息
		systemRoute.POST("/message/push", pushExtendedMessage(ctx))   // 推送消息(用户消息/好友消息/群组消息/自定义消息)
	}
}
