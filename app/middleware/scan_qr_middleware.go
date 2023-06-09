package middleware

import (
	"github.com/frchandra/ticketing-gmcgo/app/controller"
	"net/http"

	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ScanQrMiddleware struct {
	tokenUtil      *util.TokenUtil
	log            *logrus.Logger
	config         *config.AppConfig
	userService    *service.UserService
	seatController *controller.SeatController
}

func NewScanQrMiddleware(tokenUtil *util.TokenUtil, log *logrus.Logger, config *config.AppConfig, userService *service.UserService, seatController *controller.SeatController) *ScanQrMiddleware {
	return &ScanQrMiddleware{tokenUtil: tokenUtil, log: log, config: config, userService: userService, seatController: seatController}
}

func (m *ScanQrMiddleware) HandleScanQr(c *gin.Context) {
	accessDetails, err := m.tokenUtil.GetValidatedAccess(c) //get the user data from the token in the request header
	if accessDetails == nil || err != nil {                 //redirect as user
		c.Next()
		return
	}

	err = m.tokenUtil.FetchAuthn(accessDetails.AccessUuid) //check if token exist in the token storage (Check if the token is expired)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "fail",
			"error":   "your credentials are invalid. try to refresh your credentials",
		})
		m.log.
			WithField("occurrence", "AdminMiddelware@HandleAdminAcccess").
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			WithField("source_messages", err.Error()).
			Info("cannot find access token in the storage")
		c.Abort()
		return
	}

	adminUser, err := m.userService.GetById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": err.Error()})
		return
	}
	if adminUser.Name == m.config.AdminName && adminUser.Email == m.config.AdminEmail && adminUser.Phone == m.config.AdminPhone { //check if this user is admin
		if m.config.QrScanBehaviour == "default" {
			m.seatController.DetailsByLink(c) //redirect to admin controller
		} else {
			m.seatController.UpdateToStatus(c, m.config.QrScanBehaviour) //redirect to admin controller
		}
		c.Abort()
		return
	}
	//redirect as user
	c.Next()
	return

}
