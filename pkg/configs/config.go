package configs

import (
	"time"
)

type BaseConfig struct {
	Name                  string                   `yaml:"name"`
	Env                   string                   `yaml:"env"`
	Debug                 bool                     `yaml:"debug"`
	Logger                *Logger                  `yaml:"logger"`
	Port                  string                   `yaml:"port"`
	Nats                  *Nats                    `yaml:"nats"`
	CockroachDB           *CockroachDB             `yaml:"cockroachdb"`
	JWT                   *JWT                     `yaml:"jwt"`
	RemoteTrace           *RemoteTrace             `yaml:"remoteTrace"`
	StatsEnabled          bool                     `yaml:"statsEnabled"`
	RemoteProfiler        *RemoteProfiler          `yaml:"remoteProfiler"`
	DefaultPondInitDevice *DefaultPondInitDeviceV2 `yaml:"defaultPondInitDevice"`
	Redis                 *Redis                   `yaml:"redis"`
}
type Redis struct {
	Host                string `yaml:"host"`
	Password            Secret `yaml:"password"`
	ReportExpiredSecond int64  `yaml:"reportExpiredSecond"`
}
type CockroachDB struct {
	URI           Secret `yaml:"uri"`
	MigrationPath string `yaml:"migrationPath"`
}

type IoTHub struct {
	MaxIdleConnections uint `yaml:"maxIdleConnections"`
	RequestTimeout     uint `yaml:"requestTimeout"`
	ConnectTimeout     uint `yaml:"connectTimeout"`
}

type JWT struct {
	Expiry        time.Duration `yaml:"expiry"`
	EncryptionKey string        `yaml:"encryptionKey"`
	Audience      string        `yaml:"audience"`
	Issuer        string        `yaml:"issuer"`
}

type RemoteTrace struct {
	Enabled        bool    `yaml:"enabled"`
	TraceAgent     string  `yaml:"traceAgent"`
	TraceCollector string  `yaml:"traceCollector"`
	Ratio          float64 `yaml:"ratio"`
}

type Nats struct {
	ClusterID string `yaml:"clusterId"`
	Address   string `yaml:"address"`
}
type DefaultPondInitDevice struct {
	Feeder      int `yaml:"feeder"`
	PaddleWheel int `yaml:"paddleWheel"`
	SyphonPump  int `yaml:"syphonPump"`
	SumPump     int `yaml:"sumPump"`
	AirBlower   int `yaml:"airBlower"`
}

type DefaultPondInitDeviceV2 struct {
	FeederParams FeederParams                   `yaml:"feederParams"`
	OtherDevices map[string]map[string]Schedule `yaml:"otherDevices"`
}

type Schedule map[string]int

type FeederParams struct {
	FeederSpeed                     int      `yaml:"feederSpeed"`
	FeedingSchedule                 Schedule `yaml:"feedingSchedule"`
	FeedingActiveIntervalInMinute   int      `yaml:"feedingActiveIntervalInMinute"`
	FeedingDeactiveIntervalInMinute int      `yaml:"feedingDeactiveIntervalInMinute"`
	TotalFeedAmount                 int      `yaml:"totalFeedAmount"`
	Semi                            string   `yaml:"semi"`
}

type Logger struct {
	LogLevel string `yaml:"logLvl"`
	LogReq   bool   `yaml:"logReq"`
	LogResp  bool   `yaml:"logResp"`
}

type RemoteProfiler struct {
	Enabled     bool   `yaml:"enabled"`
	ProfilerURL string `yaml:"profilerURL"`
}
