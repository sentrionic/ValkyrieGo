package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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

func (h *Handler) GetGuildMembers(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

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

type createRequest struct {
	Name string `json:"name" binding:"required,gte=3"`
}

func (h *Handler) CreateGuild(c *gin.Context) {
	var req createRequest

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

	guild.Members = append(guild.Members, *authUser)

	if err := h.guildService.CreateGuild(&guild); err != nil {
		log.Printf("Failed to create guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	channel := model.Channel{
		GuildID: guild.ID,
		Name:    "general",
	}

	if err := h.guildService.CreateDefaultChannel(&channel); err != nil {
		log.Printf("Failed to create guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
	return
}

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

	if len(authUser.Guilds) >= 100 {
		e := apperrors.NewBadRequest("guild limit is 100")

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if strings.Contains(req.Link, "/") {
		req.Link = req.Link[strings.LastIndex(req.Link, "/")+1:]
	}

	ctx := context.Background()
	guildId, err := h.guildService.GetGuildIdFromInvite(ctx, req.Link)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

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

	channel, _ := h.guildService.GetDefaultChannel(guildId)

	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
	return
}

func (h *Handler) LeaveGuild(c *gin.Context) {
	c.JSON(http.StatusOK, "LeaveGuild")
}

func (h *Handler) EditGuild(c *gin.Context) {
	c.JSON(http.StatusOK, "EditGuild")
}

func (h *Handler) DeleteGuild(c *gin.Context) {
	c.JSON(http.StatusOK, "DeleteGuild")
}

func (h *Handler) GetMemberSettings(c *gin.Context) {
	c.JSON(http.StatusOK, "GetMemberSettings")
}

func (h *Handler) EditMemberSettings(c *gin.Context) {
	c.JSON(http.StatusOK, "EditMemberSettings")
}

func (h *Handler) GetBanList(c *gin.Context) {
	c.JSON(http.StatusOK, "GetBanList")
}

func (h *Handler) BanMember(c *gin.Context) {
	c.JSON(http.StatusOK, "BanMember")
}

func (h *Handler) UnbanMember(c *gin.Context) {
	c.JSON(http.StatusOK, "UnbanMember")
}

func (h *Handler) KickMember(c *gin.Context) {
	c.JSON(http.StatusOK, "KickMember")
}

// isMember checks if the given user is a member
func isMember(guild *model.Guild, userId string) bool {
	for _, v := range guild.Members {
		if v.ID == userId {
			return true
		}
	}
	return false
}
