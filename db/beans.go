package db

type User struct {
	Id       int32  `orm:"pk"`
	UserName string `orm:"unique"`
	PassWord string
}
