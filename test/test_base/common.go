package test_base

import (
	"fmt"
	"math"
	"net/http"
	"time"
)

var defaultHttpClient = http.Client{
	Transport: &http.Transport{
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   60,
		IdleConnTimeout:       20 * time.Second,
		ResponseHeaderTimeout: time.Second,
		ExpectContinueTimeout: time.Second,
	},
}

func PrintHttpResults(task *HttpTestTask) {
	var maxDuration int64 = 0
	var minDuration int64 = math.MaxInt64
	var totalBodySize int64 = 0
	durationMap := make(map[int64]int64, 0)
	statusCodeMap := make(map[int]int64, 0)
	for _, result := range task.results {
		duration := result.Duration()
		if duration <= 10 {
			durationMap[10]++
		} else if duration <= 50 {
			durationMap[50]++
		} else if duration <= 100 {
			durationMap[50]++
		} else if duration <= 500 {
			durationMap[500]++
		} else if duration <= 1000 {
			durationMap[1000]++
		} else if duration <= 5000 {
			durationMap[5000]++
		} else if duration <= 10000 {
			durationMap[10000]++
		}
		statusCodeMap[result.StatusCode()]++
		if maxDuration < duration {
			maxDuration = duration
		}
		if minDuration >= duration {
			minDuration = duration
		}
		totalBodySize += result.BodySize()
	}
	cost := task.endTime.UnixMilli() - task.startTime.UnixMilli()
	println(fmt.Sprintf("Count: %d, cost: %d ms, max_duartion: %d ms, min_duration: %d ms, total_body_size: %d",
		len(task.results), cost, maxDuration, minDuration, totalBodySize))
	for k, v := range durationMap {
		println(fmt.Sprintf("Duration less than %d ms, count: %d", k, v))
	}
	for k, v := range statusCodeMap {
		println(fmt.Sprintf("StatusCode: %d , count: %d", k, v))
	}
}
