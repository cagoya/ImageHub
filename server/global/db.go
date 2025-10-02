package global

import (
	"fmt"
	"server/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	// MySQL连接配置
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	var err error
	db, err = gorm.Open(mysql.Open(dsn))

	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败：%v", err)
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}
