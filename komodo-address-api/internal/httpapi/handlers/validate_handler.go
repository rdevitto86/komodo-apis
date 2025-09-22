package handlers

import (
	"komodo-address-api/internal/address"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidateResponse struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

func HandleValidate(ctx *gin.Context) {
	var addr address.Address

	if err := ctx.ShouldBindJSON(&addr); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	errs := address.ValidateAddress(addr)
	res := ValidateResponse{Valid: len(errs) == 0}

	if len(errs) > 0 {
		res.Errors = errs
		ctx.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}
