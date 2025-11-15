package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Team содержит информацию о команде. Модель используется для миграции	
type Team struct {
	ID    string `gorm:"type:varchar(255);primaryKey"`
	Name  string `gorm:"unique"`
	Users []User `gorm:"many2many:team_users;"`
}

// BeforeCreate хук GORM, который генерирует UUID перед созданием записи
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}



