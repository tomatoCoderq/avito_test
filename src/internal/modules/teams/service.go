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
	AddUsersToTeam(teamName string, users []models.User) (*models.Team, error)
	DeactivateUsersInTeam(teamName string, userIDs []string) error
	GetOpenPRsForReviewers(userIDs []string) ([]models.PR, error)
	GetActiveTeamMembersForReassignment(teamID string, excludeUserIDs []string) ([]models.User, error)
	BatchReassignReviewers(reassignments []models.ReassignmentData) error
	ValidateUsersInTeam(teamName string, userIDs []string) ([]string, error)
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
	exists, err := s.repo.TeamExists(team.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("team already exists")
	}


	if err := s.repo.CreateOrUpdateUsers(team.Users); err != nil {
		return nil, err
	}

	return s.repo.TeamCreate(team)
}

func (s *Service) TeamGetByName(name string) (*models.Team, error) {
	result, err := s.repo.TeamGetByName(name)
	return result, err
}

func (s *Service) AddUsersToTeam(teamName string, users []models.User) (*models.Team, error) {
	
	exists, err := s.repo.TeamExists(teamName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("team not found")
	}

	return s.repo.AddUsersToTeam(teamName, users)
}

// DeactivateTeamUsersWithPRReassignment деактивирует пользователей команды и переназначает их PR
func (s *Service) DeactivateTeamUsersWithPRReassignment(teamName string, userIDs []string) (*models.DeactivationResult, error) {
	result := &models.DeactivationResult{
		DeactivatedUsers: []string{},
		ReassignedPRs:    []models.PRReassignmentInfo{},
		Errors:           []string{},
	}

	if exists, err := s.repo.TeamExists(teamName); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New("team not found")
	}

	validUserIDs, err := s.repo.ValidateUsersInTeam(teamName, userIDs)
	if err != nil {
		return nil, err
	}

	validUserMap := make(map[string]bool)
	for _, id := range validUserIDs {
		validUserMap[id] = true
	}
	for _, id := range userIDs {
		if !validUserMap[id] {
			result.Errors = append(result.Errors, "user "+id+" is not in team "+teamName)
		}
	}

	if len(validUserIDs) == 0 {
		return result, nil
	}

	team, err := s.repo.TeamGetByName(teamName)
	if err != nil {
		return nil, err
	}

	openPRs, err := s.repo.GetOpenPRsForReviewers(validUserIDs)
	if err != nil {
		return nil, err
	}

	excludeUserIDs := append(validUserIDs, s.extractAuthorIDs(openPRs)...)
	activeCandidates, err := s.repo.GetActiveTeamMembersForReassignment(team.ID, excludeUserIDs)
	if err != nil {
		return nil, err
	}

	reassignments, reassignmentInfos := s.prepareReassignments(openPRs, validUserIDs, activeCandidates)

	if err := s.repo.DeactivateUsersInTeam(teamName, validUserIDs); err != nil {
		return nil, err
	}

	if len(reassignments) > 0 {
		if err := s.repo.BatchReassignReviewers(reassignments); err != nil {
			return nil, err
		}
	}

	result.DeactivatedUsers = validUserIDs
	result.ReassignedPRs = reassignmentInfos

	return result, nil
}

// extractAuthorIDs извлекает ID авторов из списка PR
func (s *Service) extractAuthorIDs(prs []models.PR) []string {
	authorMap := make(map[string]bool)
	for _, pr := range prs {
		authorMap[pr.AuthorID] = true
	}

	authors := make([]string, 0, len(authorMap))
	for authorID := range authorMap {
		authors = append(authors, authorID)
	}
	return authors
}

// prepareReassignments подготавливает данные для батчевого переназначения
func (s *Service) prepareReassignments(prs []models.PR, deactivatedUserIDs []string, candidates []models.User) ([]models.ReassignmentData, []models.PRReassignmentInfo) {
	deactivatedMap := make(map[string]bool)
	for _, userID := range deactivatedUserIDs {
		deactivatedMap[userID] = true
	}

	reassignments := []models.ReassignmentData{}
	reassignmentInfos := []models.PRReassignmentInfo{}
	candidateIndex := 0

	for _, pr := range prs {
		for _, reviewer := range pr.Reviewers {
			// Если ревьювер деактивирован, переназначаем
			if deactivatedMap[reviewer.ID] {
				if candidateIndex < len(candidates) {
					// Равномерно распределяем между кандидатами
					newReviewer := candidates[candidateIndex%len(candidates)]
					candidateIndex++

					reassignments = append(reassignments, models.ReassignmentData{
						PRID:          pr.ID,
						OldReviewerID: reviewer.ID,
						NewReviewerID: newReviewer.ID,
					})

					reassignmentInfos = append(reassignmentInfos, models.PRReassignmentInfo{
						PRID:         pr.ID,
						FromReviewer: reviewer.ID,
						ToReviewer:   newReviewer.ID,
					})
				}

			}
		}
	}

	return reassignments, reassignmentInfos
}
