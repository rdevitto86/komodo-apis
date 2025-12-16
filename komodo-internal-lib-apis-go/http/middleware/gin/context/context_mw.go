package context

import (
	"github.com/gin-gonic/gin"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		g.Next()
	}
}