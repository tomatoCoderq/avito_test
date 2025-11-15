package stats

// UserStats - статистика конкретного пользователя
type UserStats struct {
	UserID         string `json:"user_id"`
	Username       string `json:"username"`
	TeamName       string `json:"team_name"`
	AuthoredTotal  int    `json:"authored_total"`
	AuthoredOpen   int    `json:"authored_open"`
	AuthoredMerged int    `json:"authored_merged"`
	ReviewingTotal int    `json:"reviewing_total"`
	ReviewingOpen  int    `json:"reviewing_open"`
	ReviewingMerged int   `json:"reviewing_merged"`
}

// OverviewStats - общая статистика системы
type OverviewStats struct {
	TotalUsers    int          `json:"total_users"`
	ActiveUsers   int          `json:"active_users"`
	TotalTeams    int          `json:"total_teams"`
	TotalPRs      int          `json:"total_prs"`
	OpenPRs       int          `json:"open_prs"`
	MergedPRs     int          `json:"merged_prs"`
	TopReviewers  []TopReviewer `json:"top_reviewers"`
}

// TopReviewer - топ ревьювер
type TopReviewer struct {
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	ReviewCount int    `json:"review_count"`
}

// TeamStats - статистика команды
type TeamStats struct {
	TeamName        string           `json:"team_name"`
	TotalMembers    int              `json:"total_members"`
	ActiveMembers   int              `json:"active_members"`
	TotalPRs        int              `json:"total_prs"`
	OpenPRs         int              `json:"open_prs"`
	MergedPRs       int              `json:"merged_prs"`
	TopContributors []TopContributor `json:"top_contributors"`
}

// TopContributor - топ автор PR в команде
type TopContributor struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username"`
	AuthoredCount int    `json:"authored_count"`
}

// AuthorStats - статистика авторства PR (внутренняя структура)
type AuthorStats struct {
	Total  int
	Open   int
	Merged int
}

// ReviewStats - статистика ревью PR (внутренняя структура)
type ReviewStats struct {
	Total  int
	Open   int
	Merged int
}

// GeneralStats - общая статистика системы (внутренняя структура)
type GeneralStats struct {
	TotalUsers  int
	ActiveUsers int
	TotalTeams  int
	TotalPRs    int
	OpenPRs     int
	MergedPRs   int
}

// MemberStats - статистика участников команды (внутренняя структура)
type MemberStats struct {
	TotalMembers  int
	ActiveMembers int
}

// PRStats - статистика PR команды (внутренняя структура)
type PRStats struct {
	Total  int
	Open   int
	Merged int
}