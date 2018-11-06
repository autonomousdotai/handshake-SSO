package server

import (
	_ "github.com/ninjadotorg/handshake-dispatcher/config"
)

func Init() {
	r := NewRouter()
	r.Run(":8081")
}
