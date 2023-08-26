package handler

import (
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

const (
	tokenKey = "Token"
	uidKey   = "Uid"
)

func userTokenAuth(appCtx *app.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get(tokenKey)
		if token == "" {
			dto.ResponseUnauthorized(context)
			context.Abort()
			return
		}
		req := rpc.GetUserIdByTokenReq{Token: token}
		res, err := appCtx.RpcUserApi().GetUserIdByToken(req)
		if err != nil {
			dto.ResponseUnauthorized(context)
			context.Abort()
		} else {
			context.Set(uidKey, res.UserId)
			context.Next()
		}
	}
}

func whiteIpAuth(appCtx *app.Context) gin.HandlerFunc {
	ipWhiteList := appCtx.Config().IpWhiteList
	ips := strings.Split(ipWhiteList, ",")
	return func(context *gin.Context) {
		ip := context.ClientIP()
		appCtx.Logger().Infof("RemoteAddr: %s", ip)
		if isIpValid(ip, ips) {
			dto.ResponseForbidden(context)
			context.Abort()
		} else {
			context.Next()
		}
	}
}

func isIpValid(clientIp string, whiteIpList []string) bool {
	ip := net.ParseIP(clientIp)
	for _, whiteIp := range whiteIpList {
		_, ipNet, err := net.ParseCIDR(whiteIp)
		if err != nil {
			return false
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
