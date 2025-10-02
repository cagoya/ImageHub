package global

import (
	"github.com/spf13/viper"
)

func ReadConfigFile() error {
	viper.SetConfigName("server.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
