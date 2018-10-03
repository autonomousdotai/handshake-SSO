package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ninjadotorg/handshake-dispatcher/config"
	"github.com/ninjadotorg/handshake-dispatcher/controllers"
	"github.com/ninjadotorg/handshake-dispatcher/middlewares"
	"github.com/ninjadotorg/handshake-dispatcher/models"
)

func NewRouter() *gin.Engine {
	// Logger
	logFile, err := os.OpenFile("logs/log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter) // You may need this
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.ChainMiddleware())
	// router.Use(middlewares.IpFilterMiddleware())

	defaultController := new(controllers.DefaultController)
	router.GET("/", defaultController.Home)

	nonceController := new(controllers.NonceController)
	nonceGroup := router.Group("nonce")
	{
		nonceGroup.GET("/get", nonceController.Get)
		nonceGroup.POST("/set", nonceController.Set)
	}

	userController := new(controllers.UserController)
	verificationController := new(controllers.VerifierController)
	userGroup := router.Group("user")
	{
		userGroup.GET("/profile", middlewares.AuthMiddleware(), userController.Profile)
		userGroup.GET("/username-exist", middlewares.AuthMiddleware(), userController.UsernameExist)
		userGroup.GET("/username/:id", userController.Username)
		userGroup.POST("/profile", middlewares.AuthMiddleware(), userController.UpdateProfile)

		userGroup.POST("/sign-up", userController.SignUp)
		userGroup.POST("/free-rinkeby-eth", middlewares.AuthMiddleware(), userController.FreeRinkebyEther)

		userGroup.GET("/referred", middlewares.AuthMiddleware(), userController.Referred)
		userGroup.POST("/complete-profile", middlewares.AuthMiddleware(), userController.CompleteProfile)
		userGroup.POST("/subscribe", userController.Subscribe)
		userGroup.POST("/count_subscribed_user", userController.CountSubscribedUsers)

		userGroup.GET("/check-email-exist", middlewares.AuthMiddleware(), userController.CheckEmailExist)
		userGroup.POST("/notification", middlewares.AuthMiddleware(), userController.Notification)

		userGroup.POST("/verification/phone/start", middlewares.AuthMiddleware(), verificationController.SendPhoneVerification)
		userGroup.POST("/verification/phone/check", middlewares.AuthMiddleware(), verificationController.CheckPhoneVerification)
		userGroup.POST("/verification/email/start", middlewares.AuthMiddleware(), verificationController.SendEmailVerification)
		userGroup.POST("/verification/email/check", middlewares.AuthMiddleware(), verificationController.CheckEmailVerification)
		userGroup.POST("/verification/redeem-code/check", middlewares.AuthMiddleware(), verificationController.CheckRedeemCodeVerification)
	}

	handshakeController := new(controllers.HandshakeController)
	handshakeGroup := router.Group("handshake")
	{
		handshakeGroup.GET("/me", middlewares.AuthMiddleware(), handshakeController.Me)
		handshakeGroup.GET("/discover", middlewares.AuthMiddleware(), handshakeController.Discover)
		handshakeGroup.POST("/create", middlewares.AuthMiddleware(), handshakeController.Create)
		handshakeGroup.POST("/update", middlewares.AuthMiddleware(), handshakeController.Update)
		handshakeGroup.POST("/delete", middlewares.AuthMiddleware(), handshakeController.Delete)
	}

	systemController := new(controllers.SystemController)
	systemGroup := router.Group("system")
	{
		systemGroup.GET("/user/:id", systemController.User)
		systemGroup.POST("/betsuccess/:id", systemController.BetSuccess)
	}

	conf := config.GetConfig()
	for ex, ep := range conf.GetStringMap("forwarding") {
		endpoint := ep
		router.Any(ex, middlewares.AuthMiddleware(), func(c *gin.Context) {
			Forwarding(c, &endpoint, "")
		})
		router.Any(ex+"/*path", middlewares.AuthMiddleware(), func(c *gin.Context) {
			path := c.Param("path")
			Forwarding(c, &endpoint, path)
		})
	}

	router.NoRoute(defaultController.NotFound)

	return router
}

type ForwardingTransport struct {
}

func (t *ForwardingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	start := time.Now()
	response, err := http.DefaultTransport.RoundTrip(request)

	if err != nil {
		fmt.Println("\n\ncame in error resp here", err)
		return nil, err
	}

	elapsed := time.Since(start)

	_, err = httputil.DumpResponse(response, true)
	if err != nil {
		fmt.Println("\n\ndump response error")
	}

	log.Printf("%s - %d\n", request.Method+request.URL.Path, elapsed.Nanoseconds)
	return response, nil
}

func Forwarding(c *gin.Context, endpoint *interface{}, path string) {
	r := c.Request
	w := c.Writer

	user, _ := c.Get("User")

	url, _ := url.Parse((*endpoint).(string) + path)
	director := func(req *http.Request) {
		req.URL.Scheme = url.Scheme
		req.URL.Host = url.Host
		req.URL.Path = url.Path

		for k, _ := range r.Header {
			v := c.GetHeader(k)
			req.Header.Set(k, v)
		}
		req.Header.Set("Uid", strconv.FormatUint(uint64((user.(models.User)).ID), 10))
		req.Header.Set("Fcm-Token", (user.(models.User)).FCMToken)
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.Transport = &ForwardingTransport{}
	proxy.ServeHTTP(w, r)
}
