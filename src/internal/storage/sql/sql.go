package sql

import (
	"time"

	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(95)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.PR{}); err != nil {
		return nil, err
	}

	return db, nil
}
