package models

import (
    "time"
)

type Bookmark struct {
    ID            uint      `json:"-" gorm:"primary_key"`
    CreatedAt     time.Time `json:"-"`
    UpdatedAt     time.Time `json:"-"`
    DeletedAt    *time.Time `json:"-"`
    Hash          string    `json:"hash"`
    Name          string    `json:"name"`
    Description   string    `json:"description" gorm:"size:500"`
    URL           string    `json:"url" gorm:"size:1000"`
    ImageURL      string    `json:"image_url"`
    Tags        []Tag       `json:"tags,omitempty"`
}
