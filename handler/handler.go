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
	guildService  model.GuildService
	MaxBodyBytes  int64
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	FriendService   model.FriendService
	GuildService    model.GuildService
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
		guildService:  c.GuildService,
		MaxBodyBytes:  c.MaxBodyBytes,
	}

	// Create an account group
	ag := c.R.Group("api/account")
	ag.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)
	ag.POST("/logout", h.Logout)
	ag.POST("/forgot-password", h.ForgotPassword)
	ag.POST("/reset-password", h.ResetPassword)

	ag.Use(middleware.AuthUser())

	ag.GET("/", h.Me)
	ag.PUT("/", h.Edit)
	ag.PUT("/change-password", h.ChangePassword)

	ag.GET("/me/friends", h.GetUserFriends)
	ag.GET("/me/pending", h.GetUserRequests)
	ag.POST("/:memberId/friend", h.SendFriendRequest)
	ag.DELETE("/:memberId/friend", h.RemoveFriend)
	ag.POST("/:memberId/friend/accept", h.AcceptFriendRequest)
	ag.POST("/:memberId/friend/cancel", h.CancelFriendRequest)

	gg := c.R.Group("api/guilds")
	gg.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))
	gg.Use(middleware.AuthUser())

	gg.GET("/:guildId/members", h.GetGuildMembers)
	gg.GET("/", h.GetUserGuilds)
	gg.POST("/create", h.CreateGuild)
	gg.GET("/:guildId/invite", h.GetInvite)
	gg.DELETE("/:guildId/invite", h.DeleteGuildInvites)
	gg.POST("/join", h.JoinGuild)
	gg.GET("/:guildId/member", h.GetMemberSettings)
	gg.PUT("/:guildId/member", h.EditMemberSettings)
	gg.DELETE("/:guildId", h.LeaveGuild)
	gg.PUT("/:guildId", h.EditGuild)
	gg.DELETE("/:guildId/delete", h.DeleteGuild)
	gg.GET("/:guildId/bans", h.GetBanList)
	gg.POST("/:guildId/bans", h.BanMember)
	gg.DELETE("/:guildId/bans", h.UnbanMember)
	gg.POST("/:guildId/kick", h.KickMember)
}
