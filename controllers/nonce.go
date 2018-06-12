package controllers

import (
    "net/http"
    "strconv"
    "time"
    "github.com/gin-gonic/gin"

    "github.com/ninjadotorg/handshake-dispatcher/models"
)

type NonceController struct{}

func (u NonceController) Get(c *gin.Context) {
    address := c.DefaultQuery("address", "_")
    networkId := c.DefaultQuery("network_id", "_")
    nonce := c.DefaultQuery("nonce", "0")

    if address == "_" || networkId == "_" {
        resp := JsonResponse{0, "Invalid address or network_id.", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    // todo add new user with key
    var model models.Nonce
    db := models.Database()
    err := db.Where("address = ? AND network_id = ?", address, networkId).First(&model).Error

    if err != nil {
        model.Address = address
        model.NetworkID = networkId
        if nonce != "0" {
            ni, err := strconv.Atoi(nonce)
            if err != nil {
                resp := JsonResponse{0, "Invalid nonce.", nil}
                c.JSON(http.StatusOK, resp)
                c.Abort()
                return;
            }
            model.Nonce = ni
        }
        errDb := db.Save(&model).Error

        if errDb != nil {
            resp := JsonResponse{0, "Get nonce failed.", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    } else {
        model.Nonce += 1
        if nonce != "0" {
            ni, err := strconv.Atoi(nonce)
            if err != nil {
                resp := JsonResponse{0, "Invalid nonce.", nil}
                c.JSON(http.StatusOK, resp)
                c.Abort()
                return;
            }
            model.Nonce = ni + 1
        }
        errDb := db.Save(&model).Error

        if errDb != nil {
            resp := JsonResponse{0, "Get nonce failed.", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    }

    resp := JsonResponse{1, "", model}
    c.JSON(http.StatusOK, resp)
}

func (u NonceController) Set(c *gin.Context) {
    address := c.DefaultQuery("address", "_")
    networkId := c.DefaultQuery("network_id", "_")
    nonce := c.DefaultQuery("nonce", "0")

    if address == "_" || networkId == "_" {
        resp := JsonResponse{0, "Invalid address or network_id.", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    // todo add new user with key
    var model models.Nonce
    db := models.Database()
    err := db.Where("address = ? AND network_id = ?", address, networkId).First(&model).Error

    if err != nil {
        model.Address = address
        model.NetworkID = networkId
        model.Timestamp = time.Now().UTC().Unix()
        if nonce != "0" {
            ni, err := strconv.Atoi(nonce)
            if err != nil {
                resp := JsonResponse{0, "Invalid nonce.", nil}
                c.JSON(http.StatusOK, resp)
                c.Abort()
                return;
            }
            model.Nonce = ni
        }
        errDb := db.Save(&model).Error

        if errDb != nil {
            resp := JsonResponse{0, "Get nonce failed.", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    } else {
        model.Nonce = 0
        model.Timestamp = time.Now().UTC().Unix()
        if nonce != "0" {
            ni, err := strconv.Atoi(nonce)
            if err != nil {
                resp := JsonResponse{0, "Invalid nonce.", nil}
                c.JSON(http.StatusOK, resp)
                c.Abort()
                return;
            }
            model.Nonce = ni
        }
        errDb := db.Save(&model).Error

        if errDb != nil {
            resp := JsonResponse{0, "Get nonce failed.", nil}
            c.JSON(http.StatusOK, resp)
            c.Abort()
            return;
        }
    }

    resp := JsonResponse{1, "", model}
    c.JSON(http.StatusOK, resp)
}
