package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func OptionsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			// Handle OPTIONS request
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE, PUT")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Max-Age", "600")

			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			// Continue to next middleware or handler
			c.Next()
		}
	}
}
