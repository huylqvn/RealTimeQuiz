package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	Env       string `default:"development" envconfig:"ENV"`
	APP       string `default:"quiz-server" envconfig:"APP"`
	Port      string `default:"7000" envconfig:"PORT"`
	Version   string `default:"v1" envconfig:"VERSION"`
	LogLevel  string `default:"info" envconfig:"LOG_LEVEL"`
	LogFile   string `default:"false" envconfig:"LOG_FILE"`
	JwtSecret string `default:"test@2024" envconfig:"JWT_SECRET"`
}

var config *Config
var once = sync.Once{}

func NewFromFile() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file", err)
		return nil, err
	}
	err := viper.Unmarshal(config)
	if err != nil {
		fmt.Println("unable to decode into struct", err)
	}
	return config, nil
}

func GetEnvBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func New() (*Config, error) {
	once.Do(func() {
		godotenv.Load()
		config = new(Config)
		if err := envconfig.Process("", config); err != nil {
			panic(err)
		}
	})

	return config, nil
}
