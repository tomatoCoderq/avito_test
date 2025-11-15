package teams

// DeactivationRequestData - внутренняя структура для запроса деактивации
type DeactivationRequestData struct {
	TeamName string   `json:"team_name"`
	UserIDs  []string `json:"user_ids"`
}

// TeamMemberInfo - информация о участнике команды (внутренняя структура)
type TeamMemberInfo struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

// BatchUpdateResult - результат пакетного обновления (внутренняя структура)
type BatchUpdateResult struct {
	Updated int
	Failed  int
	Errors  []string
}