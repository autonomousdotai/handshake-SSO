package models

import (
    _ "github.com/jinzhu/gorm"
)

type User struct {
    BaseModel
    Email string `gorm:"column:email" json:"email"`
    Name string `gorm:"column:name" json:"name"`
    Avatar string `gorm:"column:avatar" json:"avatar"`
    UUID string `gorm:"column:uuid;unique;not null;" json:"uuid,omitempty"`
    Status int `gorm:"column:status;default:1;" json:"status"`
    CardID string `gorm:"column:card_id;" json:"card_id"`
    CardVerified int `gorm:"column:card_verified;default:0" json:"card_verified"`
    RewardWalletAddresses string `gorm:"column:reward_wallet_addresses;size:1000" json:"reward_wallet_addresses"`
}

func (u User) TableName() string {
    return "user"
}
