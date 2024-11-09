package domain

import "time"

const BirthDayLayout = "2006-01-02"

type UserID int

func (i UserID) Int() int {
	return int(i)
}

type BirthDay time.Time

func (b BirthDay) Time() time.Time {
	return time.Time(b)
}

func (b BirthDay) String() string {
	return time.Time(b).Format(BirthDayLayout)
}

type User struct {
	id       UserID
	name     string
	password string
	email    string
	birthDay BirthDay
}

func NewUser(name, password, email string, birthDay time.Time) User {
	return User{
		name:     name,
		password: password,
		email:    email,
		birthDay: BirthDay(birthDay),
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

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetBirthDay() BirthDay {
	return u.birthDay
}
