package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env   string      `yaml:"env"`
	SMTP  SMTPConfig  `yaml:"smtp" env_required:"true"`
	Kafka KafkaConfig `yaml:"kafka" env_required:"true"`
}

type SMTPConfig struct {
	Addr   string `yaml:"addr"`
	Host   string `yaml:"host"`
	From   string `yaml:"from"`
	Pass   string `yaml:"pass"`
	User   string `yaml:"user"`
}

type KafkaConfig struct {
	Host  string `yaml:"host"`
	Port  string `yaml:"port"`
	Topic string `yaml:"topic"`
}

func MustLoadConfig() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("can not parse config")
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	return res
}
