package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"komodo-internal-lib-apis-go/config"
	"komodo-internal-lib-apis-go/crypto/encryption"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyJWT verifies a JWT string. It supports RS256 (using PEM public key
// from env JWT_PUBLIC_KEY) and HS256 (using secret from JWT_HMAC_SECRET).
// Returns true if valid.
func VerifyJWT(tokenString string) (bool, error) {
	if tokenString == "" {
		return false, errors.New("empty token")
	}
	tokenString = strings.TrimSpace(tokenString)

	var err error
	// decrypt if encrypted; DecryptToken lives in encryption.go
	if tokenString, err = encryption.DecryptToken(tokenString); err != nil {
		return false, err
	}

	// Try HS256 if secret present
	if hmac := config.GetConfigValue("JWT_HMAC_SECRET"); hmac != "" {
		_, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(hmac), nil
		})
		if err != nil {
			return false, err
		}
		return true, nil
	}

	// Fallback to RS256 with provided public key
	pubPEM := config.GetConfigValue("JWT_PUBLIC_KEY")
	if pubPEM == "" {
		return false, errors.New("no JWT verification configuration available")
	}

	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return false, errors.New("failed to parse public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// try parsing as PKCS1
		if pk, err2 := x509.ParsePKCS1PublicKey(block.Bytes); err2 == nil {
			pub = pk
		} else {
			return false, err
		}
	}

	var rsaPub *rsa.PublicKey
	switch k := pub.(type) {
		case *rsa.PublicKey:
			rsaPub = k
		default:
			return false, errors.New("unsupported public key type")
	}

	_, err = jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return rsaPub, nil
	})
	if err != nil {
		return false, err
	}

	return true, nil
}
