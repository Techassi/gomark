package handlers

import (
    "github.com/Techassi/gomark/internal/models"
    "github.com/Techassi/gomark/internal/server/status"

    "github.com/gin-gonic/gin"
)

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
