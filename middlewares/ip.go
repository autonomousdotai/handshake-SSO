package middlewares

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func IpMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := getIPAdress(c.Request)
		log.Printf("Client IP : %s", ipAddress)
		c.Set("IpAddress", ipAddress)
		c.Next()
	}
}

func getIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if realIP == nil {
				// bad address, go to next
				continue
			} else {
				ipAddress = ip
				goto Done
			}
		}
	}
Done:
	return ipAddress
}
