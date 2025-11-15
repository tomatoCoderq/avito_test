package stats

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

func (r *Repo) GetUserStats(userID string) (*UserStats, error) {
    var stats UserStats
    
    // Получаем пользователя с командами
    var user models.User
    if err := r.db.Preload("Teams").First(&user, "id = ?", userID).Error; err != nil {
        return nil, err
    }
    
    stats.UserID = user.ID
    stats.Username = user.Name
    if len(user.Teams) > 0 {
        stats.TeamName = user.Teams[0].Name
    }
    
    var totalAuthored, openAuthored, mergedAuthored int64
    
    // Всего созданных PR
    r.db.Model(&models.PR{}).Where("author_id = ?", userID).Count(&totalAuthored)
    // Открытых PR
    r.db.Model(&models.PR{}).Where("author_id = ? AND status = ?", userID, "OPEN").Count(&openAuthored)
    // Смерженных PR
    r.db.Model(&models.PR{}).Where("author_id = ? AND status = ?", userID, "MERGED").Count(&mergedAuthored)
    
    stats.AuthoredTotal = int(totalAuthored)
    stats.AuthoredOpen = int(openAuthored)
    stats.AuthoredMerged = int(mergedAuthored)
    
    var totalReviewing, openReviewing, mergedReviewing int64
    
    // Всего на ревью
    r.db.Table("pr_reviewers").
        Joins("JOIN prs ON pr_reviewers.pr_id = prs.id").
        Where("pr_reviewers.user_id = ?", userID).
        Count(&totalReviewing)
        
    // Открытых на ревью
    r.db.Table("pr_reviewers").
        Joins("JOIN prs ON pr_reviewers.pr_id = prs.id").
        Where("pr_reviewers.user_id = ? AND prs.status = ?", userID, "OPEN").
        Count(&openReviewing)
        
    // Смерженных на ревью
    r.db.Table("pr_reviewers").
        Joins("JOIN prs ON pr_reviewers.pr_id = prs.id").
        Where("pr_reviewers.user_id = ? AND prs.status = ?", userID, "MERGED").
        Count(&mergedReviewing)
    
    stats.ReviewingTotal = int(totalReviewing)
    stats.ReviewingOpen = int(openReviewing)
    stats.ReviewingMerged = int(mergedReviewing)
    
    return &stats, nil
}

func (r *Repo) GetOverviewStats() (*OverviewStats, error) {
	var stats OverviewStats

	var totalUsers, activeUsers, totalTeams, totalPRs, openPRs, mergedPRs int64

	// Подсчеты пользователей
	r.db.Model(&models.User{}).Count(&totalUsers)
	r.db.Model(&models.User{}).Where("is_active = ?", true).Count(&activeUsers)

	// Подсчеты команд
	r.db.Model(&models.Team{}).Count(&totalTeams)

	// Подсчеты PR
	r.db.Model(&models.PR{}).Count(&totalPRs)
	r.db.Model(&models.PR{}).Where("status = ?", "OPEN").Count(&openPRs)
	r.db.Model(&models.PR{}).Where("status = ?", "MERGED").Count(&mergedPRs)

	stats.TotalUsers = int(totalUsers)
	stats.ActiveUsers = int(activeUsers)
	stats.TotalTeams = int(totalTeams)
	stats.TotalPRs = int(totalPRs)
	stats.OpenPRs = int(openPRs)
	stats.MergedPRs = int(mergedPRs)

	// Топ 5 ревьюверов
	var topReviewers []TopReviewer

	r.db.Table("users u").
		Select("u.id as user_id, u.name as username, COUNT(pr.pr_id) as review_count").
		Joins("LEFT JOIN pr_reviewers pr ON u.id = pr.user_id").
		Where("u.is_active = ?", true).
		Group("u.id, u.name").
		Order("review_count DESC").
		Limit(5).
		Find(&topReviewers)

	stats.TopReviewers = topReviewers

	return &stats, nil
}

func (r *Repo) GetTeamStats(teamName string) (*TeamStats, error) {
	var stats TeamStats

	// Проверяем, что команда существует
	var team models.Team
	if err := r.db.Where("name = ?", teamName).First(&team).Error; err != nil {
		return nil, err
	}

	stats.TeamName = teamName

	// Статистика участников команды
	var memberStats MemberStats

	if err := r.db.Raw(`
		SELECT 
			COUNT(*) as total_members,
			COUNT(*) FILTER (WHERE u.is_active = true) as active_members
		FROM team_users tu
		JOIN users u ON tu.user_id = u.id
		WHERE tu.team_id = ?`, team.ID).Scan(&memberStats).Error; err != nil {
		return nil, err
	}

	stats.TotalMembers = memberStats.TotalMembers
	stats.ActiveMembers = memberStats.ActiveMembers

	// Статистика PR команды
	var prStats PRStats

	if err := r.db.Raw(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE p.status = 'OPEN') as open,
			COUNT(*) FILTER (WHERE p.status = 'MERGED') as merged
		FROM prs p
		JOIN team_users tu ON p.author_id = tu.user_id
		WHERE tu.team_id = ?`, team.ID).Scan(&prStats).Error; err != nil {
		return nil, err
	}

	stats.TotalPRs = prStats.Total
	stats.OpenPRs = prStats.Open
	stats.MergedPRs = prStats.Merged

	// Топ 5 авторов PR в команде
	var topContributors []TopContributor
	if err := r.db.Raw(`
		SELECT 
			u.id as user_id,
			u.name as username,
			COUNT(p.id) as authored_count
		FROM users u
		JOIN team_users tu ON u.id = tu.user_id
		LEFT JOIN prs p ON u.id = p.author_id
		WHERE tu.team_id = ? AND u.is_active = true
		GROUP BY u.id, u.name
		ORDER BY authored_count DESC
		LIMIT 5
	`, team.ID).Scan(&topContributors).Error; err != nil {
		return nil, err
	}
	stats.TopContributors = topContributors

	return &stats, nil
}
