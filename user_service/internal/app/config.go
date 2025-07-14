package app

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPC     GRPC     `envPrefix:"GRPC_SERVER_"`
	Gateway  Gateway  `envPrefix:"GATEWAY_SERVER_"`
	Postgres Postgres `envPrefix:"POSTGRES_"`
	Token    Token
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

type Token struct {
	AccessTokenSecretKey  string        `env:"ACCESS_TOKEN_SECRET_KEY,required"`
	RefreshTokenSecretKey string        `env:"REFRESH_TOKEN_SECRET_KEY,required"`
	AccessTokenTTL        time.Duration `env:"ACCESS_TOKEN_TTL,required"`
	RefreshTokenTTL       time.Duration `env:"REFRESH_TOKEN_TTL,required"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)

	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	err = env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	return cfg, nil
}
