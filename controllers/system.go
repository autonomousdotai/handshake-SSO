package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/models"
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
