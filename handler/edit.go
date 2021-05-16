package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
)

type editReq struct {
	Username string                `form:"username" binding:"required,min=3,max=30"`
	Email    string                `form:"email" binding:"required,email"`
	Image    *multipart.FileHeader `form:"image" binding:"omitempty"`
}

// Edit handler
func (h *Handler) Edit(c *gin.Context) {
	authUser := c.MustGet("user").(*model.User)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	var req editReq

	if ok := bindData(c, &req); !ok {
		return
	}

	// Should be returned with current imageURL
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
			log.Println("Image is not an allowable mime-type")
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

	err := h.userService.UpdateAccount(authUser)

	if err != nil {
		log.Printf("Failed to update user: %v\n", err.Error())

		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, authUser)
}
