package repository

import (
	"fmt"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	_ "gorm.io/gorm"
)

type UserRepository struct {
	Store *database.CustomDB
}

func (s *UserRepository) CreateUser(user *model.User) (error, *model.User) {
	newUser := &model.User{
		Active:   user.Active,
		Type:     user.Type,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	result := s.Store.Db.Create(&newUser)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, newUser
}

func (s *UserRepository) GetUser() (error, []*model.User) {
	var users []*model.User
	result := s.Store.Db.Order("id").Find(&users)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, users
}

func (s *UserRepository) GetUserByID(id string) (*model.User, error) {
	var user *model.User
	result := s.Store.Db.First(&user, id)
	if id == "" {
		return nil, fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (s *UserRepository) DeleteUser(id string) error {
	result := s.Store.Db.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *UserRepository) UpdateUser(user *model.UserDTO, id string) (*model.User, error) {
	userById, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if (user.Active == true) || (user.Active == false) {
		userById.Active = user.Active
	}

	if user.Type != 0 {
		userById.Type = user.Type
	}

	if user.Username != "" {
		userById.Username = user.Username
	}

	if user.Email != "" {
		userById.Email = user.Email
	}

	if user.Password != "" {
		userById.Password = user.Password
	}

	s.Store.Db.Save(userById)
	return userById, nil
}

func NewUser(active bool, userType int, username, email, password string) *model.User {
	return &model.User{
		CustomModel: model.CustomModel{},
		Active:      active,
		Type:        model.Type(userType),
		Username:    username,
		Email:       email,
		Password:    password,
	}
}
