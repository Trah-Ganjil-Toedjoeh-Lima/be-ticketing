package controller

import (
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConfigController struct {
	config *config.AppConfig
	log    *util.LogUtil
}

func NewConfigController(config *config.AppConfig, log *util.LogUtil) *ConfigController {
	return &ConfigController{config: config, log: log}
}

func (g *ConfigController) GetAppConfig(c *gin.Context) {
	filteredConfig := make(map[string]string, 20)
	filteredConfig["AppName"] = g.config.AppName
	filteredConfig["IsProduction"] = strconv.FormatBool(g.config.IsProduction)
	filteredConfig["MidtransIsProduction"] = strconv.FormatBool(g.config.MidtransIsProduction)
	filteredConfig["IsOpenGate"] = strconv.FormatBool(g.config.IsOpenGate)
	filteredConfig["IsOpenAuth"] = strconv.FormatBool(g.config.IsOpenAuth)
	filteredConfig["QrScanBehaviour"] = g.config.QrScanBehaviour
	filteredConfig["AppUrl"] = g.config.AppUrl
	filteredConfig["AppPort"] = g.config.AppPort
	filteredConfig["AccessMinute"] = g.config.AccessMinute.String()
	filteredConfig["RefreshMinute"] = g.config.RefreshMinute.String()
	filteredConfig["TransactionMinute"] = g.config.TransactionMinute.String()
	filteredConfig["TotpPeriod"] = strconv.FormatUint(uint64(g.config.TotpPeriod), 10)

	c.JSON(http.StatusOK, gin.H{
		"app_config": filteredConfig,
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

func (g *ConfigController) OpenAuth(c *gin.Context) {
	g.config.IsOpenAuth = true
	c.Status(http.StatusOK)
	return
}

func (g *ConfigController) CloseAuth(c *gin.Context) {
	g.config.IsOpenAuth = false
	c.Status(http.StatusOK)
	return
}

func (g *ConfigController) SetToProduction(c *gin.Context) {
	g.config.MidtransIsProduction = true
	g.config.IsProduction = true
	c.Status(http.StatusOK)
	return
}

func (g *ConfigController) SetToSandbox(c *gin.Context) {
	g.config.MidtransIsProduction = false
	g.config.IsProduction = false
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
	g.config.QrScanBehaviour = inputData["qr_scan_behaviour"]
	c.Status(http.StatusOK)
	return
}
