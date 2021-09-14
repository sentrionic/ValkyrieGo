package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateChannel(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid, _ := service.GenerateId()
	user := &model.User{}
	user.ID = uid

	guildId, _ := service.GenerateId()
	guild := &model.Guild{OwnerId: uid}
	guild.ID = guildId

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("userId", uid)
	})
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("vlk", store))

	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userId", uid)
	})

	mockUserService := new(mocks.UserService)
	mockUserService.On("Get", uid).Return(user, nil)

	mockGuildService := new(mocks.GuildService)
	mockGuildService.On("GetGuild", guildId).Return(guild, nil)

	mockChannelService := new(mocks.ChannelService)
	mockSocketService := new(mocks.SocketService)

	NewHandler(&Config{
		R:              router,
		UserService:    mockUserService,
		GuildService:   mockGuildService,
		ChannelService: mockChannelService,
		SocketService:  mockSocketService,
	})

	t.Run("Name required", func(t *testing.T) {
		mockChannelService.On("CreateChannel", mock.AnythingOfType("*model.Channel")).Return(nil)
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Name too short", func(t *testing.T) {
		mockChannelService.On("CreateChannel", mock.AnythingOfType("*model.Channel")).Return(nil)
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": "Ch",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Name too long", func(t *testing.T) {
		mockChannelService.On("CreateChannel", mock.AnythingOfType("*model.Channel")).Return(nil)
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		name := ""
		for i := 0; i < 31; i++ {
			name += "a"
		}
		reqBody, err := json.Marshal(gin.H{
			"name": name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Successful channel creation", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		name := "main"
		reqBody, err := json.Marshal(gin.H{
			"name": name,
		})
		assert.NoError(t, err)

		id, _ := service.GenerateId()
		channel := &model.Channel{
			GuildID:  &guildId,
			Name:     name,
			IsPublic: true,
		}

		mockChannelService.On("CreateChannel", channel).
			Run(func(args mock.Arguments) {
				channelArgs := args.Get(0).(*model.Channel)
				channelArgs.ID = id
			}).Return(nil)
		guild.Channels = append(guild.Channels, *channel)
		mockGuildService.On("UpdateGuild", guild).Return(nil)
		mockSocketService.On("EmitNewChannel", guildId, mock.AnythingOfType("*model.ChannelResponse")).Return(nil)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockChannelService.AssertCalled(t, "CreateChannel", channel)
		mockGuildService.AssertCalled(t, "UpdateGuild", guild)
		mockSocketService.AssertCalled(t, "EmitNewChannel", guildId, mock.AnythingOfType("*model.ChannelResponse"))
		mockChannelService.AssertExpectations(t)
	})

	t.Run("Error Returned from ChannelService.CreateChannel", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		name := "main"
		reqBody, err := json.Marshal(gin.H{
			"name": name,
		})
		assert.NoError(t, err)

		router := gin.Default()
		router.Use(func(c *gin.Context) {
			c.Set("userId", uid)
		})
		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("vlk", store))

		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("userId", uid)
		})

		mockChannelService := new(mocks.ChannelService)
		NewHandler(&Config{
			R:              router,
			UserService:    mockUserService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		mockError := apperrors.NewInternal()
		mockChannelService.On("CreateChannel", mock.AnythingOfType("*model.Channel")).Return(mockError)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "CreateChannel", mock.AnythingOfType("*model.Channel"))
		mockChannelService.AssertExpectations(t)
	})

	t.Run("Successful private channel creation", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		name := "secret"
		reqBody, err := json.Marshal(gin.H{
			"name":     name,
			"isPublic": false,
		})
		assert.NoError(t, err)

		channel := &model.Channel{
			GuildID:  &guildId,
			Name:     name,
			IsPublic: false,
		}

		mockChannelService.On("CreateChannel", channel).Return(nil)
		guild.Channels = append(guild.Channels, *channel)
		mockGuildService.On("UpdateGuild", guild).Return(nil)
		mockGuildService.On("FindUsersByIds", []string{uid}, guildId).Return(&[]model.User{*user}, nil)
		channel.PCMembers = append(channel.PCMembers, *user)
		mockSocketService.On("EmitNewChannel", guildId, mock.AnythingOfType("*model.ChannelResponse")).Return(nil)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusCreated, rr.Code)
		mockChannelService.AssertCalled(t, "CreateChannel", channel)
		mockGuildService.AssertCalled(t, "UpdateGuild", guild)
		mockSocketService.AssertCalled(t, "EmitNewChannel", guildId, mock.AnythingOfType("*model.ChannelResponse"))
		mockChannelService.AssertExpectations(t)
	})

	t.Run("Guild already has 50 channels", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()
		guild.Channels = []model.Channel{}

		// Create 50 channels
		for i := 0; i < 50; i++ {
			c := model.Channel{GuildID: &guildId}
			c.ID = string(rune(i))
			guild.Channels = append(guild.Channels, c)
		}
		assert.Len(t, guild.Channels, 50)

		name := "test"
		reqBody, err := json.Marshal(gin.H{
			"name": name,
		})
		assert.NoError(t, err)

		channel := &model.Channel{
			GuildID:  &guildId,
			Name:     name,
			IsPublic: true,
		}

		mockError := apperrors.NewBadRequest("channel limit is 50")
		mockChannelService.On("CreateChannel", channel).Return(nil)
		mockGuildService.On("UpdateGuild", guild).Return(nil)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", guildId)
		mockChannelService.AssertNotCalled(t, "CreateChannel", channel)
	})

	t.Run("Member that is not the owner", func(t *testing.T) {
		uid, _ := service.GenerateId()
		user := &model.User{}
		user.ID = uid

		router := gin.Default()
		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("vlk", store))

		router.Use(func(c *gin.Context) {
			session := sessions.Default(c)
			session.Set("userId", uid)
		})

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(user, nil)

		NewHandler(&Config{
			R:              router,
			UserService:    mockUserService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": "uashdui",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/channels/"+guildId, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 401, rr.Code)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})
}
