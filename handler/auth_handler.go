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

/*
 * AuthHandler contains all routes related to account actions (/api/account)
 */

type registerReq struct {
	// Must be unique
	Email string `json:"email" binding:"required,email"`
	// Min 3, max 30 characters.
	Username string `json:"username" binding:"required,gte=3,lte=30"`
	// Min 6, max 150 characters.
	Password string `json:"password" binding:"required,gte=6,lte=150"`
} //@name RegisterRequest

// Register handler creates a new user
// Register godoc
// @Tags Account
// @Summary Create an Account
// @Accept  json
// @Produce  json
// @Param account body registerReq true "Create account"
// @Success 201 {object} model.User
// @Router /account/register [post]
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

type loginReq struct {
	// Must be unique
	Email string `json:"email" binding:"required,email"`
	// Min 6, max 150 characters.
	Password string `json:"password" binding:"required,gte=6,lte=30"`
} //@name LoginRequest

// Login used to authenticate existent user
// Login godoc
// @Tags Account
// @Summary User Login
// @Accept  json
// @Produce  json
// @Param account body loginReq true "Login account"
// @Success 200 {object} model.User
// @Router /account/login [post]
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

// Logout handler removes the current session
// Logout godoc
// @Tags Account
// @Summary User Logout
// @Accept  json
// @Produce  json
// @Param account body loginReq true "Login account"
// @Success 200 {object} model.Success
// @Router /account/logout [post]
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

type forgotRequest struct {
	Email string `json:"email" binding:"required,email"`
} //@name ForgotPasswordRequest

// ForgotPassword sends a password reset email to the requested email
// ForgotPassword godoc
// @Tags Account
// @Summary Forgot Password Request
// @Accept  json
// @Produce  json
// @Param email body forgotRequest true "Forgot Password"
// @Success 200 {object} model.Success
// @Router /account/forgot-password [post]
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

	// No user with the email found
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

type resetRequest struct {
	// The from the email provided token.
	Token string `json:"token" binding:"required"`
	// Min 6, max 150 characters.
	Password string `json:"newPassword" binding:"required"`
	// Must be the same as the password value.
	ConfirmPassword string `json:"confirmNewPassword" binding:"required"`
} //@name ResetPasswordRequest

// ResetPassword resets the users password with the provided token
// ResetPassword godoc
// @Tags Account
// @Summary Reset Password
// @Accept  json
// @Produce  json
// @Param request body resetRequest true "Reset Password"
// @Success 200 {object} model.User
// @Router /account/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetRequest

	if valid := bindData(c, &req); !valid {
		return
	}

	// Check if passwords match
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

// setUserSession saves the users ID in the session
func setUserSession(c *gin.Context, id string) {
	session := sessions.Default(c)
	session.Set("userId", id)
	if err := session.Save(); err != nil {
		fmt.Println(err)
	}
}
