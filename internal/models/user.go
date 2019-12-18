package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username           string     `json:"username"`
	Password           string     `json:"password"`
	TwoFA              bool       `json:"-"`
	TwoFAKey           string     `json:"-"`
	TempTwoFAToken     string     `json:"-"`
	TempTwoFATokenDate *time.Time `json:"-"`
}
