package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
)

type changeRequest struct {
	CurrentPassword    string `json:"currentPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required,gte=6"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required,gte=6"`
}

// ChangePassword handler
func (h *Handler) ChangePassword(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)
	var req changeRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	if req.NewPassword != req.ConfirmNewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"field":   "password",
			"message": "passwords do not match",
		})
		return
	}

	err := h.userService.ChangePassword(req.NewPassword, authUser)

	if err != nil {
		log.Printf("Failed to change password: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
