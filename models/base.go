package models

import (
    _ "github.com/jinzhu/gorm"
    "time"
)

type Model struct {
    ID uint `gorm:"primary_key"`

    DateCreated time.Time `gorm:"column:date_created"`
    DateModified time.Time `gorm:"column:date_modified"`

}

func (m *Model) BeforeCreate() (err error) {
    m.DateCreated = time.Now()
    m.DateModified = time.Now()
    return
}

func (m *Model) BeforeUpdate() (err error) {
    m.DateModified = time.Now()
    return
}
