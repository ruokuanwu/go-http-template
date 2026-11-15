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
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{TranslateError: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, oops.Wrapf(err, "open sqlite database")
	}
	connection, err := database.DB()
	if err != nil {
		return nil, oops.Wrapf(err, "get database connection")
	}
	connection.SetMaxOpenConns(1)
	if err := database.AutoMigrate(&model.Todo{}, &model.User{}); err != nil {
		_ = connection.Close()
		return nil, oops.Wrapf(err, "migrate database")
	}
	if err := database.Exec(`CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY, data BLOB NOT NULL, expiry REAL NOT NULL
	)`).Error; err != nil {
		_ = connection.Close()
		return nil, oops.Wrapf(err, "migrate sessions")
	}
	if err := database.Exec(`CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry)`).Error; err != nil {
		_ = connection.Close()
		return nil, oops.Wrapf(err, "index sessions")
	}
	return database, nil
}
