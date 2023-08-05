package mq

import (
	"context"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	RedisPublisher struct {
		name     string
		stream   string
		group    string
		clientId string
		maxLen   int64
		client   *redis.Client
		logger   *logrus.Entry
	}
)

func (d RedisPublisher) Pub(id string, msg map[string]interface{}) error {
	ctx := context.Background()
	err := d.client.XAdd(ctx, &redis.XAddArgs{
		Stream: d.stream,
		MaxLen: d.maxLen,
		ID:     id,
		Values: msg,
	}).Err()
	return err
}

func NewRedisPublisher(config conf.Mq, clientId string, logger *logrus.Entry, client *redis.Client) Publisher {
	return RedisPublisher{
		name:     config.Name,
		group:    config.Group,
		maxLen:   config.MaxLen,
		stream:   config.Topic,
		clientId: clientId,
		logger:   logger,
		client:   client,
	}
}
