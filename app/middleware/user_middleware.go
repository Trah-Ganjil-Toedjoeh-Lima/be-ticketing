package middleware

import (
	"github.com/frchandra/gmcgo/app/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserMiddleware struct {
	tokenUtil *util.TokenUtil
}

func NewUserMiddleware(tokenUtil *util.TokenUtil) *UserMiddleware {
	return &UserMiddleware{
		tokenUtil: tokenUtil,
	}
}

func (this *UserMiddleware) HandleUserAccess(c *gin.Context) {

	//get the user data from the token in the request header
	accessDetails, err := this.tokenUtil.GetValidatedAccess(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	//check if token exist in the token storage (Check if the token is expired)
	err = this.tokenUtil.FetchAuthn(accessDetails.AccessUuid)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":      "fail",
			"error":       err.Error(),
			"description": "cannot found access token. Try to refresh then token",
		})
		c.Abort()
		return
	}
	c.Set("accessDetails", accessDetails)
	c.Next()
}
