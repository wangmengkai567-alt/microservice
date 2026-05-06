package router

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 中间件
	r.Use(middleware.LoggerMiddleware()) // 添加日志中间件
	r.Use(middleware.CORS())
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 版本
	v1 := r.Group("/api/v1")

	// User Service 路由
	userHandler := handler.NewUserHandler()
	userGroup := v1.Group("/users")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}

	// Product Service 路由
	productHandler := handler.NewProductHandler()
	productGroup := v1.Group("/products")
	{
		// 公开接口
		productGroup.GET("", productHandler.ListProducts)
		productGroup.GET("/:id", productHandler.GetProduct)

		// 需要认证的接口
		productAuth := productGroup.Group("")
		productAuth.Use(middleware.AuthMiddleware())
		{
			productAuth.POST("/upload-image", productHandler.UploadImage)
			productAuth.POST("", productHandler.CreateProduct)
			productAuth.PUT("/:id", productHandler.UpdateProduct)
			productAuth.DELETE("/:id", productHandler.DeleteProduct)
		}
	}

	// Order Service 路由（所有接口需要认证）
	orderHandler := handler.NewOrderHandler()
	orderGroup := v1.Group("/orders")
	orderGroup.Use(middleware.AuthMiddleware())
	{
		orderGroup.POST("", orderHandler.CreateOrder)
		orderGroup.POST("/:id/cancel", orderHandler.CancelOrder)
		orderGroup.POST("/:id/pay", orderHandler.PayOrder)
		orderGroup.DELETE("/:id", orderHandler.DeleteOrder)
		orderGroup.GET("/:id", orderHandler.GetOrder)
		orderGroup.GET("", orderHandler.ListOrders)
	}

	return r
}
