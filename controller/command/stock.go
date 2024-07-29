package command

import (
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"

	"github.com/gin-gonic/gin"
)

type ListItem struct {
	Item []ItemStock `json:"item" binding:"required"`
}

type ItemStock struct {
	ItemId uint `json:"item_id" binding:"required"`
	Qty    int  `json:"qty" binding:"required"`
}

type StockAdjusmentInput struct {
	ItemId      uint   `json:"item_id"`
	Qty         int    `json:"qty"`
	Description string `json:"desc"`
}

func CheckStock(c *gin.Context) {
	var input ListItem
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	listItemNoStock := []map[string]interface{}{}
	for _, data := range input.Item {
		var item model.Item
		err := initializer.DB.Where("id = ?", data.ItemId).First(&item).Error
		if err != nil {
			message := map[string]interface{}{
				"item_id": data.ItemId,
				"qty":     item.Stock,
			}
			listItemNoStock = append(listItemNoStock, message)
		}
	}

	if len(listItemNoStock) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"item": listItemNoStock,
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Stock Available",
		})
		return
	}
}
