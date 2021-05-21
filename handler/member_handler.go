package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
)

type memberReq struct {
	MemberId string `json:"memberId"`
}

// GetMemberSettings handler
func (h *Handler) GetMemberSettings(c *gin.Context) {
	guildId := c.Param("guildId")
	userId := c.MustGet("userId").(string)
	_, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	settings, err := h.guildService.GetMemberSettings(userId, guildId)

	if err != nil {
		log.Printf("Unable to find settings: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// EditMemberSettings handler
func (h *Handler) EditMemberSettings(c *gin.Context) {
	guildId := c.Param("guildId")
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
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	var req model.MemberSettings

	if ok := bindData(c, &req); !ok {
		return
	}

	err = h.guildService.UpdateMemberSettings(&req, userId, guildId)

	if err != nil {
		log.Printf("Unable to update settings for user: %v\n%v", userId, err)
		e := apperrors.NewInternal()

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// GetBanList handler
func (h *Handler) GetBanList(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	bans, err := h.guildService.GetBanList(guildId)

	if err != nil {
		log.Printf("Failed to get banned members: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, bans)
}

// BanMember handler
func (h *Handler) BanMember(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	member, err := h.guildService.GetUser(req.MemberId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", req.MemberId, err)
		e := apperrors.NewNotFound("user", req.MemberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if member.ID == userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot ban yourself",
		})
		return
	}

	guild.Bans = append(guild.Bans, *member)

	if err := h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to ban member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	err = h.guildService.RemoveMember(req.MemberId, guildId)

	if err != nil {
		log.Printf("Failed to ban member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// UnbanMember handler
func (h *Handler) UnbanMember(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	if req.MemberId == userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot unban yourself",
		})
		return
	}

	if err := h.guildService.UnbanMember(req.MemberId, guildId); err != nil {
		log.Printf("Failed to unban member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// KickMember handler
func (h *Handler) KickMember(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "only the owner can do that",
		})
		return
	}

	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	member, err := h.guildService.GetUser(req.MemberId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", req.MemberId, err)
		e := apperrors.NewNotFound("user", req.MemberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if member.ID == userId {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "you cannot kick yourself",
		})
		return
	}

	err = h.guildService.RemoveMember(req.MemberId, guildId)

	if err != nil {
		log.Printf("Failed to kick member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
