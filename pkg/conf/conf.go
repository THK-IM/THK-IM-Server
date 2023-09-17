package conf

import (
	consul "github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
	"os"
)

const (
	DeployBackend = "backend"
	DeployExposed = "exposed"
)

type (
	WebSocket struct {
		Uri           string `yaml:"Uri"`
		MaxClient     int64  `yaml:"MaxClient"`
		MultiPlatform int    `yaml:"MultiPlatform"` // 0:不允许跨平台, -1:随意跨平台, 1:一个平台只能登录一台设备
	}

	Logg struct {
		IndexName   string `yaml:"IndexName"`
		Dir         string `yaml:"Dir"`
		Level       string `yaml:"Level"`
		RetainAge   int    `yaml:"RetainAge"`   // 日志文件保留时间,单位:小时
		RotationAge int    `yaml:"RotationAge"` // 日志文件翻转时间,单位:小时
	}

	Sdk struct {
		Name     string `yaml:"Name"`
		Endpoint string `yaml:"Endpoint"`
	}

	Model struct {
		Name   string `yaml:"Name""`
		Shards int64  `yaml:"Shards"`
	}

	DataSource struct {
		Endpoint        string `yaml:"Endpoint"`
		Uri             string `yaml:"Uri"`
		MaxIdleConn     int    `yaml:"MaxIdleConn"`
		MaxOpenConn     int    `yaml:"MaxOpenConn"`
		ConnMaxLifeTime int64  `yaml:"ConnMaxLifeTime"` // 单位:秒
		ConnMaxIdleTime int64  `yaml:"ConnMaxIdleTime"` // 单位:秒
	}

	RedisSource struct {
		Endpoint        string `yaml:"Endpoint"`
		Uri             string `yaml:"Uri"`
		MaxIdleConn     int    `yaml:"MaxIdleConn"`
		MaxOpenConn     int    `yaml:"MaxOpenConn"`
		ConnMaxLifeTime int64  `yaml:"ConnMaxLifeTime"` // 单位:秒
		ConnMaxIdleTime int64  `yaml:"ConnMaxIdleTime"` // 单位:秒
		ConnTimeout     int64  `yaml:"ConnTimeout"`     // 单位:秒
	}

	ObjectStorage struct {
		Endpoint string `yaml:"Endpoint"`
		Bucket   string `yaml:"Bucket"`
		AK       string `yaml:"AK"`
		SK       string `yaml:"SK"`
		Region   string `yaml:"Region"`
	}

	Metric struct {
		Endpoint     string `yaml:"Endpoint"`
		PushGateway  string `yaml:"PushGateway"`
		PushInterval int64  `yaml:"PushInterval"`
	}

	Node struct {
		MaxCount        int64 `yaml:"MaxCount"`        // 最大工作节点数
		PollingInterval int64 `yaml:"PollingInterval"` // 维持工作节点间隔
	}

	KafkaPublisher struct {
		Brokers    string `yaml:"Brokers"`
		RequireAck int    `yaml:"RequireAck"`
		BatchSize  int    `yaml:"BatchSize"`
		Async      bool   `yaml:"Async"`
	}

	RedisPublisher struct {
		MaxQueueLen int64        `yaml:"MaxQueueLen"`
		RedisSource *RedisSource `yaml:"RedisSource"`
	}

	Publisher struct {
		Topic          string          `yaml:"Topic"`
		KafkaPublisher *KafkaPublisher `yaml:"KafkaPublisher"`
		RedisPublisher *RedisPublisher `yaml:"RedisPublisher"`
	}

	RedisSubscriber struct {
		RedisSource *RedisSource `yaml:"RedisSource"`
		RetryTime   int64        `yaml:"RetryTime"`
	}

	KafkaSubscriber struct {
		Brokers           string `yaml:"Brokers"`
		Partition         int    `yaml:"Partition"`
		NumPartitions     int    `yaml:"NumPartitions"`
		ReplicationFactor int    `yaml:"ReplicationFactor"`
	}

	Subscriber struct {
		Topic           string           `yaml:"Topic"`
		Group           *string          `yaml:"Group"`
		RedisSubscriber *RedisSubscriber `yaml:"RedisSubscriber"`
		KafkaSubscriber *KafkaSubscriber `yaml:"KafkaSubscriber"`
	}

	MsgQueue struct {
		Publishers  []*Publisher  `yaml:"Publishers"`
		Subscribers []*Subscriber `yaml:"Subscribers"`
	}

	IM struct {
		OnlineTimeout       int64 `yaml:"OnlineTimeout"`
		MaxGroupMember      int   `yaml:"MaxGroupMember"`
		MaxSuperGroupMember int   `yaml:"MaxSuperGroupMember"`
	}

	Config struct {
		Name          string         `yaml:"Name"`
		Host          string         `yaml:"Host"`
		Port          string         `yaml:"Port"`
		Mode          string         `yaml:"Mode"`
		DeployMode    string         `yaml:"DeployMode"`
		IpWhiteList   string         `yaml:"IpWhiteList"`
		IM            *IM            `yaml:"IM"`
		WebSocket     *WebSocket     `yaml:"WebSocket"`
		Logg          *Logg          `yaml:"Logg"`
		Sdks          []Sdk          `yaml:"Sdks"`
		Node          *Node          `yaml:"Node"`
		ObjectStorage *ObjectStorage `yaml:"ObjectStorage"`
		DataSource    *DataSource    `yaml:"DataSource"`
		RedisSource   *RedisSource   `yaml:"RedisSource"`
		Models        []Model        `yaml:"Models"`
		Metric        *Metric        `yaml:"Metric"`
		MsgQueue      MsgQueue       `yaml:"MsgQueue"`
	}
)

func Load(f string) (c Config, err error) {
	data, e := os.ReadFile(f)
	if e != nil {
		return c, e
	}
	expanded := os.ExpandEnv(string(data))
	err = yaml.Unmarshal([]byte(expanded), &c)
	return c, err
}

func LoadString(data string) (c Config, err error) {
	expanded := os.ExpandEnv(data)
	err = yaml.Unmarshal([]byte(expanded), &c)
	return c, err
}

func LoadFromConsul(consulAddress, key string) (c Config, err error) {
	config := consul.DefaultConfig()
	config.Address = consulAddress
	client, e1 := consul.NewClient(config)
	if e1 != nil {
		panic(e1)
	}
	pair, _, e2 := client.KV().Get(key, nil)
	if e2 != nil {
		panic(e2)
	} else {
		return LoadString(string(pair.Value))
	}
}
