package config

import (
	"fmt"

	wbfconf "github.com/wb-go/wbf/config"
)

type Server struct {
	Host string
	Port int
}

type DB struct {
	DSN      string
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	Server Server
	DB     DB
}

func Load(path string) (*Config, error) {
	v := wbfconf.New()
	if err := v.Load(path); err != nil {
		return nil, err
	}
	cfg := &Config{
		Server: Server{
			Host: v.GetString("server.host"),
			Port: v.GetInt("server.port"),
		},
		DB: DB{
			DSN:      v.GetString("database.dsn"),
			Host:     v.GetString("database.host"),
			Port:     v.GetInt("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			Name:     v.GetString("database.name"),
			SSLMode:  v.GetString("database.sslmode"),
		},
	}
	return cfg, nil
}

func (d DB) DSNString() string {
	if d.DSN != "" {
		return d.DSN
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}
