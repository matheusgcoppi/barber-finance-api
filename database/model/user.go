package model

import (
	"gorm.io/gorm"
	"time"
)

type Type int

const (
	System  Type = 1
	Admin   Type = 2
	Support Type = 3
)

type CustomModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	CustomModel
	Active   bool    `json:"active"`
	Type     Type    `gorm:"not null" sql:"index" json:"type"`
	Username string  `json:"username"`
	Email    string  `gorm:"unique" json:"email"`
	Password string  `json:"password"`
	Account  Account `gorm:"ForeignKey:UserId" json:"account,omitempty"`
}

type UserDTO struct {
	Active   bool   `json:"active"`
	Type     Type   `json:"type"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
