package test

import (
	"database/sql"
	"fmt"
	"server/model"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

// 创建数据库
func CreateDB(dbConfig *DBConfig) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port))
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbConfig.Database)
	if err != nil {
		return fmt.Errorf("创建数据库失败：%v", err)
	}
	return nil
}

// 删除数据库
func DropTestDB(dbConfig *DBConfig) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port))
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}
	defer db.Close()
	_, err = db.Exec("DROP DATABASE IF EXISTS " + dbConfig.Database)
	if err != nil {
		return fmt.Errorf("删除数据库失败：%v", err)
	}
	return nil
}

// 清空数据库中的表
func ClearTables(dbConfig *DBConfig, tables []string) error {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port))
	if err != nil {
		return fmt.Errorf("连接数据库失败：%v", err)
	}
	defer db.Close()
	_, err = db.Exec("USE " + dbConfig.Database)
	if err != nil {
		return fmt.Errorf("使用数据库失败：%v", err)
	}
	for _, table := range tables {
		_, err = db.Exec("DELETE FROM " + table)
		if err != nil {
			return fmt.Errorf("清空表失败：%v", err)
		}
	}
	return nil
}

// 从配置文件加载配置
func LoadConfig() (*DBConfig, error) {
	dbConfig := &DBConfig{}
	viper.SetConfigName("test.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../config")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	dbConfig.Username = viper.GetString("db.username")
	dbConfig.Password = viper.GetString("db.password")
	dbConfig.Host = viper.GetString("db.host")
	dbConfig.Port = viper.GetString("db.port")
	dbConfig.Database = uuid.New().String()[:8] + viper.GetString("db.database")
	return dbConfig, nil
}

// 连接测试数据库
func ConnectTestDB(dbConfig *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		return nil, fmt.Errorf("连接数据库失败：%v", err)
	}

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("数据库迁移失败：%v", err)
	}

	return db, nil
}
