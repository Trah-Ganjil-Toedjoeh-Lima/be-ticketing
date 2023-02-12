package util

import "github.com/gin-gonic/gin"

func GinResponseError(c *gin.Context, code int, message string, error string) {
	c.JSON(code, gin.H{
		"message": message,
		"error":   error,
	})
}
