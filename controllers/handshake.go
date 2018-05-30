package controllers

import (
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "github.com/gin-gonic/gin"

    "github.com/autonomousdotai/handshake-dispatcher/models"
)

const LIMIT = 100

type HandshakeController struct{}

func (u HandshakeController) Me(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

    user, _ := c.Get("User")
    userModel := user.(models.User)
    
    var q, fq, s []string

    s = append(s, "def(init_at_i, 0) desc")

    search_init_user_id := fmt.Sprintf("init_user_id_i: %d", userModel.ID)
    search_shaked_user_ids := fmt.Sprintf("shaked_user_ids_is: %d", userModel.ID)  
    combine_search := fmt.Sprintf("%s %s", search_init_user_id, search_shaked_user_ids)
    
    q = append(q, combine_search)

    data, err := solrService.List("handshake", q, fq, (page - 1) * LIMIT, LIMIT, s) 
 
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
    kws := c.DefaultQuery("query", "_")
    t := c.DefaultQuery("type", "_")
    
    has_cond := false
    var q, fq, s []string
    
    s = append(s, "sum(mul(def(shake_count_i,0), 8),mul(def(comment_count_i,0), 4),mul(def(view_count_i,0), 2),if(def(last_update_at_i, 0), div(last_update_at_i, 3000000), 0)) desc")

    q = append(q, "is_private_i:0")
    fq = append(fq, "is_private_i:[* TO *]")

    if kws!= "_" {
        has_cond = true
        search_text_search := fmt.Sprintf("text_search_ss:*\"%s\"*", kws)
        has_text_search := fmt.Sprint("text_search_ss:[* TO *]")
        q = append(q, search_text_search)
        fq = append(fq, has_text_search)
    }

    if t != "_" {
        has_cond = true
        search_type := fmt.Sprintf("type_i:%s", t)
        has_type := fmt.Sprint("type_i:[* TO *]")
        q = append(q, search_type)
        fq = append(fq, has_type)
    }

    if !has_cond {
        q = append(q, "id:*")
    }

    data, err := solrService.List("handshake", q, fq, (page - 1) * LIMIT, LIMIT, s)

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

    result, _ := solrService.Create("handshake", handshake)

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

    result, _ := solrService.Update("handshake", handshake)

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
    result, _ := solrService.Delete("handshake", id)

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
