package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

/*
 * MessageHandler contains all routes related to message actions (/api/messages)
 */

// GetMessages returns messages for the given channel
// It returns the most recent 35 or the ones after the given cursor
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

	// Check if the user has access to said channel
	err = h.channelService.IsChannelMember(channel, userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
		return
	}

	// Cursor is based on the created_at field of the message
	cursor := c.Query("cursor")

	messages, err := h.messageService.GetMessages(userId, channel, cursor)

	if err != nil {
		e := apperrors.NewNotFound("messages", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the channel does not have any messages, return an empty array
	if len(*messages) == 0 {
		var empty = make([]model.MessageResponse, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// messageRequest contains all field required to create a message.
// Either text or file must be provided
type messageRequest struct {
	Text *string               `form:"text" binding:"omitempty,lte=2000"`
	File *multipart.FileHeader `form:"file" binding:"omitempty"`
}

// CreateMessage creates a message in the given channel
func (h *Handler) CreateMessage(c *gin.Context) {
	channelId := c.Param("channelId")
	userId := c.MustGet("userId").(string)
	channel, err := h.channelService.Get(channelId)

	// Check if the user has access to said channel
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

	// Either text or file must be provided
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

	// Get member settings if it is not a DM
	if !channel.IsDM {
		settings, _ := h.guildService.GetMemberSettings(userId, *channel.GuildID)
		response.User.Nickname = settings.Nickname
		response.User.Color = settings.Color
	}

	// Emit new message to the channel
	h.socketService.EmitNewMessage(channelId, &response)

	if channel.IsDM {
		// Open the DM and push it to the top
		_ = h.channelService.OpenDMForAll(channelId)
		// Post a notification
		h.socketService.EmitNewDMNotification(channelId, author)
	} else {
		// Update last activity in channel
		channel.LastActivity = time.Now()
		_ = h.channelService.UpdateChannel(channel)
		// Post a notification
		h.socketService.EmitNewNotification(*channel.GuildID, channelId)
	}

	c.JSON(http.StatusOK, true)
}

// EditMessage edits the given message with the given text
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

	if message.Text == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request parameters. See errors",
			"errors":  "Text is required",
		})
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

	response := model.MessageResponse{
		Id:         message.ID,
		Text:       message.Text,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		Attachment: message.Attachment,
		User: model.MemberResponse{
			Id: userId,
		},
	}

	// Emit edited message to the channel
	h.socketService.EmitEditMessage(message.ChannelId, &response)

	c.JSON(http.StatusOK, true)
}

// DeleteMessage deletes the given message
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

	// Check if message author or guild owner
	if !channel.IsDM {
		guild, err := h.guildService.GetGuild(*channel.GuildID)

		if err != nil {
			e := apperrors.NewNotFound("message", messageId)

			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		if message.UserId != userId && guild.OwnerId != userId {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Only the author or owner can delete the message",
			})
			return
		}
		// Only message author check required
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

	// Emit delete message to the channel
	h.socketService.EmitDeleteMessage(message.ChannelId, message.ID)

	c.JSON(http.StatusOK, true)
}
