package status

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

/////////////////////////////////////////////
////////////////// SESSION //////////////////
/////////////////////////////////////////////

func InvalidParameters(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusBadRequest, "message": "Parameters are invalid or empty", "action": "parameters_invalid"})
}

func InvalidCredentials(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusUnauthorized, "message": "Inavlid credentials", "action": "credentials_invalid"})
}

func FailedToSaveSession(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": http.StatusInternalServerError, "message": "Faied to save session", "action": "failed_to_save"})
}