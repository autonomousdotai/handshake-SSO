package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/controllers"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/utils"
)

func isWhiteEndpoint(conf *viper.Viper, url string) bool {
	for _, v := range conf.GetStringSlice("white_endpoints") {
		fmt.Println("DTHTRONG ", v, url)
		if v == url {
			return true
		}
	}
	return false
}

// AuthMiddleware : verify valid user or not
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfig()
		if !isWhiteEndpoint(conf, c.Request.URL.Path) {
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
		} else {
			c.Set("WhiteUser", 1)
		}

		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfig()
		adminHash := conf.GetString("admin_hash")
		bearer := strings.TrimSpace(c.Request.Header.Get("Bearer"))
		if len(bearer) < 1 || bearer != adminHash {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.JsonResponse{0, "Unauthorized", nil})
			return
		}

		c.Next()
	}
}
