package handler

import (
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
)

/*
 * ChannelHandler contains all routes related to channel actions (/api/channels)
 */

// GuildChannels returns the given guild's channels
func (h *Handler) GuildChannels(c *gin.Context) {
	guildId := c.Param("id")
	userId := c.MustGet("userId").(string)

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Only get the channels if the user is a member
	if !isMember(guild, userId) {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channels, err := h.channelService.GetChannels(userId, guildId)

	if err != nil {
		log.Printf("Unable to find channels for guild id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("channels", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// channelReq specifies the input form for creating a channel
// IsPublic and Members do not need to be specified if you want
// to create a public channel
type channelReq struct {
	Name     string   `json:"name" binding:"required,gte=3,lte=30"`
	IsPublic *bool    `json:"isPublic"`
	Members  []string `json:"members" binding:"omitempty"`
}

// CreateChannel creates a channel for the given guild param
func (h *Handler) CreateChannel(c *gin.Context) {
	var req channelReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	userId := c.MustGet("userId").(string)
	guildId := c.Param("id")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	// Check if the server already has 50 channels
	if len(guild.Channels) >= 50 {
		e := apperrors.NewBadRequest("channel limit is 50")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channel := model.Channel{
		Name:     req.Name,
		IsPublic: true,
		GuildID:  &guildId,
	}

	// Channel is private
	if req.IsPublic != nil && !*req.IsPublic {
		channel.IsPublic = false

		// Add the current user to the members if they are not in there
		if !containsUser(req.Members, userId) {
			req.Members = append(req.Members, userId)
		}
		members, err := h.guildService.FindUsersByIds(req.Members, guildId)

		if err != nil {
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		// Create private channel members
		for _, m := range *members {
			channel.PCMembers = append(channel.PCMembers, m)
		}
	}

	if err := h.channelService.CreateChannel(&channel); err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	guild.Channels = append(guild.Channels, channel)

	if err := h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	response := channel.SerializeChannel()

	// Emit the new channel to the guild members
	h.socketService.EmitNewChannel(guildId, &response)

	c.JSON(http.StatusCreated, response)
	return
}

// PrivateChannelMembers returns the ids of all members
// that are part of the channel
func (h *Handler) PrivateChannelMembers(c *gin.Context) {
	channelId := c.Param("id")
	userId := c.MustGet("userId").(string)

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	// Public channels do not have any private members
	if channel.IsPublic {
		var empty = make([]string, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	members, err := h.channelService.GetPrivateChannelMembers(channelId)

	if err != nil {
		log.Printf("Unable to find members for channel: %v\n%v", channelId, err)
		e := apperrors.NewNotFound("members", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

// DirectMessages returns a list of the current users DMs
func (h *Handler) DirectMessages(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	channels, err := h.channelService.GetDirectMessages(userId)

	if err != nil {
		log.Printf("Unable to find dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the user does not have any dms, return an empty array
	if len(*channels) == 0 {
		var empty = make([]model.DirectMessage, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, channels)
}

// GetOrCreateDM gets the DM with the given member and creates it
// if it does not already exist
func (h *Handler) GetOrCreateDM(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	memberId := c.Param("id")

	if userId == memberId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot dm yourself",
		})
		return
	}

	member, err := h.friendService.GetMemberById(memberId)

	if err != nil {
		log.Printf("Unable to find member for id: %v\n%v", memberId, err)
		e := apperrors.NewNotFound("member", memberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// check if dm channel already exists with these members
	dm, err := h.channelService.GetDirectMessageChannel(userId, memberId)

	if err != nil {
		log.Printf("Unable to find or create dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// dm already exists
	if dm != nil && *dm != "" {
		_ = h.channelService.SetDirectMessageStatus(*dm, userId, true)
		c.JSON(http.StatusOK, toDMChannel(member, *dm, userId))
		return
	}

	// Create the dm channel between the current user and the member
	id, _ := gonanoid.Nanoid(20)
	channel := model.Channel{
		Name:     id,
		IsPublic: false,
		IsDM:     true,
	}

	// Create the DM channel
	if err := h.channelService.CreateChannel(&channel); err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Add the users to it
	ids := []string{userId, memberId}
	err = h.channelService.AddDMChannelMembers(ids, channel.ID, userId)

	if err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, toDMChannel(member, channel.ID, userId))
}

// toDMChannel returns the DM response for the given channel and member
func toDMChannel(member *model.User, channelId string, userId string) model.DirectMessage {
	return model.DirectMessage{
		Id: channelId,
		User: model.DMUser{
			Id:       member.ID,
			Username: member.Username,
			Image:    member.Image,
			IsOnline: member.IsOnline,
			IsFriend: isFriend(member, userId),
		},
	}
}

// EditChannel edits the specified channel
func (h *Handler) EditChannel(c *gin.Context) {
	var req channelReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("guild", *channel.GuildID)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	// Used to be private and now is public
	if *req.IsPublic && !channel.IsPublic {
		err = h.channelService.CleanPCMembers(channelId)
		if err != nil {
			log.Printf("error removing pc members: %v", err)
		}
	}

	channel.IsPublic = *req.IsPublic
	channel.Name = req.Name

	// Member Changes
	if !*req.IsPublic {
		// Check if the array contains the current member
		if !containsUser(req.Members, userId) {
			req.Members = append(req.Members, userId)
		}

		// Current members of the channel
		current := make([]string, 0)
		for _, member := range channel.PCMembers {
			current = append(current, member.ID)
		}

		// Newly added members
		newMembers := difference(req.Members, current)
		// Members that got removed
		toRemove := difference(current, req.Members)

		err = h.channelService.AddPrivateChannelMembers(newMembers, channelId)
		if err != nil {
			log.Printf("Failed to add new members: %v\n", err.Error())
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		err = h.channelService.RemovePrivateChannelMembers(toRemove, channelId)
		if err != nil {
			log.Printf("Failed to add remove members: %v\n", err.Error())
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}
	}

	if err := h.channelService.UpdateChannel(channel); err != nil {
		log.Printf("Failed to update channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit the channel changes to the guild members
	response := channel.SerializeChannel()
	h.socketService.EmitEditChannel(*channel.GuildID, &response)

	c.JSON(http.StatusCreated, true)
	return
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// DeleteChannel removes the given channel from the guild
func (h *Handler) DeleteChannel(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	// Check if the guild has the minimum amount of channels
	if len(guild.Channels) == 1 {
		e := apperrors.NewBadRequest("A server needs at least one channel")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if err := h.channelService.DeleteChannel(channel); err != nil {
		log.Printf("Failed to delete channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the channel from the guild
	h.socketService.EmitDeleteChannel(channel)

	c.JSON(http.StatusCreated, true)
	return
}

// CloseDM closes the DM on the current users side
func (h *Handler) CloseDM(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	dm, err := h.channelService.GetDMByUserAndChannel(userId, channelId)

	if err != nil || dm == "" {
		log.Printf("Unable to find or create dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	_ = h.channelService.SetDirectMessageStatus(channelId, userId, false)

	c.JSON(http.StatusOK, true)
}

// containsUser checks if the array contains the user
func containsUser(members []string, userId string) bool {
	for _, m := range members {
		if m == userId {
			return true
		}
	}
	return false
}
