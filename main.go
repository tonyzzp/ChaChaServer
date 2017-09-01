package main

import (
	"github.com/tonyzzp/ChaCha_Server/server"
)

func main() {
	svr := server.NewServer("0.0.0.0:2626")
	svr.Start()
}
