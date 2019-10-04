package database

import (
    "fmt"

    "github.com/Techassi/gomark/internal/constants"
    "github.com/Techassi/gomark/internal/models"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

func Open() (*gorm.DB, error) {
    conn := fmt.Sprintf("%s:%s@/%s?%s", "root", "", "gomark", constants.STORE_PARAMS)
    db, err := gorm.Open("mysql", conn)
    if err != nil {
		return db, err
	}

    db = db.Set("gorm:table_options", constants.STORE_TABLE_OPTIONS)

    db.AutoMigrate(
		&models.Tag{},
		&models.Account{},
		&models.Bookmark{},
	)

    return db, nil
}
