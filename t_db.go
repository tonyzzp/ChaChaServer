package main

import (
	"fmt"
	"github.com/tonyzzp/ChaCha_Server/db"
)

func main() {
	ids := db.FindFriends(1)
	fmt.Println(ids)

	db.AddFriend(1, 1)
	fmt.Println(db.FindFriends(1))
	for _, u := range db.FetchUsrs(ids) {
		fmt.Println(u)
	}
}
