package websocket

import (
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

const (
	UidKey      = "uid"
	TokenKey    = "token"
	PlatformKey = "platform"
)

type OnClientChange func(client Client)
type OnClientMsgReceived func(msg string, client Client)
type UidGetter func(token string, platform string) (uid int64, err error)

type Server interface {
	Init() error
	Clients() map[int64][]Client
	ClientCount() int64
	SetUidGetter(g UidGetter)
	SetOnClientConnected(f OnClientChange)
	SetOnClientClosed(f OnClientChange)
	SetOnClientMsgReceived(r OnClientMsgReceived)
	AddClient(uid int64, client Client) (err error)
	RemoveClient(uid int64, reason string, client Client) error
	OnUserConnected(uid int64, connId int64, platform string) error
	SendMessage(uid int64, msg string) (err error)
	SendMessageToUsers(uIds []int64, msg string) (err error)
}

type WsServer struct {
	g                   *gin.Engine
	mode                string
	conf                *conf.WebSocket
	mutex               *sync.RWMutex
	logger              *logrus.Entry // 日志打印
	curCount            *atomic.Int64
	OnClientMsgReceived OnClientMsgReceived
	snowflakeNode       *snowflake.Node
	GetUidByToken       UidGetter
	userClients         map[int64][]Client
	OnClientConnected   OnClientChange
	OnClientClosed      OnClientChange
}

func NewServer(conf *conf.WebSocket, logger *logrus.Entry, g *gin.Engine, snowflakeNode *snowflake.Node, mode string) *WsServer {
	curCount := &atomic.Int64{}
	curCount.Store(0)
	mutex := &sync.RWMutex{}
	return &WsServer{
		g:             g,
		mode:          mode,
		logger:        logger,
		conf:          conf,
		curCount:      curCount,
		mutex:         mutex,
		snowflakeNode: snowflakeNode,
		userClients:   make(map[int64][]Client, 0),
	}
}

func (server *WsServer) SetUidGetter(g UidGetter) {
	server.GetUidByToken = g
}

func (server *WsServer) SetOnClientConnected(f OnClientChange) {
	server.OnClientConnected = f
}
func (server *WsServer) SetOnClientClosed(f OnClientChange) {
	server.OnClientClosed = f
}

func (server *WsServer) AddClient(uid int64, client Client) (err error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	clients, ok := server.userClients[uid]
	if !ok {
		clients = make([]Client, 0)
	}
	clients = append(clients, client)
	server.userClients[uid] = clients

	server.curCount.Add(1)
	if server.OnClientConnected != nil {
		server.OnClientConnected(client)
	}
	return
}

func (server *WsServer) RemoveClient(uid int64, reason string, client Client) (err error) {
	server.mutex.Lock()
	clients, ok := server.userClients[uid]
	if ok {
		if len(clients) == 1 {
			if clients[0].Info().Id == client.Info().Id {
				delete(server.userClients, uid)
				server.curCount.Add(-1)
			}
		} else {
			for i := 0; i < len(clients); i++ {
				if clients[i].Info().Id == client.Info().Id {
					newClients := append(clients[:i], clients[i+1:]...)
					server.userClients[uid] = newClients
					server.curCount.Add(-1)
					break
				}
			}
		}
	}
	server.mutex.Unlock()
	ok, err = client.Close(reason)
	if err == nil && ok {
		if server.OnClientConnected != nil {
			server.OnClientClosed(client)
		}
	}
	return err
}

func (server *WsServer) OnUserConnected(uid int64, connId int64, platform string) error {
	server.mutex.RLock()
	clients, ok := server.userClients[uid]
	server.mutex.RUnlock()
	if ok {
		for _, client := range clients {
			if client.Info().Id != connId {
				if server.conf.MultiPlatform == 0 {
					if err := server.RemoveClient(uid, "connect at other device", client); err != nil {
						server.logger.Error("OnUserConnected RemoveClient err: ", err)
					}
				} else if server.conf.MultiPlatform == 1 {
					if client.Info().Platform == platform {
						if err := server.RemoveClient(uid, "connect at other device", client); err != nil {
							server.logger.Error("OnUserConnected RemoveClient err: ", err)
						}
					}
				}
			}
		}
	}
	return nil
}

func (server *WsServer) SendMessage(uid int64, msg string) (err error) {
	server.mutex.RLock()
	clients, ok := server.userClients[uid]
	server.mutex.RUnlock()
	if ok {
		for _, c := range clients {
			if e := c.WriteMessage(msg); e != nil {
				server.logger.Errorf("client: %s, err, %s", c.Info(), err.Error())
			}
		}
	}
	return nil
}

func (server *WsServer) SendMessageToUsers(uIds []int64, msg string) (err error) {
	server.mutex.RLock()
	allClients := make([]Client, 0)
	for _, uid := range uIds {
		clients, ok := server.userClients[uid]
		if ok {
			allClients = append(allClients, clients...)
		}
	}
	server.mutex.RUnlock()
	server.logger.Info("SendMessageToUsers", uIds, len(allClients))
	for _, c := range allClients {
		e := c.WriteMessage(msg)
		if e != nil {
			server.logger.Errorf("client: %v, err, %s", c.Info(), err.Error())
		}
	}
	return nil
}

func (server *WsServer) Init() error {
	ws := websocket.Server{
		Handshake: func(c *websocket.Config, r *http.Request) error {
			return nil
		},
		Handler: server.onNewConn,
	}
	server.g.GET(server.conf.Uri, func(ctx *gin.Context) {
		err := server.getToken(ctx)
		if err != nil {
			ctx.Status(http.StatusForbidden)
		} else {
			ws.ServeHTTP(ctx.Writer, ctx.Request)
		}
	})
	return nil
}

func (server *WsServer) Clients() map[int64][]Client {
	return server.userClients
}

func (server *WsServer) ClientCount() int64 {
	return server.curCount.Load()
}

func (server *WsServer) SetOnClientMsgReceived(r OnClientMsgReceived) {
	server.OnClientMsgReceived = r
}

func (server *WsServer) onNewConn(ws *websocket.Conn) {
	if server.curCount.Load() >= server.conf.MaxClient {
		_ = ws.Close()
		server.logger.Infof("client count reach max count %d", server.conf.MaxClient)
		return
	}
	uid := ws.Request().Header.Get(UidKey)
	uId, err := strconv.Atoi(uid)
	if err != nil {
		_ = ws.Close()
		server.logger.Infof("uid: %s is invaild", uid)
		return
	}
	platform := ws.Request().Header.Get(PlatformKey)
	id := server.snowflakeNode.Generate()
	client := NewClient(ws, int64(id), int64(uId), platform, server)
	err = server.AddClient(int64(uId), client)
	if err != nil {
		server.logger.Error(err)
	} else {
		client.AcceptMessage()
	}
}

func (server *WsServer) getToken(ctx *gin.Context) error {
	pf := ctx.Query(PlatformKey)
	if strings.EqualFold(pf, "") {
		pf = ctx.GetHeader(PlatformKey)
		if strings.EqualFold(pf, "") {
			pf, _ = ctx.Cookie(PlatformKey)
		}
	}

	// debug 模式下直接传uid, 线上环境需要传token
	if server.mode == "debug" {
		uid := ctx.Query(UidKey)
		if strings.EqualFold(uid, "") {
			uid = ctx.GetHeader(UidKey)
			if strings.EqualFold(uid, "") {
				uid, _ = ctx.Cookie(UidKey)
			}
		}
		if uid != "" {
			ctx.Request.Header.Set(PlatformKey, pf)
			ctx.Request.Header.Set(UidKey, uid)
			return nil
		} else {
			return errors.New("token nil")
		}
	} else {
		token := ctx.Query(TokenKey)
		if strings.EqualFold(token, "") {
			token = ctx.GetHeader(TokenKey)
			if strings.EqualFold(token, "") {
				token, _ = ctx.Cookie(TokenKey)
			}
		}
		if strings.EqualFold("", token) {
			return errors.New("token nil")
		} else {
			uid, err := server.GetUidByToken(token, pf)
			if err == nil {
				server.logger.Infof("GetUidByToken, token: %s, uid: %d", token, uid)
				ctx.Request.Header.Set(PlatformKey, pf)
				ctx.Request.Header.Set(UidKey, fmt.Sprintf("%d", uid))
			}
			return err
		}
	}

}
