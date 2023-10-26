package repository

import (
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	_ "gorm.io/gorm"
)

type UserRepository struct {
	Store *database.CustomDB
}

func (s *UserRepository) CreateUser(user *model.User) error {
	newUser := model.User{
		Active:   user.Active,
		Type:     user.Type,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	result := s.Store.Db.Create(&newUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s *UserRepository) GetUser() (error, []model.User) {
	var users []model.User
	result := s.Store.Db.Find(&users)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, users
}

func NewUser(active bool, userType int, username, email, password string) *model.User {
	return &model.User{
		Active:   active,
		Type:     model.Type(userType),
		Username: username,
		Email:    email,
		Password: password,
	}
}
