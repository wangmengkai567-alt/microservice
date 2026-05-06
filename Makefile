.PHONY: help build run test clean proto docker-build docker-up docker-down

help:
	@echo "可用命令:"
	@echo "  make build         - 构建所有服务"
	@echo "  make run           - 运行所有服务"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理构建文件"
	@echo "  make proto         - 生成 proto 文件"
	@echo "  make docker-build  - 构建 Docker 镜像"
	@echo "  make docker-up     - 启动 Docker Compose"
	@echo "  make docker-down   - 停止 Docker Compose"

build:
	@echo "构建 User Service..."
	cd user-service && go build -o bin/user-service cmd/main.go
	@echo "构建 Product Service..."
	cd product-service && go build -o bin/product-service cmd/main.go
	@echo "构建 Order Service..."
	cd order-service && go build -o bin/order-service cmd/main.go

run:
	@echo "启动所有服务..."
	cd user-service && go run cmd/main.go &
	cd product-service && go run cmd/main.go &
	cd order-service && go run cmd/main.go &

test:
	@echo "运行测试..."
	cd user-service && go test ./...
	cd product-service && go test ./...
	cd order-service && go test ./...

clean:
	@echo "清理构建文件..."
	rm -rf user-service/bin
	rm -rf product-service/bin
	rm -rf order-service/bin

proto:
	@echo "生成 proto 文件..."
	cd user-service && protoc --go_out=. --go-grpc_out=. proto/user.proto
	cd product-service && protoc --go_out=. --go-grpc_out=. proto/product.proto
	cd order-service && protoc --go_out=. --go-grpc_out=. proto/order.proto

docker-build:
	@echo "构建 Docker 镜像..."
	docker-compose build

docker-up:
	@echo "启动 Docker Compose..."
	docker-compose up -d

docker-down:
	@echo "停止 Docker Compose..."
	docker-compose down

docker-logs:
	docker-compose logs -f

init-db:
	@echo "初始化数据库..."
	mysql -u root -p < init.sql
