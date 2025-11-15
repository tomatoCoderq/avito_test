package models

// DeactivationResult представляет результат операции массовой деактивации пользователей
type DeactivationResult struct {
	DeactivatedUsers []string              `json:"deactivated_users"`
	ReassignedPRs    []PRReassignmentInfo  `json:"reassigned_prs"`
	Errors           []string              `json:"errors,omitempty"`
}

// PRReassignmentInfo содержит информацию о переназначении ревьювера в PR
type PRReassignmentInfo struct {
	PRID         string `json:"pr_id"`
	FromReviewer string `json:"from_reviewer"`
	ToReviewer   string `json:"to_reviewer"`
}

// ReassignmentData используется для батчевого переназначения ревьюверов
type ReassignmentData struct {
	PRID          string
	OldReviewerID string
	NewReviewerID string
}