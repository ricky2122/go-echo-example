package domain

import "time"

const BirthDayLayout = "2006-01-02"

type UserID int

type User struct {
	id       UserID
	name     string
	password string
	email    string
	birthDay time.Time
}

func NewUser(name, password, email string, birthDay time.Time) User {
	return User{
		name:     name,
		password: password,
		email:    email,
		birthDay: birthDay,
	}
}

func (u *User) GetID() UserID {
	return u.id
}

func (u *User) SetID(id int) {
	u.id = UserID(id)
}

func (u *User) GetName() string {
	return u.name
}
