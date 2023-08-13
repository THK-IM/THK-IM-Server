package loader

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/redis/go-redis/v9"
	"os"
	"strings"
	"time"
)

func LoadRedis(source *conf.RedisSource) *redis.Client {
	if source == nil {
		return nil
	}
	endpoint := source.Endpoint
	if strings.HasPrefix(source.Endpoint, "{{") && strings.HasSuffix(source.Endpoint, "}}") {
		endpointEnvKey := strings.Replace(source.Endpoint, "{{", "", -1)
		endpointEnvKey = strings.Replace(endpointEnvKey, "}}", "", -1)
		endpoint = os.Getenv(endpointEnvKey)
	}
	opt, err := redis.ParseURL(fmt.Sprintf("%s%s", endpoint, source.Uri))
	if err != nil {
		panic(err)
	}
	opt.ConnMaxLifetime = time.Duration(source.ConnMaxLifeTime) * time.Second
	opt.ConnMaxIdleTime = time.Duration(source.ConnMaxIdleTime) * time.Second
	opt.ReadTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.WriteTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.PoolTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.MaxIdleConns = source.MaxIdleConn
	opt.PoolSize = source.MaxOpenConn

	rdb := redis.NewClient(opt)
	return rdb
}
