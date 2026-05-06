package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"order-service/config"
	"order-service/internal/client"
	"order-service/internal/handler"
	"order-service/internal/logger"
	"order-service/internal/middleware"
	"order-service/internal/model"
	"order-service/internal/processor"
	"order-service/internal/queue"
	"order-service/internal/repository"
	"order-service/internal/scheduler"
	"order-service/internal/service"
	pb "order-service/proto"

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

	logger.Info("Starting order service", zap.String("env", env))

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
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		logger.Error("Failed to migrate database", zap.Error(err))
		log.Fatal(err)
	}
	logger.Info("Database migration completed")

	// 5. 初始化商品服务客户端
	productClient, err := client.NewProductClient(config.Global.Services.ProductService)
	if err != nil {
		logger.Error("Failed to create product client", zap.Error(err))
		log.Fatal(err)
	}
	defer productClient.Close()
	logger.Info("Product service client initialized", zap.String("address", config.Global.Services.ProductService))

	// 6. 初始化 RabbitMQ
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://admin:admin123@localhost:5672/"
	}
	
	mq, err := queue.NewRabbitMQ(rabbitmqURL)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		log.Fatal(err)
	}
	defer mq.Close()
	logger.Info("RabbitMQ connected successfully", zap.String("url", rabbitmqURL))

	// 7. 初始化 Repository
	repo := repository.NewOrderRepository(db)

	// 8. 启动异步订单处理器
	orderProcessor := processor.NewOrderProcessor(repo, mq)
	if err := orderProcessor.Start(); err != nil {
		logger.Error("Failed to start order processor", zap.Error(err))
		log.Fatal(err)
	}
	logger.Info("Order processor started")

	// 9. 依赖注入（传入 MQ）
	orderService := service.NewOrderService(repo, productClient, mq)
	handler := handler.NewOrderHandler(orderService)

	// 10. 启动订单自动取消定时任务（每10秒检查一次）
	orderScheduler := scheduler.NewOrderScheduler(orderService, 10*time.Second)
	orderScheduler.Start()
	defer orderScheduler.Stop()
	logger.Info("Order auto-cancel scheduler started")

	// 11. 创建 gRPC 服务器
	lis, err := net.Listen("tcp", ":8083")
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
	pb.RegisterOrderServiceServer(grpcServer, handler)

	logger.Info("Order service running on :8083")
	fmt.Println("Order service running on 8083...")
	
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve", zap.Error(err))
		log.Fatal(err)
	}
}
