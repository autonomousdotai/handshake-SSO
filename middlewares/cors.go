package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Origin", "Payload", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	})
}
