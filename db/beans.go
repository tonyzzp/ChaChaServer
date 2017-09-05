package db

import "time"

type User struct {
	Id       int
	UserName string `orm:"unique"`
	PassWord string
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}
