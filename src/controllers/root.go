package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controllers) GetInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"henchies": "hello"})
}
