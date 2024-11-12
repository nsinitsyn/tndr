package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env             string              `yaml:"env" env-required:"true"`
	GRPC            GRPCConfig          `yaml:"grpc" env-required:"true"`
	DB              DBConfig            `yaml:"db" env-required:"true"`
	ReactionService RemoteServiceConfig `yaml:"reaction_service" env-required:"true"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-required:"true"`
}

type DBConfig struct {
	ConnectionString string `yaml:"connection_string" env-required:"true"`
}

type RemoteServiceConfig struct {
	URL string `yaml:"url" env-required:"true"`
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

// func (config *GRPCConfig) Address() string {
// 	return net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
// }
