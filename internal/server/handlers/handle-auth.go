package handlers

import (
    "fmt"
    "strings"
    "encoding/hex"
    "crypto/sha256"

    "github.com/Techassi/gomark/internal/models"
    "github.com/Techassi/gomark/internal/database"
    "github.com/Techassi/gomark/internal/server/status"

    "github.com/gin-gonic/gin"
)

var (
    reserved = []string{"login", "register", "s", "shared", "recent", "bookmarks", "b", "api", "refresh", "logout", ".", "..", "..."}
)

func AUTH_Login(c *gin.Context) {
    return
}

func AUTH_Logout(c *gin.Context) {
    return
}

func AUTH_Register(c *gin.Context) {
    db, ok := c.MustGet("DB").(database.DB)
    if !ok { return }

    new := models.User{
        Username:  c.PostForm("username"),
        Password:  c.PostForm("password"),
        Email:     c.PostForm("email"),
        Firstname: c.PostForm("firstname"),
        Lastname:  c.PostForm("lastname"),
    }

    if err := createUser(db, new); err != "" {
        switch err {
            case "USERNAME_INVALID":
                status.UsernameInvalid(c)
                return
            case "EMAIL_INVALID":
                status.EMailInvalid(c)
                return
        }

        status.AccountNotCreated(c)
        return
    }

    status.AccountCreated(c)
}

/////////////////////////////////////////////
/////////////// CREATE USER /////////////////
/////////////////////////////////////////////

func createUser(db database.DB, new models.User) (string) {
    if usernameInvalid(db, new) {
        return "USERNAME_INVALID"
    }

    if emailInvalid(db, new) {
        return "EMAIL_INVALID"
    }

    new.Email    = sanitizeEmail(new.Email)
    new.Password = hashPassword(new)
    new.Username = sanitizeUsername(new.Username)

    db.CreateUser(new)
    return ""
}

func usernameInvalid(db database.DB, new models.User) (bool) {
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
    if db.FindUsersWithUsername(new.Username) > 0 {
        return true
    }

    return false
}

func emailInvalid(db database.DB, new models.User) (bool) {
    email := strings.TrimSpace(strings.ToLower(new.Email))

    // Check for minimum length (E.g. a@b.de)
    if len(email) < 7 {
        return true
    }

    // Check if email is already taken
    if db.FindUsersWithEmail(email) > 0 {
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

func hashPassword(new models.User) (string) {
    user     := sha256.New()
    password := sha256.New()
    salted   := sha256.New()

    user.Write([]byte(new.Username))
    password.Write([]byte(new.Password))

    salted.Write([]byte(fmt.Sprintf("%s%s", user, password)))

    return hex.EncodeToString(salted.Sum(nil))
}
