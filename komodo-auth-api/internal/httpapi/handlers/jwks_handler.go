package handlers

import (
	"encoding/json"
	logger "komodo-internal-lib-apis-go/services/logger/runtime"
	"net/http"
	"os"
)

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

const JWKS_PATH = "./config/jwks.json"

func JWKSHandler(wtr http.ResponseWriter, req *http.Request) {
	wtr.Header().Set("Content-Type", "application/json")

	var (
		jwks JWKS
		data []byte
		err error
	)

	if data, err = os.ReadFile(JWKS_PATH); err == nil {
		if err = json.Unmarshal(data, &jwks); err == nil {
			// If "kid" query param is provided, filter keys
			kid := req.URL.Query().Get("kid")
			if kid != "" {
				for _, key := range jwks.Keys {
					if key.Kid == kid {
						json.NewEncoder(wtr).Encode(key)
						break
					}
				}
			} else {
				json.NewEncoder(wtr).Encode(jwks)
			}

			wtr.WriteHeader(http.StatusOK)
			return
		}
	}

	logger.Error("failed to read JWKS file: " + err.Error(), req)
	http.Error(wtr, "failed to read JWKS file", http.StatusInternalServerError)
}