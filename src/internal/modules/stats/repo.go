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
	
	// Получаем информацию о пользователе
	var user models.User
	if err := r.db.Preload("Teams").First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	
	stats.UserID = user.ID
	stats.Username = user.Name
	if len(user.Teams) > 0 {
		stats.TeamName = user.Teams[0].Name
	}
	
	// Статистика авторства PR
	var authorStats AuthorStats
	
	if err := r.db.Raw(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'OPEN') as open,
			COUNT(*) FILTER (WHERE status = 'MERGED') as merged
		FROM prs 
		WHERE author_id = ?`, userID).Scan(&authorStats).Error; err != nil {
		return nil, err
	}
	
	stats.AuthoredTotal = authorStats.Total
	stats.AuthoredOpen = authorStats.Open
	stats.AuthoredMerged = authorStats.Merged
	
	// Статистика ревью PR
	var reviewStats ReviewStats
	
	if err := r.db.Raw(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE p.status = 'OPEN') as open,
			COUNT(*) FILTER (WHERE p.status = 'MERGED') as merged
		FROM pr_reviewers pr
		JOIN prs p ON pr.pr_id = p.id
		WHERE pr.user_id = ?`, userID).Scan(&reviewStats).Error; err != nil {
		return nil, err
	}
	
	stats.ReviewingTotal = reviewStats.Total
	stats.ReviewingOpen = reviewStats.Open
	stats.ReviewingMerged = reviewStats.Merged
	
	return &stats, nil
}

func (r *Repo) GetOverviewStats() (*OverviewStats, error) {
	var stats OverviewStats
	
	// Общая статистика
	var generalStats GeneralStats
	
	if err := r.db.Raw(`
		SELECT 
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM users WHERE is_active = true) as active_users,
			(SELECT COUNT(*) FROM teams) as total_teams,
			(SELECT COUNT(*) FROM prs) as total_prs,
			(SELECT COUNT(*) FROM prs WHERE status = 'OPEN') as open_prs,
			(SELECT COUNT(*) FROM prs WHERE status = 'MERGED') as merged_prs
	`).Scan(&generalStats).Error; err != nil {
		return nil, err
	}
	
	stats.TotalUsers = generalStats.TotalUsers
	stats.ActiveUsers = generalStats.ActiveUsers
	stats.TotalTeams = generalStats.TotalTeams
	stats.TotalPRs = generalStats.TotalPRs
	stats.OpenPRs = generalStats.OpenPRs
	stats.MergedPRs = generalStats.MergedPRs
	
	// Среднее количество ревьюверов на PR
	var avgReviewers float64
	if err := r.db.Raw(`
		SELECT COALESCE(AVG(reviewer_count), 0) as avg_reviewers
		FROM (
			SELECT COUNT(pr.user_id) as reviewer_count
			FROM prs p
			LEFT JOIN pr_reviewers pr ON p.id = pr.pr_id
			GROUP BY p.id
		) subquery
	`).Scan(&avgReviewers).Error; err != nil {
		return nil, err
	}
	stats.AvgReviewers = avgReviewers
	
	// Топ 5 ревьюверов
	var topReviewers []TopReviewer
	if err := r.db.Raw(`
		SELECT 
			u.id as user_id,
			u.name as username,
			COUNT(pr.pr_id) as review_count
		FROM users u
		LEFT JOIN pr_reviewers pr ON u.id = pr.user_id
		WHERE u.is_active = true
		GROUP BY u.id, u.name
		ORDER BY review_count DESC
		LIMIT 5
	`).Scan(&topReviewers).Error; err != nil {
		return nil, err
	}
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