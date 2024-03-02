package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

const envfile = "dev/.env"

// Config represents the application configuration.
type Config struct {
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn, debug, info" env:"LOG_LEVEL" default:"info"`

	Debug   bool   `long:"debug" description:"Developer mode" env:"DEBUG"`
	PathLog string `long:"path_log" description:"Path log" env:"PATH_LOG" default:"stdout"`

	AppInfo struct {
		Name    string `long:"name" description:"App name" env:"APP_NAME" required:"true" default:"default app"`
		Version string `long:"version" description:"App version" env:"APP_VERSION" required:"true" default:"0.0.1"`
	}

	Nats struct {
		Host      string `long:"nats_host" description:"Host Nats" env:"NATS_HOST" required:"true" default:"0.0.0.0"`
		Port      int    `long:"nats_port" description:"Port Nats" env:"NATS_PORT" required:"true" default:"4222"`
		ClusterID string `long:"nats_cluster_id" description:"Nats cluster id" env:"NATS_CLUSTER_ID" required:"true" default:"nats-cluster"`
		Client1ID string `long:"nats_client_1_id" description:"Nats client id" env:"NATS_CLIENT_1_ID" required:"true" default:"nats-client-1"`
		Client2ID string `long:"nats_client_2_id" description:"Nats client id" env:"NATS_CLIENT_2_ID" required:"true" default:"nats-client-2"`
		Subject   string `long:"nats_subject" description:"Nats subject" env:"NATS_SUBJECT" required:"true" default:"test-subject"`
	}

	HttpServer struct {
		Host string `long:"http_host" description:"Host HTTP server" env:"HTTP_HOST" required:"true" default:"0.0.0.0"`
		Port int    `long:"http_port" description:"Post HTTP sever" env:"HTTP_PORT" required:"true" default:"80"`
	}

	DB struct {
		Host     string `long:"db_host" description:"Host DB" env:"DB_HOST" required:"true" default:"127.0.0.1"`
		Port     int    `long:"db_port" description:"Port DB" env:"DB_PORT" required:"true" default:"5432"`
		Name     string `long:"db_name" description:"Name DB" env:"DB_NAME" required:"true" default:"db"`
		Username string `long:"db_username" description:"Username DB" env:"DB_USER" required:"true" default:"dbuser"`
		Password string `long:"db_password" description:"Password DB" env:"DB_PASS" required:"true" default:"dbpass"`
		SSLMode  string `long:"db_sslmode" description:"SSLMode DB" env:"DB_SSLMODE" required:"true" default:"disable"`
	}
}

var (
	appConfig     *Config
	appConfigOnce sync.Once
)

// newConfig creates a new configuration instance by parsing environment variables and command-line flags.
func newConfig() (*Config, error) {
	godotenv.Load(envfile)

	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default|flags.IgnoreUnknown)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(log.Writer())
		return nil, fmt.Errorf("config parse failed: %v", err)
	}

	return &cfg, nil
}

// GetAppConfig returns the application configuration.
func GetAppConfig() (*Config, error) {
	appConfigOnce.Do(func() {
		config, err := newConfig()
		if err != nil {
			log.Fatalf("can't load config: %v", err)
		}
		appConfig = config
	})

	return appConfig, nil
}
