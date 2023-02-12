package controller

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GateController struct {
	config *config.AppConfig
}

func NewGateController(config *config.AppConfig) *GateController {
	return &GateController{config: config}
}

func (g *GateController) GetAppConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_config": g.config,
	})
	return
}

func (g *GateController) OpenGate(c *gin.Context) {
	g.config.IsOpenGate = true
	c.Status(http.StatusOK)
	return
}

func (g *GateController) CloseGate(c *gin.Context) {
	g.config.IsOpenGate = false
	c.Status(http.StatusOK)
	return
}
