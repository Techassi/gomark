package models

import (
	"github.com/jinzhu/gorm"
)

type Settings struct {
	gorm.Model
	RegisterEnabled bool `json:"register_enabled"`
	SharingEnabled  bool `json:"sharing_enabled"`
}
