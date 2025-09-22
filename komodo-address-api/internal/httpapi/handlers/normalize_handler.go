package handlers

import (
	"komodo-address-api/internal/address"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NormalizeResponse struct {
	Address address.Address `json:"address"`
}

func HandleNormalize(ctx *gin.Context) {
	var addr address.Address

	if err := ctx.ShouldBindJSON(&addr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	normalized := address.NormalizeAddress(addr)
	ctx.JSON(http.StatusOK, NormalizeResponse{Address: normalized})
}
