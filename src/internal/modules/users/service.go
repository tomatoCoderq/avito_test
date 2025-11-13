package users

import "github.com/tomatoCoderq/avito_task/src/models"

type RepositoryMethods interface {
	SetIsActive(userID string, isActive bool) (*models.User, error)
	GetUserReviews(userID string) ([]models.PR, error)
}

type Service struct {
	repo RepositoryMethods
}

func RegisterService(repo RepositoryMethods) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) SetIsActive(userID string, isActive bool) (*models.User, error) {
	return s.repo.SetIsActive(userID, isActive)
}

func (s *Service) GetUserReviews(userID string) ([]models.PR, error) {
	return s.repo.GetUserReviews(userID)
}
