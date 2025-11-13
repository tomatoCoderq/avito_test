package teams

import (
	"errors"

	"github.com/tomatoCoderq/avito_task/src/models"
)

type RepositoryMethods interface {
	TeamCreate(team *models.Team) (*models.Team, error)
	TeamGetByName(name string) (*models.Team, error)
	TeamExists(name string) (bool, error)
	CreateOrUpdateUsers(users []models.User) error
}

type Service struct {
	repo RepositoryMethods
}

func RegisterService(repo RepositoryMethods) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) TeamCreate(team *models.Team) (*models.Team, error) {
	// Проверяем, существует ли команда
	exists, err := s.repo.TeamExists(team.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("team already exists")
	}

	// Создаем или обновляем пользователей
	if err := s.repo.CreateOrUpdateUsers(team.Users); err != nil {
		return nil, err
	}

	return s.repo.TeamCreate(team)
}

func (s *Service) TeamGetByName(name string) (*models.Team, error) {
	result, err := s.repo.TeamGetByName(name)
	return result, err
}
