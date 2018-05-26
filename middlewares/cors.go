package middlewares

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    config.AllowHeaders = []string{"Content-Type", "Origin", "Payload", "*"}
    config.ExposeHeaders = []string{"Content-Length"}

    return cors.New(config)
}
