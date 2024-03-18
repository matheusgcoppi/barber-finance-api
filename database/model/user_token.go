package model

type UserToken struct {
	CustomModel
	UserID uint   `json:"user_id"`
	Token  string `json:"token"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

type UserTokenDTO struct {
	UserID uint   `gorm:"primaryKey" json:"user_id"`
	Token  string `json:"token"`
}

func (UserToken) TableName() string {
	return "users_token"
}
