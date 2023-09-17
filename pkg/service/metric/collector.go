package metric

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type Collector struct {
	p                     *Prometheus
	nodeId                int64
	onlineUserGaugeMetric *Metric
	pushMsgCntMetric      *Metric
	pushMsgBytesSumMetric *Metric
	pushMsgDurationMetric *Metric
	rcvMsgCntMetric       *Metric
	rcvMsgBytesSumMetric  *Metric
	rcvMsgDurationMetric  *Metric
	mqMsgCntMetric        *Metric
	mqMsgDurationMetric   *Metric
}

func (c Collector) OnLineUser(totalCnt, authedCnt int64) {
	c.onlineUserGaugeMetric.MetricCollector.(*prometheus.GaugeVec).WithLabelValues("total_count").Set(float64(totalCnt))
	c.onlineUserGaugeMetric.MetricCollector.(*prometheus.GaugeVec).WithLabelValues("authed_count").Set(float64(authedCnt))
}

func (c Collector) PushMsg(status, types string, size int64, millSeconds int64) {
	fmt.Println("PushMsg", status, types, size, millSeconds)
	c.pushMsgCntMetric.MetricCollector.(*prometheus.CounterVec).WithLabelValues(status, types).Inc()
	c.pushMsgDurationMetric.MetricCollector.(*prometheus.HistogramVec).WithLabelValues(status, types).Observe(float64(millSeconds))
	if size > 0 {
		c.pushMsgBytesSumMetric.MetricCollector.(*prometheus.SummaryVec).WithLabelValues(types).Observe(float64(size))
	}
}

func (c Collector) RcvClientMsg(status, types string, size, millSeconds int64) {
	fmt.Println("RcvClientMsg", status, types, size, millSeconds)
	c.rcvMsgCntMetric.MetricCollector.(*prometheus.CounterVec).WithLabelValues(status, types).Inc()
	c.rcvMsgDurationMetric.MetricCollector.(*prometheus.HistogramVec).WithLabelValues(status, types).Observe(float64(millSeconds))
	if size > 0 {
		c.rcvMsgBytesSumMetric.MetricCollector.(*prometheus.SummaryVec).WithLabelValues(types).Observe(float64(size))
	}
}

func (c Collector) RcvMqMsg(status, types string, millSeconds int64) {
	fmt.Println("RcvMqMsg", status, types, millSeconds)
	c.mqMsgCntMetric.MetricCollector.(*prometheus.CounterVec).WithLabelValues(status, types).Inc()
	c.mqMsgDurationMetric.MetricCollector.(*prometheus.HistogramVec).WithLabelValues(status, types).Observe(float64(millSeconds) / 1000)
}

func (c Collector) getNodeId() string {
	return fmt.Sprintf("node-%d", c.nodeId)
}

func NewCollector(serverName string, nodeId int64, source *conf.Metric, log *logrus.Entry, httpEngine *gin.Engine) *Collector {
	onlineUserGaugeMetric := &Metric{
		ID:          "onlineUserGauge",
		Name:        "online_user_gauge",
		Description: "How many users are online in realtime, partitioned by work_node_id status",
		Type:        "gauge_vec",
		Args:        []string{"status"},
	}
	pushMsgCntMetric := &Metric{
		ID:          "pushMessageCnt",
		Name:        "push_message_count",
		Description: "How many messages pushed to client, partitioned by work_node_id status and type.",
		Type:        "counter_vec",
		Args:        []string{"status", "type"},
	}
	pushMsgBytesSumMetric := &Metric{
		ID:          "pushMsgBytes",
		Name:        "push_message_size",
		Description: "How many bytes pushed to client, partitioned by work_node_id and type",
		Type:        "summary_vec",
		Args:        []string{"type"},
	}
	pushMsgDurationMetric := &Metric{
		ID:          "PushMsgDur",
		Name:        "push_message_duration",
		Description: "The message be pushed latencies in seconds, partitioned by work_node_id, status and type",
		Type:        "histogram_vec",
		Args:        []string{"status", "type"},
	}
	rcvMsgCntMetric := &Metric{
		ID:          "rcvClientMessageCnt",
		Name:        "rcv_message_count",
		Description: "How many client messages received, partitioned by work_node_id, status and type.",
		Type:        "counter_vec",
		Args:        []string{"status", "type"},
	}
	rcvMsgBytesSumMetric := &Metric{
		ID:          "rcvClientMsgBytes",
		Name:        "rcv_message_size",
		Description: "How many bytes received from client, partitioned by work_node_id and type",
		Type:        "summary_vec",
		Args:        []string{"type"},
	}
	rcvMsgDurationMetric := &Metric{
		ID:          "RcvClientMsgDur",
		Name:        "rcv_message_duration",
		Description: "The client message be processed latencies in seconds, partitioned by work_node_id, status and type",
		Type:        "histogram_vec",
		Args:        []string{"status", "type"},
	}
	mqMsgCntMetric := &Metric{
		ID:          "mqMessageCnt",
		Name:        "mq_message_count",
		Description: "How many mq messages received, partitioned by work_node_id status and type.",
		Type:        "counter_vec",
		Args:        []string{"status", "type"},
	}
	mqMsgDurationMetric := &Metric{
		ID:          "mqMsgDur",
		Name:        "mq_message_duration",
		Description: "The mq message be processed latencies in seconds, partitioned by work_node_id status and type.",
		Type:        "histogram_vec",
		Args:        []string{"status", "type"},
	}

	metrics := []*Metric{onlineUserGaugeMetric, pushMsgCntMetric, pushMsgDurationMetric, pushMsgBytesSumMetric,
		rcvMsgCntMetric, rcvMsgBytesSumMetric, rcvMsgDurationMetric, mqMsgCntMetric, mqMsgDurationMetric}
	p := NewPrometheus(serverName, nodeId, log, metrics)
	p.Use(httpEngine)
	if !strings.EqualFold(source.PushGateway, "") {
		p.SetPushGatewayJob(serverName)
		p.SetPushGateway(
			source.PushGateway,
			source.Endpoint,
			time.Duration(source.PushInterval),
		)
	}
	return &Collector{
		p:                     p,
		nodeId:                nodeId,
		onlineUserGaugeMetric: onlineUserGaugeMetric,
		pushMsgCntMetric:      pushMsgCntMetric,
		pushMsgDurationMetric: pushMsgDurationMetric,
		pushMsgBytesSumMetric: pushMsgBytesSumMetric,
		rcvMsgCntMetric:       rcvMsgCntMetric,
		rcvMsgBytesSumMetric:  rcvMsgBytesSumMetric,
		rcvMsgDurationMetric:  rcvMsgDurationMetric,
		mqMsgCntMetric:        mqMsgCntMetric,
		mqMsgDurationMetric:   mqMsgDurationMetric,
	}
}
