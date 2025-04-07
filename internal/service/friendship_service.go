package service

import (
	"damapp-server/internal/domain"
	"damapp-server/internal/repository"

	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FriendshipService struct {
	repo repository.FriendshipRepository
}

func NewFriendshipService(repo repository.FriendshipRepository) *FriendshipService {
	return &FriendshipService{repo: repo}
}

func (s *FriendshipService) SendFriendRequest(senderID, receiverID uint64) error {
	if senderID == receiverID {
		return errors.New("you cannot send a friend request to yourself")
	}

	existing, err := s.repo.GetByUserIDs(senderID, receiverID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("could not check existing friendship: %w", err)
	}

	if existing != nil {
		switch existing.Status {
		case domain.FriendshipStatusPending:
			return errors.New("a friend request is already pending")
		case domain.FriendshipStatusAccepted:
			return errors.New("you are already friends")
		case domain.FriendshipStatusRejected:
			return errors.New("friend request was previously rejected")
		}
	}
	friendship := &domain.Friendship{
		UserID1:   senderID,
		UserID2:   receiverID,
		Status:    domain.FriendshipStatusPending,
		CreatedAt: uint64(time.Now().Unix()),
	}

	if err := s.repo.Create(friendship); err != nil {
		return fmt.Errorf("could not send friend request: %w", err)
	}

	return nil
}

func (s *FriendshipService) AcceptFriendRequest(friendshipID uint64) error {
	err := s.repo.AcceptFriendRequest(friendshipID)
	if err != nil {
		return errors.New("could not accept the friend request")
	}
	return nil
}

func (s *FriendshipService) RejectFriendRequest(friendshipID uint64) error {
	err := s.repo.RejectFriendRequest(friendshipID)
	if err != nil {
		return errors.New("could not reject the friend request")
	}
	return nil
}

func (s *FriendshipService) GetPendingRequests(userID2 uint64) ([]domain.Friendship, error) {
	friendships, err := s.repo.GetPendingRequests(userID2)
	if err != nil {
		return nil, errors.New("failed to fetch pending requests")
	}
	return friendships, nil
}

func (s *FriendshipService) GetSentRequests(userID1 uint64) ([]domain.Friendship, error) {
	friendships, err := s.repo.GetSentFriendRequests(userID1)
	if err != nil {
		return nil, errors.New("failed to fetch sent requests")
	}
	return friendships, nil
}

func (s *FriendshipService) AreFriends(userID1, userID2 uint64) (bool, error) {
	friendship, err := s.repo.GetByUserIDs(userID1, userID2)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, fmt.Errorf("failed to check friendship: %w", err)
	}
	if friendship == nil {
		return false, nil
	}
	return friendship.Status == domain.FriendshipStatusAccepted, nil
}
