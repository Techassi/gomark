package models

import (
    "fmt"

    "github.com/Techassi/gomark/internal/constants"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

func OpenMySQL() {
    conn := fmt.Sprintf("%s:%s@/%s?%s", "root", "", "gomark", constants.STORE_PARAMS)
    db, err = gorm.Open("mysql", conn)
    if err != nil {
        panic(err)
    }

    db = db.Set("gorm:table_options", constants.STORE_TABLE_OPTIONS)

    db.AutoMigrate(
		&User{},
		&Entity{},
        &NoteData{},
		&FolderData{},
		&BookmarkData{},
	)
}

func FindUsersWithUsername(username string) (int) {
    users := []User{}
    db.Where("username = ?", username).Find(&users)

    return len(users)
}

func FindUsersWithEmail(email string) (int) {
    users := []User{}
    db.Where("email = ?", email).Find(&users)

    return len(users)
}
