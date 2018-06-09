package main

import (
	"github.com/jinzhu/gorm"
	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

func main() {
	config.Init()

	//
	var db *gorm.DB = models.Database()
	defer db.Close()

	db.AutoMigrate(&models.User{})
    db.AutoMigrate(&models.Nonce{})
}
