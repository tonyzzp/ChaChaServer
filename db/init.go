package db

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tonyzzp/ChaCha_Server/cfg"
)

var ORM orm.Ormer

func init() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", cfg.DBPATH)
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
	ORM = orm.NewOrm()

	_, e := ORM.Raw(`create table if not exists friends(
		id integer primary key autoincrement,
		user_id integer,
		friend_id integer)
	`).Exec()
	if e != nil {
		panic(e)
	}
	_, e = ORM.Raw("create index if not exists index_friends on friends (user_id,friend_id)").Exec()
	if e != nil {
		panic(e)
	}
}
