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

    var refMd map[string]interface{}
    if ref.Metadata != "" { 
        json.Unmarshal([]byte(ref.Metadata), &refMd)   
    } else {
        refMd = map[string]interface{}{}
    }
    
    referrals, hasReferrals := refMd["referrals"]
    if !hasReferrals {
        referrals = map[string]interface{}{}
    }
    
    aReferrals := referrals.(map[string]interface{})

    bonusKey := fmt.Sprintf("firstbet%d", user.ID)
    
    _, hasBonus := aReferrals[bonusKey]
    if hasBonus {
        log.Println("Received first bet bonus")
        resp := JsonResponse{0, "Received first bet bonus", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    var refWallets map[string]interface{}
    if ref.RewardWalletAddresses == "" {
        log.Println("The referrer have empty reward wallet address")
        resp := JsonResponse{0, "The referrer have empty reward wallet address", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    json.Unmarshal([]byte(ref.RewardWalletAddresses), &refWallets)
    ethWallet, hasEthWallet := refWallets["ETH"]

    if !hasEthWallet {
        log.Println("The referrer don't have eth reward wallet")
        resp := JsonResponse{0, "The referrer don't have eth reward wallet", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    amount := "20"
    address := ((ethWallet.(map[string]interface{}))["address"]).(string)
    status, hash := ethereumService.FreeToken(fmt.Sprint(ref.ID), address, amount, network)
    log.Println("status", status, hash)
    if !status {
        log.Println("Receive first bet bonus failed")
        resp := JsonResponse{0, "Receive first bet bonus failed", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }
    aReferrals[bonusKey] = map[string]interface{}{
        "address": address,
        "amount": amount,
        "hash": hash,
        "time": time.Now().UTC().Unix(), 
    }

    refMd["referrals"] = aReferrals        
    metadata, _ := json.Marshal(refMd)
    ref.Metadata = string(metadata)
    errDb = models.Database().Save(&ref).Error
    if errDb != nil {
        log.Println("Mark referrer received bonus failed", errDb.Error())
        resp := JsonResponse{0, "Mark referrer received bonus failed", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;

    }

    log.Println(ref)
    
    resp := JsonResponse{1, "Receive first bet bonus success", hash}
    c.JSON(http.StatusOK, resp)
}
