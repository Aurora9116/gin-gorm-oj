package middlewares

import (
	"gin-gorm-oj/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthAdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		//
		auth := c.GetHeader("Authorization")
		if strings.Contains(auth, "bearer") {
			auth = strings.ReplaceAll(auth, "bearer ", "")
		}
		userClaim, err := helper.AnalyseToken(auth)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Unauthorized Authorization",
			})
			c.Abort()
			return
		}
		if userClaim == nil || userClaim.IsAdmin != 1 {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "Unauthorized Admin",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
