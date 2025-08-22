package config

import (
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

func GetConfig() (*Config, error) {
	var c Config
	content, err := os.ReadFile("/etc/app.yaml")

	if err != nil {
		zap.Error(err)

		return nil, err
	}

	err = yaml.Unmarshal(content, &c)

	return &c, err
}

func GetLogginConfig() *LoggingConfig {
	return nil
}

type SentryConfig struct {
	Dsn   string `yaml:"dsn,omitempty"`
	Debug bool   `yaml:"debug"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password"`
	Database uint   `yaml:"database"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}
