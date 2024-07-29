package query

import (
	"net/http"
	"test_backend_esb/initializer"

	"github.com/gin-gonic/gin"
)

type TransactionResponse struct {
	ID                uint   `json:"id"`
	TransactionNumber string `json:"transaction_number"`
	TransactionDate   string `json:"transaction_date"`
	Total             int64  `json:"total"`
	TotalQty          int    `json:"total_qty"`
	Status            string `json:"status"`
	OrderMethodName   string `json:"order_method_name"`
	OrderMethodPrice  int64  `json:"order_method_price"`
	PaymentMethodName string `json:"payment_method_name"`
}

type userID struct {
	ID int `uri:"id"`
}

func GetTransaction(c *gin.Context) {
	var input userID
	err := c.ShouldBindUri(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var result []TransactionResponse
	err = initializer.DB.Select(`t.id, t.transaction_number, DATE_FORMAT(t.transaction_date, '%Y-%m-%d %H:%i:%s') AS transaction_date, t.total, t.total_qty, t.status, o.name AS order_method_name, o.price AS order_method_price, p.name AS payment_method_name`).
		Table("transactions t").Joins("INNER JOIN order_methods o ON t.order_method_id = o.id").Joins("INNER JOIN payment_methods p ON t.payment_method_id = p.id").
		Where("t.user_id = ?", input.ID).Order("t.id DESC").Find(&result).Error
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
