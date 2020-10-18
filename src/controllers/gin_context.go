package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func WriteInternalErrorResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
func WriteInvalidRequestResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
func WriteAuthenticationErrorResponse(c * gin.Context, err error)  {
	c.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": err.Error()})
}