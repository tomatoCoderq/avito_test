package stats

type RepositoryMethods interface {
	GetUserStats(userID string) (*UserStats, error)
	GetOverviewStats() (*OverviewStats, error)
	GetTeamStats(teamName string) (*TeamStats, error)
}

type Service struct {
	repo RepositoryMethods
}

func RegisterService(repo RepositoryMethods) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUserStats(userID string) (*UserStats, error) {
	return s.repo.GetUserStats(userID)
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	return s.repo.GetOverviewStats()
}

func (s *Service) GetTeamStats(teamName string) (*TeamStats, error) {
	return s.repo.GetTeamStats(teamName)
}