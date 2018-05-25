package middlewares

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins: []string{"http://*, https://*"},
        AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"},
        AllowHeaders: []string{"Content-Type", "Origin", "Device-Type", "Device-Id", "*"},
        ExposeHeaders: []string{"Content-Length"},
        AllowCredentials: true,
        AllowOriginFunc: func (origin string) bool {
            return true
        },
        MaxAge: 12 * time.Hour,
    })
}
