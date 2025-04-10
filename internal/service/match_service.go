package service

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/domain"
	"damapp-server/internal/repository"
	"damapp-server/utils"

	"errors"
	"time"

	"gorm.io/gorm"
)

type MatchService struct {
	repo repository.MatchRepository
}

func NewMatchService(repo repository.MatchRepository) *MatchService {
	return &MatchService{repo: repo}
}

func (s *MatchService) CreateMatch(playerID uint64) (*domain.Match, error) {
	// check if the player already in a match (each player can be in one match at a time)
	_, err := s.repo.CurrentMatch(playerID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to check if match already exist")
	}

	code, err := utils.Generate6DigitCode()
	if err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to generate match code")
	}

	match := &domain.Match{
		ID:        utils.GenerateID(),
		Code:      code,
		Player1ID: playerID,
		Status:    domain.MatchStatusWaiting,
		CreatedAt: time.Now().Unix(),
	}

	if err := s.repo.Create(match); err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to create match")
	}

	return match, nil
}

func (s *MatchService) GetByID(id uint64) (*domain.Match, error) {
	match, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrCodeDataNotFound, "match not found")
		}
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve match")
	}
	return match, nil
}

func (s *MatchService) GetByCode(code int) (*domain.Match, error) {
	match, err := s.repo.GetByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrCodeDataNotFound, "match not found")
		}
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve match")
	}
	return match, nil
}

func (s *MatchService) JoinMatch(code int, playerID uint64) (*domain.Match, error) {
	match, err := s.GetByCode(code)
	if err != nil {
		return nil, err
	}

	if match.Status != domain.MatchStatusWaiting {
		return nil, apperror.New(apperror.ErrCodeConflict, "match already started")
	}

	if match.Player1ID == playerID {
		return nil, apperror.New(apperror.ErrCodeConflict, "player already in the match")
	}

	if err := s.repo.JoinMatch(code, playerID); err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to join match")
	}

	return match, nil
}

func (s *MatchService) FinishMatch(matchID uint64, winnerID uint64) (*domain.Match, error) {
	match, err := s.GetByID(matchID)
	if err != nil {
		return nil, err
	}

	if match.Status == domain.MatchStatusFinished || match.Status == domain.MatchStatusStopped {
		return nil, apperror.New(apperror.ErrCodeConflict, "match has already finished or stopped")
	}

	if err := s.repo.FinishMatch(matchID, winnerID, time.Now().UnixMilli()); err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to finish match")
	}

	return match, nil
}

func (s *MatchService) StopMatch(matchID uint64, winnerID uint64) (*domain.Match, error) {
	match, err := s.GetByID(matchID)
	if err != nil {
		return nil, err
	}

	if match.Status == domain.MatchStatusFinished || match.Status == domain.MatchStatusStopped {
		return nil, apperror.New(apperror.ErrCodeConflict, "match has already finished or stopped")
	}

	if err := s.repo.StopMatch(matchID, time.Now().UnixMilli()); err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to stop match")
	}

	return match, nil
}
