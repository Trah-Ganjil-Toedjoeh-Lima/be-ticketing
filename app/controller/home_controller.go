package controller

import (
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HomeController struct {
	config *config.AppConfig
}

func NewHomeController(config *config.AppConfig) *HomeController {
	return &HomeController{config: config}
}

func (u *HomeController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"app_name":  u.config.AppName,
		"app_url":   u.config.AppUrl,
		"time_unix": time.Now().Unix(),
		"time":      time.Now(),
	})
}
