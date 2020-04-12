package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username           string     `json:"username"`
	Password           string     `json:"password"`
	Role               Role       `json:"role"`
	Firstname          string     `json:"firstname"`
	Lastname           string     `json:"lastname"`
	EMail              string     `json:"email"`
	TwoFA              bool       `json:"-"`
	TwoFAKey           string     `json:"-"`
	TempTwoFAToken     string     `json:"-"`
	TempTwoFATokenDate *time.Time `json:"-"`
}

type Role int

const (
	Admin Role = iota
	Normal
)

func (r Role) String() string {
	return [...]string{"Admin", "Normal"}[r]
}
