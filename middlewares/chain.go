package middlewares

import (
    "strings"
    "strconv"
    "github.com/gin-gonic/gin"
)

func ChainMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        cid := c.Request.Header.Get("ChainId")

        cid = strings.TrimSpace(cid)

        if len(cid) > 0 {
            chainId, _ := strconv.Atoi(cid)
            c.Set("ChainId", chainId)
        }
       
        c.Next()
    }
}
