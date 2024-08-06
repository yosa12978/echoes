package config

import (
	"encoding/json"
	"fmt"
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
		Addr       string `yaml:"addr" envconfig:"SERVER_ADDR" json:"addr"`
		SessionKey string `yaml:"session_key" envconfig:"SERVER_SESSION_KEY" json:"session_key"`
		RootPass   string `yaml:"root_pass" envconfig:"SERVER_ROOT_PASS" json:"root_pass"`
	} `yaml:"server" json:"server"`
	Postgres struct {
		User string `yaml:"username" envconfig:"POSTGRES_USER" json:"username"`
		Pass string `yaml:"password" envconfig:"POSTGRES_PASS" json:"password"`
		Addr string `yaml:"addr" envconfig:"POSTGRES_ADDR" json:"addr"`
	} `yaml:"postgres" json:"postgres"`
	Redis struct {
		Addr     string `yaml:"addr" envconfig:"REDIS_ADDR" json:"addr"`
		Db       int    `yaml:"db" envconfig:"REDIS_DB" json:"db"`
		Password string `yaml:"password" envconfig:"REDIS_PASSWORD" json:"password"`
	} `yaml:"redis" json:"redis"`
	Feed struct {
		Title      string `yaml:"title" envconfig:"FEED_TITLE" json:"title"`
		Desc       string `yaml:"desc" envconfig:"FEED_DESC" json:"desc"`
		Link       string `yaml:"link" envconfig:"FEED_LINK" json:"link"`
		DetailLink string `yaml:"detail_link" envconfig:"FEED_DETAIL_LINK" json:"detail_link"`
		Author     string `yaml:"author" envconfig:"FEED_AUTHOR" json:"author"`
		Email      string `yaml:"email" envconfig:"FEED_EMAIL" json:"email"`
	} `yaml:"feed" json:"feed"`
	Profile struct {
		Name    string `yaml:"name" envconfig:"PROFILE_NAME" json:"name"`
		Bio     string `yaml:"bio" envconfig:"PROFILE_BIO" json:"bio"`
		Picture string `yaml:"picture" envconfig:"PROFILE_PICTURE" json:"picture"`
	} `yaml:"profile" json:"profile"`
}

func Get() Config {
	once.Do(func() {
		if err := readFile("config.yaml", &c); err != nil {
			panic(err)
		}
		if err := readEnv(&c); err != nil {
			panic(err)
		}

		foo, _ := json.MarshalIndent(c, "", "    ")
		fmt.Println(string(foo))
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
