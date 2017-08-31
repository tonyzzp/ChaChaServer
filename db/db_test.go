package db

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {
	u := FindUserByUserName("zzp")
	fmt.Println(u)

	u = new(User)
	u.UserName = "zzp3"
	u.PassWord = "abc"
	id, e := ORM.Insert(u)
	fmt.Println(e)
	fmt.Println(id)
	fmt.Println(u)
}
