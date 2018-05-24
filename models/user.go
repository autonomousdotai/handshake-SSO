package models

import (
    _ "github.com/jinzhu/gorm"
)

type User struct {
    Model
    Email string `gorm:"column:email"`
    Name string `gorm:"column:name"`
    Avatar string `gorm:"column:avatar"`
    UUID string `gorm:"column:uuid;unique;not null;"`
}

func (u User) TableName() string {
    return "user"
}
