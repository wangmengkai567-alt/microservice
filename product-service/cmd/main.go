package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"product-service/config"
	"product-service/internal/handler"
	"product-service/internal/logger"
	"product-service/internal/middleware"
	"product-service/internal/model"
	"product-service/internal/repository"
	"product-service/internal/service"
	pb "product-service/proto"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	logger.Info("Starting product service", zap.String("env", env))

	// 2. 初始化配置
	if err := config.Init(); err != nil {
		logger.Error("Failed to initialize config", zap.Error(err))
		log.Fatal(err)
	}

	// 3. 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Global.Database.User,
		config.Global.Database.Password,
		config.Global.Database.Host,
		config.Global.Database.Port,
		config.Global.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		log.Fatal(err)
	}
	logger.Info("Database connected successfully")

	// 4. 自动迁移
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		logger.Error("Failed to migrate database", zap.Error(err))
		log.Fatal(err)
	}
	logger.Info("Database migration completed")

	// 5. 依赖注入
	repo := repository.NewProductRepository(db)
	service := service.NewProductService(repo)
	handler := handler.NewProductHandler(service)

	// 6. 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		logger.Error("Failed to listen", zap.Error(err))
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.TraceInterceptor(),
			middleware.AuthInterceptor(),
		),
	)
	pb.RegisterProductServiceServer(grpcServer, handler)

	logger.Info("Product service running on :8082")
	fmt.Println("Product service running on 8082...")
	
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve", zap.Error(err))
		log.Fatal(err)
	}
}
