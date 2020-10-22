package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func writeInternalErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
}
func writeInvalidRequestResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("%v", err)})
}
func writeAuthenticationErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusNetworkAuthenticationRequired, gin.H{"error": fmt.Sprintf("%v", err)})
}

func writeSuccessIfNoErrors(ctx *gin.Context, errors ...error) {
	for _, err := range errors {
		if err != nil {
			writeInternalErrorResponse(ctx, err)
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
