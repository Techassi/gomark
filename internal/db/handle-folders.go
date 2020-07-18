package db

import (
	"errors"

	m "github.com/Techassi/gomark/internal/models"
)

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
