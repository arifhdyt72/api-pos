package model

import (
	"time"
)

type Transaction struct {
	MavisModel
	TransactionNumber string        `json:"transaction_number" gorm:"index"`
	TransactionDate   time.Time     `json:"transaction_date"`
	Total             int64         `json:"total"`
	TotalQty          int           `json:"total_qty"`
	PaymentMethodID   uint          `json:"payment_method_id"`
	PaymentMethod     PaymentMethod `json:"payment_method"`
	OrderMethodID     uint          `json:"order_method_id"`
	OrderMethod       OrderMethod   `json:"order_method"`
	Status            string        `json:"status"`
	StoreID           uint          `json:"store_id"`
	Store             *Store
	UserID            uint `json:"user_id"`
	User              *User
	TransactionDetail []TransactionDetail `json:"transaction_detail" binding:"dive"`
}

type TransactionDetail struct {
	MavisModel
	ItemID        uint `json:"item_id"`
	Item          *Item
	TransactionID uint `json:"transaction_id"`
	Transaction   *Transaction
	Qty           int `json:"qty"`
}
