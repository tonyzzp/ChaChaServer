package main

import (
	"fmt"
	"github.com/tonyzzp/ChaCha_Server/db"
	"github.com/tonyzzp/ChaCha_Server/server"
)

func main() {
	svr := server.NewServer("0.0.0.0:2626")
	svr.Start()
	fmt.Println(db.ORM)
}
