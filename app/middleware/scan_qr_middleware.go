package middleware

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ScanQrMiddleware struct {
	tokenUtil   *util.TokenUtil
	log         *logrus.Logger
	config      *config.AppConfig
	userService *service.UserService
}

func NewScanQrMiddleware(tokenUtil *util.TokenUtil, log *logrus.Logger, config *config.AppConfig, userService *service.UserService) *ScanQrMiddleware {
	return &ScanQrMiddleware{tokenUtil: tokenUtil, log: log, config: config, userService: userService}
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

	adminUser, _ := m.userService.GetById(accessDetails.UserId)
	if adminUser.Name == m.config.AdminName && adminUser.Email == m.config.AdminEmail && adminUser.Phone == m.config.AdminPhone { //check if this user is admin
		//redirect as admin
		if m.config.QrScanBehaviour != "default" {
			c.Redirect(http.StatusFound, "/api/v1/admin/seat/"+c.Param("link")+"/"+m.config.QrScanBehaviour)
		} else {
			c.Redirect(http.StatusFound, "/api/v1/admin/seat/"+c.Param("link"))
		}
		return
	}
	//redirect as user
	c.Next()
	return

}
