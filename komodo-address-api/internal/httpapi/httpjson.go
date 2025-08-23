package httpapi

import (
	encodingjson "encoding/json"
	"fmt"
	"komodo-address-api/internal/address"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = encodingjson.NewEncoder(w).Encode(v)
}

func ParseAddress(r *http.Request) (address.Address, error) {
	defer r.Body.Close()
	dec := encodingjson.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var a address.Address

	if err := dec.Decode(&a); err != nil {
		return address.Address{}, fmt.Errorf("invalid JSON body: %w", err)
	}

	return a, nil
}

func Method(m string, h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != m {
			w.Header().Set("Allow", m)
			WriteJSON(w, http.StatusMethodNotAllowed, ErrorObj("method not allowed"))
			return
		}
		h(w, r)
	}
}
