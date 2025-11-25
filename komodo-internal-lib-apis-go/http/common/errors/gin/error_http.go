package ginerrors

import (
	"fmt"
	errorTypes "komodo-internal-lib-apis-go/http/common/errors"
	httptypes "komodo-internal-lib-apis-go/http/types"
	"time"

	"github.com/gin-gonic/gin"
)

// Writes a standardized error response for Gin
func WriteErrorResponse(c *gin.Context, status int, message string, errCode string) {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	c.JSON(status, errorTypes.ErrorStandard{
		Status:    status,
		Code:      errCode,
		Message:   message,
		RequestId: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Writes a verbose error response with API details for Gin
func WriteErrorVerboseResponse(c *gin.Context, status int, message string, errCode string, apiError any) {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	c.JSON(status, errorTypes.ErrorVerbose{
		Status:    status,
		Code:      errCode,
		Message:   message,
		APIName:   c.Request.URL.Path,
		APIError:  fmt.Sprintf("%v", apiError),
		RequestId: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Forwards an existing APIResponse error to Gin context
func ForwardErrorResponse(c *gin.Context, res *httptypes.APIResponse) {
	c.JSON(res.Status, errorTypes.ErrorStandard{
		Status:    res.Status,
		Code:      res.Error.Code,
		Message:   res.Error.Message,
		RequestId: res.RequestID,
		Timestamp: res.Timestamp,
	})
}
