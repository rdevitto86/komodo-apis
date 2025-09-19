package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// EncryptToken encrypts plain using AES-GCM with key from JWT_ENC_KEY (base64).
// If JWT_ENC_KEY is empty, returns the input unchanged.
func EncryptToken(plain string) (string, error) {
    keyB64 := os.Getenv("JWT_ENC_KEY")
    if keyB64 == "" {
      return plain, nil
    }

    key, err := base64.StdEncoding.DecodeString(keyB64)
    if err != nil {
      return "", err
    }
    if len(key) != 16 && len(key) != 24 && len(key) != 32 {
      return "", errors.New("invalid JWT_ENC_KEY length (must be 16/24/32 bytes after base64)")
    }

    block, err := aes.NewCipher(key)
    if err != nil {
      return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
      return "", err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
      return "", err
    }

    ct := gcm.Seal(nil, nonce, []byte(plain), nil)
    out := append(nonce, ct...)
    return base64.StdEncoding.EncodeToString(out), nil
}

// DecryptToken decrypts a token encrypted with EncryptToken.
// If JWT_ENC_KEY is empty, returns the input unchanged.
func DecryptToken(enc string) (string, error) {
    keyB64 := os.Getenv("JWT_ENC_KEY")
    if keyB64 == "" {
      return enc, nil
    }

    key, err := base64.StdEncoding.DecodeString(keyB64)
    if err != nil {
      return "", err
    }

    data, err := base64.StdEncoding.DecodeString(enc)
    if err != nil {
    	return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
      return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
      return "", err
    }

    nonceSize := gcm.NonceSize()
    if len(data) < nonceSize {
      return "", errors.New("malformed ciphertext")
    }

    nonce, ct := data[:nonceSize], data[nonceSize:]
    plain, err := gcm.Open(nil, nonce, ct, nil)
    if err != nil {
      return "", err
    }
	
    return string(plain), nil
}

// SignJWT signs claims with either HMAC (JWT_HMAC_SECRET) or RSA (JWT_PRIVATE_KEY PEM).
// ttl is optional; pass 0 to not set exp.
func SignJWT(claims jwt.MapClaims, ttl time.Duration) (string, error) {
    // set exp if ttl provided
    if ttl > 0 {
      claims["exp"] = time.Now().Add(ttl).Unix()
    }

    // prefer HMAC if provided
    if secret := os.Getenv("JWT_HMAC_SECRET"); secret != "" {
      tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
      return tok.SignedString([]byte(secret))
    }

    // fallback to RSA private key
    privPEM := os.Getenv("JWT_PRIVATE_KEY")
    if privPEM == "" {
      return "", errors.New("no signing configuration (JWT_HMAC_SECRET or JWT_PRIVATE_KEY)")
    }

    privKey, err := parseRSAPrivateKeyFromPEM([]byte(privPEM))
    if err != nil {
      return "", err
    }

    tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    return tok.SignedString(privKey)
}

func parseRSAPrivateKeyFromPEM(pemBytes []byte) (*rsa.PrivateKey, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil {
      return nil, errors.New("failed to parse private key PEM")
    }
    if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
      return key, nil
    }
    priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
      return nil, err
    }
    if rsaKey, ok := priv.(*rsa.PrivateKey); ok {
      return rsaKey, nil
    }
    return nil, errors.New("unsupported private key type")
}
