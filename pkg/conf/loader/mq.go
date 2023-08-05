package loader

import (
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/mq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func LoadPublishers(pubConfigs []conf.Mq, clientId string, logger *logrus.Entry, client *redis.Client) map[string]mq.Publisher {
	publisherMap := make(map[string]mq.Publisher, 0)
	for _, pubConfig := range pubConfigs {
		publisherMap[pubConfig.Name] = mq.NewRedisPublisher(pubConfig, clientId, logger, client)
	}
	return publisherMap
}

func LoadSubscribers(subConfigs []conf.Mq, clientId string, logger *logrus.Entry, client *redis.Client) map[string]mq.Subscriber {
	subscriberMap := make(map[string]mq.Subscriber, 0)
	for _, subConfig := range subConfigs {
		subscriberMap[subConfig.Name] = mq.NewRedisSubscribe(subConfig, clientId, logger, client)
	}
	return subscriberMap
}
