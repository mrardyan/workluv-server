package domain

import "strings"

type Account struct {
	ID       uint
	Username string
	Email    Email
	Password Password
}

// TableName overrides the table name for GORM
func (Account) TableName() string {
	return "accounts"
}

type Email string

func (e Email) IsValid() bool {
	return strings.Contains(string(e), "@")
}

type Password string

func (e Password) IsValid() bool {
	return len(e) >= 8
}
