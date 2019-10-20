package models

import (
    "fmt"
    "strings"
    "encoding/hex"
    "crypto/sha256"
)

type User struct {
    ID          uint     `json:"-" gorm:"primary_key"`
    Username    string   `json:"username,omitempty"`
    Password    string   `json:"-"`
    Email       string   `json:"enail"`
    Firstname   string   `json:"firstname"`
    Lastname    string   `json:"lastname"`
    // Bookmarks []Bookmark `json:"bookmarks,omitempty"`
}

var (
    reserved = []string{"login", "register", "s", "shared", "recent", "bookmarks", "b", "api", "refresh", "logout", ".", "..", "..."}
)

/////////////////////////////////////////////
/////////////// CREATE USER /////////////////
/////////////////////////////////////////////

func CreateUser(new User) (string) {
    if usernameInvalid(new) {
        return "USERNAME_INVALID"
    }

    if emailInvalid(new) {
        return "EMAIL_INVALID"
    }

    new.Email    = sanitizeEmail(new.Email)
    new.Password = hashPassword(new.Username, new.Password)
    new.Username = sanitizeUsername(new.Username)

    db.Create(&new)
    return ""
}

func usernameInvalid(new User) (bool) {
    username := strings.TrimSpace(strings.ToLower(new.Username))

    // Check if username shorter than 3 characters
    if len(username) < 3 {
        return true
    }

    // Check username against some reserved names
    for _, name := range reserved {
        if name == username {
            return true
        }
    }

    // Check if username is already taken
    if findUsersWithUsername(new.Username) > 0 {
        return true
    }

    return false
}

func emailInvalid(new User) (bool) {
    email := strings.TrimSpace(strings.ToLower(new.Email))

    // Check for minimum length (E.g. a@b.de)
    if len(email) < 7 {
        return true
    }

    // Check if email is already taken
    if findUsersWithEmail(email) > 0 {
        return true
    }

    // Add check for valid email adress

    return false
}

func sanitizeEmail(email string) (string) {
    return strings.TrimSpace(strings.ToLower(email))
}

func sanitizeUsername(username string) (string) {
    return strings.Replace(strings.TrimSpace(username), " ", "_", -1)
}

func hashPassword(un, pw string) (string) {
    user     := sha256.New()
    password := sha256.New()
    salted   := sha256.New()

    user.Write([]byte(un))
    password.Write([]byte(pw))

    salted.Write([]byte(fmt.Sprintf("%s%s", user, password)))

    return hex.EncodeToString(salted.Sum(nil))
}
