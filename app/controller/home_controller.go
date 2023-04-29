package controller

import (
	"net/http"
	"time"

	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
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

func (u *HomeController) checkAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":   "success",
		"time_unix": time.Now().Unix(),
		"time":      time.Now(),
	})
}
