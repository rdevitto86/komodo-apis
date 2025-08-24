package httpapi

import (
	encodingjson "encoding/json"
	"fmt"
	"komodo-address-api/internal/address"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, val any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = encodingjson.NewEncoder(w).Encode(val)
}

func ParseAddress(req *http.Request) (address.Address, error) {
	defer req.Body.Close()

	dec := encodingjson.NewDecoder(req.Body)
	dec.DisallowUnknownFields()

	var addr address.Address

	if err := dec.Decode(&addr); err != nil {
		return address.Address{}, fmt.Errorf("invalid JSON body: %w", err)
	}
	return addr, nil
}

func Method(method string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			WriteJSON(w, http.StatusMethodNotAllowed, Error405("method not allowed"))
			return
		}
		handler(w, r)
	}
}
