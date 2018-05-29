package controllers

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/config"
    "github.com/autonomousdotai/handshake-dispatcher/models"
    "github.com/autonomousdotai/handshake-dispatcher/utils"
)

type UserController struct{}

func (u UserController) SignUp(c *gin.Context) {
    config := config.GetConfig()
    UUID, passpharse, err := utils.HashNewUID(config.GetString("secret_key"))
   
    if err != nil {
        resp := JsonResponse{0, "Sign up failed", nil}
        c.JSON(http.StatusOK, resp)
        return
    }

    // todo add new user with key
    db := models.Database()
    user := models.User{UUID: UUID}

    errDb := db.Create(&user).Error;

    if errDb != nil {
        resp := JsonResponse{0, "Sign up failed", nil}
        c.JSON(http.StatusOK, resp)
        return
    }

    // implement another logic
    go ExchangeSignUp(user.ID)

    resp := JsonResponse{1, "", map[string]interface{}{"passpharse": passpharse}}
    c.JSON(http.StatusOK, resp)
    return
}

func (u UserController) Profile(c *gin.Context) {  
    var userModel models.User
    
    user, _ := c.Get("User")
    userModel = user.(models.User)
    userModel.UUID = ""
    
    resp := JsonResponse{1, "", userModel}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) UpdateProfile(c *gin.Context) {
    var userModel models.User
    
    user, _ := c.Get("User")
    userModel = user.(models.User)
    
    email := c.DefaultPostForm("email", "_")
    name := c.DefaultPostForm("name", "_")
    avatar, _ := c.FormFile("avatar")
    
    if email != "_" {
        userModel.Email = email
    }
    if name != "_" {
        userModel.Name = name
    }
    
    if avatar != nil {
        userModel.Avatar = avatar.Filename
    }

    db := models.Database()
    dbErr := db.Save(&userModel).Error

    if dbErr != nil {
        resp := JsonResponse{0, "Update profile failed.", nil}
        c.JSON(http.StatusOK, resp)
        return
    }

    userModel.UUID = ""

    resp := JsonResponse{1, "", userModel}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) ExportPassphrase(c *gin.Context) {
    resp := JsonResponse{1, "", "Export passpharse"}
    c.JSON(http.StatusOK, resp)
}

func ExchangeSignUp(userId uint) {
    jsonData := make(map[string]interface{})
    jsonData["id"] = userId

    endpoint, found := utils.GetForwardingEndpoint("exchange")
    fmt.Println(endpoint, found)
    jsonValue, _ := json.Marshal(jsonData)
  
    endpoint = fmt.Sprintf("%s/%s", endpoint, "/user/profile")
    
    request, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
    request.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{}
    _, err := client.Do(request)
    if err != nil {
        fmt.Println("call exchange failed ", err)
    } else {
        fmt.Println("call exchange on SignUp success")
    }
}
