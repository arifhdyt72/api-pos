package helper

import (
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func AuthUser(c *gin.Context) *model.User {
	bearer := c.Request.Header.Get("Authorization")
	token := strings.Split(bearer, " ")
	if len(token) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request token",
		})
		return nil
	}

	var user model.User
	initializer.DB.Where("token = ?", token[1]).Preload(clause.Associations).First(&user)
	return &user
}
