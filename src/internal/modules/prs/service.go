package prs

import (
	"errors"
	"math/rand"

	"github.com/tomatoCoderq/avito_task/src/models"
)

type RepositoryMethods interface {
	CreatePR(pr *models.PR) (*models.PR, error)
	GetPRByID(prID string) (*models.PR, error)
	MergePR(prID string) (*models.PR, error)
	ReassignReviewer(prID string, oldUserID string, newUserID string) (*models.PR, error)
	GetUserByID(userID string) (*models.User, error)
	GetActiveTeamMembers(teamID string, excludeUserID string) ([]models.User, error)
}

type Service struct {
	repo RepositoryMethods
}

func RegisterService(repo RepositoryMethods) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetPRByID(prID string) (*models.PR, error) {
	return s.repo.GetPRByID(prID)
}

func (s *Service) CreatePR(prID, prName, authorID string) (*models.PR, error) {
	author, err := s.repo.GetUserByID(authorID)
	if err != nil {
		return nil, err
	}

	if len(author.Teams) == 0 {
		return nil, errors.New("author has no team")
	}

	teamMembers, err := s.repo.GetActiveTeamMembers(author.Teams[0].ID, author.ID)
	if err != nil {
		return nil, err
	}

	// Выбираем до 2 ревьюверов случайным образом из активных членов команды
	reviewers := s.selectReviewers(teamMembers, 2)

	pr := &models.PR{
		ID:        prID,
		Name:      prName,
		AuthorID:  author.ID,
		Status:    "OPEN",
		Reviewers: reviewers,
	}

	return s.repo.CreatePR(pr)
}

func (s *Service) MergePR(prID string) (*models.PR, error) {
	return s.repo.MergePR(prID)
}

func (s *Service) ReassignReviewer(prID, oldUserID string) (*models.PR, string, error) {
	pr, err := s.repo.GetPRByID(prID)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == "MERGED" {
		return nil, "", errors.New("PR_MERGED: cannot reassign on merged PR")
	}

	oldUser, err := s.repo.GetUserByID(oldUserID)
	if err != nil {
		return nil, "", err
	}

	isAssigned := false
	for _, reviewer := range pr.Reviewers {
		if reviewer.ID == oldUserID {
			isAssigned = true
			break
		}
	}

	if !isAssigned {
		return nil, "", errors.New("NOT_ASSIGNED: reviewer is not assigned to this PR")
	}

	// Получаем команду старого пользователя
	if len(oldUser.Teams) == 0 {
		return nil, "", errors.New("user has no team")
	}

	// Получаем активных членов команды (исключая автора и текущих ревьюверов)
	teamMembers, err := s.repo.GetActiveTeamMembers(oldUser.Teams[0].ID, pr.AuthorID)
	if err != nil {
		return nil, "", err
	}

	// Фильтруем текущих ревьюверов
	candidates := make([]models.User, 0)
	for _, member := range teamMembers {
		isCurrentReviewer := false
		for _, reviewer := range pr.Reviewers {
			if member.ID == reviewer.ID {
				isCurrentReviewer = true
				break
			}
		}
		if !isCurrentReviewer {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil, "", errors.New("NO_CANDIDATE: no active replacement candidate in team")
	}

	// Выбираем случайного кандидата
	newReviewer := candidates[rand.Intn(len(candidates))]

	// Переназначаем
	updatedPR, err := s.repo.ReassignReviewer(prID, oldUserID, newReviewer.ID)
	if err != nil {
		return nil, "", err
	}

	return updatedPR, newReviewer.ID, nil
}

func (s *Service) selectReviewers(candidates []models.User, maxCount int) []models.User {
	if len(candidates) == 0 {
		return []models.User{}
	}

	count := maxCount
	if len(candidates) < count {
		count = len(candidates)
	}

	shuffled := make([]models.User, len(candidates))
	copy(shuffled, candidates)

	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	// Возвращаем до двух случайных ревьюверов
	return shuffled[:count]
}
