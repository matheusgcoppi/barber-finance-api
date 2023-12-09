package model

type Account struct {
	CustomModel
	UserId  uint       `json:"user_id"`
	Balance float64    `json:"balance"`
	User    *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Income  *[]Income  `gorm:"ForeignKey:AccountID" json:"income,omitempty"`
	Expense *[]Expense `gorm:"ForeignKey:AccountID" json:"expense,omitempty"`
}

type AccountDTO struct {
	CustomModel
	UserId  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
}

func (Account) tableName() string {
	return "accounts"
}
