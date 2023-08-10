package conf

import (
	"gopkg.in/yaml.v3"
	"os"

	consul "github.com/hashicorp/consul/api"
)

type WebSocket struct {
	Uri           string `yaml:"Uri"`
	MaxClient     int64  `yaml:"MaxClient"`
	MultiPlatform int    `yaml:"MultiPlatform"` // 0:不允许跨平台, -1:随意跨平台, 1:一个平台只能登录一台设备
}

type Logg struct {
	IndexName   string `yaml:"IndexName"`
	Dir         string `yaml:"Dir"`
	Level       string `yaml:"Level"`
	RetainAge   int    `yaml:"RetainAge"`   // 日志文件保留时间,单位:小时
	RotationAge int    `yaml:"RotationAge"` // 日志文件翻转时间,单位:小时
}

type Sdk struct {
	Name     string `yaml:"Name"`
	Endpoint string `yaml:"Endpoint"`
}

type Model struct {
	Name   string `yaml:"Name""`
	Shards int64  `yaml:"Shards"`
}

type Database struct {
	Endpoint        string `yaml:"Endpoint"`
	Uri             string `yaml:"Uri"`
	MaxIdleConn     int    `yaml:"MaxIdleConn"`
	MaxOpenConn     int    `yaml:"MaxOpenConn"`
	ConnMaxLifeTime int64  `yaml:"ConnMaxLifeTime"` // 单位:秒
	ConnMaxIdleTime int64  `yaml:"ConnMaxIdleTime"` // 单位:秒
}

type RedisCache struct {
	Endpoint        string `yaml:"Endpoint"`
	Uri             string `yaml:"Uri"`
	MaxIdleConn     int    `yaml:"MaxIdleConn"`
	MaxOpenConn     int    `yaml:"MaxOpenConn"`
	ConnMaxLifeTime int64  `yaml:"ConnMaxLifeTime"` // 单位:秒
	ConnMaxIdleTime int64  `yaml:"ConnMaxIdleTime"` // 单位:秒
	ConnTimeout     int64  `yaml:"ConnTimeout"`     // 单位:秒
}

type Metric struct {
	Endpoint     string `yaml:"Endpoint"`
	PushGateway  string `yaml:"PushGateway"`
	PushInterval int64  `yaml:"PushInterval"`
}

type Node struct {
	MaxCount        int64 `yaml:"MaxCount"`        // 最大工作节点数
	PollingInterval int64 `yaml:"PollingInterval"` // 维持工作节点间隔
}

type Mq struct {
	Engine    string `yaml:"Engine"`
	Name      string `yaml:"Name"`
	Topic     string `yaml:"Topic"`
	Group     string `yaml:"Group"`
	MaxLen    int64  `yaml:"MaxLen"`
	RetryTime int64  `yaml:"RetryTime"`
}

type IM struct {
	MaxGroupMember      int `yaml:"MaxGroupMember"`
	MaxSuperGroupMember int `yaml:"MaxSuperGroupMember"`
}

type Config struct {
	Name          string      `yaml:"Name"`
	Host          string      `yaml:"Host"`
	Port          string      `yaml:"Port"`
	Mode          string      `yaml:"Mode"`
	OnlineTimeout int64       `yaml:"OnlineTimeout"` // 单位 秒
	IM            *IM         `yaml:"IM"`
	WebSocket     *WebSocket  `yaml:"WebSocket"`
	Logg          *Logg       `yaml:"Logg"`
	Sdks          []Sdk       `yaml:"Sdks"`
	Node          *Node       `yaml:"Node"`
	Database      *Database   `yaml:"DataSource"`
	RedisCache    *RedisCache `yaml:"RedisCache"`
	Models        []Model     `yaml:"Models"`
	Metric        *Metric     `yaml:"Metric"`
	Publishes     []Mq        `yaml:"Publishes"`
	Subscribers   []Mq        `yaml:"Subscribers"`
}

func Load(f string) (c Config, err error) {
	data, e := os.ReadFile(f)
	if e != nil {
		return c, e
	}
	err = yaml.Unmarshal(data, &c)
	return c, err
}

func LoadString(data string) (c Config, err error) {
	err = yaml.Unmarshal([]byte(data), &c)
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
