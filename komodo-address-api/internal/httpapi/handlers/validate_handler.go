package handlers

import (
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/httpapi"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidateResponse struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

func HandleValidate(c *gin.Context) {
	var addr address.Address

	if err := c.ShouldBindJSON(&addr); err != nil {
		c.JSON(http.StatusBadRequest, httpapi.Error400(err.Error()))
		return
	}

	errs := address.ValidateAddress(addr)
	res := ValidateResponse{Valid: len(errs) == 0}

	if len(errs) > 0 {
		res.Errors = errs
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	c.JSON(http.StatusOK, res)
}
