package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-default:"local"`
	StaragePath string     `yaml:"storage_path" env-required:"true"`
	HttpServer  HTTPServer `yaml:"http_server"`
}
type HTTPServer struct {
	Addres       string        `yaml:"addres" env-default:"localhost:8080"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	Idle_timeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	config_path := flag.String("CONFIG_PATH", "./config/local.yaml", "")
	flag.Parse()
	if *config_path == "" {
		log.Fatalf("Config path is empty")
	}
	if _, err := os.Stat(*config_path); os.IsNotExist(err) {
		log.Fatalf("Config path is incorrect: %s", *config_path)
	}
	var cfg Config
	err := cleanenv.ReadConfig(*config_path, &cfg)
	if err != nil {
		log.Fatal("Cant read config file")
	}
	return &cfg
}
