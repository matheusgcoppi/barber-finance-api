package model

import "time"

type Payment int

const (
	Cash       Payment = 1
	CreditCard Payment = 2
	DebitCard  Payment = 3
	Pix        Payment = 4
	Bill       Payment = 5
	Other      Payment = 6
)

type Income struct {
	CustomModel
	AccountID   uint      `gorm:"not null" json:"account_id"`
	Price       float64   `gorm:"not null" json:"price"`
	Description string    `gorm:"not null" json:"description"`
	When        time.Time `gorm:"not null" json:"when"`
	Payment     Payment   `gorm:"not null" json:"payment"`
}

type IncomeDTO struct {
	AccountID   uint      `gorm:"not null" json:"account_id"`
	Price       float64   `gorm:"not null" json:"price"`
	Description string    `gorm:"not null" json:"description"`
	When        time.Time `gorm:"not null" json:"when"`
	Payment     Payment   `gorm:"not null" json:"payment"`
}

func (Income) tableName() string {
	return "incomes"
}
