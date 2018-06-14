package controllers

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
    "strings"
    "time"
    "github.com/gin-gonic/gin"

    "github.com/ninjadotorg/handshake-dispatcher/config"
    "github.com/ninjadotorg/handshake-dispatcher/models"
    "github.com/ninjadotorg/handshake-dispatcher/utils"
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

    ref := c.Query("ref")

    db := models.Database()
    
    exist := true
    username := ""
    
    for exist {
        count := 0
        username = utils.RandomNinjaName()
        errDb := db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error
        if errDb == nil && count == 0 {
            exist = false
        }
    }

    user := models.User{UUID: UUID, Username: username}
    if ref != "" {
        refUser := models.User{}
        refErr := db.Where("username = ?", ref).First(&refUser).Error

        if refErr == nil {
            user.RefID = refUser.ID
        }
    }

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

func (u UserController) UsernameExist(c *gin.Context) {
    username := c.DefaultQuery("username", "_")

    if username == "_" {
        resp := JsonResponse{0, "Invalid Username", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    var userModel models.User 
    user, _ := c.Get("User")
    userModel = user.(models.User)

    var _u models.User
    errDb := models.Database().Where("username = ? AND id != ?", username, userModel.ID).First(&_u).Error
  
    var result bool

    if errDb != nil {
        fmt.Println("Error", errDb.Error())
        result = false
    } else {
        result = true
    }

    resp := JsonResponse{1, "", result}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) UpdateProfile(c *gin.Context) {
    var userModel models.User
    
    user, _ := c.Get("User")
    userModel = user.(models.User)
    
    email := c.DefaultPostForm("email", "_")
    name := c.DefaultPostForm("name", "_")
    username := c.DefaultPostForm("username", "_")
    rwas := c.DefaultPostForm("reward_wallet_addresses", "_")
    phone := c.DefaultPostForm("phone", "_")
    ft := c.DefaultPostForm("fcm_token", "_")
    avatar, avatarErr := c.FormFile("avatar")
    
    if email != "_" {
        userModel.Email = email
    }
    if username != "_" {
        userModel.Username = username
    }
    if name != "_" {
        userModel.Name = name
    }
    if rwas != "_" {
        userModel.RewardWalletAddresses = rwas
    }
    if phone != "_" {
        userModel.Phone = phone
    }
    if ft != "_" {
        userModel.FCMToken = ft
    }
    
    if avatarErr == nil {
        uploadImageFolder := "user"
        fileName := avatar.Filename
        imageExt := strings.Split(fileName, ".")[1]
        fileNameImage := fmt.Sprintf("avatar-%d-image-%s.%s", userModel.ID, time.Now().Format("20060102150405"), imageExt)
        path := uploadImageFolder + "/" + fileNameImage 

        success, _ := uploadService.Upload(path, avatar)
        if !success {
            resp := JsonResponse{0, "Update profile failed: upload file error", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return  
        }

        userModel.Avatar = path
    }

    db := models.Database()
    dbErr := db.Save(&userModel).Error

    if dbErr != nil {
        fmt.Println("Error", dbErr.Error())
        resp := JsonResponse{0, "Update profile failed.", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return
    }

    userModel.UUID = ""

    go afterUpdateProfile(fmt.Sprint(userModel.ID))

    resp := JsonResponse{1, "", userModel}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) FreeRinkebyEther(c *gin.Context) {  
    var userModel models.User
    user, _ := c.Get("User")
    userModel = user.(models.User)
   
    address := c.DefaultQuery("address", "_")

    if address == "_" {
        resp := JsonResponse{0, "Invalid address", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    var md map[string]interface{}
    if userModel.Metadata != "" { 
        json.Unmarshal([]byte(userModel.Metadata), &md)   
    }


    var status bool
    var message string
    shouldRequest := false

    rinkeby, ok := md["free-rinkeby"]
    if ok {
        status = false
        message = fmt.Sprintf("Your free eth transaction is %s", rinkeby.(map[string]interface{})["hash"])
    } else {
        shouldRequest = true
    }

    if shouldRequest {
        value := "1"
        status, message := ethereumService.FreeEther(fmt.Sprint(userModel.ID), address, value, "rinkeby")
        if status {
            md["free-rinkeby"] = map[string]interface{}{
                "address": address,
                "value": value,
                "hash": message,
                "time": time.Now().UTC().Unix(), 
            }
        
            metadata, _ := json.Marshal(md)
            userModel.Metadata = string(metadata)
            dbErr := models.Database().Save(&userModel).Error
            if dbErr != nil {
                status = false
                message = dbErr.Error()
            } else {
                status = true
            }   
        }
    }
   
    resp := JsonResponse{1, message, status}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) Referred(c *gin.Context) {
    var userModel models.User
    var count int 

    user, _ := c.Get("User")
    userModel = user.(models.User)

    db := models.Database()
    errDb := db.Model(&models.User{}).Where("ref_id = ?", userModel.ID).Count(&count).Error
   
    if errDb != nil {
        panic(errDb)
    }
    
    resp := JsonResponse{1, "", count}
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
  
    endpoint = fmt.Sprintf("%s/%s", endpoint, "user/profile")
    
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

func afterUpdateProfile(userId string) {
    //todo check & bonus

    fmt.Println("todo after update propfile", userId)
}
