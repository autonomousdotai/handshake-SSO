package models

import (
	_ "github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-dispatcher/services"
)

type User struct {
	BaseModel
	Username              string `gorm:"column:username;unique;default:NULL" json:"username"`
	Email                 string `gorm:"column:email" json:"email"`
	Name                  string `gorm:"column:name" json:"name"`
	Phone                 string `gorm:"column:phone" json:"phone"`
	Avatar                string `gorm:"column:avatar" json:"avatar"`
	UUID                  string `gorm:"column:uuid;unique;not null;" json:"uuid,omitempty"`
	Verified              int    `gorm:"column:verified;default:0;" json:"verified"`
	CardID                string `gorm:"column:card_id;" json:"card_id"`
	CardVerified          int    `gorm:"column:card_verified;default:0" json:"card_verified"`
	RewardWalletAddresses string `gorm:"column:reward_wallet_addresses;size:1000" json:"reward_wallet_addresses"`
	WalletAddresses       string `gorm:"column:wallet_addresses;size:1000" json:"wallet_addresses"`
	Metadata              string `gorm:"column:metadata;size:5000" json:"-"`
	FCMToken              string `gorm:"column:fcm_token;size:200" json:"fcm_token"`
	IDVerified            int    `gorm:"column:id_verified;default:0" json:"id_verified"`
	IDVerificationLevel   int    `gorm:"column:id_verification_level;default:0" json:"id_verification_level"`
	RefID                 uint   `gorm:"column:ref_id;" json:"-"`
	Type                  int    `gorm:"column:type;default:0;" json:"type"`
	Password              string `gorm:"column:password;size:200" json:"password"`
}

func (u User) TableName() string {
	return "user"
}

var hookService = new(services.HookService)

// AfterUpdate :
func (u *User) AfterUpdate() {
	if len(u.Email) > 0 {
		go hookService.UserModelHooks("Update", u.ID, u.Metadata, u.Email, u.Name)
	}
}
