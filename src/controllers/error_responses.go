package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func writeInternalErrorResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
func writeInvalidRequestResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
func writeAuthenticationErrorResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": err.Error()})
}