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

func HandleValidate(w http.ResponseWriter, r *http.Request) {
	addr, err := httpapi.ParseAddress(r)

	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorObj(err.Error()))
		return
	}

	errs := address.ValidateAddress(addr)
	resp := ValidateResponse{Valid: len(errs) == 0}

	if len(errs) > 0 {
		resp.Errors = errs
		httpapi.WriteJSON(w, http.StatusUnprocessableEntity, resp)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, resp)
}
