package main

import (
	"fmt"
	"github.com/tonyzzp/ChaCha_Server/db"
	"testing"
)

func Test_db(t *testing.T) {
	u := db.FindUserByUserName("zzp")
	fmt.Println(u)
}
