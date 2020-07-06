package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Log        Log        `json:"log" mapstructure:"log"`
	HttpListen HttpListen `json:"http" mapstructure:"http"`
	DBConfig   DBConfig   `json:"db_config" mapstructure:"db"`
	GrpcServer GrpcServer `json:"grpc_server" mapstructure:"grpc"`
	Ampq       Ampq       `json:"ampq" mapstructure:"ampq"`
	Scheduler  Scheduler  `json:"scheduler" mapstructure:"scheduler"`
}

type Log struct {
	LogFile  string `json:"log_file" mapstructure:"log_file"`
	LogLevel string `json:"log_level" mapstructure:"log_level"`
}

type HttpListen struct {
	Ip   string `json:"ip" mapstructure:"ip"`
	Port string `json:"port" mapstructure:"port"`
}

type DBConfig struct {
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Host     string `json:"host" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	Database string `json:"database" mapstructure:"database"`
}

type GrpcServer struct {
	Host string `json:"ip" mapstructure:"host"`
	Port string `json:"port" mapstructure:"port"`
}

type Ampq struct {
	Host     string `json:"ip" mapstructure:"host"`
	Port     string `json:"port" mapstructure:"port"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
	Queue    string `json:"queue" mapstructure:"queue"`
}

type Scheduler struct {
	Period     string `json:"period" mapstructure:"period"`
	BeforeTime string `json:"before_time" mapstructure:"before_time"`
	EventTime  string `json:"event_time" mapstructure:"event_time"`
}

func GetConfigFromFile(filePath string) *Config {
	viper.SetConfigFile(filePath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file: %s", err.Error())
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	var C Config
	err := viper.Unmarshal(&C)
	if err != nil {
		log.Fatal(err)
	}
	return &C
}
