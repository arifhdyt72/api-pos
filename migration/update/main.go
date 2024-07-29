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
}
