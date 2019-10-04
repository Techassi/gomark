package models

import (
    "time"
)

type Tag struct {
    ID            uint      `json:"-" gorm:"primary_key"`
    CreatedAt     time.Time `json:"-"`
    UpdatedAt     time.Time `json:"-"`
    DeletedAt    *time.Time `json:"-"`
    Hash          string    `json:"hash"`
    Name          string    `json:"name"`
    BookmarkID    uint      `json:"bookmark_id"`
}
