package db

import (
	"errors"
	"fmt"
	"time"

	m "github.com/Techassi/gomark/internal/models"
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
func (d *DB) Init(c *m.Config) {
	db, err := gorm.Open("mysql", Connection(c))
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&m.User{},
		&m.Settings{},
		&m.Bookmark{},
		&m.Folder{},
		&m.Note{},
		&m.Entity{},
		&m.EntityRelation{},
	)

	d.Conn = db
}

// Connection creates the connection string and returns it.
func Connection(c *m.Config) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Database)
}

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

////////////////////////////////////////////////////////////////////////////////
////////////////////////////// SETTINGS FUNCTIONS //////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) DefaultSettings() m.Settings {
	if settings := d.LoadSettings(); settings.ID != 0 {
		return settings
	}

	s := m.Settings{
		RegisterEnabled: true,
		SharingEnabled:  false,
	}
	d.Conn.Create(&s)

	return s
}

func (d *DB) LoadSettings() m.Settings {
	s := m.Settings{}
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

func (d *DB) SaveEntity(e m.Entity) {
	d.Conn.Save(&e)
}

func (d *DB) SaveBookmarkToFolder(parentHash string, b m.Bookmark) error {
	p := m.Folder{}

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

func (d *DB) GetBookmarksByUserID(id uint) []m.Entity {
	e := []m.Entity{}

	d.Conn.Set("gorm:auto_preload", true).Where("owner_id = ? AND type = ?", id, "bookmark").Find(&e)
	return e
}

func (d *DB) ShareBookmark(hash string) (string, error) {
	var (
		e         m.Entity
		shareHash string
		err       error
	)

	for e.ID != 0 || shareHash == "" {
		shareHash, err = util.RandomCryptoString(16)
		if err != nil {
			return "", err
		}

		d.Conn.Where("share_hash = ?", shareHash).Find(&e)
	}

	d.Conn.Model(&e).Where("hash = ? AND type = ?", hash, "bookmark").Update("share_hash", shareHash)
	return shareHash, nil
}

func (d *DB) GetShared(shareHash string) m.Entity {
	var e m.Entity
	d.Conn.Set("gorm:auto_preload", true).Where("share_hash = ?", shareHash).First(&e)
	return e
}

func (d *DB) GetBookmarkByHash(hash string) m.Entity {
	var e m.Entity
	d.Conn.Set("gorm:auto_preload", true).Where("hash = ?", hash).First(&e)
	return e
}

func (d *DB) GetBookmarks() []m.Entity {
	var e []m.Entity
	d.Conn.Set("gorm:auto_preload", true).Where("type = ?", "bookmark").Find(&e)
	return e
}

////////////////////////////////////////////////////////////////////////////////
/////////////////////////////// FOLDER FUNCTIONS ///////////////////////////////
////////////////////////////////////////////////////////////////////////////////

func (d *DB) SaveFolder(f m.Folder) {
	d.Conn.Save(&f)
}

func (d *DB) SaveSubFolder(parentHash string, sub m.Folder) error {
	p := m.Folder{}

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

func (d *DB) GetFolders(id uint) []m.Folder {
	f := []m.Folder{}

	d.Conn.Where("owner_id = ? AND has_parent = ?", id, false).Find(&f)
	return f
}

func (d *DB) GetSubFolders(hash string, id uint) []m.Folder {
	f := []m.Folder{}
	d.Conn.Raw(folderJoin(hash, id)).Scan(&f)
	return f
}

func (d *DB) GetSubEntities(hash string, id uint) []m.Folder {
	f := []m.Folder{}
	d.Conn.Raw(entityJoin(hash, id)).Scan(&f)
	return f
}
