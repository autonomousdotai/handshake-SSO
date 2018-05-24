package main

import (
    _ "fmt"
    "github.com/autonomousdotai/handshake-dispatcher/config"
    "github.com/autonomousdotai/handshake-dispatcher/server"
)

func main() {
    config.Init()
    //db.Init()
    server.Init()
}
