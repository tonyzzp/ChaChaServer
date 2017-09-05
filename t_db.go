package main

import (
	"fmt"
	"github.com/tonyzzp/ChaCha_Server/db"
)

func main() {
	ids := db.FindFriends(1)
	fmt.Println(ids)
}
