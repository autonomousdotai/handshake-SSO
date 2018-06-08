package middlewares

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

func IpFilterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ipAddress := getIPAdress(c.Request)
		log.Printf("Client IP : %s", ipAddress)
		c.Set("Ip", ipAddress)
		if checkIp(net.ParseIP(ipAddress)) {
			panic("Your country is not supported")
			c.Abort()
			return
		}
		c.Next()
	}
}

type IpRange struct {
	start net.IP
	end   net.IP
}

func inRange(r IpRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

func getIPAdress(r *http.Request) string {
	var ipAddress string
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		for _, ip := range strings.Split(r.Header.Get(h), ",") {
			// header can contain spaces too, strip those out.
			ip = strings.TrimSpace(ip)
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() {
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

var (
	ipRangeArr = MakeIpFilter()
)

func MakeIpFilter() []IpRange {
	ipRangeArrTmp := make([]IpRange, 0)
	file, err := os.Open("./config/ip_ranges.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	ipArr := []map[string]string{}
	err = decoder.Decode(&ipArr)
	if err != nil {
		panic(err)
	}

	for _, v := range ipArr {
		ipRange := IpRange{
			start: net.ParseIP(v["start"]),
			end:   net.ParseIP(v["end"]),
		}
		ipRangeArrTmp = append(ipRangeArrTmp, ipRange)
	}
	sort.Slice(ipRangeArrTmp, func(i, j int) bool {
		ai := ipRangeArrTmp[i]
		aj := ipRangeArrTmp[j]
		if bytes.Compare(ai.start, aj.start) >= 0 {
			return false
		}
		return true
	})
	for _, v := range ipRangeArrTmp {
		log.Println(v.start.String() + " - " + v.end.String())
	}
	log.Println("MakeIpFilter() Finished")

	return ipRangeArrTmp
}

func checkIp(ipAddress net.IP) bool {
	index := 0
	sort.Search(len(ipRangeArr), func(i int) bool {
		if bytes.Compare(ipRangeArr[i].end, ipAddress) >= 0 {
			if bytes.Compare(ipRangeArr[i].start, ipAddress) <= 0 {
				index = i
				return false
			}
			return true
		}
		return false
	})
	if index < len(ipRangeArr) && bytes.Compare(ipAddress, ipRangeArr[index].start) >= 0 && bytes.Compare(ipAddress, ipRangeArr[index].end) <= 0 {
		return true
	}
	return false
}
