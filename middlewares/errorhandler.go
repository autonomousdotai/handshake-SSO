package middlewares

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func(c *gin.Context) {
            if rec := recover(); rec != nil {
                c.JSON(http.StatusOK, gin.H{"status": 0, "message": rec})
                c.Abort()
                return;
            }
        }(c)
        c.Next()
    }
}
