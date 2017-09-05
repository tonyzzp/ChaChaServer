package db

func FindUserById(id int) *User {
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

func FetchUsrs(ids []int) []*User {
	r := make([]*User, len(ids))
	for i := 0; i < len(ids); i++ {
		id := ids[i]
		u := FindUserById(id)
		r[i] = u
	}
	return r
}
