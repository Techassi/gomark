package database

import (
    "github.com/Techassi/gomark/internal/models"
)

type DB interface {
    FindUsersWithUsername(username string) (int)
    FindUsersWithEmail(email string) (int)
    CreateUser(new models.User)
}
