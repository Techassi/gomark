package db

import (
	"errors"

	m "github.com/Techassi/gomark/internal/models"
	"github.com/Techassi/gomark/internal/util"
)

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

func (d *DB) GetBookmarks(id uint) []m.Entity {
	e := []m.Entity{}

	d.Conn.Set("gorm:auto_preload", true).Where("owner_id = ? AND type = ?", id, "bookmark").Order("created_at DESC").Find(&e)
	return e
}

func (d *DB) GetRecentBookmarks(id uint) []m.Entity {
	e := []m.Entity{}

	d.Conn.Set("gorm:auto_preload", true).Limit(6).Where("owner_id = ? AND type = ?", id, "bookmark").Order("created_at DESC").Find(&e)
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

func (d *DB) Archived(hash string) {
	d.Conn.Set("gorm:auto_preload", true).Model(m.Entity{}).Where("type = ? AND hash = ?", "bookmark", hash).Update("archived", true)
}
