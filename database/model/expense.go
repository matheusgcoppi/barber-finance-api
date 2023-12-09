package model

import "time"

type Expense struct {
	CustomModel
	AccountID   uint      `gorm:"not null" json:"account_id"`
	Price       float64   `gorm:"not null" json:"price"`
	Description string    `gorm:"not null" json:"description"`
	When        time.Time `gorm:"not null" json:"when"`
	Payment     Payment   `gorm:"not null" json:"payment"`
}

type ExpenseDTO struct {
	AccountID   uint      `gorm:"not null" json:"account_id"`
	Price       float64   `gorm:"not null" json:"price"`
	Description string    `gorm:"not null" json:"description"`
	When        time.Time `gorm:"not null" json:"when"`
	Payment     Payment   `gorm:"not null" json:"payment"`
}

func (Expense) tableName() string {
	return "expenses"
}
