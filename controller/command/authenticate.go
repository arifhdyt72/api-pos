package command

import (
	"errors"
	"fmt"
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"
	"test_backend_esb/tools"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authRequest struct {
	Username string `form:"username" json:"username" binding:"required" gorm:"unique"`
	Password string `form:"password" json:"password" binding:"required"`
}

func AuthHandler(c *gin.Context) {
	var req authRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "parameter tidak sesuai",
		})
		return
	}

	// SEARCH USER BY USERNAME & PASSWORD
	var user model.User
	user.Username = req.Username
	res := initializer.DB.Where("status = 1").First(&user, user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username tidak ditemukan",
		})
		return
	}

	// VALIDATE PASSWORD
	fmt.Println(user.HashedPassword)
	fmt.Println(req.Password)
	result := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if result != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password salah",
		})
		return
	}

	//RETURN TOKEN TO USER
	token := tools.GenerateToken(64)
	user.Token = token
	initializer.DB.Save(&user)

	var store model.Store
	initializer.DB.Where("id = ?", user.StoreID).First(&store)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"store": store,
		"user":  user,
	})
}
