package models

type User struct {
	ID       string `gorm:"type:varchar(255);primaryKey"`
	Name     string
	IsActive bool
	Teams    []Team `gorm:"many2many:team_users;"`
}