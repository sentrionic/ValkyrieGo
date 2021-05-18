package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
