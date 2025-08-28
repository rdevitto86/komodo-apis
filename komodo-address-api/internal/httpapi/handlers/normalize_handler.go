package handlers

import (
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/httpapi"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NormalizeResponse struct {
	Address address.Address `json:"address"`
}

func HandleNormalize(c *gin.Context) {
	var addr address.Address

	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, httpapi.Error400(err.Error()))
		return
	}

	normalized := address.NormalizeAddress(addr)
	c.JSON(http.StatusOK, NormalizeResponse{Address: normalized})
}
