package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type DefaultController struct{}


func (d DefaultController) Home(c *gin.Context) {
    resp := JsonResponse{1, "Handshake REST API", nil}
    c.JSON(http.StatusOK, resp)
}

func (d DefaultController) Notify(c *gin.Context) {
    jsonData := map[string]interface{}{
        "data": map[string]interface{}{
            "notification": map[string]interface{}{
                "title":"Demo Notification",
                "body":"Demo Notification description",
                "click_action":"https://stag-handshake.autonomous.ai",
            },
            "to": "eFY_FQvPAd8:APA91bFPBjSWgh82s5mgg2IKwnBHCFu2StaT8D7H9uIk-zLfU0xPxWOIIoZ4RCVeU_m6tk9oyGseddGQZdFhGPDEnNq883Qo7OGmuyslUEE9tadBgX9vcJpQG2LVwLlQFW8cc4azNASq",
        },
    }

    result, err := fcmService.Notify(jsonData)
    
    var status int
    var message string

    if result {
        status = 1
    } else {
        status = 0
        if err != nil {
            message = err.Error()
        }
    }
    
    resp := JsonResponse{status, message, nil}
    c.JSON(http.StatusOK, resp)
}

func (d DefaultController) NotFound(c *gin.Context) {
    resp := JsonResponse{0, "Page not found", nil}
    c.JSON(http.StatusOK, resp)
}
