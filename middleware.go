package goniplus

import (
	"github.com/gin-gonic/gin"
)

// GinMiddleware returns goniplus middleware for Gin Framework
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := startHTTPTrack(c.Request)
		c.Next()
		finishHTTPTrack(t, c.Writer.Status())
	}
}
