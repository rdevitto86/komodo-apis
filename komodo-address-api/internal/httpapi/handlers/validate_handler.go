package handlers

import (
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/httpapi"
	"net/http"
)

type ValidateResponse struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

func HandleValidate(writer http.ResponseWriter, req *http.Request) {
	addr, err := httpapi.ParseAddress(req)

	if err != nil {
		httpapi.WriteJSON(writer, http.StatusBadRequest, httpapi.Error400(err.Error()))
		return
	}

	errs := address.ValidateAddress(addr)
	res := ValidateResponse{Valid: len(errs) == 0}

	if len(errs) > 0 {
		res.Errors = errs
		httpapi.WriteJSON(writer, http.StatusUnprocessableEntity, res)
		return
	}

	httpapi.WriteJSON(writer, http.StatusOK, res)
}
