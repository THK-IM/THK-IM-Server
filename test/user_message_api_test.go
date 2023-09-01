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

func TestSendUserMessage(t *testing.T) {
	uri := "/message"
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
		senMessageReq := &dto.SendMessageReq{
			CId:       snowNode.Generate().Int64(),
			SId:       sessionIds[index],
			Type:      1,
			FUid:      1696519117500059648,
			CTime:     time.Now().UnixMilli(),
			Body:      "This is a text message",
			RMsgId:    nil,
			AtUsers:   nil,
			Receivers: nil,
		}
		dataBytes, errJson := json.Marshal(senMessageReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("POST", url, body)
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

func TestDelUserMessage(t *testing.T) {
	uri := "/message"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 605394647632969758)
	msgIds := make([][]int64, 0)
	msgIds = append(msgIds, []int64{1697446767454326784, 1697448240481308672})
	timesFrom := make([]int64, 0)
	timesFrom = append(timesFrom, 1693536903991)
	timesTo := make([]int64, 0)
	timesTo = append(timesTo, 1693539619506)

	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		delMessageReq := &dto.DeleteMessageReq{
			UId:        uIds[index],
			SId:        sessionIds[index],
			MessageIds: msgIds[index],
			TimeFrom:   &timesFrom[index],
			TimeTo:     &timesTo[index],
		}
		dataBytes, errJson := json.Marshal(delMessageReq)
		if errJson == nil {
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

func TestAckUserMessage(t *testing.T) {
	uri := "/message/ack"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 605394647632969758)
	msgIds := make([][]int64, 0)
	msgIds = append(msgIds, []int64{1697442921013317632, 1697446152502251520})
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		ackUserMessageReq := &dto.AckUserMessagesReq{
			UId:    uIds[index],
			SId:    sessionIds[index],
			MsgIds: msgIds[index],
		}
		dataBytes, errJson := json.Marshal(ackUserMessageReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("POST", url, body)
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

func TestReadUserMessage(t *testing.T) {
	uri := "/message/read"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 605394647632969758)
	msgIds := make([][]int64, 0)
	msgIds = append(msgIds, []int64{1697442921013317632, 1697446152502251520})
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		readUserMessageReq := &dto.ReadUserMessageReq{
			UId:    uIds[index],
			SId:    sessionIds[index],
			MsgIds: msgIds[index],
		}
		dataBytes, errJson := json.Marshal(readUserMessageReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("POST", url, body)
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

func TestRevokeUserMessage(t *testing.T) {
	uri := "/message/revoke"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696519117500059648)
	msgIds := make([]int64, 0)
	msgIds = append(msgIds, 1697450778714705920)
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		revokeUserMessageReq := &dto.RevokeUserMessageReq{
			UId:   uIds[index],
			SId:   sessionIds[index],
			MsgId: msgIds[index],
		}
		dataBytes, errJson := json.Marshal(revokeUserMessageReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("POST", url, body)
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

func TestReeditUserMessage(t *testing.T) {
	uri := "/message/reedit"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	uIds := make([]int64, 0)
	uIds = append(uIds, 1696519117500059648)
	msgIds := make([]int64, 0)
	msgIds = append(msgIds, 1697450778714705920)
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		reeditUserMessageReq := &dto.ReeditUserMessageReq{
			UId:     uIds[index],
			SId:     sessionIds[index],
			MsgId:   msgIds[index],
			Content: fmt.Sprintf("Reedit message %d", msgIds[index]),
		}
		dataBytes, errJson := json.Marshal(reeditUserMessageReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("POST", url, body)
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
