package models

import (
	_ "github.com/jinzhu/gorm"
)

// SubscribedUser : this is used for collecting email
type SubscribedUser struct {
	BaseModel
	Email   string `gorm:"column:email" json:"email"`
	Product string `gorm:"column:product" json:"product"`
}

// TableName : SubscribeEmail
func (u SubscribedUser) TableName() string {
	return "subscribed_user"
}
