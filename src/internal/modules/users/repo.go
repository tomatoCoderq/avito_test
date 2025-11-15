package users

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

func (r *Repo) SetIsActive(userID string, isActive bool) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	user.IsActive = isActive
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	if err := r.db.Preload("Teams").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetUserReviews(userID string) ([]models.PR, error) {
	var prs []models.PR

	
	if err := r.db.
		Joins("JOIN pr_reviewers ON pr_reviewers.pr_id = prs.id").
		Where("pr_reviewers.user_id = ?", userID).
		Preload("Author").
		Find(&prs).Error; err != nil {
		return nil, err
	}

	return prs, nil
}

