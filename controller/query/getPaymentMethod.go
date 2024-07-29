package query

import (
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"github.com/gin-gonic/gin"
)

func GetPaymentMethod(c *gin.Context) {
	var result []model.OrderMethod
	err := initializer.DB.Table("payment_methods o").Where("deleted_at IS NULL").
		Order("ID DESC").Find(&result).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": result,
	})
}
