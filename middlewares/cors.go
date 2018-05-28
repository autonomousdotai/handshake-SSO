package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Content-Type", "Origin", "Payload", "*"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}
	return cors.New(config)
}
