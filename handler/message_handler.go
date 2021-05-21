package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetMessages(c *gin.Context)  {
	c.JSON(http.StatusOK, "GetMessages")
}

func (h *Handler) CreateMessage(c *gin.Context)  {
	c.JSON(http.StatusOK, "CreateMessage")
}

func (h *Handler) EditMessage(c *gin.Context)  {
	c.JSON(http.StatusOK, "EditMessage")
}

func (h *Handler) DeleteMessage(c *gin.Context)  {
	c.JSON(http.StatusOK, "DeleteMessage")
}