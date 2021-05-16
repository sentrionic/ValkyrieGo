package handler

import (
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,gte=3,lte=30"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Register handler
func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	u := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := h.userService.Register(u)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, u.ID)

	c.JSON(http.StatusCreated, u)
}
