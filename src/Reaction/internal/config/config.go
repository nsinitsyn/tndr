package config

import (
	"log"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service ServiceConfig `yaml:"service" env-required:"true"`
	GRPC    GRPCConfig    `yaml:"grpc" env-required:"true"`
	HTTP    HTTPConfig    `yaml:"http" env-required:"true"`
	Storage StorageConfig `yaml:"storage" env-required:"true"`
	Tracing TracingConfig `yaml:"tracing" env-required:"true"`
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
	Hostname  string `yaml:"hostname" env-required:"true"`
	Port      int    `yaml:"port" env-required:"true"`
	Namespace string `yaml:"namespace" env-required:"true"`
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
