package service

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/domain"
	"damapp-server/internal/repository"

	"errors"
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
		return apperror.New(apperror.ErrCodeBadRequest, "you cannot send friend request to yourself")
	}

	existing, err := s.repo.GetByUserIDs(senderID, receiverID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.New(apperror.ErrCodeInternal, "could not check friendship existence")
	}

	if existing != nil {
		switch existing.Status {
		case domain.FriendshipStatusPending:
			return apperror.New(apperror.ErrCodeBadRequest, "a friend request is already pending")
		case domain.FriendshipStatusAccepted:
			return apperror.New(apperror.ErrCodeBadRequest, "you are already friends")
		case domain.FriendshipStatusRejected:
			return apperror.New(apperror.ErrCodeBadRequest, "friend request was previously rejected")
		}
	}

	friendship := &domain.Friendship{
		UserID1:   senderID,
		UserID2:   receiverID,
		Status:    domain.FriendshipStatusPending,
		CreatedAt: uint64(time.Now().Unix()),
	}

	if err := s.repo.Create(friendship); err != nil {
		return apperror.New(apperror.ErrCodeInternal, "could not send friend request")
	}

	return nil
}

func (s *FriendshipService) AcceptFriendRequest(friendshipID uint64) error {
	err := s.repo.AcceptFriendRequest(friendshipID)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "could not accept the friend request")
	}
	return nil
}

func (s *FriendshipService) RejectFriendRequest(friendshipID uint64) error {
	err := s.repo.RejectFriendRequest(friendshipID)
	if err != nil {
		return apperror.New(apperror.ErrCodeInternal, "could not reject the friend request")
	}
	return nil
}

func (s *FriendshipService) GetPendingRequests(userID2 uint64) ([]domain.Friendship, error) {
	friendships, err := s.repo.GetPendingRequests(userID2)
	if err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve pending requests")
	}
	return friendships, nil
}

func (s *FriendshipService) GetSentRequests(userID1 uint64) ([]domain.Friendship, error) {
	friendships, err := s.repo.GetSentFriendRequests(userID1)
	if err != nil {
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve sent requests")
	}
	return friendships, nil
}

func (s *FriendshipService) AreFriends(userID1, userID2 uint64) (bool, error) {
	friendship, err := s.repo.GetByUserIDs(userID1, userID2)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, apperror.New(apperror.ErrCodeInternal, "failed to check friendship")
	}
	if friendship == nil {
		return false, nil
	}
	return friendship.Status == domain.FriendshipStatusAccepted, nil
}
