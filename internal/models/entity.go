package models

import (
	"github.com/jinzhu/gorm"
)

type Bookmark struct {
	gorm.Model  `json:"-"`
	OwnerID     uint   `json:"owner_id"`
	Hash        string `json:"hash"`
	Name        string `json:"name"`
	Shared      bool   `json:"shared"`
	Pinned      bool   `json:"pinned"`
	ClickedOn   uint   `json:"clicked_on"`
	Description string `json:"description" gorm:"size:500"`
	URL         string `json:"url" gorm:"size:1000"`
	ImageURL    string `json:"image_url"`
	HasParent   bool   `json:"has_parent"`
}

type Note struct {
	gorm.Model `json:"-"`
	OwnerID    uint   `owner_id`
	Hash       string `json:"hash"`
	Name       string `json:"name"`
	ClickedOn  uint   `json:"clicked_on"`
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
	OwnerID       uint   `json:"owner_id"`
	Hash          string `json:"hash"`
	Name          string `json:"name"`
	Shared        bool   `json:"shared"`
	ClickedOn     uint   `json:"clicked_on"`
	ChildrenCount uint   `json:"children_count"`
	HasParent     bool   `json:"has_parent"`
}

type EntityRelation struct {
	ID       uint   `json:"id" gorm:"UNIQUE;AUTO_INCREMENT"`
	ParentID uint   `json:"parent_id"`
	ChildID  uint   `json:"child_id"`
	Type     string `type`
}
