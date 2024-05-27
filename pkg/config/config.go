package config

import (
	"log"
	"sync"

	"github.com/BurntSushi/toml"
)

// Config ...
type Config struct {
	ServiceId               string `toml:"service_id"`
	Service                 string `toml:"service"`
	BindAddr                string `toml:"bind_addr"`
	LogLevel                string `toml:"log_level"`
	DatabaseURL             string `toml:"database_url"`
	TokenAuthSecurityKey    string `toml:"token_auth_secret"`
	TokenRefreshSecurityKey string `toml:"token_refresh_secret"`
	StaticPath              string `toml:"static_path"`
	FilesPath               string `toml:"files_path"`
	FilesUrlPrefix          string `toml:"files_url_prefix"`
	IsDevelopment           bool   `toml:"is_development"`
	GuidePath               string `toml:"guide_path"`
	SSL                     struct {
		CertPath string `toml:"cert_path"`
		KeyPath  string `toml:"key_path"`
	} `toml:"ssl"`
	Telegram struct {
		ChatId string `toml:"chat_id"`
		BotId  string `toml:"bot_id"`
	} `toml:"telegram"`
	RabbitMQ struct {
		Url       string `toml:"url"`
		QueueName string `toml:"queue_name"`
	} `toml:"rabbitmq"`
	Database struct {
		MainServer  string   `toml:"main"`
		ReadServers []string `toml:"read"`
	} `toml:"database"`
}

const configPath = "configs/apiserver.toml"

var (
	singleInstance *Config
	lock           = &sync.Mutex{}
)

// NewConfig ...
func newConfig() *Config {
	config := &Config{
		BindAddr:      ":8081",
		LogLevel:      "debug",
		IsDevelopment: false,
	}
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func GetInstance() *Config {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = newConfig()
		}
	}
	return singleInstance
}
