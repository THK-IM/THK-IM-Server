package handler

import (
	"errors"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/gin-gonic/gin"
)

func RegisterGroupApiHandlers(ctx *app.Context) {
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

	groupRoute := httpEngine.Group("group")
	groupRoute.Use(authMiddleware)
	{
		groupRoute.POST("", createGroup(ctx))                           // 创建群
		groupRoute.GET("/:id", queryGroup(ctx))                         // 查询群资料
		groupRoute.PUT("/:id", updateGroup(ctx))                        // 修改群信息
		groupRoute.DELETE("/:id", deleteGroup(ctx))                     // 删除群
		groupRoute.POST("/:id/transfer", transferGroup(ctx))            // 转让群
		groupRoute.POST("/:id/admin", transferGroup(ctx))               // 添加管理员
		groupRoute.DELETE("/:id/admin", transferGroup(ctx))             // 删除管理员
		groupRoute.GET("/apply", queryJoinGroupApply(ctx))              // 查询群申请列表
		groupRoute.POST("/apply", createJoinGroupApply(ctx))            // 申请加入群
		groupRoute.POST("/apply/:id/review", reviewJoinGroupApply(ctx)) // 审核群加入申请
		groupRoute.POST("/invite", inviteJoinGroup(ctx))
		groupRoute.POST("/invite", inviteJoinGroup(ctx))
	}
}
