package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	rawConf *viper.Viper
	once    sync.Once
	App     appConf
)

type appConf struct {
	Env     string      `mapstructure:"env" json:"env" yaml:"env"`
	Mode    string      `mapstructure:"mode" json:"mode" yaml:"mode"`
	Http    httpConf    `mapstructure:"http" json:"http" yaml:"http"`
	Data    data        `mapstructure:"data" json:"data" yaml:"data"`
	Log     logConf     `mapstructure:"log" json:"log" yaml:"log"`
	Gateway gatewayConf `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
}

type gatewayConf struct {
	Apisix apisixConf `mapstructure:"apisix" json:"apisix" yaml:"apisix"`
}

type apisixConf struct {
	Url    string `mapstructure:"url" json:"url" yaml:"url"`
	Token  string `mapstructure:"token" json:"token" yaml:"token"`
	Prefix string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
}

type postgres struct {
	Url                string `mapstructure:"url" json:"url" yaml:"url"`
	AutoMigrate        bool   `mapstructure:"migrate" json:"migrate" yaml:"migrate"`
	LogFile            string `mapstructure:"log_file_name" json:"log_file_name" yaml:"log_file_name"`
	LogLevel           string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	MaxIdleConnections int    `mapstructure:"max_idle_connections" json:"max_idle_connections" yaml:"max_idle_connections"`
	MaxOpenConnections int    `mapstructure:"max_open_connections" json:"max_open_connections" yaml:"max_open_connections"`
}

type dbConf struct {
	Postgres postgres `mapstructure:"postgres" json:"postgres" yaml:"postgres"`
}

type data struct {
	Database dbConf `mapstructure:"db" json:"db" yaml:"db"`
}

type httpConf struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type logConf struct {
	LogLevel    string `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	Encoding    string `mapstructure:"encoding" json:"encoding" yaml:"encoding"`
	LogFileName string `mapstructure:"log_file_name" json:"log_file_name" yaml:"log_file_name"`
	MaxBackups  int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge      int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	MaxSize     int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	Compress    bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

func InitConfig(p string) {
	once.Do(
		func() {
			rawConf = newConfig(p)
			if err := rawConf.Unmarshal(&App); err != nil {
				panic(err)
			}
		},
	)
}

func newConfig(p string) *viper.Viper {
	envConf := os.Getenv("APP_CONF")
	if envConf == "" {
		envConf = p
	}
	fmt.Println("load conf file:", envConf)
	return getConfig(envConf)
}

func getConfig(path string) *viper.Viper {
	conf := viper.New()
	conf.SetConfigFile(path)
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	return conf
}
