package configs

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/configs"
)

type Config struct {
	configs.BaseConfig `yaml:",inline"`
	Firebase           *Firebase      `yaml:"firebase"`
	SmsBrand           *SmsBrand      `yaml:"smsBrandName"`
	S3Storage          *S3Storage     `yaml:"s3Storage"`
	FCMServerKey       configs.Secret `yaml:"fcmServerKey"`
}

type Firebase struct {
	ConfigPath string `yaml:"configPath"`
}

type SmsBrand struct {
	Token             configs.Secret `yaml:"token"`
	BaseURL           string         `yaml:"baseUrl"`
	Name              string         `yaml:"name"`
	TemplateOTP       string         `yaml:"templateOTP"`
	OTPLiveTime       time.Duration  `yaml:"otpLiveTime"`
	DebugPhoneNumbers []string       `yaml:"debugPhoneNumbers"`
}

type S3Storage struct {
	Region    string `yaml:"region"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Endpoint  string `yaml:"endpoint"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	rawContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rawContent, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
