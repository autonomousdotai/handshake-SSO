package server

import (
    "fmt"
    "net/http"
    "net/http/httputil"
    "net/url"
    "time"
    "log"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/controllers"
    "github.com/autonomousdotai/handshake-dispatcher/middlewares"
    "github.com/autonomousdotai/handshake-dispatcher/config"
)

func NewRouter() *gin.Engine {
    router := gin.Default()
    router.Use(gin.Logger())
    router.Use(gin.Recovery())

    api := router.Group("api")
    {
        userController := new(controllers.UserController)
        userGroup := api.Group("user")
        {
            userGroup.GET("/profile", middlewares.AuthMiddleware(), userController.Profile)

            userGroup.POST("/profile", middlewares.AuthMiddleware(), userController.UpdateProfile)

            userGroup.POST("/sign-up", userController.SignUp)
        }
        
        conf := config.GetConfig()
        for ex, ep := range conf.GetStringMap("forwarding") { 
            endpoint := ep
            api.Any(ex, middlewares.AuthMiddleware(), func(c *gin.Context) {
                Forwarding(c, &endpoint, "")
            })
            api.Any(ex + "/*path", middlewares.AuthMiddleware(), func(c *gin.Context) {
                path := c.Param("path")
                Forwarding(c, &endpoint, path)
            })
        }
    }

    return router
}

type ForwardingTransport struct {
}

func (t *ForwardingTransport) RoundTrip(request *http.Request) (*http.Response, error) {
    start := time.Now()
    response, err := http.DefaultTransport.RoundTrip(request)

    if err != nil {
        fmt.Println("\n\ncame in error resp here", err)
    }

    elapsed := time.Since(start)
    
    body, err := httputil.DumpResponse(response, true)
    if err != nil {
        fmt.Println("\n\ndump response error");
    }

    log.Printf("%s - %d\n", request.Method + request.URL.Path, elapsed.Nanoseconds)
    log.Println("Response Body:", string(body))
    return response, nil
}

func Forwarding(c *gin.Context, endpoint *interface{}, path string) { 
    r := c.Request
    w := c.Writer
    
    url, _ := url.Parse((*endpoint).(string) + path)
    director := func(req *http.Request) {
        req.URL.Scheme = url.Scheme
        req.URL.Host = url.Host
        req.URL.Path = url.Path

        for k, _ := range r.Header {
            v := c.GetHeader(k)
            req.Header.Set(k, v)
        }
    }
    proxy := &httputil.ReverseProxy{Director: director} 
    proxy.Transport = &ForwardingTransport{}
    proxy.ServeHTTP(w, r)
}
