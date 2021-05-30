package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
 * GuildHandler contains all routes related to guild actions (/api/guilds)
 */

// GetUserGuilds returns the current users guilds
func (h *Handler) GetUserGuilds(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	guilds, err := h.guildService.GetUserGuilds(userId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, guilds)
}

// GetGuildMembers returns the given guild's members
func (h *Handler) GetGuildMembers(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if a member
	if !isMember(guild, userId) {
		e := apperrors.NewAuthorization("not a member")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	members, err := h.guildService.GetGuildMembers(userId, guildId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

type createGuildRequest struct {
	Name string `json:"name" binding:"required,gte=3,lte=30"`
}

// CreateGuild creates a guild
func (h *Handler) CreateGuild(c *gin.Context) {
	var req createGuildRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	userId := c.MustGet("userId").(string)

	authUser, err := h.guildService.GetUser(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user is already in 100 guilds
	if len(authUser.Guilds) >= 100 {
		e := apperrors.NewBadRequest("guild limit is 100")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild := model.Guild{
		Name:    req.Name,
		OwnerId: userId,
	}

	// Add the current user as a member
	guild.Members = append(guild.Members, *authUser)

	if err := h.guildService.CreateGuild(&guild); err != nil {
		log.Printf("Failed to create guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Create the default 'general' channel for the guild
	channel := model.Channel{
		GuildID:  &guild.ID,
		Name:     "general",
		IsPublic: true,
	}

	if err := h.channelService.CreateChannel(&channel); err != nil {
		log.Printf("Failed to create guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
	return
}

// editGuildRequest specifies the form to edit the guild.
// If Image is not nil then the guilds icon got changed.
// If Icon is not nil then the guild kept its old one.
// If both are nil then the icon got reset.
type editGuildRequest struct {
	Name  string                `form:"name" binding:"required,gte=3,lte=30"`
	Image *multipart.FileHeader `form:"image" binding:"omitempty"`
	Icon  *string               `form:"icon" binding:"omitempty"`
}

// EditGuild edits the given guild
func (h *Handler) EditGuild(c *gin.Context) {
	var req editGuildRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

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
			"error": "must be the owner for that",
		})
		return
	}

	guild.Name = req.Name

	// Guild icon got changed
	if req.Image != nil {
		// Validate image mime-type is allowable
		mimeType := req.Image.Header.Get("Content-Type")

		if valid := isAllowedImageType(mimeType); !valid {
			log.Println("Image is not an allowable mime-type")
			e := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		directory := fmt.Sprintf("valkyrie_go/guilds/%s", guild.ID)
		url, err := h.userService.ChangeAvatar(req.Image, directory)

		if err != nil {
			fmt.Println(err)
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		if guild.Icon != nil {
			_ = h.userService.DeleteImage(*guild.Icon)
		}
		guild.Icon = &url
		// Guild kept its old icon
	} else if req.Icon != nil {
		guild.Icon = req.Icon
		// Guild reset its icon
	} else {
		guild.Icon = nil
	}

	if err := h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to update guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit guild changes to guild members
	h.socketService.EmitEditGuild(guild)

	c.JSON(http.StatusCreated, true)
	return
}

// GetInvite creates an invite for the given channel
// The isPermanent query parameter specifies if the invite
// should not be deleted after it got used
func (h *Handler) GetInvite(c *gin.Context) {
	guildId := c.Param("guildId")
	permanent := c.Query("isPermanent")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)
	// Must be a member to create an invite
	if !isMember(guild, userId) {
		e := apperrors.NewBadRequest("must be a member to fetch an invite")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	isPermanent := false
	if permanent != "" {
		isPermanent, err = strconv.ParseBool(permanent)

		if err != nil {
			e := apperrors.NewBadRequest("isPermanent is not a boolean")

			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}
	}

	ctx := context.Background()
	link, err := h.guildService.GenerateInviteLink(ctx, guild.ID, isPermanent)

	if isPermanent {
		guild.InviteLinks = append(guild.InviteLinks, link)
		_ = h.guildService.UpdateGuild(guild)
	}

	origin := os.Getenv("CORS_ORIGIN")
	c.JSON(http.StatusOK, fmt.Sprintf("%s/%s", origin, link))
}

// DeleteGuildInvites removes all permanent invites from the given guild
func (h *Handler) DeleteGuildInvites(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

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
			"error": "only the owner can invalidate invites",
		})
		return
	}

	ctx := context.Background()
	h.guildService.InvalidateInvites(ctx, guild)
	guild.InviteLinks = make(pq.StringArray, 0)

	if err := h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to join guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
	return
}

type joinReq struct {
	Link string `json:"link" binding:"required"`
}

// JoinGuild adds the current user to invited guild
func (h *Handler) JoinGuild(c *gin.Context) {
	var req joinReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	userId := c.MustGet("userId").(string)

	authUser, err := h.guildService.GetUser(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user has reached the guild limit
	if len(authUser.Guilds) >= 100 {
		e := apperrors.NewBadRequest("guild limit is 100")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the link contains the domain, remove it
	if strings.Contains(req.Link, "/") {
		req.Link = req.Link[strings.LastIndex(req.Link, "/")+1:]
	}

	ctx := context.Background()
	guildId, err := h.guildService.GetGuildIdFromInvite(ctx, req.Link)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid Link or the server got deleted",
		})
		return
	}

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Invalid Link or the server got deleted",
		})
		return
	}

	// Check if the user is banned from the guild
	if isBanned(guild, authUser.ID) {
		e := apperrors.NewBadRequest("You are banned from this server")

		c.JSON(e.Status(), gin.H{
			"message": e,
		})
		return
	}

	// Check if the user is already a member
	if isMember(guild, authUser.ID) {
		e := apperrors.NewBadRequest("already a member")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild.Members = append(guild.Members, *authUser)

	if err := h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to join guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit new member to the guild
	h.socketService.EmitAddMember(guild.ID, authUser)

	channel, _ := h.guildService.GetDefaultChannel(guildId)

	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
	return
}

// LeaveGuild leaves the given guild
func (h *Handler) LeaveGuild(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId == userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "the owner cannot leave their server",
		})
		return
	}

	if err := h.guildService.RemoveMember(userId, guildId); err != nil {
		log.Printf("Failed to leave guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the member from the guild
	h.socketService.EmitRemoveMember(guild.ID, userId)

	c.JSON(http.StatusOK, true)
}

// DeleteGuild deletes the given guild
func (h *Handler) DeleteGuild(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

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
			"error": "only the owner can delete their server",
		})
		return
	}

	members := make([]string, 0)
	for _, member := range guild.Members {
		members = append(members, member.ID)
	}

	if err := h.guildService.DeleteGuild(guildId); err != nil {
		log.Printf("Failed to leave guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the guild to its members
	h.socketService.EmitDeleteGuild(guildId, members)

	c.JSON(http.StatusOK, true)
}

// isMember checks if the given user is a member of the guild
func isMember(guild *model.Guild, userId string) bool {
	for _, v := range guild.Members {
		if v.ID == userId {
			return true
		}
	}
	return false
}

// isBanned checks if the given user is banned from the guild
func isBanned(guild *model.Guild, userId string) bool {
	for _, v := range guild.Bans {
		if v.ID == userId {
			return true
		}
	}
	return false
}
