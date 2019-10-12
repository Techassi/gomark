package status

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func UsernameInvalid(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Username is invalid or taken", "action": "username_invalid"})
}

func EMailInvalid(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "E-Mail is invalid or taken", "action": "email_invalid"})
}

func AccountCreated(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "New account was created", "action": "account_created"})
}

func AccountNotCreated(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "E-Mail is invalid or taken", "action": "account_not_created"})
}
