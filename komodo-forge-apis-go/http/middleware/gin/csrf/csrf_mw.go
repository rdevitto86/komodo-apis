package csrf

import (
	ctxKeys "komodo-forge-apis-go/http/common/context"
	errCodes "komodo-forge-apis-go/http/common/errors"
	errors "komodo-forge-apis-go/http/common/errors/gin"
	hdrSrv "komodo-forge-apis-go/http/headers/eval"
	hdrTypes "komodo-forge-apis-go/http/types"
	httpUtils "komodo-forge-apis-go/http/utils/http"
	logger "komodo-forge-apis-go/logging/runtime"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFMiddleware validates CSRF tokens for state-changing requests from browser clients
// API clients are exempt from CSRF validation
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
			// Check if client type is already set in context
			clientTypeVal, exists := c.Get(string(ctxKeys.CLIENT_TYPE_KEY))
			var clientType string
			if exists {
				clientType, _ = clientTypeVal.(string)
			} else {
				clientType = httpUtils.GetClientType(c.Request)
			}

			// API clients are exempt from CSRF validation
			if clientType == "api" {
				c.Set(string(ctxKeys.CSRF_TOKEN_KEY), "api-client-exempt")
				c.Set(string(ctxKeys.CSRF_VALID_KEY), true)
				c.Next()
				return
			}

			// Browser client - require CSRF token
			if ok, err := hdrSrv.ValidateHeaderValue(hdrTypes.HEADER_X_CSRF_TOKEN, c.Request); !ok || err != nil {
				logger.Error("invalid or missing CSRF token for browser client", err)
				errors.WriteErrorResponse(c, http.StatusBadRequest, "invalid CSRF token", errCodes.ERR_INVALID_REQUEST)
				c.Abort()
				return
			}
		}

		// Set CSRF validation flags in context
		c.Set(string(ctxKeys.CSRF_TOKEN_KEY), "")
		c.Set(string(ctxKeys.CSRF_VALID_KEY), true)

		c.Next()
	}
}
