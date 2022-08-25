package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	SSLCertificateKey   string      `yaml:"ssl_certificate_key"` // 	https 需要的密钥
	Location            []*Location `yaml:"location"`
	Schema              string      `yaml:"schema"`
	Port                int         `yaml:"port"`
	SSLCertificate      string      `yaml:"ssl_certificate"` // https 需要的证书
	HealthCheck         bool        `yaml:"health_check"`
	HealthCheckInterval uint        `yaml:"health_check_interval"`
	MaxAllowed          uint        `yaml:"max_allowed"` // 最大并发请求
}

type Location struct {
	Pattern     string   `yaml:"pattern"`
	ProxyPass   []string `yaml:"proxy_pass"`
	BalanceMode string   `yaml:"balance_mode"`
}

// Print 打印配置
func (c *Config) Print() {
	fmt.Printf("Schema: %s\nPort: %d\nHealth Check: %v\nLocation:\n", c.Schema, c.Port, c.HealthCheck)
	for _, l := range c.Location {
		fmt.Printf("\tRoute: %s\n\tProxy Pass: %s\n\tMode: %s\n\n", l.Pattern, l.ProxyPass, l.BalanceMode)
	}
}

// Validation 验证配置
func (c *Config) Validation() error {
	if c.Schema != "http" && c.Schema != "https" {
		return fmt.Errorf("unkonw schemd '%s'", c.Schema)
	}
	// todo 其他字段校验
	return nil
}

// ReadConfig 读取配置 -- 从文件中读取
func ReadConfig(filename string) (*Config, error) {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	// 解析yaml
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
