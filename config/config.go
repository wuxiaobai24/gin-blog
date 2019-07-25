package config

import (
	"github.com/spf13/viper"
)

// Config config
type Config struct {
	Base     BaseConf     `mapstructure:"base"`
	Database DatabaseConf `mapstructure:"database"`
}

// BaseConf base cinfig
type BaseConf struct {
	Mode string
	Url  string
	Port int
}

// DatabaseConf config about database
type DatabaseConf struct {
	Type         string `mapstructure:"type"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DatabaseName string `mapstructure:"name"` // if Type == 'sqlite3' then DatabaseName should set path
}

var (
	Conf *Config
)

// InitConfig init config(load config file)
func InitConfig() error {
	Conf = NewConfig()

	viper.SetConfigName("conf")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	//viper.AddConfigPath("./conf")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(Conf); err != nil {
		return err
	}

	return nil
}

// NewConfig return a pointer to empty Config
func NewConfig() *Config {
	return &Config{}
}
