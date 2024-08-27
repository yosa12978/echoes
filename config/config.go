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
		Addr       string `yaml:"addr" envconfig:"ECHOES_ADDR" json:"addr"`
		SessionKey string `yaml:"session_key" envconfig:"ECHOES_SESSION_KEY" json:"session_key"`
		RootPass   string `yaml:"root_pass" envconfig:"ECHOES_ROOT_PASS" json:"root_pass"`
	} `yaml:"server" json:"server"`
	Postgres struct {
		User string `yaml:"username" envconfig:"ECHOES_POSTGRES_USER" json:"username"`
		Pass string `yaml:"password" envconfig:"ECHOES_POSTGRES_PASS" json:"password"`
		Addr string `yaml:"addr" envconfig:"ECHOES_POSTGRES_ADDR" json:"addr"`
	} `yaml:"postgres" json:"postgres"`
	Redis struct {
		Addr     string `yaml:"addr" envconfig:"ECHOES_REDIS_ADDR" json:"addr"`
		Db       int    `yaml:"db" envconfig:"ECHOES_REDIS_DB" json:"db"`
		Password string `yaml:"password" envconfig:"ECHOES_REDIS_PASSWORD" json:"password"`
	} `yaml:"redis" json:"redis"`
	Feed struct {
		Title      string `yaml:"title" envconfig:"ECHOES_FEED_TITLE" json:"title"`
		Desc       string `yaml:"desc" envconfig:"ECHOES_FEED_DESC" json:"desc"`
		Link       string `yaml:"link" envconfig:"ECHOES_EED_LINK" json:"link"`
		DetailLink string `yaml:"detail_link" envconfig:"ECHOES_FEED_DETAIL_LINK" json:"detail_link"`
		Author     string `yaml:"author" envconfig:"ECHOES_FEED_AUTHOR" json:"author"`
		Email      string `yaml:"email" envconfig:"ECHOES_FEED_EMAIL" json:"email"`
	} `yaml:"feed" json:"feed"`
	Profile struct {
		Name    string `yaml:"name" envconfig:"ECHOES_PROFILE_NAME" json:"name"`
		Bio     string `yaml:"bio" envconfig:"ECHOES_PROFILE_BIO" json:"bio"`
		Picture string `yaml:"picture" envconfig:"ECHOES_PROFILE_PICTURE" json:"picture"`
	} `yaml:"profile" json:"profile"`
	Website struct {
		Title string `yaml:"title" envconfig:"ECHOES_WEBSITE_TITLE" json:"title"`
	} `yaml:"website" json:"website"`
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
