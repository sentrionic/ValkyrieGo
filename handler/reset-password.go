package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"net/http"
)

type resetRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmNewPassword" binding:"required"`
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetRequest

	if valid := bindData(c, &req); !valid {
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Passwords do not match",
		})
	}

	ctx := c.Request.Context()
	user, err := h.userService.ResetPassword(ctx, req.Password, req.Token)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
	return
}
