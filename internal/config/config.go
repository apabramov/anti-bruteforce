package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger   LoggerConf
	GrpsServ GrpcServerConf
	Storage  StorageConf
	Redis    RedisConf
	Limit    LimitConf
}

type LoggerConf struct {
	Level string
}

type GrpcServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type string
	Dsn  string
}

type RedisConf struct {
	Host string
	Port string
	Pass string
}

type LimitConf struct {
	LimitLogin int
	LimitPass  int
	LimitIP    int
}

func NewConfig(cfg string) (Config, error) {
	var conf Config
	f, err := os.ReadFile(cfg)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(f), &conf); err != nil {
		return Config{}, err
	}
	return conf, nil
}
