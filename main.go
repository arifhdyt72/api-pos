package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"test_backend_esb/controller/command"
	"test_backend_esb/controller/master"
	"test_backend_esb/controller/middleware"
	"test_backend_esb/controller/query"
	"test_backend_esb/helper"
	"test_backend_esb/initializer"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	for _, k := range helper.TypeRegistry {
		fmt.Println(k)
	}
	// INIT ENV VARIABLE
	initializer.LoadEnv()

	// CONNECT TO DATABASE
	initializer.ConnectDB()

	os.Setenv("TZ", "Asia/Jakarta")

	// INIT CUSTOM LOG
	initializer.InitLogger()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	// CREATE HTTP SERVER USING GIN
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	masterRoutes := r.Group("master")
	{
		masterRoutes.Use(middleware.AuthMiddlware())
		masterRoutes.GET("/*any", master.MasterHandle)
		masterRoutes.POST("/*any", master.MasterHandle)
		masterRoutes.PATCH("/*any", middleware.LogoutMiddleware(), master.MasterHandle)
		masterRoutes.DELETE("/*any", master.MasterHandle)
	}

	r.Static("/images/foto", "images/foto")
	r.Static("/images/icon", "images/icon")

	r.GET("/info", func(c *gin.Context) {
		routeInfo := r.Routes()
		fmt.Println(routeInfo[0].Method)
		paths := make([]string, len(routeInfo))

		for i, k := range routeInfo {
			paths[i] = k.Method + " : " + k.Path
		}

		c.JSON(http.StatusOK, gin.H{
			"routes": paths,
		})
	})

	apiV1 := r.Group("/api/v1")
	{
		apiV1.Use(middleware.AuthMiddlware())
		apiV1.GET("/auth-user", query.GetAuthUser)
		apiV1.POST("/menu", query.GetMenu)
		apiV1.POST("/transaction", command.CreateTransaction)
		apiV1.GET("/transaction/:id", query.GetTransaction)
		apiV1.POST("/auth", command.AuthHandler)
		apiV1.GET("/category_limit/:id", query.GetLimitCategory)
		apiV1.GET("/category/:id", query.GetAllCategory)
		apiV1.GET("/order_method", query.GetOrderMethod)
		apiV1.GET("/payment_method", query.GetPaymentMethod)
	}

	r.Run("0.0.0.0:8082")
}
