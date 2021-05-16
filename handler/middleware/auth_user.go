package middleware

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
)

func AuthUser(s model.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("userId")

		if id == nil {
			err := errors.New("provided session is invalid")
			c.JSON(401, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		userId := id.(string)

		user, err := s.Get(userId)

		if err != nil {
			c.JSON(401, gin.H{
				"error": err,
			})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}