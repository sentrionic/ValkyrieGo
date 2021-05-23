package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/handler"
	"github.com/sentrionic/valkyrie/repository"
	"github.com/sentrionic/valkyrie/service"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func inject(d *dataSources) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	/*
	 * repository layer
	 */
	userRepository := repository.NewUserRepository(d.DB)
	friendRepository := repository.NewFriendRepository(d.DB)
	guildRepository := repository.NewGuildRepository(d.DB)
	channelRepository := repository.NewChannelRepository(d.DB)
	messageRepository := repository.NewMessageRepository(d.DB)

	bucketName := os.Getenv("AWS_STORAGE_BUCKET_NAME")
	fileRepository := repository.NewFileRepository(d.S3Session, bucketName)
	redisRepository := repository.NewRedisRepository(d.RedisClient)

	gmailUser := os.Getenv("GMAIL_USER")
	gmailPassword := os.Getenv("GMAIL_PASSWORD")
	origin := os.Getenv("CORS_ORIGIN")
	mailRepository := repository.NewMailRepository(gmailUser, gmailPassword, origin)

	/*
	 * service layer
	 */
	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		FileRepository:  fileRepository,
		RedisRepository: redisRepository,
		MailRepository:  mailRepository,
	})

	friendService := service.NewFriendService(&service.FSConfig{
		UserRepository:   userRepository,
		FriendRepository: friendRepository,
	})

	guildService := service.NewGuildService(&service.GSConfig{
		UserRepository:    userRepository,
		FileRepository:    fileRepository,
		RedisRepository:   redisRepository,
		GuildRepository:   guildRepository,
		ChannelRepository: channelRepository,
	})

	channelService := service.NewChannelService(&service.CSConfig{
		ChannelRepository: channelRepository,
		GuildRepository:   guildRepository,
	})

	messageService := service.NewMessageService(&service.MSConfig{
		MessageRepository: messageRepository,
		FileRepository:    fileRepository,
	})

	// initialize gin.Engine
	router := gin.Default()

	// set cors settings
	router.Use(CORSMiddleware(origin))

	redisURL := os.Getenv("REDIS_URL")
	secret := os.Getenv("SECRET")
	store, _ := redis.NewStore(10, "tcp", redisURL, "", []byte(secret))

	store.Options(sessions.Options{
		Domain:   "",
		MaxAge:   60 * 60 * 24 * 7, // 7 days
		Secure:   false,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions("vlk", store))

	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	maxBodyBytes := os.Getenv("MAX_BODY_BYTES")
	mbb, err := strconv.ParseInt(maxBodyBytes, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse MAX_BODY_BYTES as int: %w", err)
	}

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		FriendService:   friendService,
		GuildService:    guildService,
		ChannelService:  channelService,
		MessageService:  messageService,
		TimeoutDuration: time.Duration(ht) * time.Second,
		MaxBodyBytes:    mbb,
	})

	return router, nil
}

func CORSMiddleware(origin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
