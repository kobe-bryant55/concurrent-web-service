package config

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/MehmetTalhaSeker/concurrent-web-service/assets"
)

type Config struct {
	Rest struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		BaseURL string `yaml:"base_url"`
		Version string `yaml:"version"`
	} ` yaml:"rest"`

	DB struct {
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
	} `yaml:"db"`
	Env string `yaml:"env"`
	JWT struct {
		Secret string `yaml:"secret"`
		Exp    string `yaml:"exp"`
	} `yaml:"jwt"`

	Worker struct {
		MaxWorker int   `mapstructure:"max_worker"`
		MaxQueue  int   `mapstructure:"max_queue"`
		MaxLength int64 `mapstructure:"max_length"`
	} `yaml:"worker"`

	Version bool `yaml:"version"`
}

func Init() (*Config, error) {
	env := os.Getenv("cws-api-env")
	if env == "" {
		env = "local"
	}

	file, err := assets.EmbeddedFiles.ReadFile(fmt.Sprintf("configs/env.%s.yaml", env))
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	viper.AddConfigPath("./")
	viper.SetConfigName(fmt.Sprintf("env.%s", env))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.ReadConfig(bytes.NewReader(file))
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var c Config

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	log.Println("Config initialize success!")

	return &c, nil
}
