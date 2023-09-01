package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
	"github.com/bwmarrin/snowflake"
	"net/http"
	"testing"
	"time"
)

func TestUserOnlineQuery(t *testing.T) {
	uri := "/system/user/online"
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696502911565041665)

	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		getUsersOnlineStatusReq := &dto.GetUsersOnlineStatusReq{
			UIds: uIds,
		}
		url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
		dataBytes, errJson := json.Marshal(getUsersOnlineStatusReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		req, errReq := http.NewRequest("GET", url, body)
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

func TestUserOnlinePost(t *testing.T) {
	uri := "/system/user/online"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696502911565041665)
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		postUserOnline := &dto.PostUserOnlineReq{
			NodeId: 1,
			ConnId: 1,
			Online: true,
			UId:    uIds[index],
		}
		dataBytes, errJson := json.Marshal(postUserOnline)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		response, errHttp := client.Post(url, contentType, body)
		if errHttp != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
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

func TestKickoffUser(t *testing.T) {
	uri := "/system/user/kickoff"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696502911565041665)
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		kickoffReq := &dto.KickUserReq{UId: uIds[index]}
		dataBytes, errJson := json.Marshal(kickoffReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		response, errHttp := client.Post(url, contentType, body)
		if errHttp != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
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

func TestSendSystemMessage(t *testing.T) {
	uri := "/system/message/send"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	snowNode, err := snowflake.NewNode(1)
	if err != nil {
		t.Fail()
	}
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		sendMessageReq := &dto.SendMessageReq{
			CId:       snowNode.Generate().Int64(),
			SId:       sessionIds[index],
			Type:      1,
			FUid:      0,
			CTime:     time.Now().UnixMilli(),
			Body:      "This is a system message",
			RMsgId:    nil,
			AtUsers:   nil,
			Receivers: nil,
		}
		dataBytes, errJson := json.Marshal(sendMessageReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		response, errHttp := client.Post(url, contentType, body)
		if errHttp != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
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

func TestPushSystemMessage(t *testing.T) {
	uri := "/system/message/push"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	uIds := make([][]int64, 0)
	uIds = append(uIds, []int64{894385949183117216})

	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		pushMessageReq := &dto.PushMessageReq{
			UIds:        uIds[index],
			Type:        0,
			SubType:     11,
			Body:        "11111",
			OfflinePush: false,
		}
		dataBytes, errJson := json.Marshal(pushMessageReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		response, errHttp := client.Post(url, contentType, body)
		if errHttp != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
		} else {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, response.StatusCode, response.ContentLength, duration, nil)
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
