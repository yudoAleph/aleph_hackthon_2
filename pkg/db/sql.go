package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewSQLConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
