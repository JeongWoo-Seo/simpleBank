package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver      string `mapstructure:"db.driver"`
	DBSource      string `mapstructure:"db.source"`
	ServerAddress string `mapstructure:"server.address"`
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
