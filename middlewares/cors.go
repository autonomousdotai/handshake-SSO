package middlewares

import (
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"https://stag-handshake.autonomous.ai", "http://localhost:8080", "http://0.0.0.0:8080"}
    config.AllowHeaders = []string{"Content-Type", "Origin", "Payload", "*"}
    config.ExposeHeaders = []string{"Content-Length"}
    config.AllowCredentials = true
    return cors.New(config)
}
