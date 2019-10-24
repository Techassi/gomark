package models

import (

)

func GetEntities() ([]Entity) {
    var entities []Entity

    db.Set("gorm:auto_preload", true).Find(&entities)

    return entities
}

func GetNotes() ([]Entity) {
    var entities []Entity

    db.Set("gorm:auto_preload", true).Where("type = ?", "note").Find(&entities)

    return entities
}

func GetSharedEntities() ([]Entity) {
    var entities []Entity

    db.Set("gorm:auto_preload", true).Where("shared = ?", true).Find(&entities)

    return entities
}
