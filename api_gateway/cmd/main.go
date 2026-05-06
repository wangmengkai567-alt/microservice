package main

import (
	"fmt"
	"log"
	"os"

	"api-gateway/config"
	"api-gateway/internal/client"
	"api-gateway/internal/logger"
	"api-gateway/internal/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 初始化日志
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	if err := logger.Init(env); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting API Gateway", zap.String("env", env))

	// 2. 初始化配置
	if err := config.Init(); err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		log.Fatal("Failed to load config:", err)
	}

	// 3. 设置 Gin 模式
	gin.SetMode(config.Global.Server.Mode)

	// 4. 初始化 gRPC 客户端
	if err := client.InitClients(); err != nil {
		logger.Error("Failed to init gRPC clients", zap.Error(err))
		log.Fatal("Failed to init gRPC clients:", err)
	}
	defer client.Clients.Close()
	logger.Info("gRPC clients initialized successfully")

	// 5. 设置路由
	r := router.SetupRouter()

	// 6. 启动服务器
	addr := fmt.Sprintf(":%d", config.Global.Server.Port)
	logger.Info("API Gateway running", zap.String("address", addr))
	fmt.Printf("API Gateway running on %s...\n", addr)
	
	if err := r.Run(addr); err != nil {
		logger.Error("Failed to start server", zap.Error(err))
		log.Fatal("Failed to start server:", err)
	}
}
