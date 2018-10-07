package models

import (
	_ "github.com/jinzhu/gorm"
)

type IDVerification struct {
	BaseModel
	UserID      uint   `gorm:"column:user_id" json:"user_id"`
	IDType      int    `gorm:"column:id_type;default:0" json:"id_type"`
	FrontImage  string `gorm:"column:front_image" json:"front_image"`
	BackImage   string `gorm:"column:back_image" json:"back_image"`
	SelfieImage string `gorm:"column:selfie_image" json:"selfie_image"`
	Status      int    `gorm:"column:status;default:0" json:"status"`
}

func (u IDVerification) TableName() string {
	return "id_verification"
}
