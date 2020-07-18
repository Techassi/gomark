package db

import (
	"fmt"

	m "github.com/Techassi/gomark/internal/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB is the top-level Database instance.
type DB struct {
	Conn *gorm.DB
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
