package sql

import (
	"gorm.io/driver/postgres"
	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/gorm"
)

func New(connectionString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&models.User{}, &models.Team{}, &models.PR{}); err != nil {
		return nil, err
	}

	return db, nil
}