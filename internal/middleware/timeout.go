package middleware

import (
	"net/http"
	"time"
	"user-service/internal/logger"

	"github.com/gin-gonic/gin"
)

// TimeoutMiddleware adds timeout handling to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a channel to signal timeout
		timeoutChan := make(chan struct{})

		// Start a goroutine that will close the channel after timeout
		go func() {
			time.Sleep(timeout)
			close(timeoutChan)
		}()

		// Create a channel to signal request completion
		doneChan := make(chan struct{})

		// Start the request processing in a goroutine
		go func() {
			defer close(doneChan)
			c.Next()
		}()

		// Wait for either completion or timeout
		select {
		case <-doneChan:
			// Request completed normally
			return
		case <-timeoutChan:
			// Request timed out
			if !c.Writer.Written() {
				// Log timeout error
				logger.LogEndpointTimeout(c, "TimeoutMiddleware", timeout, map[string]interface{}{
					"middleware": "timeout",
				})

				c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{
					"status":     0,
					"statusCode": http.StatusRequestTimeout,
					"message":    "Request timeout",
					"data": gin.H{
						"error":           "Request processing took too long",
						"timeout_seconds": timeout.Seconds(),
					},
				})
			}
			return
		}
	}
}
