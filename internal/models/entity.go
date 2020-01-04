package models

import (
	"github.com/jinzhu/gorm"
)

type Entity struct {
	gorm.Model   `json:"-"`
	EntityID     uint          `json:"-"`
	Type         string        `json:"type"`
	Hash         string        `json:"hash"`
	Name         string        `json:"name"`
	Shared       bool          `json:"shared"`
	ClickedOn    uint          `json:"clicked_on"`
	FolderData   *FolderData   `json:"folder_data,omitempty"`
	NoteData     *NoteData     `json:"note_data,omitempty"`
	BookmarkData *BookmarkData `json:"bookmark_data,omitempty"`
}

type FolderData struct {
	gorm.Model         `json:"-"`
	EntityID           uint     `json:"-"`
	ChildEntitiesCount uint     `json:"child_entities_count"`
	ChildEntities      []Entity `json:"child_entities"`
}

type NoteData struct {
	gorm.Model `json:"-"`
	EntityID   uint   `json:"entity_id"`
	Content    string `json:"content"`
}

type BookmarkData struct {
	gorm.Model  `json:"-"`
	EntityID    uint   `json:"-"`
	Description string `json:"description" gorm:"size:500"`
	URL         string `json:"url" gorm:"size:1000"`
	ImageURL    string `json:"image_url"`
}

type Tag struct {
	gorm.Model `json:"-"`
	Hash       string `json:"hash"`
	Name       string `json:"name"`
	BookmarkID uint   `json:"bookmark_id"`
}
