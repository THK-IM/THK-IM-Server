package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func getTestEndPoint() string {
	return "http://127.0.0.1:10000"
}

func TestCreateSingleSession(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 10
	concurrent := 2
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		userId := int64(1)
		entityId := int64(index + 1)
		createSessionReq := &dto.CreateSessionReq{
			Type:     model.SingleSessionType,
			EntityId: nil,
			Members:  []int64{userId, entityId},
		}
		dataBytes, errJson := json.Marshal(createSessionReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			response, errHttp := client.Post(url, contentType, body)
			if errHttp != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			} else {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
			}
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
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

func TestCreateGroupSession(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1000
	concurrent := 10
	successChan := make(chan bool)
	snowNode, err := snowflake.NewNode(1)
	if err != nil {
		t.Fail()
	}
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		entityId := snowNode.Generate().Int64()
		members := make([]int64, 0)
		for i := 0; i < 100; i++ {
			members = append(members, snowNode.Generate().Int64())
		}
		createSessionReq := &dto.CreateSessionReq{
			Type:     model.GroupSessionType,
			EntityId: &entityId,
			Members:  members,
		}
		dataBytes, errJson := json.Marshal(createSessionReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			response, errHttp := client.Post(url, contentType, body)
			if errHttp != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			} else {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
			}
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
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

func TestCreateSuperGroupSession(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1000
	concurrent := 10
	successChan := make(chan bool)
	snowNode, err := snowflake.NewNode(1)
	if err != nil {
		t.Fail()
	}
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		entityId := snowNode.Generate().Int64()
		members := make([]int64, 0)
		for i := 0; i < 100; i++ {
			members = append(members, rand.Int63())
		}
		createSessionReq := &dto.CreateSessionReq{
			Type:     model.SuperGroupSessionType,
			EntityId: &entityId,
			Members:  members,
		}
		dataBytes, errJson := json.Marshal(createSessionReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			response, errHttp := client.Post(url, contentType, body)
			if errHttp != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			} else {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
			}
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
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

func TestUpdateSession(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696506637453365248)
	mutes := make([]*int, 0)
	names := make([]*string, 0)
	remarks := make([]*string, 0)
	for i := 0; i < len(sessionIds); i++ {
		mute := 0
		name := fmt.Sprintf("name:%d", sessionIds[i])
		remark := fmt.Sprintf("remark:%d", sessionIds[i])
		mutes = append(mutes, &mute)
		names = append(names, &name)
		remarks = append(remarks, &remark)
	}
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		updateSessionReq := &dto.UpdateSessionReq{
			Id:     sessionIds[index],
			Mute:   mutes[index],
			Name:   names[index],
			Remark: remarks[index],
		}
		dataBytes, errJson := json.Marshal(updateSessionReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("PUT", fmt.Sprintf("%s/%d", url, sessionIds[index]), body)
			req.Header.Set("Content-Type", contentType)
			if errReq != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, -2, 0, duration, errReq)
			}
			response, errHttp := client.Do(req)
			if errHttp != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			} else {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
			}
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
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
