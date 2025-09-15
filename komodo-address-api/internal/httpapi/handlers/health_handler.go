package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status string `json:"status"`
}

func HandleHealth(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, healthResponse{Status: "OK"})
}