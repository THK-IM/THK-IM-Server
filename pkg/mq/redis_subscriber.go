package mq

import (
	"context"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	redisSubscriber struct {
		name      string
		stream    string
		group     string
		clientId  string
		maxLen    int64
		retryTime int64
		client    *redis.Client
		logger    *logrus.Entry
	}
)

func (d redisSubscriber) Sub(
	onMessageReceived OnMessageReceived) {
	if d.group == "" {
		d.subscribe(onMessageReceived)
	} else {
		d.groupSubscribe(onMessageReceived)
	}
}

func (d redisSubscriber) subscribe(onMessageReceived OnMessageReceived) {
	go func() {
		ctx := context.Background()
		lastId := "$"
		for {
			if entries, err := d.client.XRead(ctx, &redis.XReadArgs{
				Streams: []string{d.stream, lastId},
				Count:   10,
				Block:   0,
			}).Result(); err != nil {
				d.logger.Error(err)
			} else {
				d.consumeXStreams(entries, onMessageReceived)
				if len(entries) > 0 && len(entries[0].Messages) > 0 {
					lastId = entries[0].Messages[len(entries[0].Messages)-1].ID
				}
			}
		}
	}()
}

func (d redisSubscriber) infoGroups() ([]redis.XInfoGroup, error) {
	ctx := context.Background()
	return d.client.XInfoGroups(ctx, d.stream).Result()
}

func (d redisSubscriber) createGroupIfNeeded() {
	groups, _ := d.infoGroups()
	existed := false
	for _, group := range groups {
		if group.Name == d.group {
			existed = true
			break
		}
	}
	if !existed {
		ctx := context.Background()
		err := d.client.XGroupCreateMkStream(ctx, d.stream, d.group, "0").Err()
		if err != nil {
			d.logger.Error(err)
			panic(err)
		}
	}
}

func (d redisSubscriber) groupSubscribe(onMessageReceived OnMessageReceived) {
	d.createGroupIfNeeded()
	ctx := context.Background()
	lastId := ">"
	go func() {
		for {
			if entries, err := d.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    d.group,
				Consumer: d.clientId,
				Streams:  []string{d.stream, lastId},
				Count:    1,
				Block:    0,
				NoAck:    false,
			}).Result(); err != nil {
				d.logger.Error(err)
			} else {
				d.consumeXStreams(entries, onMessageReceived)
			}
		}
	}()
	if d.retryTime > 0 {
		d.pendingMessage(onMessageReceived)
	}
}

func (d redisSubscriber) pendingMessage(onMessageReceived OnMessageReceived) {
	ctx := context.Background()
	ticker := time.NewTicker(time.Minute * time.Duration(d.retryTime))
	go func() {
		for range ticker.C {
			if pendingInfos, err := d.client.XPendingExt(ctx, &redis.XPendingExtArgs{
				Stream:   d.stream,
				Group:    d.group,
				Consumer: d.clientId,
				Start:    "-",
				End:      "+",
				Count:    10,
			}).Result(); err != nil {
				d.logger.Error(err)
			} else {
				for _, pendingInfo := range pendingInfos {
					d.retryConsume(pendingInfo.ID, onMessageReceived)
				}
			}
		}
	}()
}

func (d redisSubscriber) retryConsume(id string, onMessageReceived OnMessageReceived) {
	ctx := context.Background()
	if messages, err := d.client.XRangeN(ctx, d.stream, id, "+", 1).Result(); err != nil {
		d.logger.Error(err)
	} else {
		d.consumeMessages(messages, onMessageReceived)
	}
}

func (d redisSubscriber) consumeMessages(messages []redis.XMessage, onMessageReceived OnMessageReceived) {
	ctx := context.Background()
	for _, msg := range messages {
		// messageID := msg.ID
		// eventDescription := fmt.Sprintf("%v", values["whatHappened"])
		// ticketID := fmt.Sprintf("%v", values["ticketID"])
		// ticketData := fmt.Sprintf("%v", values["ticketData"])
		if err := onMessageReceived(msg.Values); err == nil {
			if d.group != "" {
				d.client.XAck(ctx, d.stream, d.group, msg.ID)
			}
		} else {
			d.logger.Error(
				fmt.Sprintf("subscribe group: %s, client id: %s, msgId: %s, values: %v",
					d.group, d.clientId, msg.ID, msg.Values,
				),
			)
		}
	}
}

func (d redisSubscriber) consumeXStreams(entries []redis.XStream, onMessageReceived OnMessageReceived) {
	for _, entry := range entries {
		d.consumeMessages(entry.Messages, onMessageReceived)
	}
}

func NewRedisSubscribe(config conf.Mq, clientId string, logger *logrus.Entry, client *redis.Client) Subscriber {
	return redisSubscriber{
		name:      config.Name,
		group:     config.Group,
		maxLen:    config.MaxLen,
		stream:    config.Topic,
		retryTime: config.RetryTime,
		clientId:  clientId,
		logger:    logger,
		client:    client,
	}
}
