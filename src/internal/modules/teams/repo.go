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

func (r *Repo) AddUsersToTeam(teamName string, users []models.User) (*models.Team, error) {
	var team models.Team
	if err := r.db.Where("name = ?", teamName).First(&team).Error; err != nil {
		return nil, err
	}

	if err := r.CreateOrUpdateUsers(users); err != nil {
		return nil, err
	}

	if err := r.db.Model(&team).Association("Users").Append(users); err != nil {
		return nil, err
	}

	if err := r.db.Preload("Users").First(&team, "id = ?", team.ID).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// DeactivateUsersInTeam деактивирует пользователей в команде (batch операция)
func (r *Repo) DeactivateUsersInTeam(teamName string, userIDs []string) error {
	var team models.Team
	if err := r.db.Where("name = ?", teamName).First(&team).Error; err != nil {
		return err
	}

	var validUserIDs []string
	err := r.db.Table("team_users").
		Select("user_id").
		Where("team_id = ? AND user_id IN ?", team.ID, userIDs).
		Pluck("user_id", &validUserIDs).Error

	if err != nil {
		return err
	}

	if len(validUserIDs) == 0 {
		return nil
	}

	result := r.db.Model(&models.User{}).
		Where("id IN ?", validUserIDs).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetOpenPRsForReviewers получает все открытые PR для указанных ревьюверов
func (r *Repo) GetOpenPRsForReviewers(userIDs []string) ([]models.PR, error) {
	var prs []models.PR

	err := r.db.
		Joins("JOIN pr_reviewers ON pr_reviewers.pr_id = prs.id").
		Where("pr_reviewers.user_id IN ? AND prs.status = 'OPEN'", userIDs).
		Preload("Author").
		Preload("Reviewers").
		Find(&prs).Error

	if err != nil {
		return nil, err
	}

	return prs, nil
}

// GetActiveTeamMembersForReassignment получает активных участников команды для переназначения
func (r *Repo) GetActiveTeamMembersForReassignment(teamID string, excludeUserIDs []string) ([]models.User, error) {
	var users []models.User

	query := r.db.
		Joins("JOIN team_users ON team_users.user_id = users.id").
		Where("team_users.team_id = ? AND users.is_active = ?", teamID, true)

	if len(excludeUserIDs) > 0 {
		query = query.Where("users.id NOT IN ?", excludeUserIDs)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// BatchReassignReviewers выполняет батчевое переназначение ревьюверов
func (r *Repo) BatchReassignReviewers(reassignments []models.ReassignmentData) error {
	if len(reassignments) == 0 {
		return nil
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, reassignment := range reassignments {
			if err := tx.Exec(
				"DELETE FROM pr_reviewers WHERE pr_id = ? AND user_id = ?",
				reassignment.PRID, reassignment.OldReviewerID,
			).Error; err != nil {
				return err
			}

			if err := tx.Exec(
				"INSERT INTO pr_reviewers (pr_id, user_id) VALUES (?, ?)",
				reassignment.PRID, reassignment.NewReviewerID,
			).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ValidateUsersInTeam проверяет, что все указанные пользователи состоят в команде
func (r *Repo) ValidateUsersInTeam(teamName string, userIDs []string) ([]string, error) {
	var team models.Team
	if err := r.db.Where("name = ?", teamName).First(&team).Error; err != nil {
		return nil, err
	}

	var existingUserIDs []string
	err := r.db.
		Table("team_users").
		Select("user_id").
		Where("team_id = ? AND user_id IN ?", team.ID, userIDs).
		Pluck("user_id", &existingUserIDs).Error

	if err != nil {
		return nil, err
	}

	return existingUserIDs, nil
}
