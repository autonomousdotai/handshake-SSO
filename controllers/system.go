package controllers

import (
    "time"
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"
    
    "github.com/ninjadotorg/handshake-dispatcher/models"
    "github.com/ninjadotorg/handshake-dispatcher/config"
)

type SystemController struct{}

func (s SystemController) User(c *gin.Context) {  
    userId := c.Param("id")
    user := models.User{}
    err := models.Database().Where("id = ?", userId).First(&user).Error

    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return; 
    }   

    resp := JsonResponse{1, "", user}
    c.JSON(http.StatusOK, resp)
}

func (s SystemController) BetSuccess(c *gin.Context) {
    userId := c.Param("id")
    betType := c.DefaultPostForm("type", "1")

    user := models.User{}
    errDb := models.Database().Where("id = ?", userId).First(&user).Error

    if errDb != nil {
        log.Println("Not exist user", errDb.Error())
        resp := JsonResponse{0, errDb.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    ref := models.User{}
    errDb = models.Database().Where("id = ?", user.RefID).First(&ref).Error

    if errDb != nil {
        log.Println("Not exist referrer", errDb.Error())
        resp := JsonResponse{0, errDb.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    conf := config.GetConfig()
    
    env := conf.GetString("env")
    network := "rinkeby"
    if env == "prod" {
        network = "mainnet"
    }

    if betType == "1" {
        BonusFirstBet(c, &user, &ref, network)
    } else {
        BonusFreeBet(c, &user, network)
    }

    c.Abort()
    return;
}

func BonusFirstBet(c *gin.Context, user *models.User, ref *models.User, network string) {
    userMd := GetUserMetadata(user)

    // todo send user free ether
    bonusKey := fmt.Sprintf("firstbet")

    amount := "80"
    status, hash, address := SendFreeShuriken(user, userMd, bonusKey, amount, network)

    if status {
        userMd[bonusKey] = map[string]interface{}{
            "address": address,
            "amount": amount,
            "hash": hash,
            "time": time.Now().UTC().Unix(), 
        }

        metadata, _ := json.Marshal(userMd)
        user.Metadata = string(metadata)
        errDb := models.Database().Save(&user).Error
        if errDb != nil {
            log.Println("Mark user send bonus failed", errDb.Error())
            resp := JsonResponse{0, "Mark user send bonus failed", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    }

    go mailService.SendFirstBet(user.Email, user.Username, hash)


    refMd := GetUserMetadata(ref)
    refReferrals, hasRefReferrals := refMd["referrals"]
    if !hasRefReferrals {
        refReferrals = map[string]interface{}{}
    }
    
    rbReferrals := refReferrals.(map[string]interface{})
    rbBonusKey := fmt.Sprintf("firstbet%d", user.ID)
    
    rbAmount := "20"
    rbStatus, rbHash, rbAddress := SendFreeShuriken(ref, rbReferrals, rbBonusKey, rbAmount, network)
   
    if rbStatus {
        rbReferrals[bonusKey] = map[string]interface{}{
            "address": rbAddress,
            "amount": rbAmount,
            "hash": rbHash,
            "time": time.Now().UTC().Unix(), 
        }

        refMd["referrals"] = rbReferrals        
        metadata, _ := json.Marshal(refMd)
        ref.Metadata = string(metadata)
        errDb := models.Database().Save(&ref).Error
        if errDb != nil {
            log.Println("Mark referrer send bonus failed", errDb.Error())
            resp := JsonResponse{0, "Mark referrer send bonus failed", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;

        }
    }

    log.Println("after BetSuccess", user, ref)

    go mailService.SendFirstBetReferrer(ref.Email, ref.Username, rbHash)
    resp := JsonResponse{1, "Send first bet bonus success", hash}
    c.JSON(http.StatusOK, resp)
}

func BonusFreeBet(c *gin.Context, user *models.User, network string) {
    userMd := GetUserMetadata(user)

    // todo send user free ether
    bonusKey := fmt.Sprintf("firstbet")

    amount := "20"
    status, hash, address := SendFreeShuriken(user, userMd, bonusKey, amount, network)

    if status {
        userMd[bonusKey] = map[string]interface{}{
            "address": address,
            "amount": amount,
            "hash": hash,
            "time": time.Now().UTC().Unix(), 
        }

        metadata, _ := json.Marshal(userMd)
        user.Metadata = string(metadata)
        errDb := models.Database().Save(&user).Error
        if errDb != nil {
            log.Println("Mark user send bonus failed", errDb.Error())
            resp := JsonResponse{0, "Mark user send bonus failed", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    }

    log.Println("after FreeBetSuccess", user)

    go mailService.SendFreeBet(user.Email, user.Username, hash)

    resp := JsonResponse{1, "Send free bet bonus success", hash}
    c.JSON(http.StatusOK, resp)
}

func GetUserMetadata(user *models.User) map[string]interface{} {
    var md map[string]interface{}
    if user.Metadata != "" { 
        json.Unmarshal([]byte(user.Metadata), &md)   
    } else {
        md = map[string]interface{}{}
    }
    return md;
}

func SendFreeShuriken(user *models.User, mark map[string]interface{}, bonusKey string, amount string, network string) (bool, string, string) {
    rtStatus := false
    rtHash := ""
    rtAddress := ""

    _, hasBonus := mark[bonusKey]
    if hasBonus {
        rtHash = "Received first bet bonus"
        return rtStatus, rtHash, rtAddress
    }

    var userWallets map[string]interface{}
    if user.RewardWalletAddresses == "" {
        rtHash = "The referrer have empty reward wallet address"
        return rtStatus, rtHash, rtAddress
    }

    json.Unmarshal([]byte(user.RewardWalletAddresses), &userWallets)
    ethWallet, hasEthWallet := userWallets["ETH"]

    if !hasEthWallet {
        rtHash = "The referrer don't have eth reward wallet"
        return rtStatus, rtHash, rtAddress
    }

    address := ((ethWallet.(map[string]interface{}))["address"]).(string)
    status, hash := ethereumService.FreeToken(fmt.Sprint(user.ID), address, amount, network)
    log.Println("status", status, hash)
    if !status {
        rtHash = "Receive first bet bonus failed"
        return rtStatus, rtHash, rtAddress
    }

    rtStatus = true
    rtHash = hash
    rtAddress = address

    return rtStatus, rtHash, rtAddress
}
