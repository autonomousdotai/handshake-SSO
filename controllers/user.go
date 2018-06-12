package controllers

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"
    "io/ioutil"
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

    // todo add new user with key
    db := models.Database()
    user := models.User{UUID: UUID, Username: UUID[len(UUID)-8:]}

    errDb := db.Create(&user).Error;

    if errDb != nil {
        resp := JsonResponse{0, "Sign up failed", nil}
        c.JSON(http.StatusOK, resp)
        return
    }

    user.Username = fmt.Sprintf("%s%d", UUID[len(UUID) - 8:], user.ID) 

    errDb = db.Save(&user).Error

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

    resp := JsonResponse{1, "", userModel}
    c.JSON(http.StatusOK, resp)
}

func (u UserController) FreeRinkebyEth(c *gin.Context) {  
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

    var status bool
    var message string

    if userModel.Metadata != "" { 
        var md map[string]interface{}
        json.Unmarshal([]byte(userModel.Metadata), &md)   
        rinkeby, ok := md["free-rinkeby"]
        if ok {
            status = false
            message = fmt.Sprintf("Your free eth transaction is %s", rinkeby.(map[string]interface{})["hash"])
        } else {
            status, message = UserRequestRinkebyFreeEth(userModel, md, address)
        }
    } else {
        md := make(map[string]interface{})
        status, message = UserRequestRinkebyFreeEth(userModel, md, address)
    }
   
    resp := JsonResponse{1, message, status}
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

func UserRequestRinkebyFreeEth(user models.User, metadata map[string]interface{}, address string) (bool, string) {
    value := "1"
    status, hash := RequestRinkebyFreeEth(user.ID, address, value)
    
    if status {
        metadata["free-rinkeby"] = map[string]interface{}{
            "address": address,
            "value": value,
            "hash": hash,
            "time": time.Now().UTC().Unix(), 
        }
        
        md, _ := json.Marshal(metadata)
        user.Metadata = string(md)
        dbErr := models.Database().Save(&user).Error
        if dbErr != nil {
            return false, dbErr.Error()
        } else {
            return true, ""
        }
    } else {
        return false, hash
    }
}

func RequestRinkebyFreeEth(userId uint, address string, value string) (bool, string) {
    endpoint, _ := utils.GetServicesEndpoint("ethereum")
    
    endpoint = fmt.Sprintf("%s/rinkeby/free-ether?to_address=%s&value=%s", endpoint, address, value)

    request, _ := http.NewRequest("POST", endpoint, nil)
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Uid", fmt.Sprint(userId))
    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        fmt.Println(err.Error())
        return false, err.Error()
    }

    b, _ := ioutil.ReadAll(response.Body)

    var data map[string]interface{}
    json.Unmarshal(b, &data)

    status, ok := data["status"]
    message, _ := data["message"]

    if ok && (float64(1) == status) {
        rData := data["data"].(map[string]interface{})
        return true, rData["hash"].(string)
    } else {
        return false, message.(string)
    }
}

