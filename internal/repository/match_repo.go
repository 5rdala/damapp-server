package repository

import (
	"damapp-server/internal/domain"

	"gorm.io/gorm"
)

type MatchRepository interface {
	Create(match *domain.Match) error
	GetByID(id uint64) (*domain.Match, error)
	CurrentMatch(playerID uint64) (*domain.Match, error)
	GetByCode(code int) (*domain.Match, error)
	JoinMatch(code int, player2ID uint64) error
	UpdateStatus(id uint64, status domain.MatchStatus) error
	FinishMatch(id uint64, winnerID uint64, finishedAt int64) error
	StopMatch(id uint64, finishedAt int64) error
}

type MatchRepo struct {
	DB *gorm.DB
}

func NewMatchRepo(db *gorm.DB) MatchRepository {
	return &MatchRepo{DB: db}
}

func (repo *MatchRepo) Create(match *domain.Match) error {
	if err := repo.DB.Create(match).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MatchRepo) CurrentMatch(playerID uint64) (*domain.Match, error) {
	var match domain.Match

	err := repo.DB.Where("(player_id_1 = ? OR player_id_2 = ?) AND status = ?",
		playerID, playerID, domain.MatchStatusWaiting).First(&match).Error
	if err != nil {
		return nil, err
	}
	return &match, nil
}

func (repo *MatchRepo) GetByID(id uint64) (*domain.Match, error) {
	var match domain.Match
	if err := repo.DB.First(&match, id).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (repo *MatchRepo) GetByCode(code int) (*domain.Match, error) {
	var match domain.Match
	if err := repo.DB.Where("code = ?", code).First(&match).Error; err != nil {
		return nil, err
	}
	return &match, nil
}

func (repo *MatchRepo) JoinMatch(code int, player2ID uint64) error {
	return repo.DB.Model(&domain.Match{}).
		Where("code = ?", code).
		Update("player_2_id", player2ID).Error
}

func (repo *MatchRepo) UpdateStatus(id uint64, status domain.MatchStatus) error {
	return repo.DB.Model(&domain.Match{}).Where("id = ?", id).
		Update("status", status).Error
}

func (repo *MatchRepo) FinishMatch(id uint64, winnerID uint64, finishedAt int64) error {
	return repo.DB.Model(&domain.Match{}).Where("id = ?", id).
		Updates(map[string]any{
			"code":        0,
			"winner":      winnerID,
			"status":      domain.MatchStatusFinished,
			"finished_at": finishedAt,
		}).Error
}

func (repo *MatchRepo) StopMatch(id uint64, finishedAt int64) error {
	return repo.DB.Model(&domain.Match{}).Where("id = ?", id).
		Updates(map[string]any{
			"code":        0,
			"status":      domain.MatchStatusStopped,
			"finished_at": finishedAt,
		}).Error
}
