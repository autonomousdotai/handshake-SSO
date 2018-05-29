package utils

import (
    "fmt"
    "github.com/autonomousdotai/handshake-dispatcher/config"
)

func GetForwardingEndpoint(t string) (string, bool) {
    conf := config.GetConfig()
    var endpoint string
    found := false

    for n, ep := range conf.GetStringMap("forwarding") {
        fmt.Println(n, ep)
        if n == t {
            endpoint = ep.(string)
            found = true
            break
        }
    }

    return endpoint, found
}

func GetServicesEndpoint(t string) (string, bool) {
    conf := config.GetConfig()
    var endpoint string
    found := false
    
    for n, ep := range conf.GetStringMap("services") {
        if n == t {
            endpoint = ep.(string)
            found = true
            break
        }
    }

    return endpoint, found
}
