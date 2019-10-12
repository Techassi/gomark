package middlewares

import (
    "github.com/Techassi/gomark/internal/database"

    "github.com/gin-gonic/gin"
)

func InjectDB(db database.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("DB", db)
        c.Next()
    }
}
