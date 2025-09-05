package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// config 파일 읽기
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No config file found: %v (falling back to env only)", err)
	}

	// config struct에 매핑
	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
