package server

import (
    _ "github.com/autonomousdotai/handshake-dispatcher/config"
)

func Init() {
    r := NewRouter()
    r.Run(":8080")
}
