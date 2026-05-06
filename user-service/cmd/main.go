//程序入口
package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"user-service/config"
	"user-service/internal/handler"
	"user-service/internal/logger"
	"user-service/internal/middleware"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
	pb "user-service/proto"

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

	logger.Info("Starting user service", zap.String("env", env))

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
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Error("Failed to migrate database", zap.Error(err))
		log.Fatal(err)
	}
	logger.Info("Database migration completed")

	// 5. 依赖注入
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	handler := handler.NewUserHandler(service)

	// 6. 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":8081")
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
	pb.RegisterUserServiceServer(grpcServer, handler)

	logger.Info("User service running on :8081")
	fmt.Println("User service running on 8081...")
	
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve", zap.Error(err))
		log.Fatal(err)
	}
}