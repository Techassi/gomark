package models

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB is the top-level Database instance.
type DB struct {
	Conn *gorm.DB
}

// Init sets up the Database connection. This function will panic if the connection
// isn't possible.
func (d *DB) Init(c *Config) {
	db, err := gorm.Open("mysql", Connection(c))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&Entity{},
		&User{},
	)

	d.Conn = db
}

// Connection creates the connection string and returns it.
func Connection(c *Config) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Database)
}

func (d *DB) ValidCredentials(u *User) bool {
	var user User

	if !d.Conn.HasTable(&user) {
		return false
	}

	r := d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user).RecordNotFound()
	u.TwoFA = user.TwoFA

	return !r
}

// SetTempTwoFAToken sets a temp 2FA token and its expiration timestamp
func (d *DB) SetTempTwoFAToken(u *User, t string, currTime time.Time) error {
	var user User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TempTwoFAToken = t
	user.TempTwoFATokenDate = &currTime

	d.Conn.Save(&user)
	return nil
}

// CheckTempTwoFAToken checks the temp 2FA token and its expiration timestamp
// TODO: Check expiration timestamp
func (d *DB) CheckTempTwoFAToken(t string) string {
	var user User

	d.Conn.Where("temp_two_fa_token = ? AND temp_two_fa_token_date > ?", t, time.Now()).First(&user)

	return user.TwoFAKey
}

func (d *DB) Update2FA(u *User, key string) {
	var user User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TwoFA = true
	user.TwoFAKey = key

	d.Conn.Save(&user)
}
