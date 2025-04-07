package service

import (
	"damapp-server/internal/domain"
	"damapp-server/internal/repository"

	"errors"

	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, password string) error {
	// Check if the username already exists
	existingUser, err := s.repo.GetByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to check username availability")
	}

	if existingUser != nil {
		return errors.New("username already used")
	}

	err = s.repo.Create(&domain.User{Username: username, Password: password})
	if err != nil {
		return errors.New("internal server err")
	}
	return nil
}

func (s *UserService) GetUserByID(id uint64) (*domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetUserByUserName(username string) (*domain.User, error) {
	return s.repo.GetByUserName(username)
}
