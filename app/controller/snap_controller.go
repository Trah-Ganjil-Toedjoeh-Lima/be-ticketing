package controller

import (
	"github.com/frchandra/gmcgo/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SnapController struct {
	snapService *service.SnapService
}

func NewSnapController(snapService *service.SnapService) *SnapController {
	return &SnapController{snapService: snapService}
}

func (s *SnapController) HandleCallback(c *gin.Context) {
	message := make(map[string]interface{})
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": message,
	})
	return

}
