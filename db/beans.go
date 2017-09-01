package db

import "time"

type User struct {
	Id       int32
	UserName string    `orm:"type(text);unique"`
	PassWord string    `orm:"type(text)"`
	Created  time.Time `orm:"auto_now_add;type(datetime)"`
}
