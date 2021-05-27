package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/handler/ws"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func (h *Handler) GetMessages(c *gin.Context) {
	channelId := c.Param("channelId")
	userId := c.MustGet("userId").(string)

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	err = h.channelService.IsChannelMember(channel, userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}

	cursor := c.Query("cursor")

	messages, err := h.messageService.GetMessages(userId, channel, cursor)

	if err != nil {
		e := apperrors.NewNotFound("messages", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if len(*messages) == 0 {
		var empty = make([]model.MessageResponse, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, messages)
}

type messageRequest struct {
	Text *string               `form:"text" binding:"omitempty,lte=2000"`
	File *multipart.FileHeader `form:"file" binding:"omitempty"`
}

func (h *Handler) CreateMessage(c *gin.Context) {
	channelId := c.Param("channelId")
	userId := c.MustGet("userId").(string)
	channel, err := h.channelService.Get(channelId)

	err = h.channelService.IsChannelMember(channel, userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}

	var req messageRequest
	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	if req.Text == nil && req.File == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Either a message pr a file is required",
		})
		return
	}

	author, _ := h.userService.Get(userId)

	message := model.Message{
		UserId:    userId,
		ChannelId: channel.ID,
	}

	if req.Text != nil {
		message.Text = req.Text
	}

	if req.File != nil {
		mimeType := req.File.Header.Get("Content-Type")

		if valid := isAllowedFileType(mimeType); !valid {
			log.Println("File is not an allowable mime-type")
			e := apperrors.NewBadRequest("imageFile must be 'image' or 'audio'")
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		attachment, err := h.messageService.UploadFile(req.File, channel.ID)

		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		message.Attachment = attachment
	}

	if err := h.messageService.CreateMessage(&message); err != nil {
		log.Printf("Failed to create message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	response := model.MessageResponse{
		Id:         message.ID,
		Text:       message.Text,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		Attachment: message.Attachment,
		User: model.MemberResponse{
			Id:        author.ID,
			Username:  author.Username,
			Image:     author.Image,
			IsOnline:  author.IsOnline,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
			IsFriend:  false,
		},
	}

	if !channel.IsDM {
		settings, _ := h.guildService.GetMemberSettings(userId, *channel.GuildID)
		response.User.Nickname = settings.Nickname
		response.User.Color = settings.Color
	}

	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.NewMessageAction,
		Data:   response,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	h.WsServer.broadcastToRoom(data, channelId)

	if channel.IsDM {
		// Open the DM and push it to the top
		_ = h.channelService.OpenDMForAll(channelId)
		//TODO: Emit new_dm_notification event
	} else {
		// Update last activity in channel
		channel.LastActivity = time.Now()
		_ = h.channelService.UpdateChannel(channel)
		//TODO: Emit new_notification event
	}

	c.JSON(http.StatusOK, true)
}

func (h *Handler) EditMessage(c *gin.Context) {
	messageId := c.Param("messageId")
	userId := c.MustGet("userId").(string)
	message, err := h.messageService.Get(messageId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if message.UserId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Only the author can edit the message",
		})
		return
	}

	var req messageRequest
	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	message.Text = req.Text

	if err := h.messageService.UpdateMessage(message); err != nil {
		log.Printf("Failed to edit message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	//TODO: Emit edit_message event

	c.JSON(http.StatusOK, true)
}

func (h *Handler) DeleteMessage(c *gin.Context) {
	messageId := c.Param("messageId")
	userId := c.MustGet("userId").(string)
	message, err := h.messageService.Get(messageId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channel, err := h.channelService.Get(message.ChannelId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if !channel.IsDM {
		guild, err := h.guildService.GetGuild(*channel.GuildID)

		if err != nil {
			e := apperrors.NewNotFound("message", messageId)

			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		if message.UserId != userId || guild.OwnerId != userId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Only the author can delete the message",
			})
			return
		}
	} else {
		if message.UserId != userId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Only the author can delete the message",
			})
			return
		}
	}

	if err := h.messageService.DeleteMessage(message); err != nil {
		log.Printf("Failed to delete message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	//TODO: Emit delete_message event

	c.JSON(http.StatusOK, true)
}
