package handlers

import (
	"komodo-address-api/internal/address"
	"komodo-address-api/internal/httpapi"
	"net/http"
)

type NormalizeResponse struct {
	Address address.Address `json:"address"`
}

func HandleNormalize(writer http.ResponseWriter, req *http.Request) {
	addr, err := httpapi.ParseAddress(req)

	if err != nil {
		httpapi.WriteJSON(writer, http.StatusBadRequest, httpapi.Error400(err.Error()))
		return
	}

	n := address.NormalizeAddress(addr)
	httpapi.WriteJSON(writer, http.StatusOK, NormalizeResponse{Address: n})
}
