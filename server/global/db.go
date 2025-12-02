package global

import (
	"fmt"
	"server/model"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() error {
	// PostgreSQL连接配置
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, username, password, database, port)

	var err error
	db, err = gorm.Open(postgres.Open(dsn))

	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{}, &model.Image{}, &model.Tag{}, &model.Chat{}, &model.Message{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败：%v", err)
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}
