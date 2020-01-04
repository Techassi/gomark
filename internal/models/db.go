package models

import (
	"fmt"
	"time"

	"github.com/Techassi/gomark/internal/util"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB is the top-level Database instance.
type DB struct {
	Conn *gorm.DB
}

type DBModel struct {
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// GENERAL FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// Init sets up the Database connection. This function will panic if the
// connection isn't possible.
func (d *DB) Init(c *Config) {
	db, err := gorm.Open("mysql", Connection(c))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&Entity{},
		&User{},
		&Settings{},
		&BookmarkData{},
		&FolderData{},
		&NoteData{},
	)

	d.Conn = db
}

// Connection creates the connection string and returns it.
func Connection(c *Config) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Database)
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// AUTH FUNCTIONS ////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) ValidateCredentials(username, inputPass string) (bool, User) {
	user := User{}

	// Check if the 'User' table exists
	if !d.Conn.HasTable(&user) {
		return false, User{}
	}

	// Check if the user exists via the username
	r := d.Conn.Where("username = ?", username).First(&user).RecordNotFound()
	if r {
		return false, User{}
	}

	// Compare the provided and saved password
	_, err := util.ComparePassword(inputPass, user.Password)
	if err != nil {
		return false, User{}
	}

	return true, user
}

func (d *DB) ValidateNewCredentials(u, p string) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// REGISTER FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) Register(user, pass, last, first, mail string) {
	hashedPass, _ := util.HashPassword(pass)

	d.Conn.Create(&User{
		Username:  user,
		Password:  hashedPass,
		Lastname:  last,
		Firstname: first,
		EMail:     mail,
	})
}

////////////////////////////////////////////////////////////////////////////////
//////////////////////////////// 2FA FUNCTIONS /////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// SetTempTwoFAToken sets a temp 2FA token and its expiration timestamp
func (d *DB) SetTemp2FAToken(u *User, t string, expirationTime time.Time) error {
	var user User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TempTwoFAToken = t
	user.TempTwoFATokenDate = &expirationTime

	d.Conn.Save(&user)
	return nil
}

// CheckTempTwoFAToken checks the temp 2FA token and its expiration timestamp
// TODO: Check expiration timestamp
func (d *DB) CheckTemp2FAToken(t string) *User {
	var user User

	d.Conn.Where("temp_two_fa_token = ? AND temp_two_fa_token_date > ?", t, time.Now()).First(&user)

	return &user
}

func (d *DB) Update2FA(u *User, key string) {
	var user User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TwoFA = true
	user.TwoFAKey = key

	d.Conn.Save(&user)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// SETTINGS FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) DefaultSettings() Settings {
	if settings := d.LoadSettings(); settings.ID != 0 {
		return settings
	}

	s := Settings{
		RegisterEnabled: true,
		SharingEnabled:  false,
	}
	d.Conn.Create(&s)

	return s
}

func (d *DB) LoadSettings() Settings {
	s := Settings{}
	d.Conn.Order("created_at DESC").First(&s)

	return s
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// ENTITY FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) SaveEntity(e Entity, id uint) {
	u := User{}
	d.Conn.First(&u, id)
	d.Conn.Model(&u).Association("Entities").Append(e)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) GetBookmarksByUserID(id uint) []Entity {
	u := User{}
	r := []Entity{}
	d.Conn.First(&u, id)
	d.Conn.Set("gorm:auto_preload", true).Model(&u).Where("type = ?", "bookmark").Association("Entities").Find(&r)
	return r
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) GetFoldersByUserID(id uint) []Entity {
	u := User{}
	r := []Entity{}
	d.Conn.First(&u, id)
	d.Conn.Set("gorm:auto_preload", true).Model(&u).Where("type = ?", "folder").Association("Entities").Find(&r)
	return r
}

func (d *DB) SaveEntityToFolder(e Entity, hash string) {
	r := Entity{}
	d.Conn.Set("gorm:auto_preload", true).Where("type = ? AND hash = ?", "folder", hash).First(&r)

	r.FolderData.ChildEntities = append(r.FolderData.ChildEntities, e)
	d.Conn.Save(&r)
}
