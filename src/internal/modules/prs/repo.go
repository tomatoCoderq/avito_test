package prs

import (
	"errors"

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

func (r *Repo) CreatePR(pr *models.PR) (*models.PR, error) {
	if err := r.db.Create(pr).Error; err != nil {
		return nil, err
	}
	
	// Загружаем связанные данные
	if err := r.db.Preload("Author").Preload("Reviewers").First(pr, "id = ?", pr.ID).Error; err != nil {
		return nil, err
	}
	
	return pr, nil
}

func (r *Repo) GetPRByID(prID string) (*models.PR, error) {
	var pr models.PR
	if err := r.db.Preload("Author").Preload("Reviewers").First(&pr, "id = ?", prID).Error; err != nil {
		return nil, err
	}
	
	return &pr, nil
}

func (r *Repo) MergePR(prID string) (*models.PR, error) {
	var pr models.PR
	if err := r.db.First(&pr, "id = ?", prID).Error; err != nil {
		return nil, err
	}

	pr.Status = "MERGED"
	
	if err := r.db.Save(&pr).Error; err != nil {
		return nil, err
	}

	// Загружаем связанные данные
	if err := r.db.Preload("Author").Preload("Reviewers").First(&pr, "id = ?", prID).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *Repo) ReassignReviewer(prID string, oldUserID string, newUserID string) (*models.PR, error) {
	var pr models.PR
	if err := r.db.Preload("Reviewers").First(&pr, "id = ?", prID).Error; err != nil {
		return nil, err
	}

	// Удаляем старого ревьювера
	found := false
	newReviewers := make([]models.User, 0)
	for _, reviewer := range pr.Reviewers {
		if reviewer.ID == oldUserID {
			found = true
			continue
		}
		newReviewers = append(newReviewers, reviewer)
	}

	if !found {
		return nil, errors.New("reviewer not assigned to this PR")
	}

	// Добавляем нового ревьювера
	var newReviewer models.User
	if err := r.db.First(&newReviewer, "id = ?", newUserID).Error; err != nil {
		return nil, err
	}
	
	newReviewers = append(newReviewers, newReviewer)
	pr.Reviewers = newReviewers

	if err := r.db.Model(&pr).Association("Reviewers").Replace(newReviewers); err != nil {
		return nil, err
	}

	// Загружаем обновленные данные
	if err := r.db.Preload("Author").Preload("Reviewers").First(&pr, "id = ?", prID).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

func (r *Repo) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Teams").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetActiveTeamMembers(teamID string, excludeUserID string) ([]models.User, error) {
	var users []models.User
	
	// Получаем активных пользователей из команды, исключая указанного пользователя (автора PR)
	// Условия: team_id = teamID AND is_active = true AND user_id != excludeUserID
	if err := r.db.
		Joins("JOIN team_users ON team_users.user_id = users.id").
		Where("team_users.team_id = ? AND users.is_active = ? AND users.id != ?", teamID, true, excludeUserID).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

