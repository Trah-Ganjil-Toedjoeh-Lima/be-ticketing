package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ConfigController struct {
	config *config.AppConfig
	log    *util.LogUtil
}

func NewConfigController(config *config.AppConfig, log *util.LogUtil) *ConfigController {
	return &ConfigController{config: config, log: log}
}

func (g *ConfigController) GetAppConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"app_config": g.config,
	})
	return
}

func (g *ConfigController) OpenGate(c *gin.Context) {
	g.config.IsOpenGate = true
	c.Status(http.StatusOK)
	return
}

func (g *ConfigController) CloseGate(c *gin.Context) {
	g.config.IsOpenGate = false
	c.Status(http.StatusOK)
	return
}

func (g *ConfigController) UpdateQrScanBehaviour(c *gin.Context) {
	var inputData map[string]string //get the data in request body
	if err := c.ShouldBindJSON(&inputData); err != nil {
		g.log.BasicLog(err, "ConfigController@UpdateQrScanBehaviour")
		util.GinResponseError(c, http.StatusBadRequest, "error when processing the request data", err.Error())
		return
	}
	if inputData["qr_scan_behaviour"] == "open_gate" {
		g.config.QrScanBehaviour = inputData["qr_scan_behaviour"]
		c.Status(http.StatusOK)
		return
	}
	if inputData["qr_scan_behaviour"] == "ticket_exchanging" {
		g.config.QrScanBehaviour = inputData["qr_scan_behaviour"]
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusBadRequest)
	return

}
