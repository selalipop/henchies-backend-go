package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetInfo a default information route to ensure server is running
func (c *Controllers) GetInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"henchies": "hello 👋🏾"})
}
