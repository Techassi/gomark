package models

import (
	"github.com/jinzhu/gorm"
)

type Entity struct {
	gorm.Model `json:"-"`
	OwnerID    uint      `json:"-"`
	Type       string    `json:"type"`
	Hash       string    `json:"hash"`
	ShareHash  string    `json:"share_hash"`
	Name       string    `json:"name"`
	Shared     bool      `json:"shared"`
	Pinned     bool      `json:"pinned"`
	Bookmark   *Bookmark `json:"bookmark"`
	Folder     *Folder   `json:"folder"`
	Note       *Note     `json:"note"`
	HasParent  bool      `json:"has_parent"`
}

type Bookmark struct {
	gorm.Model  `json:"-"`
	EntityID    uint   `json:"-"`
	Description string `json:"description" gorm:"size:500"`
	URL         string `json:"url" gorm:"size:1000"`
	ImageURL    string `json:"image_url"`
}

type Note struct {
	gorm.Model `json:"-"`
	EntityID   uint   `json:"-"`
	Content    string `json:"content" gorm:"size:20000"`
}

type Tag struct {
	gorm.Model `json:"-"`
	OwnerID    uint   `owner_id`
	Hash       string `json:"hash"`
	Name       string `json:"name"`
}

type Folder struct {
	gorm.Model    `json:"-"`
	EntityID      uint `json:"-"`
	ChildrenCount uint `json:"children_count"`
	HasParent     bool `json:"has_parent"`
}

type EntityRelation struct {
	ID       uint   `json:"id" gorm:"UNIQUE;AUTO_INCREMENT"`
	ParentID uint   `json:"parent_id"`
	ChildID  uint   `json:"child_id"`
	Type     string `type`
}
