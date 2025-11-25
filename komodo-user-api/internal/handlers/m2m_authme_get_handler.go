package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserHandler(gctx *gin.Context) {
	gctx.JSON(http.StatusOK, gin.H{
		"user_id": gctx.GetString("user_id"), // Assuming user_id is set in context by auth middleware
		"message": "Authenticated user details retrieved successfully",
	})
}
