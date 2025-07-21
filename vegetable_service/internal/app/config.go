package app

import (
	"net"
	"path"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/escoutdoor/vegetable_store/common/pkg/errwrap"
)

type Config struct {
	AppName                  string        `env:"APP_NAME,required"`
	GracefullShutdownTimeout time.Duration `env:"GRACEFULL_SHUTDOWN_TIMEOUT,required"`

	GRPC       GRPC       `envPrefix:"GRPC_SERVER_"`
	Gateway    Gateway    `envPrefix:"GATEWAY_SERVER_"`
	Postgres   Postgres   `envPrefix:"POSTGRES_"`
	Prometheus Prometheus `envPrefix:"PROMETHEUS_SERVER_"`
	Jaeger     Jaeger     `envPrefix:"JAEGER_"`
	Swagger    Swagger    `envPrefix:"SWAGGER_"`
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
	Dsn string `env:"DSN,required"`
}

type Prometheus struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (c *Prometheus) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type Jaeger struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (c *Jaeger) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

type Swagger struct {
	Path     string `env:"PATH,required"`
	FileName string `env:"FILENAME,required"`
}

func (c *Swagger) FilePath() string {
	return path.Join(c.Path, c.FileName)
}

func LoadConfig(path string) (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, errwrap.Wrap("parse env", err)
	}

	return cfg, nil
}
