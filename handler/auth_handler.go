package handler

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
)

// REGISTER

// registerReq is not exported, hence the lowercase name
// it is used for validation and json marshalling
type registerReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,gte=3,lte=30"`
	Password string `json:"password" binding:"required,gte=6,lte=150"`
}

// Register handler
func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	user := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	err := h.userService.Register(user)

	if err != nil {
		log.Printf("Failed to sign up user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusCreated, user)
}

// Login

// loginReq is not exported
type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=6,lte=30"`
}

// Login used to authenticate extant user
func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if ok := bindData(c, &req); !ok {
		return
	}

	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.userService.Login(user)

	if err != nil {
		log.Printf("Failed to sign in user: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
}

// LOGOUT

// Logout handler
func (h *Handler) Logout(c *gin.Context) {
	c.Set("user", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()

	if err != nil {
		fmt.Printf("error clearing session: %v", err)
	}

	c.JSON(http.StatusOK, true)
	return
}

// FORGOT PASSWORD

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

// RESET PASSWORD

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

func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
}
