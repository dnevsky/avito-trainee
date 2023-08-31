package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

const (
	DevEnv  = "local"
	ProdEnv = "production"
)

type Config struct {
	Debug     bool   `envconfig:"DEBUG"`
	Env       string `envconfig:"ENV"`
	SentryDSN string `envconfig:"SENTRY_DSN"`

	HTTPConfig struct {
		Port               string        `yaml:"port"`
		ReadTimeout        time.Duration `yaml:"readTimeout"`
		WriteTimeout       time.Duration `yaml:"writeTimeout"`
		MaxHeaderMegabytes int           `yaml:"maxHeaderBytes"`
	} `yaml:"http"`

	Limiter struct {
		RPS   int           `yaml:"rps"`
		Burst int           `yaml:"burst"`
		TTL   time.Duration `yaml:"ttl"`
	} `yaml:"limiter"`

	DB struct {
		Host     string `envconfig:"DB_HOST"`
		Port     string `envconfig:"DB_PORT"`
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
		DBName   string `envconfig:"DB_DATABASE"`
		SSLMode  string `envconfig:"DB_SSL_MODE"`
	}
}

func Init(configPath string) (*Config, error) {
	var cfg Config

	readEnv(&cfg)
	readFile(&cfg, configPath)

	return &cfg, nil
}

func processError(err error) {
	fmt.Printf("Error while read config: %s\n", err.Error())
	os.Exit(2)
}

func readFile(cfg *Config, configPath string) {
	f, err := os.Open(configPath)
	if err != nil {
		processError(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil && !errors.Is(err, io.EOF) {
		processError(err)
	}
}

func readEnv(cfg *Config) {
	err := godotenv.Load()
	if err != nil {
		// do nothing
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		processError(err)
	}
}
