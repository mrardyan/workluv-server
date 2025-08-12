package domain

import "strings"

type User struct {
	ID       uint
	Username string
	Email    Email
	Password Password
}

type Email string

func (e Email) IsValid() bool {
	return strings.Contains(string(e), "@")
}

type Password string

func (e Password) IsValid() bool {
	return len(e) >= 8
}
