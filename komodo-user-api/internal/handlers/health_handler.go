package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthHandler(g *gin.Context) {
	g.Header("Content-Type", "application/json")
	g.JSON(http.StatusOK, gin.H{"status": "OK"})
}
