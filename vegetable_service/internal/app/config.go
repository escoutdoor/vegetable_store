package app

import (
	"net"

	"github.com/caarlos0/env/v11"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPC     GRPC     `envPrefix:"GRPC_SERVER_"`
	Gateway  Gateway  `envPrefix:"GATEWAY_SERVER_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
}

type GRPC struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (c *GRPC) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type Gateway struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (c *Gateway) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type Postgres struct {
	Dsn string `env:"DSN"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := godotenv.Load(path)
	if err != nil {
		return nil, errwrap.Wrap("load config", err)
	}

	err = env.Parse(cfg)
	if err != nil {
		return nil, errwrap.Wrap("parse env", err)
	}

	return cfg, nil
}
