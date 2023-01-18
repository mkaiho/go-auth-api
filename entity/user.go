package entity

type User struct {
	ID    ID
	Name  string
	Email Email
}

type Users []*User
