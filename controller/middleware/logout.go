package middleware

import (
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"
	"test_backend_esb/tools"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func LogoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/master/user" && c.Request.Method == "PATCH" {
			var req map[string]interface{}
			err := c.ShouldBindBodyWith(&req, binding.JSON)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "invalid params",
				})
				return
			}
			var user model.User
			if err := initializer.DB.First(&user, req["ID"].(float64)).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "invalid user id",
				})
				return
			}
			user.Token = tools.GenerateToken(64)
			c.Set("binder", req)
			initializer.DB.Save(&user)
		}
		c.Next()
	}
}
