package database

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/CPtung/mattercontroller/internal/config"
	"github.com/CPtung/mattercontroller/pkg/model"
	"github.com/sirupsen/logrus"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once  sync.Once
	appdb *gorm.DB
	mutex sync.Mutex
)

func dbClient() *gorm.DB {
	once.Do(func() {
		var (
			err    error
			dbPath = filepath.Join(config.LibPath, "app.db")
		)
		// Ensure database directory exists
		dbDir := filepath.Dir(dbPath)
		if _, err := os.Stat(dbDir); os.IsNotExist(err) {
			logrus.Infof("App db directory not found, creating: %s", dbDir)
			if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
				logrus.Errorf("Failed to create db directory: %s", err.Error())
			}
		}

		// Open database connection
		if appdb, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}); err != nil {
			logrus.Errorf("Failed to open database: %s", err.Error())
		}

		// Auto migrate tables
		if err := appdb.AutoMigrate(&model.MatterDevice{}); err != nil {
			logrus.Errorf("Failed to auto migrate database: %s", err.Error())
		}
	})

	return appdb
}
