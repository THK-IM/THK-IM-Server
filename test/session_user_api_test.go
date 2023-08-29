package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
	"github.com/bwmarrin/snowflake"
	"net/http"
	"testing"
	"time"
)

func TestSessionUserAdd(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	entityIds := make([]int64, 0)
	entityIds = append(entityIds, 1696502911565041665)
	successChan := make(chan bool)
	snowNode, err := snowflake.NewNode(1)
	if err != nil {
		t.Fail()
	}
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		members := make([]int64, 0)
		for i := 0; i < 1; i++ {
			members = append(members, snowNode.Generate().Int64())
		}
		sessionAddUserReq := &dto.SessionAddUserReq{
			EntityId: entityIds[index],
			UIds:     members,
			Role:     model.SessionMember,
		}
		dataBytes, errJson := json.Marshal(sessionAddUserReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			response, errHttp := client.Post(fmt.Sprintf("%s/%d/user", url, sessionIds[index]), contentType, body)
			if errHttp != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration)
			} else {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration)
			}
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration)
		}
	}, func(task *test_base.HttpTestTask) {
		test_base.PrintHttpResults(task)
		for _, result := range task.Results() {
			if result.StatusCode() != http.StatusOK {
				successChan <- false
				return
			}
		}
		successChan <- true
		return
	})
	task.Start()

	for {
		select {
		case success, opened := <-successChan:
			if !opened {
				t.Fail()
			}
			if success {
				t.Skip()
			} else {
				t.Fail()
			}
			return
		}
	}
}
