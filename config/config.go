package config

import (
	"os"
	"sync"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var (
	c    Config
	once sync.Once
)

type Config struct {
	Server struct {
		Addr       string `yaml:"addr" envconfig:"SERVER_ADDR"`
		SessionKey string `yaml:"session_key" envconfig:"SERVER_SESSION_KEY"`
		RootPass   string `yaml:"root_pass" envconfig:"SERVER_ROOT_PASS"`
	} `yaml:"server"`
	Postgres struct {
		Addr string `yaml:"addr" envconfig:"POSTGRES_ADDR"`
	} `yaml:"postgres"`
	Redis struct {
		Addr     string `yaml:"addr" envconfig:"REDIS_ADDR"`
		Db       int    `yaml:"db" envconfig:"REDIS_DB"`
		Password string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	} `yaml:"redis"`
	Feed struct {
		Title      string `yaml:"title" envconfig:"FEED_TITLE"`
		Desc       string `yaml:"desc" envconfig:"FEED_DESC"`
		Link       string `yaml:"link" envconfig:"FEED_LINK"`
		DetailLink string `yaml:"detail_link" envconfig:"FEED_DETAIL_LINK"`
		Author     string `yaml:"author" envconfig:"FEED_AUTHOR"`
		Email      string `yaml:"email" envconfig:"FEED_EMAIL"`
	} `yaml:"feed"`
	Profile struct {
		Name    string `yaml:"name" envconfig:"PROFILE_NAME"`
		Bio     string `yaml:"bio" envconfig:"PROFILE_BIO"`
		Picture string `yaml:"picture" envconfig:"PROFILE_PICTURE"`
	} `yaml:"profile"`
}

func Get() Config {
	once.Do(func() {
		if err := readFile("config.yaml", &c); err != nil {
			panic(err)
		}
		if err := readEnv(&c); err != nil {
			panic(err)
		}
	})
	return c
}

func readFile(filename string, cfg *Config) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return yaml.NewDecoder(f).Decode(cfg)
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}
