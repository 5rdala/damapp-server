package repository

import (
	"damapp-server/internal/domain"

	"gorm.io/gorm"
)

type FriendshipRepository interface {
	Create(friendship *domain.Friendship) error
	GetPendingRequests(userID2 uint64) ([]domain.Friendship, error)
	GetSentFriendRequests(userID1 uint64) ([]domain.Friendship, error)
	AcceptFriendRequest(friendshipID uint64) error
	RejectFriendRequest(friendshipID uint64) error
	GetByUserIDs(userID1, userID2 uint64) (*domain.Friendship, error)
}

type FriendshipRepo struct {
	DB *gorm.DB
}

func NewFriendshipRepo(db *gorm.DB) FriendshipRepository {
	return &FriendshipRepo{DB: db}
}

func (repo *FriendshipRepo) Create(friendship *domain.Friendship) error {
	return repo.DB.Create(friendship).Error
}

func (repo *FriendshipRepo) GetPendingRequests(userID2 uint64) ([]domain.Friendship, error) {
	var friendships []domain.Friendship
	err := repo.DB.Where("status = ? AND user_id_2 = ?", userID2, domain.FriendshipStatusPending).
		Find(&friendships).Error
	if err != nil {
		return nil, err
	}
	return friendships, nil
}

func (repo *FriendshipRepo) GetSentFriendRequests(userID1 uint64) ([]domain.Friendship, error) {
	var friendships []domain.Friendship
	err := repo.DB.Where("user_id_1 = ? AND status = ?", userID1, domain.FriendshipStatusPending).
		Find(&friendships).Error
	if err != nil {
		return nil, err
	}
	return friendships, nil
}

func (repo *FriendshipRepo) GetAcceptedFriendships(userID uint64) ([]domain.Friendship, error) {
	var friendships []domain.Friendship
	err := repo.DB.Where("status = ? AND (user_id_1 = ? OR user_id_2 = ?)",
		domain.FriendshipStatusAccepted, userID, userID).Find(&friendships).Error
	if err != nil {
		return nil, err
	}
	return friendships, nil
}

func (repo *FriendshipRepo) AcceptFriendRequest(friendshipID uint64) error {
	return repo.DB.Model(&domain.Friendship{}).Where("id = ?", friendshipID).
		Update("status", domain.FriendshipStatusAccepted).Error
}

func (repo *FriendshipRepo) RejectFriendRequest(friendshipID uint64) error {
	return repo.DB.Model(&domain.Friendship{}).Where("id = ?", friendshipID).
		Update("status", domain.FriendshipStatusRejected).Error
}

func (repo *FriendshipRepo) GetByUserIDs(userID1, userID2 uint64) (*domain.Friendship, error) {
	var friendship domain.Friendship
	err := repo.DB.Where(
		"(user_id_1 = ? AND user_id_2 = ?) OR (user_id_1 = ? AND user_id_2 = ?)",
		userID1, userID2, userID2, userID1).First(&friendship).Error
	if err != nil {
		return nil, err
	}
	return &friendship, nil
}
