package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/app"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/conf/loader"
	"github.com/THK-IM/THK-IM-Server/pkg/handler"
	"github.com/gin-gonic/gin"
	"os"
)

var consulEndpoint = flag.String("config-consul-endpoint", "", "consul address")
var consulKey = flag.String("config-consul-key", "", "consul key")
var configFile = flag.String("config-file", "etc/msg_api_server.yaml", "the config file")

func getConsul() (endpoint, key string) {
	if *consulEndpoint != "" && *consulKey != "" {
		return *consulEndpoint, *consulKey
	} else {
		return os.Getenv("config-consul-endpoint"), os.Getenv("config-consul-key")
	}
}

func initConfig() conf.Config {
	var (
		config conf.Config
		err    error
	)
	cAddress, cKey := getConsul()
	if cAddress != "" && cKey != "" {
		config, err = conf.LoadFromConsul(cAddress, cKey)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
	} else {
		config, err = conf.Load(*configFile)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
	}
	return config
}

func initMsgApiServer(appCtx *app.Context) {
	handler.RegisterApiHandlers(appCtx)
	if err := loader.LoadTables(appCtx.Config().Models, appCtx.Database()); err != nil {
		panic(err)
	}
}

func initMsgPushServer(appCtx *app.Context) {
	if appCtx.WebsocketServer() != nil {
		handler.RegisterMsgPushHandlers(appCtx)
		if e := appCtx.WebsocketServer().Init(); e != nil {
			panic(e)
		}
	}
}

func initMsgDBServer(appCtx *app.Context) {
	handler.RegisterSaveMsgHandlers(appCtx)
}

func main() {
	flag.Parse()
	config := initConfig()
	gin.SetMode(config.Mode)
	httpEngine := gin.Default()
	appCtx := app.NewAppContext(config, httpEngine)

	serverName := config.Name
	if serverName == "msg_api_server" {
		initMsgApiServer(appCtx)
	} else if serverName == "msg_push_server" {
		initMsgPushServer(appCtx)
	} else if serverName == "msg_db_server" {
		initMsgDBServer(appCtx)
	} else {
		panic(errors.New("server name must be one of { 'msg_api_server', 'msg_push_server', 'msg_db_server'} "))
	}
	appCtx.Start()
}
