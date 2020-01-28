package models

import (
	"errors"
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
		&User{},
		&Settings{},
		&Bookmark{},
		&Folder{},
		&EntityRelation{},
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

	if username == "" || inputPass == "" {
		return false, User{}
	}

	// Check if the 'User' table exists
	if !d.Conn.HasTable(&user) {
		return false, User{}
	}

	// Check if the user exists via the username
	d.Conn.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		return false, User{}
	}

	// Compare the provided and saved password
	correct, err := util.ComparePassword(inputPass, user.Password)
	if err != nil || !correct {
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

func (d *DB) Register(user, pass, last, first, mail string) error {
	hashedPass, err := util.HashPassword(pass)
	if err != nil {
		return err
	}

	d.Conn.Create(&User{
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
func (d *DB) SetTemp2FAToken(u *User, t string, expirationTime time.Time) error {
	var user User

	d.Conn.Where("username = ? AND password = ?", u.Username, u.Password).First(&user)
	user.TempTwoFAToken = t
	user.TempTwoFATokenDate = &expirationTime

	d.Conn.Save(&user)
	return nil
}

// CheckTemp2FAToken checks the temp 2FA token and its expiration timestamp
func (d *DB) CheckTemp2FAToken(t string) bool {
	var user User

	d.Conn.Where("temp_two_fa_token = ? AND temp_two_fa_token_date > ?", t, time.Now()).First(&user)
	if user.ID == 0 {
		return false
	}

	return true
}

func (d *DB) Update2FA(username, key string) error {
	var user User
	d.Conn.Where("username = ?", username).First(&user)
	if user.ID == 0 {
		errors.New("User with username not found")
	}

	user.TwoFA = true
	user.TwoFAKey = key

	d.Conn.Save(&user)
	return nil
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

func folderJoin(hash string, id uint) string {
	return fmt.Sprintf("SELECT c.* FROM folders JOIN entity_relations ON folders.id = entity_relations.parent_id JOIN folders c ON entity_relations.child_id = c.id WHERE folders.owner_id = %d AND folders.hash = '%s' AND entity_relations.type = 'folder'", id, hash)
}

func entityJoin(hash string, id uint) string {
	return fmt.Sprintf("SELECT c.* FROM folders JOIN entity_relations ON folders.id = entity_relations.parent_id JOIN folders c ON entity_relations.child_id = c.id WHERE folders.owner_id = %d AND folders.hash = '%s'", id, hash)
}

func subItem(parentID, childID uint, t string) string {
	return fmt.Sprintf("INSERT INTO entity_relations (parent_id,child_id,type) VALUES (%d,%d,'%s')", parentID, childID, t)
}

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// BOOKMARK FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) SaveBookmark(b Bookmark) {
	d.Conn.Save(&b)
}

func (d *DB) SaveBookmarkToFolder(parentHash string, b Bookmark) error {
	p := Folder{}

	d.Conn.Where("hash = ?", parentHash).First(&p)
	if p.ID == 0 {
		return errors.New("No such parent folder")
	}

	// Update the ChildrenCount of the parent folder. Save both the parent and
	// child folder. Create an entry in the relation table
	p.ChildrenCount++
	d.Conn.Save(&p)
	d.Conn.Save(&b)
	d.Conn.Exec(subItem(p.ID, b.ID, "bookmark"))
	return nil
}

func (d *DB) GetBookmarksByUserID(id uint) []Bookmark {
	b := []Bookmark{}

	d.Conn.Set("gorm:auto_preload", true).Where("owner_id = ?", id).Find(&b)
	return b
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) SaveFolder(f Folder) {
	d.Conn.Save(&f)
}

func (d *DB) SaveSubFolder(parentHash string, sub Folder) error {
	p := Folder{}

	d.Conn.Where("hash = ?", parentHash).First(&p)
	if p.ID == 0 {
		return errors.New("No such parent folder")
	}

	// Update the ChildrenCount of the parent folder. Save both the parent and
	// child folder. Create an entry in the relation table
	p.ChildrenCount++
	d.Conn.Save(&p)
	d.Conn.Save(&sub)
	d.Conn.Exec(subItem(p.ID, sub.ID, "folder"))
	return nil
}

func (d *DB) GetFolders(id uint) []Folder {
	f := []Folder{}

	d.Conn.Where("owner_id = ? AND has_parent = ?", id, false).Find(&f)
	return f
}

func (d *DB) GetSubFolders(hash string, id uint) []Folder {
	f := []Folder{}
	d.Conn.Raw(folderJoin(hash, id)).Scan(&f)
	return f
}

func (d *DB) GetSubEntities(hash string, id uint) []Folder {
	f := []Folder{}
	d.Conn.Raw(entityJoin(hash, id)).Scan(&f)
	return f
}
