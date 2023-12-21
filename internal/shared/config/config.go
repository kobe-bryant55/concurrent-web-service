package config

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"

	"github.com/MehmetTalhaSeker/mts-blog-api/assets"
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
		MaxWorker int   `yaml:"max_worker"`
		MaxQueue  int   `yaml:"max_queue"`
		MaxLength int64 `yaml:"max_length"`
	} `yaml:"worker"`

	Version bool `yaml:"version"`
}

func Init() (*Config, error) {
	file, err := assets.EmbeddedFiles.ReadFile("configs/env.development.yaml")
	if err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	viper.AddConfigPath("./")
	viper.SetConfigName("env.development")
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
