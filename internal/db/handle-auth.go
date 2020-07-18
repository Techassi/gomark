package db

import (
	"errors"
	"time"

	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/util"
)

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// AUTH FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) ValidateCredentials(username, inputPass string) (bool, m.User) {
	user := m.User{}

	if username == "" || inputPass == "" {
		return false, m.User{}
	}

	// Check if the 'User' table exists
	if !d.Conn.HasTable(&user) {
		return false, m.User{}
	}

	// Check if the user exists via the username
	d.Conn.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return false, m.User{}
	}

	// Compare the provided and saved password
	correct, err := util.ComparePassword(inputPass, user.Password)
	if err != nil || !correct {
		return false, m.User{}
	}

	return true, user
}

func (d *DB) ValidateNewCredentials(u, p string) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// REGISTER FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) Register(user, pass, last, first, mail string) error {
	hashedPass, err := util.HashPassword(pass)
	if err != nil {
		return err
	}

	d.Conn.Create(&m.User{
		Username:  user,
		Password:  hashedPass,
		Lastname:  last,
		Firstname: first,
		EMail:     mail,
	})
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// 2FA FUNCTIONS /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SetTempTwoFAToken sets a temp 2FA token and its expiration timestamp
func (d *DB) SetTemp2FAToken(u *m.User, t string, expirationTime time.Time) error {
	var user m.User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TempTwoFAToken = t
	user.TempTwoFATokenDate = &expirationTime

	d.Conn.Save(&user)
	return nil
}

// CheckTemp2FAToken checks the temp 2FA token and its expiration timestamp
func (d *DB) CheckTemp2FAToken(t string) bool {
	var user m.User

	d.Conn.Where("temp_two_fa_token = ? AND temp_two_fa_token_date > ?", t, time.Now()).First(&user)
	if user.ID == 0 {
		return false
	}

	return true
}

func (d *DB) Update2FA(username, key string) error {
	var user m.User
	d.Conn.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		errors.New("User with username not found")
	}

	user.TwoFA = true
	user.TwoFAKey = key

	d.Conn.Save(&user)
	return nil
}
