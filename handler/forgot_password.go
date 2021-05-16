package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type forgotRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req forgotRequest
	if valid := bindData(c, &req); !valid {
		return
	}

	user, err := h.userService.GetByEmail(req.Email)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Try again later",
		})
		return
	}

	if user.ID == "" {
		c.JSON(http.StatusOK, true)
		return
	}

	ctx := c.Request.Context()
	err = h.userService.ForgotPassword(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong. Try again later",
		})
		return
	}

	c.JSON(http.StatusCreated, true)
	return
}
