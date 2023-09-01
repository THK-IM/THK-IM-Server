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

func TestUserSessionQuery(t *testing.T) {
	uri := "/user_session"
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696519117500059648)
	cTimes := make([]int64, 0)
	cTimes = append(cTimes, time.Now().UnixMilli())
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		url := fmt.Sprintf("%s%s/%d/%d", getTestEndPoint(), uri, uIds[index], sessionIds[index])
		req, errReq := http.NewRequest("GET", url, nil)
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

func TestUserLatestSessionQuery(t *testing.T) {
	uri := "/user_session/latest"
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696519117500059648)
	mTimes := make([]int64, 0)
	mTimes = append(mTimes, time.Now().UnixMilli())
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		getLatestSessionReq := dto.GetUserSessionsReq{
			UId:    uIds[index],
			Offset: 0,
			Count:  100,
			MTime:  mTimes[index],
		}
		url := fmt.Sprintf("%s%s?u_id=%d&offset=%d&count=%d&m_time=%d", getTestEndPoint(), uri, getLatestSessionReq.UId,
			getLatestSessionReq.Offset, getLatestSessionReq.Count, getLatestSessionReq.MTime)
		req, errReq := http.NewRequest("GET", url, nil)
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

func TestUserSessionUpdate(t *testing.T) {
	uri := "/user_session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696506637453365248)
	userIds := make([]int64, 0)
	userIds = append(userIds, 9140382960051230588)
	topTimes := make([]int64, 0)
	topTimes = append(topTimes, time.Now().UnixMilli())
	statuses := make([]int, 0)
	statuses = append(statuses, 1)

	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		sessionUpdateUserReq := &dto.UpdateUserSessionReq{
			UId:    userIds[index],
			SId:    sessionIds[index],
			Top:    &topTimes[index],
			Status: &statuses[index],
		}
		dataBytes, errJson := json.Marshal(sessionUpdateUserReq)
		if errJson != nil {
			duration := time.Now().UnixMilli() - startTime
			return test_base.NewHttpTestResult(index, -1, 0, duration, errJson)
		}
		body := bytes.NewReader(dataBytes)
		req, errReq := http.NewRequest("PUT", url, body)
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
