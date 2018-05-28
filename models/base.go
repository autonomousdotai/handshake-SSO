package models

import (
    _ "github.com/jinzhu/gorm"
    "time"
)

type BaseModel struct {
    ID uint `gorm:"primary_key" json:"id"`
    DateCreated time.Time `gorm:"column:date_created" json:"date_created"`
    DateModified time.Time `gorm:"column:date_modified" json:"date_modified"`
}

func (m *BaseModel) BeforeCreate() (err error) {
    m.DateCreated = time.Now()
    m.DateModified = time.Now()
    return
}

func (m *BaseModel) BeforeUpdate() (err error) {
    m.DateModified = time.Now()
    return
}
