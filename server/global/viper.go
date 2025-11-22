package global

import (
	"os"

	"github.com/spf13/viper"
)

func ReadConfigFile() error {
	// 从环境变量读取配置文件名，如果没有设置则使用默认值
	configName := os.Getenv("CONFIG_FILE")
	if configName == "" {
		configName = "server"
	} else {
		if len(configName) > 5 && configName[len(configName)-5:] == ".yaml" {
			configName = configName[:len(configName)-5]
		} else if len(configName) > 4 && configName[len(configName)-4:] == ".yml" {
			configName = configName[:len(configName)-4]
		}
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}
