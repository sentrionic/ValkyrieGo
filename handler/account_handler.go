package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
)

// CURRENT USER

// Me handler calls services for getting
// a user's details
func (h *Handler) Me(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	u, err := h.userService.Get(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, u)
}

// EDIT USER

type editReq struct {
	Username string                `form:"username" binding:"required,min=3,max=30"`
	Email    string                `form:"email" binding:"required,email"`
	Image    *multipart.FileHeader `form:"image" binding:"omitempty"`
}

// Edit handler
func (h *Handler) Edit(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	var req editReq

	if ok := bindData(c, &req); !ok {
		return
	}

	authUser, err := h.userService.Get(userId)

	if err != nil {
		err := errors.New("provided session is invalid")
		c.JSON(401, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}

	authUser.Username = req.Username

	if authUser.Email != req.Email {
		inUse := h.userService.CheckEmail(req.Email)

		if inUse {
			c.JSON(http.StatusBadRequest, gin.H{
				"field":   "Email",
				"message": "email already in use",
			})
			return
		} else {
			authUser.Email = req.Email
		}
	}

	if req.Image != nil {

		// Validate image mime-type is allowable
		mimeType := req.Image.Header.Get("Content-Type")

		if valid := isAllowedImageType(mimeType); !valid {
			e := apperrors.NewBadRequest("imageFile must be 'image/jpeg' or 'image/png'")
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		directory := fmt.Sprintf("valkyrie_go/users/%s", authUser.ID)
		url, err := h.userService.ChangeAvatar(req.Image, directory)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err,
			})
			return
		}

		_ = h.userService.DeleteImage(authUser.Image)

		authUser.Image = url
	}

	err = h.userService.UpdateAccount(authUser)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, authUser)
}

// CHANGE PASSWORD

type changeRequest struct {
	CurrentPassword    string `json:"currentPassword" binding:"required"`
	NewPassword        string `json:"newPassword" binding:"required,gte=6"`
	ConfirmNewPassword string `json:"confirmNewPassword" binding:"required,gte=6"`
}

// ChangePassword handler
func (h *Handler) ChangePassword(c *gin.Context) {
	userId := c.MustGet("userId").(string)
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

	authUser, err := h.userService.Get(userId)

	if err != nil {
		err := errors.New("provided session is invalid")
		c.JSON(401, gin.H{
			"error": err,
		})
		c.Abort()
		return
	}

	err = h.userService.ChangePassword(req.NewPassword, authUser)

	if err != nil {
		log.Printf("Failed to change password: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}
