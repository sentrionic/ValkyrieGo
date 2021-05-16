package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/handler/middleware"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"time"
)

// Handler struct holds required services for handler to function
type Handler struct {
	userService  model.UserService
	MaxBodyBytes int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	// Create a handler (which will later have injected services)
	h := &Handler{
		userService:  c.UserService,
		MaxBodyBytes: c.MaxBodyBytes,
	}

	// Create an account group
	g := c.R.Group("api/account")

	if gin.Mode() != gin.TestMode {
		g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
		g.GET("/", middleware.AuthUser(h.userService), h.Me)
		g.PUT("/", middleware.AuthUser(h.userService), h.Edit)
		g.PUT("/change-password", middleware.AuthUser(h.userService), h.ChangePassword)
	} else {
		g.GET("/", h.Me)
	}

	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
	g.POST("/forgot-password", h.ForgotPassword)
	g.POST("/reset-password", h.ResetPassword)
}

func setUserSession(c *gin.Context, id string) {
	if gin.Mode() != gin.TestMode {
		session := sessions.Default(c)
		session.Set("userId", id)
		if err := session.Save(); err != nil {
			fmt.Println(err)
		}
	}
}
