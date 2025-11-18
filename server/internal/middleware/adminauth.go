package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func AdminAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        if v, err := c.Cookie("admin_session"); err == nil && v == "ok" {
            c.Next()
            return
        }
        c.Redirect(http.StatusFound, "/admin/login")
        c.Abort()
    }
}