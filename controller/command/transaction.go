package command

import (
	"fmt"
	"net/http"
	"test_backend_esb/initializer"
	"test_backend_esb/model"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionInput struct {
	PaymentMethodID   uint                     `json:"payment_method_id"`
	OrderMethodID     uint                     `json:"order_method_id"`
	StoreID           uint                     `json:"store_id"`
	UserID            uint                     `json:"user_id"`
	TotalQty          int                      `json:"total_qty"`
	Total             int64                    `json:"total"`
	TransactionDetail []TransactionDetailInput `json:"transaction_detail"`
}

type TransactionDetailInput struct {
	ItemID uint `json:"item_id"`
	Qty    int  `json:"qty"`
}

func CreateTransaction(c *gin.Context) {
	var input TransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	currentTime := time.Now()
	layout := "02-01-2006"
	formattedTime := currentTime.Format(layout)

	var count int64
	initializer.DB.Model(&model.Transaction{}).Count(&count)

	trfNumber := fmt.Sprintf("#TRF-%s-%d", formattedTime, count+1)

	var transaction model.Transaction
	transaction.TransactionNumber = trfNumber
	transaction.TransactionDate = time.Now()
	transaction.TotalQty = input.TotalQty
	transaction.Total = input.Total
	transaction.PaymentMethodID = input.PaymentMethodID
	transaction.OrderMethodID = input.OrderMethodID
	transaction.Status = "Paid"
	transaction.StoreID = input.StoreID
	transaction.UserID = input.UserID

	err := initializer.DB.Create(&transaction).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, detail := range input.TransactionDetail {
		var detailTransaction model.TransactionDetail
		detailTransaction.TransactionID = transaction.ID
		detailTransaction.ItemID = detail.ItemID
		detailTransaction.Qty = detail.Qty

		initializer.DB.Create(&detailTransaction)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Transaction has been created",
	})
}
