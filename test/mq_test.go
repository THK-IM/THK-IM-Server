package test

import (
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/conf/loader"
	"github.com/THK-IM/THK-IM-Server/pkg/mq"
	"testing"
	"time"
)

func TestRedisBroadcastSubscribe(t *testing.T) {
	broadcastSubscribe(t, 1000, 5, 10, "../etc/test_redis_mq.yaml")
}

func TestKafkaBroadcastSubscribe(t *testing.T) {
	broadcastSubscribe(t, 1000, 5, 10, "../etc/test_kafka_mq.yaml")
}

func broadcastSubscribe(t *testing.T, msgCount, publisherCount, consumerCount int, configPath string) {
	c, err := conf.Load(configPath)
	if err != nil {
		t.Failed()
		return
	}
	logger := loader.LoadLogg(c.Name, c.Logg)

	msgChannel := make(chan map[string]interface{}, 0)
	subscribers := make([]mq.Subscriber, 0)
	subscribersConsumerCount := make(map[int]int, 0)
	for i := 0; i < consumerCount; i++ {
		subscriberMap := loader.LoadSubscribers(c.MsgQueue.Subscribers, fmt.Sprintf("%d", i), logger)
		subscribers = append(subscribers, subscriberMap["push_msg"])
		subscribersConsumerCount[i] = 0
	}
	for i, subscriber := range subscribers {
		func(index int) {
			subscriber.Sub(func(m map[string]interface{}) error {
				m["index"] = index
				msgChannel <- m
				return nil
			})
		}(i)
	}
	time.Sleep(time.Second) // 等待消费者准备好
	publishers := make([]mq.Publisher, 0)
	for i := 0; i < publisherCount; i++ {
		publisherMap := loader.LoadPublishers(c.MsgQueue.Publishers, fmt.Sprintf("%d", i), logger)
		publishers = append(publishers, publisherMap["push_msg"])
	}
	for _, publisher := range publishers {
		go func(p mq.Publisher) {
			fmt.Println("write start:", time.Now().UnixMilli())
			for i := 0; i < msgCount/publisherCount; i++ {
				msg := map[string]interface{}{"content": fmt.Sprintf("data-%d, %v", i, time.Now())}
				err = p.Pub(fmt.Sprintf("%d", i), msg)
				if err != nil {
					t.Failed()
				}
			}
			fmt.Println("write end:", time.Now().UnixMilli())
		}(publisher)
	}

	messages := make([]map[string]interface{}, 0)
	for {
		select {
		case msg, isOpen := <-msgChannel:
			if !isOpen {
				t.Error(errors.New("channel closed"))
				t.Failed()
			}
			messages = append(messages, msg)
			count := len(messages)
			if count%(msgCount/10) == 0 {
				fmt.Println("consumer count: ", count)
			}
			if count == msgCount*consumerCount {
				for _, m := range messages {
					index := m["index"].(int)
					subscribersConsumerCount[index] += 1
				}
				for k, v := range subscribersConsumerCount {
					fmt.Println("consumer map", k, v)
				}
				t.Skip()
				break
			}
		}
	}
}

func TestRedisGroupSubscribe(t *testing.T) {
	groupSubscribe(t, 10000, 10, 10, "../etc/test_redis_mq.yaml")
}

func TestKafkaGroupSubscribe(t *testing.T) {
	groupSubscribe(t, 10000, 100, 10, "../etc/test_kafka_mq.yaml")
}

func groupSubscribe(t *testing.T, msgCount, publisherCount, consumerCount int, configPath string) {
	c, err := conf.Load(configPath)
	if err != nil {
		t.Error(err)
		t.Failed()
		return
	}

	logger := loader.LoadLogg(c.Name, c.Logg)

	msgChannel := make(chan map[string]interface{}, 0)
	subscribers := make([]mq.Subscriber, 0)
	subscribersConsumerCount := make(map[int]int, 0)
	for i := 0; i < consumerCount; i++ {
		subscriberMap := loader.LoadSubscribers(c.MsgQueue.Subscribers, fmt.Sprintf("%d", i), logger)
		subscribers = append(subscribers, subscriberMap["save_msg"])
		subscribersConsumerCount[i] = 0
	}

	for i, subscriber := range subscribers {
		func(index int, s mq.Subscriber) {
			s.Sub(func(m map[string]interface{}) error {
				m["index"] = index
				msgChannel <- m
				return nil
			})
		}(i, subscriber)
	}

	time.Sleep(30 * time.Second) // 等待重平衡结束

	publishers := make([]mq.Publisher, 0)
	for i := 0; i < publisherCount; i++ {
		publisherMap := loader.LoadPublishers(c.MsgQueue.Publishers, fmt.Sprintf("%d", i), logger)
		publishers = append(publishers, publisherMap["save_msg"])
	}
	for _, publisher := range publishers {
		go func(p mq.Publisher) {
			fmt.Println("write start:", time.Now().UnixMilli())
			for i := 0; i < msgCount/publisherCount; i++ {
				msg := map[string]interface{}{"content": fmt.Sprintf("data-%d, %v", i, time.Now())}
				err = p.Pub(fmt.Sprintf("%d", i), msg)
				if err != nil {
					t.Failed()
				}
			}
			fmt.Println("write end:", time.Now().UnixMilli())
		}(publisher)
	}

	messages := make([]map[string]interface{}, 0)
	for {
		select {
		case msg, isOpen := <-msgChannel:
			if !isOpen {
				t.Error(errors.New("channel closed"))
				t.Failed()
			}
			messages = append(messages, msg)
			count := len(messages)
			if count%(msgCount/10) == 0 {
				fmt.Println("consumer count: ", count)
			}
			if count == msgCount {
				for _, m := range messages {
					index := m["index"].(int)
					subscribersConsumerCount[index] += 1
				}
				for k, v := range subscribersConsumerCount {
					fmt.Println("consumer map", k, v)
				}
				time.Sleep(time.Second)
				t.Skip()
				break
			}
		}
	}
}
