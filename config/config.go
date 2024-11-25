package config

import "time"

type Config struct {
	ServiceName string  `env:"SERVICE_NAME"`
	Env         string  `env:"ENV" envDefault:"local"`
	HTTP        HTTP    `env:"HTTP"`
	Storage     Storage `env:"STORAGE"`
}

type HTTP struct {
	Host    string        `env:"HOST" required:"true"`
	Port    int           `env:"PORT" required:"true"`
	Timeout time.Duration `env:"TIMEOUT" envDefault:"5s"`
}

type Storage struct {
	Path string `env:"PATH" required:"true"`
}
