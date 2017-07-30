package main

type User struct {
	name   string
	email  string
	tenant string
	cockie string
}

func (u *User) getTenant() string {
	return u.tenant
}
