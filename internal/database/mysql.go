package database

import (
    "fmt"

    "github.com/Techassi/gomark/internal/constants"
    "github.com/Techassi/gomark/internal/models"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type MySQLDatabase struct {
	gorm.DB
}

func OpenMySQL() (mysqlDB *MySQLDatabase) {
    conn := fmt.Sprintf("%s:%s@/%s?%s", "root", "", "gomark", constants.STORE_PARAMS)
    db, err := gorm.Open("mysql", conn)
    if err != nil {
        panic(err)
    }

    db = db.Set("gorm:table_options", constants.STORE_TABLE_OPTIONS)

    db.AutoMigrate(
		&models.Tag{},
		&models.User{},
		&models.Bookmark{},
	)

    mysqlDB = &MySQLDatabase{*db}
    return mysqlDB
}

func (db *MySQLDatabase) FindUsersWithUsername(username string) (int) {
    users := []models.User{}
    db.Where("username = ?", username).Find(&users)

    return len(users)
}

func (db *MySQLDatabase) FindUsersWithEmail(email string) (int) {
    users := []models.User{}
    db.Where("email = ?", email).Find(&users)

    return len(users)
}

func (db *MySQLDatabase) CreateUser(new models.User) {
    db.Create(&new)
}
