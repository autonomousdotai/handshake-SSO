package controllers

import (
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/models"
    "github.com/autonomousdotai/handshake-dispatcher/services"
)

const LIMIT = 100

type HandshakeController struct{}

func (u HandshakeController) Me(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

    user, _ := c.Get("User")
    userModel := user.(models.User)
    

    solr := new(services.SolrService)

    init_id := fmt.Sprintf("init_user_id_i: %d", userModel.ID)
    //shaked_ids := []string{"shaked_ids_is:\"[", userId, "]\""}
    data, err := solr.List("handshake", []string{init_id}, (page - 1) * LIMIT, LIMIT) 

    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    data["page"] = page
    data["page_size"] = LIMIT

    resp := JsonResponse{1, "", data}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Discover(c *gin.Context) {  
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    
    solr := new (services.SolrService)
    data, err := solr.List("handshake", []string{"id:*"}, (page - 1) * LIMIT, LIMIT)

    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }
   
    data["page"] = page
    data["page_size"] = LIMIT

    resp := JsonResponse{1, "", data}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Create(c *gin.Context) {
    data := c.PostForm("data")

    var handshake map[string]interface{}
    json.Unmarshal([]byte(data), &handshake)

    solr := new(services.SolrService)
    result, _ := solr.Create("handshake", handshake)

    if !result {
        resp := JsonResponse{0, "Create handshake fail", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    resp := JsonResponse{1, "", handshake}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Update(c *gin.Context) {
    data := c.PostForm("data")

    var handshake map[string]interface{}
    json.Unmarshal([]byte(data), &handshake)

    solr := new(services.SolrService)
    result, _ := solr.Update("handshake", handshake)

    if !result {
        resp := JsonResponse{0, "Update handshake fail", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    resp := JsonResponse{1, "", handshake}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Delete(c *gin.Context) {
    id := c.PostForm("id")
    fmt.Println("delete id", id)
    solr := new(services.SolrService)
    result, _ := solr.Delete("handshake", id)

    if !result {
        resp := JsonResponse{0, "Delete handshake fail", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    resp := JsonResponse{1, "", result}
    c.JSON(http.StatusOK, resp)
    return
}
