package storage

import (
	"fmt"

	"github.com/goclaw/goclaw/config"
	"github.com/goclaw/goclaw/logs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB     *gorm.DB
	logger logs.Logger
}

func NewDatabase(logger logs.Logger, cfg *config.Config) (*Database, error) {
	dbPath := cfg.Database.Path
	if dbPath == "" {
		dbPath = "goclaw.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	logger.Infof("Database connected: %s", dbPath)

	return &Database{
		DB:     db,
		logger: logger,
	}, nil
}

func (d *Database) Migrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}
