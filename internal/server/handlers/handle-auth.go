package handlers

import (
    "strings"

    "github.com/Techassi/gomark/internal/models"
    "github.com/Techassi/gomark/internal/server/status"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

func AUTH_Login(c *gin.Context) {
    session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		status.InvalidParameters(c)
		return
	}

	// Check for username and password match
	if !AUTH_CheckCredentials(username, password) {
		status.InvalidCredentials(c)
		return
	}

	// Save the username in the session
	session.Set("test", username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		status.FailedToSaveSession(c)
		return
	}

	status.LoggedIn(c)
}

func AUTH_Logout(c *gin.Context) {
    return
}

func AUTH_Register(c *gin.Context) {
    new := models.User{
        Username:  c.PostForm("username"),
        Password:  c.PostForm("password"),
        Email:     c.PostForm("email"),
        Firstname: c.PostForm("firstname"),
        Lastname:  c.PostForm("lastname"),
    }

    if err := models.CreateUser(new); err != "" {
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

func AUTH_CheckCredentials(username, password string) (bool) {
    return models.CheckIfValidCredentials(username, password)
}
