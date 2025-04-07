package service

import (
	"damapp-server/internal/apperror"
	"damapp-server/internal/domain"
	"damapp-server/internal/repository"
	"damapp-server/utils"

	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Authenticate(username, password string) (string, error) {
	user, err := s.repo.GetByUserName(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.New(apperror.ErrCodeDataNotFound, "user not found")
		}
		return "", apperror.New(apperror.ErrCodeInternal, "failed to retrieve user")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", apperror.New(apperror.ErrCodeInvalidCredential, "invalid username or password")
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		fmt.Println(err)
		return "", apperror.New(apperror.ErrCodeInternal, "could not create token")
	}

	return token, nil
}

func (s *UserService) CreateUser(username, password string) (string, error) {
	existingUser, err := s.repo.GetByUserName(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", apperror.New(apperror.ErrCodeInternal, "failed to check username availability")
	}

	if existingUser != nil {
		return "", apperror.New(apperror.ErrCodeInvalidCredential, "username is already taken")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperror.New(apperror.ErrCodeInternal, "failed to hash password")
	}

	userID := utils.GenerateID()

	err = s.repo.Create(&domain.User{ID: userID, Username: username, Password: string(hashedPass)})
	if err != nil {
		return "", apperror.New(apperror.ErrCodeInternal, "failed to create user")
	}

	token, err := utils.GenerateJWT(userID, username)
	if err != nil {
		return "", apperror.New(apperror.ErrCodeInternal, "could not create token")
	}

	return token, nil
}

func (s *UserService) GetUserByID(id uint64) (*domain.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrCodeDataNotFound, "user not found")
		}
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve user")
	}
	return user, nil
}

func (s *UserService) GetUserByUserName(username string) (*domain.User, error) {
	user, err := s.repo.GetByUserName(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.New(apperror.ErrCodeDataNotFound, "user not found")
		}
		return nil, apperror.New(apperror.ErrCodeInternal, "failed to retrieve user")
	}
	return user, nil

}
