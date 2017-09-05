package db

import "github.com/astaxie/beego/orm"

func FindFriends(userid int) []int {
	var m []orm.Params
	_, e := ORM.QueryTable("friends").Values(&m)
	if e != nil {
		panic(e)
	}
	l := len(m)
	r := make([]int, l)
	for i := 0; i < l; i++ {
		id := m[i]["friend_id"].(int)
		r[i] = id
	}
	return r
}
