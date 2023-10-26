package model

import "gorm.io/gorm"

type Type int

const (
	System  Type = 1
	Admin   Type = 2
	Support Type = 3
)

type User struct {
	gorm.Model
	Active   bool
	Type     Type `gorm:"not null" sql:"index"`
	Username string
	Email    string `gorm:"unique"`
	Password string
}

type UserDTO struct {
	Type     Type   `json:"type"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
