package query

import (
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Limit   int
	Page    int
	Payload interface{}
}

type storeID struct {
	ID int `uri:"id"`
}

type MenuRequest struct {
	Name string `json:"name"`
}

func GetLimitCategory(c *gin.Context) {
	var input storeID
	err := c.ShouldBindUri(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var category []model.Category

	err = initializer.DB.Table("categories c").Where("deleted_at IS NULL").Where("store_id", input.ID).
		Limit(8).Order("ID DESC").Find(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

func GetAllCategory(c *gin.Context) {
	var input storeID
	err := c.ShouldBindUri(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var category []model.Category

	err = initializer.DB.Table("categories c").Where("deleted_at IS NULL").
		Where("store_id", input.ID).Order("ID DESC").Find(&category).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": category,
	})
}

func GetMenu(c *gin.Context) {
	var req MenuRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var item []model.Item

	err = initializer.DB.Table("items i").Where("deleted_at IS NULL").
		Where("name LIKE ?", "%"+req.Name+"%").Order("ID DESC").Find(&item).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": item,
	})
}
