package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultController struct{}

func (d DefaultController) Home(c *gin.Context) {

	fmt.Println("r: %+v\n", c.Request)
	resp := JsonResponse{1, "Handshake REST API", c.Request.RemoteAddr}
	c.JSON(http.StatusOK, resp)
}

func (d DefaultController) NotFound(c *gin.Context) {
	resp := JsonResponse{0, "Page not found", nil}
	c.JSON(http.StatusOK, resp)
}
