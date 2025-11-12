package jwt

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"komodo-internal-lib-apis-go/config"
	logger "komodo-internal-lib-apis-go/logging/runtime"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	cachedPrivateKey *rsa.PrivateKey
	cachedPublicKey  *rsa.PublicKey
	keyMutex         sync.RWMutex
	keysInitialized  bool
)

const (
	MinTokenTTL     = 300    // 5 min
	DefaultTokenTTL = 3600   // 1 hour
	MaxTokenTTL     = 172800 // 2 days
)

// Loads and caches RSA keys at startup
func InitializeKeys() error {
	logger.Info("initializing JWT keys")

	keyMutex.Lock()
	defer keyMutex.Unlock()

	// Load private key
	privateKeyPEM := config.GetConfigValue("JWT_PRIVATE_KEY")
	if privateKeyPEM == "" {
		return fmt.Errorf("JWT_PRIVATE_KEY not configured")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	cachedPrivateKey = privateKey

	// Load public key
	publicKeyPEM := config.GetConfigValue("JWT_PUBLIC_KEY")
	if publicKeyPEM == "" {
		return fmt.Errorf("JWT_PUBLIC_KEY not configured")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	cachedPublicKey = publicKey

	keysInitialized = true
	logger.Info("jwt keys initialized successfully")

	return nil
}

// Returns the cached private key
func GetPrivateKey() (*rsa.PrivateKey, error) {
	keyMutex.RLock()
	defer keyMutex.RUnlock()

	if !keysInitialized {
		return nil, fmt.Errorf("jwt keys not initialized")
	}
	return cachedPrivateKey, nil
}

// Returns the cached public key
func GetPublicKey() (*rsa.PublicKey, error) {
	keyMutex.RLock()
	defer keyMutex.RUnlock()

	if !keysInitialized {
		return nil, fmt.Errorf("jwt keys not initialized")
	}
	return cachedPublicKey, nil
}

// Extracts JWT token from Authorization header or request body
func ExtractTokenFromRequest(req *http.Request) (string, error) {
	// Try Authorization header first
	if auth := req.Header.Get("Authorization"); auth != "" && strings.HasPrefix(auth, "Bearer ") {
		if token := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer ")); token != "" {
			return token, nil
		}
	}

	// Fallback to request body
	var body struct {
		Token string `json:"token,omitempty"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err == nil && body.Token != "" {
		return body.Token, nil
	}
	return "", fmt.Errorf("no token found in Authorization header or request body")
}

// Parses and validates a JWT token using the cached public key
func VerifyToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	publicKey, err := GetPublicKey()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Parse and validate token with leeway for clock skew
	parsedToken, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("token parsing failed: %w", err)
	}
	if !parsedToken.Valid {
		return nil, nil, fmt.Errorf("token is invalid")
	}

	// Extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse token claims")
	}

	return parsedToken, claims, nil
}

// Clamps requested TTL to valid range and logs warnings
func ClampTTL(requested int, clientID string) int {
	if requested < MinTokenTTL {
		logger.Warn("Client %s requested TTL %d below minimum, clamping to %d", clientID, requested, MinTokenTTL)
		return MinTokenTTL
	}
	if requested > MaxTokenTTL {
		logger.Warn("Client %s requested TTL %d above maximum, clamping to %d", clientID, requested, MaxTokenTTL)
		return MaxTokenTTL
	}
	return requested
}

// Creates and signs a JWT token with the provided claims
func SignToken(claims jwt.MapClaims) (string, error) {
	privateKey, err := GetPrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("token signing failed: %w", err)
	}
	return signedToken, nil
}

// Extracts multiple claims and returns them as a map with type-safe conversions.
// Converts float64 to int64 when the value has no decimal component to preserve precision.
// Returns strings, int64, bool, or original value for other types.
func ExtractStringClaims(claims jwt.MapClaims, keys []string) map[string]any {
	clms := map[string]any{}
	for _, k := range keys {
		if value, exists := claims[k]; exists {
			switch v := value.(type) {
				case float64:
					if v == float64(int64(v)) {
						clms[k] = int64(v)
					} else {
						clms[k] = v
					}
				case string, bool:
					clms[k] = v
				default:
					clms[k] = value
			}
		}
	}
	return clms
}

// Extracts a single string claim safely
func ExtractStringClaim(claims jwt.MapClaims, key string) string {
	if value, exists := claims[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}	

// Extracts an int64 claim (useful for timestamps)
func ExtractInt64Claim(claims jwt.MapClaims, key string) (int64, bool) {
	if value, exists := claims[key]; exists {
		if num, ok := value.(float64); ok {
			return int64(num), true
		}
	}
	return 0, false
}

// Creates a standard set of JWT claims with optional custom claims
// Automatically generates a unique JTI (JWT ID) for token tracking and revocation
func CreateStandardClaims(
	issuer string,
	subject string,
	audience string,
	expiresInSeconds int64,
	customClaims map[string]interface{},
) jwt.MapClaims {
	now := time.Now()

	claims := jwt.MapClaims{
		"iss": issuer,
		"sub": subject,
		"aud": audience,
		"exp": now.Add(time.Duration(expiresInSeconds) * time.Second).Unix(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"jti": uuid.NewString(), // Auto-generate unique token ID
	}
	
	// Add custom claims (can override jti if needed)
	for key, value := range customClaims {
		claims[key] = value
	}
	return claims
}

// Checks if a token has expired based on claims
func IsTokenExpired(claims jwt.MapClaims) bool {
	if exp, ok := ExtractInt64Claim(claims, "exp"); ok {
		return time.Now().Unix() > exp
	}
	return false // No expiration claim = assume not expired
}

// Calculates remaining time to live for a token in seconds
func GetTokenTTL(claims jwt.MapClaims) int64 {
	if exp, ok := ExtractInt64Claim(claims, "exp"); ok {
		now := time.Now().Unix()
		ttl := exp - now
		if ttl < 0 { return 0 }
		return ttl
	}
	return 0
}
