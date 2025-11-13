package teams

import (
	"github.com/tomatoCoderq/avito_task/src/models"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) TeamExists(name string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Team{}).Where("name = ?", name).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repo) CreateOrUpdateUsers(users []models.User) error {
	for _, user := range users {
		// Используем Upsert (создать или обновить)
		if err := r.db.Save(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) TeamCreate(team *models.Team) (*models.Team, error) {
	if err := r.db.Create(team).Error; err != nil {
		return nil, err
	}
	
	// Загружаем связанных пользователей
	if err := r.db.Preload("Users").First(team, "id = ?", team.ID).Error; err != nil {
		return nil, err
	}
	
	return team, nil
}

func (r *Repo) TeamGetByName(name string) (*models.Team, error) {
	var team models.Team

	if err := r.db.Preload("Users").Where("name = ?", name).First(&team).Error; err != nil {
		return nil, err
	}
	return &team, nil
}
