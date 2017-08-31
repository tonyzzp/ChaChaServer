package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

var ORM orm.Ormer

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "/db.db")
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
	ORM = orm.NewOrm()
}
