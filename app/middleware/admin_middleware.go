package middleware

import (
	"github.com/frchandra/ticketing-gmcgo/app/service"
	"github.com/frchandra/ticketing-gmcgo/app/util"
	"github.com/frchandra/ticketing-gmcgo/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type AdminMiddleware struct {
	tokenUtil   *util.TokenUtil
	log         *logrus.Logger
	config      *config.AppConfig
	userService *service.UserService
}

func NewAdminMiddleware(tokenUtil *util.TokenUtil, log *logrus.Logger, config *config.AppConfig, userService *service.UserService) *AdminMiddleware {
	return &AdminMiddleware{tokenUtil: tokenUtil, log: log, config: config, userService: userService}
}

func (u *AdminMiddleware) AdminAccess(c *gin.Context) {
	accessDetails, err := u.tokenUtil.GetValidatedAccess(c) //get the user data from the token in the request header
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "fail",
			"error":   "your credentials are invalid",
		})
		u.log.WithField("occurrence", "AdminMiddelware@HandleAdminAcccess").
			WithField("source_messages", err.Error()).
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			Info("cannot find token in the http request")
		c.Abort()
		return
	}

	err = u.tokenUtil.FetchAuthn(accessDetails.AccessUuid) //check if token exist in the token storage (Check if the token is expired)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "fail",
			"error":   "your credentials are invalid. try to refresh your credentials",
		})
		u.log.
			WithField("occurrence", "AdminMiddelware@HandleAdminAcccess").
			WithField("client_ip", c.ClientIP()).
			WithField("endpoint", c.FullPath()).
			WithField("source_messages", err.Error()).
			Info("cannot find access token in the storage")
		c.Abort()
		return
	}

	adminUser, err := u.userService.GetById(accessDetails.UserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "error", "error": err.Error()})
		c.Abort()
		return
	}

	if adminUser.Name == u.config.AdminName && adminUser.Email == u.config.AdminEmail && adminUser.Phone == u.config.AdminPhone { //check if this user is admin
		c.Set("accessDetails", accessDetails)
		c.Next()
		return
	}
	c.Abort()
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "fail",
		"error":   "you are not authorized",
	})
	return

}
