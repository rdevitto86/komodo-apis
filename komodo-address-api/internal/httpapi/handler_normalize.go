package httpapi

import (
	"komodo-address-api/internal/address"
	"net/http"
)

type NormalizeResponse struct {
	Address address.Address `json:"address"`
}

func HandleNormalize(w http.ResponseWriter, r *http.Request) {
	addr, err := ParseAddress(r)

	if err != nil {
		WriteJSON(w, http.StatusBadRequest, errorObj(err.Error()))
		return
	}

	n := address.NormalizeAddress(addr)
	WriteJSON(w, http.StatusOK, NormalizeResponse{Address: n})
}
