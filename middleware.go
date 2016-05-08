package goniplus

import (
	"github.com/gin-gonic/gin"
)

// GinMiddleware returns goniplus middleware for Gin Framework
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := startRequestTrack(c.Request)
		defer func() {
			if err := recover(); err != nil {
				r.finishRequestTrack(500, true)
			}
		}()
		c.Next()
		r.finishRequestTrack(c.Writer.Status(), false)
	}
}
