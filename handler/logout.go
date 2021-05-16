package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Logout handler
func (h *Handler) Logout(c *gin.Context) {
	c.Set("user", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()

	if err != nil {
		fmt.Printf("error clearing session: %v", err)
	}

	c.JSON(http.StatusOK, true)
	return
}