package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
	"net/http"
	"testing"
	"time"
)

func TestSessionMessageQuery(t *testing.T) {
	uri := "/session"
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	cTimes := make([]int64, 0)
	cTimes = append(cTimes, time.Now().UnixMilli())
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		getSessionMessageReq := &dto.GetSessionMessageReq{
			SId:    sessionIds[index],
			CTime:  cTimes[index],
			Offset: 0,
			Count:  100,
		}
		url := fmt.Sprintf("%s%s/%d/message", getTestEndPoint(), uri, sessionIds[index])
		query := fmt.Sprintf("?s_id=%d&offset=%d&count=%d&c_time=%d", getSessionMessageReq.SId, getSessionMessageReq.Offset,
			getSessionMessageReq.Count, getSessionMessageReq.CTime)
		req, errReq := http.NewRequest("GET", fmt.Sprintf("%s%s", url, query), nil)
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

func TestSessionMessageDelete(t *testing.T) {
	uri := "/session"
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	msgIds := make([][]int64, 0)
	msgIds = append(msgIds, []int64{1697518819884404736, 1697518185135214592})
	cTimes := make([]int64, 0)
	cTimes = append(cTimes, time.Now().UnixMilli())
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		delSessionMessageReq := &dto.DelSessionMessageReq{
			SId:      sessionIds[index],
			MsgIds:   msgIds[index],
			TimeFrom: 0,
			TimeTo:   cTimes[index],
		}
		dataBytes, errJson := json.Marshal(delSessionMessageReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -2, 0, duration, errJson)
		}
		url := fmt.Sprintf("%s%s/%d/message", getTestEndPoint(), uri, sessionIds[index])
		body := bytes.NewReader(dataBytes)
		req, errReq := http.NewRequest("DELETE", url, body)
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
