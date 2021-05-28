package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sentrionic/valkyrie/handler/middleware"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"time"
)

// Handler struct holds required services for handler to function
type Handler struct {
	userService    model.UserService
	friendService  model.FriendService
	guildService   model.GuildService
	channelService model.ChannelService
	messageService model.MessageService
	MaxBodyBytes   int64
	WsServer       *WsServer
}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
	R               *gin.Engine
	UserService     model.UserService
	FriendService   model.FriendService
	GuildService    model.GuildService
	ChannelService  model.ChannelService
	MessageService  model.MessageService
	TimeoutDuration time.Duration
	MaxBodyBytes    int64
	Redis           *redis.Client
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {

	// Websocket Setup
	wsServer := NewWebsocketServer(c)
	go wsServer.Run()

	c.R.GET("/ws", middleware.AuthUser(), func(c *gin.Context) {
		ServeWs(wsServer, c)
	})

	// Create a handler (which will later have injected services)
	h := &Handler{
		userService:    c.UserService,
		friendService:  c.FriendService,
		guildService:   c.GuildService,
		channelService: c.ChannelService,
		messageService: c.MessageService,
		MaxBodyBytes:   c.MaxBodyBytes,
		WsServer:       wsServer,
	}

	c.R.Use(middleware.Timeout(c.TimeoutDuration, apperrors.NewServiceUnavailable()))

	// Create an account group
	ag := c.R.Group("api/account")

	ag.POST("/register", h.Register)
	ag.POST("/login", h.Login)
	ag.POST("/logout", h.Logout)
	ag.POST("/forgot-password", h.ForgotPassword)
	ag.POST("/reset-password", h.ResetPassword)

	ag.Use(middleware.AuthUser())

	ag.GET("", h.Me)
	ag.PUT("", h.Edit)
	ag.PUT("/change-password", h.ChangePassword)

	ag.GET("/me/friends", h.GetUserFriends)
	ag.GET("/me/pending", h.GetUserRequests)
	ag.POST("/:memberId/friend", h.SendFriendRequest)
	ag.DELETE("/:memberId/friend", h.RemoveFriend)
	ag.POST("/:memberId/friend/accept", h.AcceptFriendRequest)
	ag.POST("/:memberId/friend/cancel", h.CancelFriendRequest)

	// Create a guild group
	gg := c.R.Group("api/guilds")
	gg.Use(middleware.AuthUser())

	gg.GET("/:guildId/members", h.GetGuildMembers)
	gg.GET("", h.GetUserGuilds)
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

	// Create a channels group
	cg := c.R.Group("api/channels")
	cg.Use(middleware.AuthUser())

	//Route parameters cause conflicts so they have to use the same parameter name
	cg.GET("/:id", h.GuildChannels)                 // id -> guildId
	cg.POST("/:id", h.CreateChannel)                // id -> guildId
	cg.GET("/:id/members", h.PrivateChannelMembers) // id -> channelId
	cg.POST("/:id/dm", h.GetOrCreateDM)             // id -> memberId
	cg.GET("/me/dm", h.DirectMessages)              //
	cg.PUT("/:id", h.EditChannel)                   // id -> channelId
	cg.DELETE("/:id", h.DeleteChannel)              // id -> channelId
	cg.DELETE("/:id/dm", h.CloseDM)                 // id -> channelId

	// Create a messages group
	mg := c.R.Group("api/messages")
	mg.Use(middleware.AuthUser())

	mg.GET("/:channelId", h.GetMessages)
	mg.POST("/:channelId", h.CreateMessage)
	mg.PUT("/:messageId", h.EditMessage)
	mg.DELETE("/:messageId", h.DeleteMessage)
}
