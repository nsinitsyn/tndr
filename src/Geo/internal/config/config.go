package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service         ServiceConfig       `yaml:"service" env-required:"true"`
	GRPC            GRPCConfig          `yaml:"grpc" env-required:"true"`
	HTTP            HTTPConfig          `yaml:"http" env-required:"true"`
	Storage         StorageConfig       `yaml:"storage" env-required:"true"`
	ReactionService RemoteServiceConfig `yaml:"reaction_service" env-required:"true"`
	Messaging       MessagingConfig     `yaml:"messaging" env-required:"true"`
	Tracing         TracingConfig       `yaml:"tracing" env-required:"true"`
}

type ServiceConfig struct {
	Env        string `yaml:"env" env-required:"true"`
	Name       string `yaml:"name" env-required:"true"`
	Version    string `yaml:"version" env-required:"true"`
	InstanceID string `yaml:"instance" env-required:"true"`
}

type GRPCConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type HTTPConfig struct {
	Port int `yaml:"port" env-required:"true"`
}

type StorageConfig struct {
	Addr     string `yaml:"addr" env-required:"true"`
	Password string `yaml:"password" env-default:""`
	DB       int    `yaml:"db" env-default:"0"`
}

type RemoteServiceConfig struct {
	Endpoint string `yaml:"endpoint" env-required:"true"`
}

type MessagingConfig struct {
	Servers string `yaml:"servers" env-required:"true"`
	Group   string `yaml:"group" env-required:"true"`
	Topic   string `yaml:"topic" env-required:"true"`
}

type TracingConfig struct {
	Enabled  bool   `yaml:"enabled" env-default:"false"`
	Endpoint string `yaml:"endpoint" env-default:""`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		path, ok := os.LookupEnv("CONFIG_PATH")
		if !ok {
			path = "config.yaml"
		}

		instance = &Config{}
		if err := cleanenv.ReadConfig(path, instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
