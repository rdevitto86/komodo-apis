package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(gctx *gin.Context) {
	gctx.Header("Content-Type", "application/json")
	gctx.JSON(http.StatusOK, gin.H{"status": "OK"})
}
