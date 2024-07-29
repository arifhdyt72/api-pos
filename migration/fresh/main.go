package main

import (
	"os"
	"test_backend_esb/initializer"
	"test_backend_esb/model"
)

func init() {
	// INIT ENV VARIABLE
	initializer.LoadEnv()

	// CONNECT TO DATABASE
	initializer.ConnectDB()

	// SET TIMEZONE
	os.Setenv("TZ", "Asia/Jakarta")
}

func main() {
	initializer.DB.AutoMigrate(
		&model.Store{},
		&model.Item{},
		&model.OrderMethod{},
		&model.User{},
		&model.Category{},
		&model.Transaction{},
		&model.TransactionDetail{},
		&model.PaymentMethod{},
	)

	truest := true
	initializer.DB.Create(&model.User{
		Name:     "Admin",
		Username: "superadmin",
		Password: "$2y$10$su9PyzRwtYt7liFFmx/uVurBHnWUPGFu91MWYNkHj6Mw.ek9FMswi",
		Email:    "arif.hidayat@testmail.com",
		StoreID:  nil,
		Status:   &truest,
	})

	initializer.DB.Create(&model.Store{
		Name:    "Store Model",
		Address: "",
	})
}
