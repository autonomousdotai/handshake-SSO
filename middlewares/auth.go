package middlewares

import (
	"net/http"
	"strings"
	"log"
	"fmt"
    	"encoding/json"
	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/controllers"
	"github.com/ninjadotorg/handshake-dispatcher/models"
	"github.com/ninjadotorg/handshake-dispatcher/utils"
)

func isWhiteEndpoint(conf *viper.Viper, url string) bool {
	for _, v := range conf.GetStringSlice("white_endpoints") {
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

func AdminAuthMiddleware1() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfig()
		adminHash := conf.GetString("admin_hash")
		bearer := strings.TrimSpace(c.Request.Header.Get("AdminHash"))
		if len(bearer) < 1 || bearer != adminHash {
			c.AbortWithStatusJSON(http.StatusUnauthorized, controllers.JsonResponse{0, "Unauthorized", nil})
			return
		}

		c.Next()
	}
}
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		conf := config.GetConfig()
		if !isWhiteEndpoint(conf, c.Request.URL.Path) {


			payload := c.Request.Header.Get("AdminHash")

			p := strings.TrimSpace(payload)

			if len(p) == 0 {
				panic("Invalid admin.")
			}

			bkey := []byte(conf.GetString("secret_key"))
			uuid, err := utils.HashDecrypt(bkey, p)

			log.Println("err", err)


			if err != nil {
				panic("Invalid admin.")
			}

			user := models.User{}
			errDb := models.Database().Where("uuid = ?", uuid).First(&user).Error

			if errDb != nil {
				panic("Invalid admin.")
			}
			c.Set("User", user)
		} else {
			c.Set("WhiteUser", 1)
		}

		c.Next()

		go saveLogs(c);
	}
}

func saveLogs(c *gin.Context)  {

	query, _ := json.Marshal(c.Request.URL.Query());

	queryStr := fmt.Sprintf("%s", query)

	log.Println("c.Request.URL.Query()", queryStr)

	method := c.Request.Method
	fmt.Println("method", method)

	handleName := c.HandlerName();
	fmt.Println("HandlerName", handleName)

	controllerName := ""
	action := ""

	ss := strings.Split(handleName, "/")
	if len(ss) > 0 {
		cl :=  strings.Split(ss[len(ss)-1], ".")
		if len(cl) > 2{
			controllerName, action = cl [1], cl[2]
			controllerName = strings.Replace(controllerName, "Controller", "", -1)
			action = strings.Replace(action, "-fm", "", -1)
		}

	}
	fmt.Println("controllerName, action", controllerName, action)

	dataFormStr := "{}"
	out, err := json.Marshal(c.Request.PostForm)
	if err != nil {
		fmt.Println("c.PostForm err", err.Error())
	} else{
		if(string(out) != "null") {
			dataFormStr = string(out)
		}
		fmt.Println("c.PostForm", string(out))
	}

	//for k, v := range c.Request.PostForm {
	//    fmt.Printf("key[%s] value[%s]\n", k, v)
	//}

	path := c.Request.URL.Path
	log.Println("path", path)

	IP := c.ClientIP()
	userAgent := c.Request.Header.Get("User-Agent")

	log.Println("IP", IP)
	log.Println("userAgent", userAgent)

	// save Logs:
	db := models.Database()

	var userModel models.User

	user, _ := c.Get("User")
	userModel = user.(models.User)

	activity_log := models.ActivityLog{
				Name: controllerName,
				Action: action,
				Description: "query:" + queryStr + ", data: " + dataFormStr,
				Path: path,
				Host: IP,
				Method: method,
				UserAgent: userAgent,
				UserID: userModel.ID,
	}

	errDb := db.Save(&activity_log).Error

	if errDb != nil {
		log.Println("save log fail")
		return
	}
	log.Println("save log ok")
}