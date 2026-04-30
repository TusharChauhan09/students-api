package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServerConfig `yaml:"http_server"`
}

type HTTPServerConfig struct {
	Address string `yaml:"address" env-required:"true"`
}

func MustLoad() *Config {
	var configPath string

	// 1. env variable
	configPath = os.Getenv("CONFIG_PATH")

	// 2. cli flag
	if configPath == "" {
		configFlag := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *configFlag

		if configPath == "" {
			log.Fatal("config path is not set")
		}
	}

	// 3. file exists?
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file error: %v", err)
	}

	// 4. read config
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	return &cfg
}