package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Team struct {
	ID    string `gorm:"type:varchar(255);primaryKey"`
	Name  string `gorm:"unique"`
	Users []User `gorm:"many2many:team_users;"`
}

func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}



