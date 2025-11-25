package context

import (
	"github.com/gin-gonic/gin"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		gctx.Next()
	}
}