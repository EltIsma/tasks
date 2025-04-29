package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string `env:"ENV" env-default:"local"`
	Postgres PostgresConfig
	Redis    RedisConfig
	Server   ServerConfig
	Kafka    KafkaConfig
}

type ServerConfig struct {
	Port            string        `yaml:"port" env:"PORT"  env-default:"8090"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"READ_TIMEOUT" env-default:"30s"`
	WriteTimeout    time.Duration `yaml:"wtite_timeout" env:"WRITE_TIMEOUT" env-default:"30s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type PostgresConfig struct {
	PostgresURL string `env:"POSTGRES_URL" env-required:"true"`
}

type RedisConfig struct {
	Hosts    []string `env:"REDIS_HOSTS" yaml:"hosts" env-required:"true"`
	Password string   `env:"REDIS_PASSWORD" env-required:"true"`
}

type KafkaConfig struct {
	BrokerList []string `yaml:"brokers" env-required:"true"`
	Topic      string   `yaml:"topic" env-required:"true"`
}

func InitConfig() (*Config, error) {
	envPath, configPath := fetchConfigPath()

	if envPath == "" {
		return nil, fmt.Errorf("'.env' file path is empty")
	}

	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("no .env file found")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("can not read config and parse it: %w", err)
	}

	return &cfg, nil
}

func fetchConfigPath() (string, string) {
	var envPath, configPath string

	flag.StringVar(&envPath, "env", "", "application configuration file")
	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.Parse()

	if envPath == "" {
		envPath = os.Getenv("ENV_PATH")
	}

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	return envPath, configPath
}
