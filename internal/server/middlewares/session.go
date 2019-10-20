package middlewares

import (
    "net/http"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

func ValidateSession(c *gin.Context) {
    session := sessions.Default(c)
	user := session.Get("test")
	if user == nil {
		// Abort the request with the appropriate error code
		c.Redirect(http.StatusMovedPermanently, "/login")
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
