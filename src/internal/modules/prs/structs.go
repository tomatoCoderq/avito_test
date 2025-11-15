package prs

// PRCreationData - данные для создания PR (внутренняя структура)
type PRCreationData struct {
	PRID     string   `json:"pr_id"`
	Name     string   `json:"name"`
	AuthorID string   `json:"author_id"`
	TeamID   string   `json:"team_id"`
	ReviewerIDs []string `json:"reviewer_ids"`
}

// ReassignmentCandidate - кандидат для переназначения (внутренняя структура)
type ReassignmentCandidate struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamID   string `json:"team_id"`
	IsActive bool   `json:"is_active"`
	ReviewCount int `json:"review_count"` // Для балансировки нагрузки
}

// PRMergeResult - результат слияния PR (внутренняя структура)
type PRMergeResult struct {
	PRID          string `json:"pr_id"`
	PreviousStatus string `json:"previous_status"`
	NewStatus      string `json:"new_status"`
	MergedAt       string `json:"merged_at"`
}