package middlewares

import (
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/autonomousdotai/handshake-dispatcher/config"
    "github.com/autonomousdotai/handshake-dispatcher/utils"
    "github.com/autonomousdotai/handshake-dispatcher/models"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        conf := config.GetConfig()
        payload := c.Request.Header.Get("Payload")

        p := strings.TrimSpace(payload)

        if len(p) == 0 {
            panic("Invalid user.")
        }
       
        bkey := []byte(conf.GetString("secret_key"))
        uuid, err := utils.HashDecrypt(bkey, p)
       
        if err != nil {
            panic("Invalid user.")
        }

        user := models.User{}
        errDb := models.Database().Where("uuid = ?", uuid).First(&user).Error
        
        if errDb != nil {
            panic("Invalid user.")    
        }
        
        c.Set("User", user)
        c.Next()
    }
}
