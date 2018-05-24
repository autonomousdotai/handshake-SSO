package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/config"
    "github.com/autonomousdotai/handshake-dispatcher/models"
    "github.com/autonomousdotai/handshake-dispatcher/utils"
)

type UserController struct{}

var userModel = new(models.User)

func (u UserController) SignUp(c *gin.Context) {
    config := config.GetConfig()
    UUID, passpharse, err := utils.HashNewUID(config.GetString("secret_key"))
   
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Sign Up failed!"})
        return
    }

    // todo add new user with key
    tx := models.Database()
    user := models.User{UUID: UUID}

    errDb := tx.Create(&user).Error;

    if errDb != nil {
        c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Sign Up failed!"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"passpharse": passpharse})
    return
}

func (u UserController) Profile(c *gin.Context) {  
    user, _ := c.Get("User")
    c.JSON(http.StatusOK, gin.H{"page": "Profile", "user": user})
}

func (u UserController) UpdateProfile(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"page": "Update Profile"})
}

func (u UserController) ExportPassphrase(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"page": "Export File Passphrase"})
}
