package models

import (
	_ "github.com/jinzhu/gorm"
)

type Store struct {
	BaseModel
	UserID       uint   `gorm:"column:user_id" json:"user_id"`
	StoreID      string `gorm:"column:store_id;unique;default:NULL" json:"store_id"`
	StoreName  string `gorm:"column:store_name" json:"store_name"`
	WalletsReceive  string `gorm:"column:wallets_receive" json:"wallets_receive"`
	ConfirmURL	string `gorm:"column:confirm_url" json:"confirm_url"`
	Status      int    `gorm:"column:status;default:0" json:"status"`
}

func (u Store) TableName() string {
	return "store"
}
