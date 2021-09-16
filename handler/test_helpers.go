package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func getAuthenticatedTestRouter(uid string) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("vlk", store))

	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userId", uid)
		c.Set("userId", uid)
	})

	return router
}

func getTestRouter() *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("vlk", store))
	return router
}
