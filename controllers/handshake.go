package controllers

import (
    "fmt"
    "net/http"
    "strconv"
    "encoding/json"
    "strings"
    "github.com/gin-gonic/gin"

    "github.com/ninjadotorg/handshake-dispatcher/models"
    "github.com/ninjadotorg/handshake-dispatcher/services"
)

const LIMIT = 100

type HandshakeController struct{}

func (u HandshakeController) Me(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

    user, _ := c.Get("User")
    userModel := user.(models.User)

    chainId, hasChain := c.Get("ChainId")

    if !hasChain {
        resp := JsonResponse{0, "Invalid Chain Id", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    var q, fq, s string

    // sort
    s = "def(last_update_at_i, 0) desc"

    // query
    q = "id:*"

    // filter query
    search_init_user_id := fmt.Sprintf("init_user_id_i: %d", userModel.ID)
    search_shaked_user_ids := fmt.Sprintf("shake_user_ids_is: %d", userModel.ID)
    search_chain_id := fmt.Sprintf("chain_id_i: %d", chainId)
    fq = fmt.Sprintf("(%s OR %s) AND %s", search_init_user_id, search_shaked_user_ids, search_chain_id)

    data, err := solrService.List("handshake", q, fq, (page - 1) * LIMIT, LIMIT, s, nil)

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
    cq := c.DefaultQuery("custom_query", "_")
    t := c.DefaultQuery("type", "_")
    pt := c.DefaultQuery("pt", "0,0")
    sfield := c.DefaultQuery("sfield", "location_p")
    d := c.DefaultQuery("d", "10")

    chainId, hasChain := c.Get("ChainId")

    if !hasChain {
        resp := JsonResponse{0, "Invalid Chain Id", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return
    }

    user, _ := c.Get("User")
    userModel := user.(models.User)

    var q, fq, s string

    // sort
    // The last condition is for sort by distance
    s = "sum(mul(def(shake_count_i,0), 8),mul(def(comment_count_i,0), 4),mul(def(view_count_i,0), 2),if(def(last_update_at_i, 0), div(last_update_at_i, 3000000), 0)) desc"

    // filter query
    //fq = fmt.Sprintf("is_private_i:0 AND chain_id_i:%d AND -init_user_id_i:%d", chainId, userModel.ID)
    fq = fmt.Sprintf("is_private_i:0 AND chain_id_i:%d", chainId, userModel.ID)

    // query
    if kws != "_" {
        words := strings.Fields(kws)
        fmt.Println(words, len(words))
        search_text_search := ""
        for _,word := range words {
            if len(search_text_search) > 0 {
                search_text_search = fmt.Sprintf("%s *%s*", search_text_search, word)
            } else {
                search_text_search = fmt.Sprintf("*%s*", word)
            }
        }
        search_text_search = fmt.Sprintf("text_search_ss:(%s)", search_text_search)
        if len(q) > 0 {
            q = fmt.Sprintf("%s AND %s", q, search_text_search)
        } else {
            q = search_text_search
        }
    }

    if cq != "_" {
        if len(q) > 0 {
            q = fmt.Sprintf("%s AND %s", q, cq)
        } else {
            q = cq
        }
    }

    if t != "_" {
        search_type := fmt.Sprintf("type_i:%s", t)
        fq = fmt.Sprintf("%s AND %s", fq, search_type)
    }

    if len(q) == 0 {
        q = "id:*"
    }

    var ss *services.SolrSpatial
    if pt != "0,0" {
        ss = &services.SolrSpatial{
            Pt: pt,
            SField: sfield,
            D: d,
        }
    }

    data, err := solrService.List("handshake", q, fq, (page - 1) * LIMIT, LIMIT, s, ss)
    if err != nil {
        resp := JsonResponse{0, err.Error(), nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return
    }

    data["page"] = page
    data["page_size"] = LIMIT

    resp := JsonResponse{1, "", data}
    c.JSON(http.StatusOK, resp)
    return
}

func (u HandshakeController) Create(c *gin.Context) {
    chainId, hasChain := c.Get("ChainId")

    if !hasChain {
        resp := JsonResponse{0, "Invalid Chain Id", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return
    }

    data := c.PostForm("data")

    var handshake map[string]interface{}
    json.Unmarshal([]byte(data), &handshake)

    handshake["chain_id_i"] = chainId

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
    chainId, hasChain := c.Get("ChainId")

    if !hasChain {
        resp := JsonResponse{0, "Invalid Chain Id", nil}
        c.JSON(http.StatusOK, resp)
        c.Abort()
        return;
    }

    data := c.PostForm("data")

    var handshake map[string]interface{}
    json.Unmarshal([]byte(data), &handshake)

    handshake["chain_id_i"] = chainId

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
