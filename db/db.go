package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/mrgThang/flashcard-be/config"
)

func MustConnectMysql(config *config.MysqlConfig) *gorm.DB {
	dsn := config.DSN()
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}

	if err := db.AutoMigrate(); err != nil {
		panic("failed to migrate database: " + err.Error())
	}

	return db
}
