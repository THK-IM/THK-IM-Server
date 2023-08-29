package test_base

import (
	"net/http"
	"time"
)

type RequestFunc func(index, channelIndex int, client http.Client) *HttpTestResult
type ResultFunc func(task *HttpTestTask)

type HttpTestTask struct {
	count         int // 总测试次数
	concurrent    int // 测试并发数
	startTime     time.Time
	endTime       time.Time
	results       []*HttpTestResult
	resultChan    chan *HttpTestResult
	taskIndexChan chan int
	requestFunc   RequestFunc
	resultFunc    ResultFunc
}

func (t *HttpTestTask) StartTime() time.Time {
	return t.startTime
}

func (t *HttpTestTask) EndTime() time.Time {
	return t.endTime
}

func (t *HttpTestTask) Results() []*HttpTestResult {
	return t.results
}

func NewHttpTestTask(count, concurrent int, request RequestFunc, result ResultFunc) *HttpTestTask {
	return &HttpTestTask{
		count:         count,
		concurrent:    concurrent,
		results:       make([]*HttpTestResult, 0),
		resultChan:    make(chan *HttpTestResult),
		taskIndexChan: make(chan int),
		requestFunc:   request,
		resultFunc:    result,
	}
}

func (t *HttpTestTask) Start() {
	// 准备接受测试指标协程启动
	t.prepare()

	t.startTime = time.Now()
	// 开始往channel发送数据进行压测
	go func() {
		for i := 0; i < t.count; i++ {
			t.taskIndexChan <- i
		}
	}()
}

func (t *HttpTestTask) prepare() {
	// 创建并发协程
	for i := 0; i < t.concurrent; i++ {
		t.startGoroutines(i)
	}
	go func() {
		for {
			select {
			case httpResult, opened := <-t.resultChan:
				if !opened {
					return
				}
				t.results = append(t.results, httpResult)
				if len(t.results) >= t.count {
					t.endTime = time.Now()
					close(t.taskIndexChan)
					close(t.resultChan)
					if t.resultFunc != nil {
						t.resultFunc(t)
					}
				}
			}
		}
	}()
}

func (t *HttpTestTask) startGoroutines(routineIndex int) {
	go func() {
		for {
			select {
			case index, opened := <-t.taskIndexChan:
				if !opened {
					return
				}
				if t.requestFunc != nil {
					httpResult := t.requestFunc(index, routineIndex, defaultHttpClient)
					t.resultChan <- httpResult
				}
			}
		}
	}()
}
