package middleware

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GateMiddleware struct {
	config *config.AppConfig
}

func NewGateMiddleware(config *config.AppConfig) *GateMiddleware {
	return &GateMiddleware{config: config}
}

func (g *GateMiddleware) HandleTransactionAccess(c *gin.Context) {
	if g.config.IsOpenGate == true {
		c.Next()
		return
	} else {
		c.JSON(http.StatusTooEarly, gin.H{
			"error": "the transaction gate has not been opened",
		})
		c.Abort()
	}
}

func (g GateMiddleware) HandleAuthAccess(c *gin.Context) {
	if g.config.IsOpenAuth == true {
		c.Next()
		return
	} else {
		c.JSON(http.StatusTooEarly, gin.H{
			"error": "the authentication gate has not been opened",
		})
		c.Abort()
	}
}
