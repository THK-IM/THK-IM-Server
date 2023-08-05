package test

import (
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/conf/loader"
	"os"
	"testing"
)

func TestRedisMqBroadcastSubscribe(t *testing.T) {
	dir, err := os.Getwd()
	fmt.Println(dir, err)
	c, e1 := conf.Load("../etc/msg_api_server.yaml")
	if e1 != nil {
		t.Error(e1)
		t.Failed()
		return
	}

	logger := loader.LoadLogg(c.Name, c.Logg)
	rdb := loader.LoadDataCache(c.RedisCache)

	id, _ := loader.LoadNodeId(c, rdb)
	pushMsgSubMap := loader.LoadSubscribers(c.Subscribers, fmt.Sprintf("%d", id), logger, rdb)
	count := 100

	msgChannel := make(chan map[string]interface{}, 0)
	pushMsgSubMap["push_msg"].Sub(func(m map[string]interface{}) error {
		msgChannel <- m
		return nil
	})

	pushMsgPubMap := loader.LoadPublishers(c.Publishes, fmt.Sprintf("%d", id), logger, rdb)
	for i := 0; i < count; i++ {
		msg := map[string]interface{}{"content": fmt.Sprintf("data-%d", i)}
		err := pushMsgPubMap["push_msg"].Pub("", msg)
		if err != nil {
			t.Error(err)
			t.Failed()
		}
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
			if len(messages) == 20 {
				t.Skip()
				break
			}
		}
	}
}
