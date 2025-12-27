package errors

import (
	"fmt"
	httpApi "komodo-forge-apis-go/http"
	"komodo-forge-apis-go/http/common/errors"
	"time"

	"github.com/gin-gonic/gin"
)

// Writes a standardized error response for Gin
func WriteErrorResponse(gctx *gin.Context, status int, message string, errCode string) {
	requestID := gctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	gctx.JSON(status, errors.ErrorStandard{
		Status:    status,
		Code:      errCode,
		Message:   message,
		RequestId: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Writes a verbose error response with API details for Gin
func WriteErrorVerboseResponse(gctx *gin.Context, status int, message string, errCode string, apiError any) {
	requestID := gctx.GetHeader("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	gctx.JSON(status, errors.ErrorVerbose{
		Status:    status,
		Code:      errCode,
		Message:   message,
		APIName:   gctx.Request.URL.Path,
		APIError:  fmt.Sprintf("%v", apiError),
		RequestId: requestID,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Forwards an existing APIResponse error to Gin context
func ForwardErrorResponse(gctx *gin.Context, res *httpApi.APIResponse) {
	gctx.JSON(res.Status, errors.ErrorStandard{
		Status:    res.Status,
		Code:      res.Error.Code,
		Message:   res.Error.Message,
		RequestId: res.RequestID,
		Timestamp: res.Timestamp,
	})
}
