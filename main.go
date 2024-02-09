package main

import (
	"fmt"
	"go_pos_v1_2/config"
	"go_pos_v1_2/helper"
	"go_pos_v1_2/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello pos")
	config.InitDB()
	InitEnv()

	r := gin.Default()
	r.POST("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "PING",
		})
	})

	route := r.Group("/")
	{
		user := route.Group("/user")
		{
			user.POST("/register", routes.RegisterUser)
			user.POST("/login", routes.GenerateToken)
			user.GET("/", helper.ValidateJWT(), routes.GetUser)
			user.GET("/:id", helper.ValidateJWT(), routes.GetUserById)
			user.DELETE("/delete/:id", helper.ValidateJWT(), routes.DeleteUser)
			user.PUT("/update/:id", helper.ValidateJWT(), routes.UpdateUser)
		}

		role := route.Group("/role").Use(helper.ValidateJWT())
		{
			role.POST("/insert", routes.InsertRole)
			role.GET("/", routes.GetRole)
			role.GET("/:id", routes.GetRoleById)
			role.DELETE("/delete/:id", routes.DeleteRole)
			role.PUT("/update/:id", routes.UpdateRole)
		}

		product := route.Group("/product").Use(helper.ValidateJWT())
		{
			product.POST("/insert", routes.InsertProduct)
			product.GET("/", routes.GetProduct)
			product.GET("/:id", routes.GetProductById)
			product.DELETE("/delete/:id", routes.DeleteProduct)
			product.PUT("/update/:id", routes.UpdateProduct)
		}

		member := route.Group("/member").Use(helper.ValidateJWT())
		{
			member.POST("/insert", routes.InsertMember)
			member.GET("/", routes.GetMember)
			member.GET("/:id", routes.GetMemberById)
			member.DELETE("/delete/:id", routes.DeleteMember)
			member.PUT("/update/:id", routes.UpdateMember)
		}

		inventory := route.Group("/inventory").Use(helper.ValidateJWT())
		{
			inventory.POST("/insert", routes.InsertInventory)
			inventory.GET("/", routes.GetInventory)
			inventory.GET("/:id", routes.GetInventoryById)
			inventory.DELETE("/delete/:id", routes.DeleteInventory)
			inventory.PUT("/update/:id", routes.UpdateInventory)
		}

		distributor := route.Group("/distributor").Use(helper.ValidateJWT())
		{
			distributor.POST("/insert", routes.InsertDistributor)
			distributor.GET("/", routes.GetDistributor)
			distributor.GET("/:id", routes.GetDistributorById)
			distributor.DELETE("/delete/:id", routes.DeleteDistributor)
			distributor.PUT("/update/:id", routes.UpdateDistributor)
		}

		chasier := route.Group("/chasier").Use(helper.ValidateJWT())
		{
			chasier.POST("/insert", routes.InsertChasier)
			chasier.GET("/", routes.GetChasier)
			chasier.GET("/:id", routes.GetChasierById)
			chasier.DELETE("/delete/:id", routes.DeleteChasier)
			chasier.PUT("/update/:id", routes.UpdateChasier)
		}

		purchase := route.Group("/purchase").Use(helper.ValidateJWT())
		{
			purchase.POST("/insert", routes.InsertPurchase)
			purchase.GET("/", routes.GetPurchase)
			purchase.GET("/:id", routes.GetPurchaseById)
			purchase.DELETE("/delete/:id", routes.DeletePurchase)
			purchase.PUT("/update/:id", routes.UpdatePurchase)
		}

		sales := route.Group("/sales").Use(helper.ValidateJWT())
		{
			sales.POST("/insert", routes.InsertSales)
			sales.GET("/", routes.GetSales)
			sales.GET("/:id", routes.GetSalesById)
			sales.DELETE("/delete/:id", routes.DeleteSales)
			sales.PUT("/update/:id", routes.UpdateSales)
		}
	}

	r.Run(":8011")
}

func InitEnv() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
