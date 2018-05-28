package controllers

import (
    "net/http"
    "strconv"
    "strings"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/models"
    "github.com/autonomousdotai/handshake-dispatcher/services"
)

type HandshakeController struct{}

func (u HandshakeController) Me(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit := 100

    user, _ := c.Get("User")
    userModel := user.(models.User)
    
    userId := strconv.FormatUint(uint64(userModel.ID), 10)

    solr := new(services.SolrService)

    init_id := []string{"init_id:", userId}
    shaked_ids := []string{"shaked_ids:\"[", userId, "]\""}
    data, err := solr.List("handshake", []string{"id:*", strings.Join(init_id, ""), strings.Join(shaked_ids, "")}, (page - 1) * limit, limit) 

    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    data["page"] = page

    resp := JsonResponse{1, "", data}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Discover(c *gin.Context) {  
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit := 100
    
    solr := new (services.SolrService)
    data, err := solr.List("handshake", []string{"id:*"}, (page - 1) * limit, limit)

    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }
   
    data["page"] = page

    resp := JsonResponse{1, "", data}
    c.JSON(http.StatusOK, resp)
    return
}
