package middlewares

import (
    "fmt"
    "strings"
    "net/http"
    "strconv"

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
            fmt.Println("return with 401")
            // Response message Unauthorized
            c.AbortWithStatus(401)
            return
        }
       
        bkey := []byte(conf.GetString("secret_key"))
        uuid, err := utils.HashDecrypt(bkey, p)
       
        if err != nil {
            c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Invalid user!"})
            return
        }

        fmt.Println(p, uuid)

        user := models.User{}
        errDb := models.Database().Where("uuid = ?", uuid).First(&user).Error
        
        if errDb != nil {
            c.JSON(http.StatusOK, gin.H{"status": 0, "message": "Invalid user!"})
            c.Abort()
            return;
        }
        
        c.Set("User", user)
        c.Writer.Header().Set("UID", strconv.FormatUint(uint64(user.ID), 10))
        c.Next()
    }
}
