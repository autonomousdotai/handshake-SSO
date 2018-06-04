package main

import (
    _ "fmt"
    "github.com/ninjadotorg/handshake-dispatcher/config"
    "github.com/ninjadotorg/handshake-dispatcher/server"
)

func main() {
    config.Init()
    //db.Init()
    server.Init()
}
