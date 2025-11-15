package users

// UserActivationResult - результат операции активации/деактивации пользователя
type UserActivationResult struct {
	UserID      string `json:"user_id"`
	OldStatus   bool   `json:"old_status"`
	NewStatus   bool   `json:"new_status"`
	Updated     bool   `json:"updated"`
	ErrorReason string `json:"error_reason,omitempty"`
}

// UserReviewInfo - информация о ревью пользователя (внутренняя структура)
type UserReviewInfo struct {
	PRID       string `json:"pr_id"`
	PRName     string `json:"pr_name"`
	AuthorID   string `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     string `json:"status"`
}