package handlers

import (
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/httpapi"
	"net/http"
)

type NormalizeResponse struct {
	Address address.Address `json:"address"`
}

func HandleNormalize(w http.ResponseWriter, r *http.Request) {
	addr, err := httpapi.ParseAddress(r)

	if err != nil {
		httpapi.WriteJSON(w, http.StatusBadRequest, httpapi.ErrorObj(err.Error()))
		return
	}

	n := address.NormalizeAddress(addr)
	httpapi.WriteJSON(w, http.StatusOK, NormalizeResponse{Address: n})
}
