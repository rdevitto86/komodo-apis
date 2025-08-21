package httpapi

import (
	"komodo-address-api/internal/address"
	"net/http"
)

type ValidateResponse struct {
	Valid  bool              `json:"valid"`
	Errors map[string]string `json:"errors,omitempty"`
}

func HandleValidate(w http.ResponseWriter, r *http.Request) {
	addr, err := ParseAddress(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, errorObj(err.Error()))
		return
	}

	errs := address.ValidateAddress(addr)
	resp := ValidateResponse{Valid: len(errs) == 0}

	if len(errs) > 0 {
		resp.Errors = errs
		WriteJSON(w, http.StatusUnprocessableEntity, resp)
		return
	}

	WriteJSON(w, http.StatusOK, resp)
}
