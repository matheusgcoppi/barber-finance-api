package repository

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/matheusgcoppi/barber-finance-api/database"
	"github.com/matheusgcoppi/barber-finance-api/database/model"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/gorm"
	"strings"
)

type DbRepository struct {
	Store *database.CustomDB
}

func (s *DbRepository) CreateUser(user *model.User) (*model.User, *model.Account, error) {
	newUser := &model.User{
		Active:   user.Active,
		Type:     user.Type,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	result := s.Store.Db.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "users_email_key") {
			return nil, nil, fmt.Errorf("email address '%s' is already in use", user.Email)
		} else {
			return nil, nil, result.Error
		}
	}
	if newUser.Type == model.Admin {
		newAccount := &model.Account{
			UserId:  newUser.ID,
			Balance: 0,
			User:    newUser,
		}

		account := s.Store.Db.Create(&newAccount)
		if account.Error != nil {
			return nil, nil, account.Error
		}
		return newUser, newAccount, nil
	}

	return newUser, nil, nil
}

func (s *DbRepository) LoginUser(email string, password string) (*model.User, error) {
	var user model.User
	result := s.Store.Db.First(&user, "email = ?", email)
	if gorm.IsRecordNotFoundError(result.Error) {
		return nil, fmt.Errorf("user not found")
	} else if result.Error != nil {
		return nil, result.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func (s *DbRepository) GetUser() (error, []*model.User) {
	var users []*model.User
	result := s.Store.Db.Order("id").Find(&users)
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, users
}

func (s *DbRepository) GetUserByID(id string) (*model.User, error) {
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

func (s *DbRepository) DeleteUser(id string) error {
	result := s.Store.Db.Delete(&model.User{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("User with id = " + id + " not found")
	}
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *DbRepository) UpdateUser(user *model.UserDTO, id string) (*model.User, error) {
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
