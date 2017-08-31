package db

func FindUserById(id int32) *User {
	u := new(User)
	u.Id = id
	e := ORM.Read(u)
	if e == nil {
		return u
	} else {
		return nil
	}
}

func FindUserByUserName(userName string) *User {
	u := new(User)
	u.UserName = userName
	e := ORM.Read(u, "UserName")
	if e == nil {
		return u
	} else {
		return nil
	}
}

func InsertUser(u *User) bool {
	_, e := ORM.Insert(u)
	return e != nil
}
