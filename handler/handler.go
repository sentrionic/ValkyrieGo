package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/handler/middleware"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"time"
)

// Handler struct holds required services for handler to function
type Handler struct {
	userService   model.UserService
	friendService model.FriendService
	MaxBodyBytes  int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	FriendService   model.FriendService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
	// Create a handler (which will later have injected services)
	h := &Handler{
		userService:   c.UserService,
		friendService: c.FriendService,
		MaxBodyBytes:  c.MaxBodyBytes,
	}

	// Create an account group
	g := c.R.Group("api/account")
	g.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/logout", h.Logout)
	g.POST("/forgot-password", h.ForgotPassword)
	g.POST("/reset-password", h.ResetPassword)

	g.Use(middleware.AuthUser())

	g.GET("/", h.Me)
	g.PUT("/", h.Edit)
	g.PUT("/change-password", h.ChangePassword)

	g.GET("/me/friends", h.GetUserFriends)
	g.GET("/me/pending", h.GetUserRequests)
	g.POST("/:memberId/friend", h.SendFriendRequest)
	g.DELETE("/:memberId/friend", h.RemoveFriend)
	g.POST("/:memberId/friend/accept", h.AcceptFriendRequest)
	g.POST("/:memberId/friend/cancel", h.CancelFriendRequest)
}