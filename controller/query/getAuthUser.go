package query

import (
	"net/http"
	"test_backend_esb/helper"

	"github.com/gin-gonic/gin"
)

func GetAuthUser(c *gin.Context) {
	user := helper.AuthUser(c)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Ok!",
		"data": user,
	})
}
