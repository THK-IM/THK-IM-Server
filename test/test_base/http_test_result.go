package test_base

type HttpTestResult struct {
	serialNo   int   // 请求流水号
	statusCode int   // 状态码
	bodySize   int64 // 响应body大小
	duration   int64 // 耗时 单位ms
}

func (h HttpTestResult) SerialNo() int {
	return h.serialNo
}

func (h HttpTestResult) StatusCode() int {
	return h.statusCode
}

func (h HttpTestResult) BodySize() int64 {
	return h.bodySize
}

func (h HttpTestResult) Duration() int64 {
	return h.duration
}

func NewHttpTestResult(serialNo, statusCode int, bodySize, duration int64) *HttpTestResult {
	return &HttpTestResult{
		serialNo:   serialNo,
		statusCode: statusCode,
		bodySize:   bodySize,
		duration:   duration,
	}
}
