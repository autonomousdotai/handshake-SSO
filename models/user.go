package models

import (
	"fmt"

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
	Status                int    `gorm:"column:status;default:1;" json:"status"`
	CardID                string `gorm:"column:card_id;" json:"card_id"`
	CardVerified          int    `gorm:"column:card_verified;default:0" json:"card_verified"`
	RewardWalletAddresses string `gorm:"column:reward_wallet_addresses;size:1000" json:"reward_wallet_addresses"`
	WalletAddresses       string `gorm:"column:wallet_addresses;size:1000" json:"wallet_addresses"`
	Metadata              string `gorm:"column:metadata;size:5000" json:"-"`
	FCMToken              string `gorm:"column:fcm_token;size:200" json:"fcm_token"`
	RefID                 uint   `gorm:"column:ref_id;" json:"-"`
}

func (u User) TableName() string {
	return "user"
}

var cryptosignService = new(services.CryptosignService)

// AfterCreate :
func (u *User) AfterCreate() {
	fmt.Println("Create")
	cryptosignService.UserModelHooks("Create", u.ID, u.Metadata, u.Email)
}

// AfterUpdate :
func (u *User) AfterUpdate() {
	fmt.Println("Update")
	cryptosignService.UserModelHooks("Update", u.ID, u.Metadata, u.Email)
}

// AfterDelete :
func (u *User) AfterDelete() (err error) {
	fmt.Println("Delete")
	cryptosignService.UserModelHooks("Delete", u.ID, u.Metadata, u.Email)
	return
}
