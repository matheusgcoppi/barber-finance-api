package model

type Account struct {
	CustomModel
	UserId  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
	User    *User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

type AccountDTO struct {
	CustomModel
	UserId  uint    `json:"user_id"`
	Balance float64 `json:"balance"`
}

func (Account) tableName() string {
	return "accounts"
}
