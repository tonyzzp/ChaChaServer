package db

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
)

var friends = make(map[int][]int)

//  查询 userid 的好友列表
func FindFriends(userid int) []int {
	r := friends[userid]
	if r == nil {
		var m []orm.Params
		_, e := ORM.Raw("select friend_id from friends where user_id=?", userid).Values(&m)
		if e != nil {
			panic(e)
		}
		l := len(m)
		r = make([]int, l)
		for i := 0; i < l; i++ {
			id, _ := strconv.Atoi(m[i]["friend_id"].(string))
			r[i] = id
		}
		friends[userid] = r
	}
	return r
}

func AddFriend(userid int, friendid int) {
	if userid == friendid {
		logs.Warn("不能添加自己为好友")
		return
	}
	count := 0
	ORM.Raw("select count(*) from friends where user_id=? and friend_id=?",
		userid, friendid).QueryRow(&count)
	if count > 0 {
		logs.Warn("好友已存在 %v %v", userid, friendid)
		return
	}
	_, e := ORM.Raw("insert into friends(user_id,friend_id) values(?,?)", userid, friendid).Exec()
	if e != nil {
		panic(e)
	}
	f := friends[userid]
	if f != nil {
		f = append(f, friendid)
		friends[userid] = f
	}
}

func RemoveFriend(userid int, friendid int) {
	ORM.Raw("delete from friends where user_id=? and friend_id=?", userid, friendid).Exec()
	ORM.Raw("delete from friends where user_id=? and friend_id=?", friendid, userid).Exec()
	friends[userid] = nil
	friends[friendid] = nil
}
