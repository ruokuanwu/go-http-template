package db

import (
	"os"
	"path/filepath"

	"example.com/go-http-template/internal/model"
	"github.com/samber/oops"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(path string) (*gorm.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, oops.Wrapf(err, "create database directory")
	}
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, oops.Wrapf(err, "open sqlite database")
	}
	connection, err := database.DB()
	if err != nil {
		return nil, oops.Wrapf(err, "get database connection")
	}
	connection.SetMaxOpenConns(1)
	if err := database.AutoMigrate(&model.Todo{}); err != nil {
		_ = connection.Close()
		return nil, oops.Wrapf(err, "migrate database")
	}
	return database, nil
}
