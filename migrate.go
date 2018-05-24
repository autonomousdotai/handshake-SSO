package main

import (
    "github.com/jinzhu/gorm" 
    "github.com/autonomousdotai/handshake-dispatcher/models"
    "github.com/autonomousdotai/handshake-dispatcher/config"
)

func main() {
    config.Init()

    //
    var db *gorm.DB = models.Database()
    defer db.Close()

    db.AutoMigrate(&models.User{})
}



