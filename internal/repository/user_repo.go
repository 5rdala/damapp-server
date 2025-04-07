package repository

import (
	"damapp-server/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint64) (*domain.User, error)
	GetByUserName(username string) (*domain.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{DB: db}
}

func (repo *UserRepo) Create(user *domain.User) error {
	if err := repo.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *UserRepo) GetByID(id uint64) (*domain.User, error) {
	var user domain.User
	if err := repo.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepo) GetByUserName(username string) (*domain.User, error) {
	var user domain.User
	if err := repo.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
