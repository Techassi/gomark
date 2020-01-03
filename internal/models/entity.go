package models

import (
	"time"
)

type Entity struct {
	ID           uint         `json:"-" gorm:"primary_key"`
	CreatedAt    time.Time    `json:"-"`
	UpdatedAt    time.Time    `json:"-"`
	DeletedAt    *time.Time   `json:"-"`
	Type         string       `json:"type"`
	Hash         string       `json:"hash"`
	Name         string       `json:"name"`
	Shared       bool         `json:"shared"`
	FolderData   FolderData   `json:"folder_data,omitempty"`
	NoteData     NoteData     `json:"note_data,omitempty"`
	BookmarkData BookmarkData `json:"bookmark_data,omitempty"`
}

type FolderData struct {
	EntityID           uint       `json:"entity_id"`
	ID                 uint       `json:"-" gorm:"primary_key"`
	CreatedAt          time.Time  `json:"-"`
	UpdatedAt          time.Time  `json:"-"`
	DeletedAt          *time.Time `json:"-"`
	ChildEntitiesCount uint       `json:"child_entities_count"`
	ChildEntities      []Entity   `json:"child_entities"`
}

type NoteData struct {
	EntityID  uint       `json:"entity_id"`
	ID        uint       `json:"-" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Content   string     `json:"content"`
}

type BookmarkData struct {
	EntityID    uint       `json:"entity_id"`
	ID          uint       `json:"-" gorm:"primary_key"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
	Description string     `json:"description" gorm:"size:500"`
	URL         string     `json:"url" gorm:"size:1000"`
	ImageURL    string     `json:"image_url"`
}

type Tag struct {
	ID         uint       `json:"-" gorm:"primary_key"`
	CreatedAt  time.Time  `json:"-"`
	UpdatedAt  time.Time  `json:"-"`
	DeletedAt  *time.Time `json:"-"`
	Hash       string     `json:"hash"`
	Name       string     `json:"name"`
	BookmarkID uint       `json:"bookmark_id"`
}
