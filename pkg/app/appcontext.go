package app

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/conf/loader"
	"github.com/THK-IM/THK-IM-Server/pkg/metric"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/THK-IM/THK-IM-Server/pkg/mq"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"github.com/THK-IM/THK-IM-Server/pkg/websocket"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Context struct {
	startTime       int64
	nodeId          int64
	config          *conf.Config
	logger          *logrus.Entry
	redisCache      *redis.Client
	database        *gorm.DB
	collector       *metric.Collector
	snowflakeNode   *snowflake.Node
	httpEngine      *gin.Engine
	websocketServer websocket.Server
	rpcMap          map[string]interface{}
	modelMap        map[string]interface{}
	publisherMap    map[string]mq.Publisher
	subscriberMap   map[string]mq.Subscriber
}

func (c *Context) ServerEventPublisher() mq.Publisher {
	return c.publisherMap["server_event"]
}

func (c *Context) MsgPusherPublisher() mq.Publisher {
	return c.publisherMap["push_msg"]
}

func (c *Context) MsgSaverPublisher() mq.Publisher {
	return c.publisherMap["save_msg"]
}

func (c *Context) MsgPusherSubscriber() mq.Subscriber {
	return c.subscriberMap["push_msg"]
}

func (c *Context) MsgSaverSubscriber() mq.Subscriber {
	return c.subscriberMap["save_msg"]
}

func (c *Context) ServerEventSubscriber() mq.Subscriber {
	return c.subscriberMap["server_event"]
}

func (c *Context) StartTime() int64 {
	return c.startTime
}

func (c *Context) NodeId() int64 {
	return c.nodeId
}

func (c *Context) Config() *conf.Config {
	return c.config
}

func (c *Context) RedisCache() *redis.Client {
	return c.redisCache
}

func (c *Context) Database() *gorm.DB {
	return c.database
}

func (c *Context) SessionModel() model.SessionModel {
	return c.modelMap["session"].(model.SessionModel)
}

func (c *Context) SessionMessageModel() model.SessionMessageModel {
	return c.modelMap["session_message"].(model.SessionMessageModel)
}

func (c *Context) SessionUserModel() model.SessionUserModel {
	return c.modelMap["session_user"].(model.SessionUserModel)
}

func (c *Context) UserMessageModel() model.UserMessageModel {
	return c.modelMap["user_message"].(model.UserMessageModel)
}

func (c *Context) UserSessionModel() model.UserSessionModel {
	return c.modelMap["user_session"].(model.UserSessionModel)
}

func (c *Context) UserOnlineStatusModel() model.UserOnlineStatusModel {
	return c.modelMap["user_online_status"].(model.UserOnlineStatusModel)
}

func (c *Context) RpcMsgApi() rpc.MsgApi {
	api, ok := c.rpcMap["msg-api"].(rpc.MsgApi)
	if ok {
		return api
	} else {
		return nil
	}
}

func (c *Context) RpcUserApi() rpc.UserApi {
	api, ok := c.rpcMap["user-api"].(rpc.UserApi)
	if ok {
		return api
	} else {
		return nil
	}
}

func (c *Context) Collector() *metric.Collector {
	return c.collector
}

func (c *Context) SnowflakeNode() *snowflake.Node {
	return c.snowflakeNode
}

func (c *Context) HttpEngine() *gin.Engine {
	return c.httpEngine
}

func (c *Context) WebsocketServer() websocket.Server {
	return c.websocketServer
}

func (c *Context) Logger() *logrus.Entry {
	return c.logger
}

func NewAppContext(config *conf.Config, httpEngine *gin.Engine) *Context {
	logger := loader.LoadLogg(config.Name, config.Logg)
	// gin.DefaultWriter = logger.WriterLevel(logrus.InfoLevel)
	// gin.DefaultErrorWriter = logger.WriterLevel(logrus.ErrorLevel)

	redisCache := loader.LoadRedis(config.RedisSource)
	id, startTime := loader.LoadNodeId(config, redisCache)
	snowflakeNode, err := snowflake.NewNode(id)
	if err != nil {
		panic(err)
	}

	ctx := &Context{
		httpEngine:    httpEngine,
		nodeId:        id,
		startTime:     startTime,
		config:        config,
		logger:        logger,
		redisCache:    redisCache,
		snowflakeNode: snowflakeNode,
	}
	if config.DataSource != nil {
		ctx.database = loader.LoadDataBase(logger, config.DataSource)
		if config.Models != nil {
			ctx.modelMap = loader.LoadModels(config.Models, ctx.database, logger, snowflakeNode)
		}
	}
	nodeIdStr := fmt.Sprintf("%d", id)
	if config.MsgQueue.Publishers != nil {
		ctx.publisherMap = loader.LoadPublishers(config.MsgQueue.Publishers, nodeIdStr, logger)
	}
	if config.MsgQueue.Subscribers != nil {
		ctx.subscriberMap = loader.LoadSubscribers(config.MsgQueue.Subscribers, nodeIdStr, logger)
	}
	if config.Sdks != nil {
		ctx.rpcMap = loader.LoadSdks(config.Sdks, logger)
	}
	if config.WebSocket != nil {
		ctx.websocketServer = websocket.NewServer(config.WebSocket, logger, httpEngine, snowflakeNode, config.Mode)
	}
	if config.Metric != nil {
		ctx.collector = metric.NewCollector(config.Name, id, config.Metric, logger, httpEngine)
	}
	return ctx
}

func (c *Context) Start() {
	address := fmt.Sprintf("%s:%s", c.config.Host, c.config.Port)
	if e := c.httpEngine.Run(address); e != nil {
		panic(e)
	} else {
		c.logger.Infof("%s server start at: %s", c.config.Name, address)
	}
}
