package middlewares

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowHeaders = []string{"Content-Type", "Origin", "Payload", "*"}
    config.ExposeHeaders = []string{"Content-Length"}
    config.AllowCredentials = true
    config.AllowAllOrigins = true
    return cors.New(config)
}
