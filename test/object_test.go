package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/dto"
	"github.com/THK-IM/THK-IM-Server/test/test_base"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestUpload(t *testing.T) {
	uri := "/object/upload_params"
	url := fmt.Sprintf("%s%s", getTestEndPoint(), uri)
	contentType := "application/json"
	count := 1
	concurrent := 1
	successChan := make(chan bool)
	task := test_base.NewHttpTestTask(count, concurrent, func(index, channelIndex int, client http.Client) *test_base.HttpTestResult {
		startTime := time.Now().UnixMilli()
		reqParams := &dto.GetUploadParamsReq{
			SId:      1696502911665704960,
			UId:      1874068156324778273,
			FileName: "image1.png",
		}
		query := fmt.Sprintf("?s_id=%d&u_id=%d&fn=%s", reqParams.SId, reqParams.UId, reqParams.FileName)
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
			resBytes, err := io.ReadAll(response.Body)
			if err != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			}
			getUploadRes := &dto.GetUploadParamsRes{}
			err = json.Unmarshal(resBytes, getUploadRes)
			if err != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
			}
			bodyBuffer := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuffer)
			for k, v := range getUploadRes.Params {
				w, errParams := bodyWriter.CreateFormField(k)
				if errParams != nil {
					duration := time.Now().UnixMilli() - startTime
					return test_base.NewHttpTestResult(index, 500, 0, duration, errParams)
				}
				_, errParams = w.Write([]byte(v))
				if errParams != nil {
					duration := time.Now().UnixMilli() - startTime
					return test_base.NewHttpTestResult(index, 500, 0, duration, errParams)
				}
			}
			w, errParams := bodyWriter.CreateFormFile("file", "image1.png")
			if errParams != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errParams)
			}
			readBytes, errRead := os.ReadFile("./media/image1.png")
			if errRead != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errRead)
			}
			_, errRead = w.Write(readBytes)
			if errRead != nil {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errRead)
			}
			contentType = bodyWriter.FormDataContentType()
			defer bodyWriter.Close()
			if getUploadRes.Method != "POST" {
				duration := time.Now().UnixMilli() - startTime
				return test_base.NewHttpTestResult(index, 500, 0, duration, errors.New("error http Method"))
			} else {
				postRes, errPost := http.Post(getUploadRes.Url, contentType, bodyBuffer)
				if errPost != nil {
					duration := time.Now().UnixMilli() - startTime
					return test_base.NewHttpTestResult(index, 500, 0, duration, errHttp)
				} else {
					duration := time.Now().UnixMilli() - startTime
					if postRes.StatusCode == 200 {
						return test_base.NewHttpTestResult(index, 200, 0, duration, nil)
					} else {
						resBytes, err = io.ReadAll(postRes.Body)
						fmt.Println(resBytes)
						return test_base.NewHttpTestResult(index, postRes.StatusCode, 0, duration, errors.New(string(resBytes)))
					}
				}
			}
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
