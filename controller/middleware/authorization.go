package middleware

import (
	"net/http"
	"strings"
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"github.com/gin-gonic/gin"
)

func AuthMiddlware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path != "/api/v1/auth" && !strings.Contains(c.Request.URL.Path, "/images/") {
			header := c.Request.Header.Get("Authorization")
			if header == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unathorized",
				})
				return
			}

			authorization := strings.Split(header, " ")
			if len(authorization) < 2 || authorization[1] == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unathorized",
				})
				return
			}

			var user model.User
			if authorization[0] == "Bearer" {
				rs := initializer.DB.Model(&model.User{}).Where("token = ?", authorization[1]).First(&user)
				if rs.Error != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "unathorized",
					})
					return
				}

				// if user.AppVersion == "" && user.RoleID == 6 {
				// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				// 		"error": "update applikasi ke versi terbaru",
				// 	})
				// 	return
				// }

				// if user.AppVersion != os.Getenv("APP_VERSION") && user.RoleID == 4 && c.Request.URL.Path != "/api/v1/lkm" && os.Getenv("LOGIN_BYPASS") != "BYPASS" {
				// 	c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{
				// 		"error": "outdated application version",
				// 	})
				// 	return
				// }
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unathorized",
				})
				return
			}

		}
		c.Next()
	}
}
