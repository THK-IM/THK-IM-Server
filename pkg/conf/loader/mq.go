package loader

import (
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	mq2 "github.com/THK-IM/THK-IM-Server/pkg/service/mq"
	"github.com/sirupsen/logrus"
)

func LoadPublishers(pubConfigs []*conf.Publisher, clientId string, logger *logrus.Entry) map[string]mq2.Publisher {
	publisherMap := make(map[string]mq2.Publisher, 0)
	for _, pubConfig := range pubConfigs {
		if pubConfig.RedisPublisher != nil {
			client := LoadRedis(pubConfig.RedisPublisher.RedisSource)
			publisherMap[pubConfig.Topic] = mq2.NewRedisPublisher(pubConfig, clientId, logger, client)
		} else if pubConfig.KafkaPublisher != nil {
			publisherMap[pubConfig.Topic] = mq2.NewKafkaPublisher(pubConfig, clientId, logger)
		}
	}
	return publisherMap
}

func LoadSubscribers(subConfigs []*conf.Subscriber, clientId string, logger *logrus.Entry) map[string]mq2.Subscriber {
	subscriberMap := make(map[string]mq2.Subscriber, 0)
	for _, subConfig := range subConfigs {
		if subConfig.RedisSubscriber != nil {
			client := LoadRedis(subConfig.RedisSubscriber.RedisSource)
			subscriberMap[subConfig.Topic] = mq2.NewRedisSubscribe(subConfig, clientId, logger, client)
		} else if subConfig.KafkaSubscriber != nil {
			subscriberMap[subConfig.Topic] = mq2.NewKafkaSubscriber(subConfig, clientId, logger)
		}
	}
	return subscriberMap
}
