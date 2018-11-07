package models

import (
	_ "github.com/jinzhu/gorm"
)

type ActivityLog struct {
	BaseModel
	UserID       	uint   `gorm:"column:user_id" json:"user_id"`
	Action  	string `gorm:"column:action" json:"action"`
	Name  		string `gorm:"column:name" json:"name"`
	Description  	string `gorm:"column:description" json:"description"`
	Path  		string `gorm:"column:path" json:"path"`
	Host		string `gorm:"column:host" json:"host"`
	Method  	string `gorm:"column:method" json:"method"`
	UserAgent	string `gorm:"column:user_agent" json:"user_agent"`
	Status      	int    `gorm:"column:status;default:1" json:"status"`

}

func (u ActivityLog) TableName() string {
	return "activity_log"
}
