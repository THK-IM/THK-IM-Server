package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/pkg/model"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
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
	members := make([][]int64, 0)
	// snowNode, err := snowflake.NewNode(1)
	// if err != nil {
	// 	t.Fail()
	// }
	// members = append(members, []int64{snowNode.Generate().Int64()})
	members = append(members, []int64{1696519117500059648, 1696519117500059649})
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()

		sessionAddUserReq := &dto.SessionAddUserReq{
			EntityId: entityIds[index],
			UIds:     members[index],
			Role:     model.SessionMember,
		}
		dataBytes, errJson := json.Marshal(sessionAddUserReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			response, errHttp := client.Post(fmt.Sprintf("%s/%d/user", url, sessionIds[index]), contentType, body)
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

func TestSessionUserUpdate(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696506637453365248)
	userIds := make([][]int64, 0)
	userIds = append(userIds, []int64{9029434827551400299, 9140382960051230588})
	roles := make([]int, 0)
	roles = append(roles, model.SessionSuperAdmin)
	mutes := make([]*int, 0)
	for i := 0; i < len(roles); i++ {
		mute := 0
		mutes = append(mutes, &mute)
	}
	// snowNode, err := snowflake.NewNode(1)
	// if err != nil {
	// 	t.Fail()
	// }
	// members = append(members, []int64{snowNode.Generate().Int64()})
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		sessionUpdateUserReq := &dto.SessionUserUpdateReq{
			SId:  sessionIds[index],
			UIds: userIds[index],
			Role: &roles[index],
			Mute: mutes[index],
		}
		dataBytes, errJson := json.Marshal(sessionUpdateUserReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, errReq := http.NewRequest("PUT", fmt.Sprintf("%s/%d/user", url, sessionIds[index]), body)
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

func TestSessionUserRemove(t *testing.T) {
	uri := "/session"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	sessionIds := make([]int64, 0)
	sessionIds = append(sessionIds, 1696502911586013184)
	successChan := make(chan bool)
	members := make([][]int64, 0)
	members = append(members, []int64{1696519117500059648, 1696519117500059649})

	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		sessionAddUserReq := &dto.SessionDelUserReq{
			UIds: members[index],
		}
		dataBytes, errJson := json.Marshal(sessionAddUserReq)
		if errJson == nil {
			body := bytes.NewReader(dataBytes)
			req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%d/user", url, sessionIds[index]), body)
			req.Header.Set("Content-Type", contentType)
			if err != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, -2, 0, duration, err)
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
