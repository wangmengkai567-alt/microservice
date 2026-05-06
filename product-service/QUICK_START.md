# Product Service 快速启动指南

## 前置条件

- Go 1.25.5 或更高版本
- MySQL 5.7 或更高版本
- Protocol Buffers 编译器（可选，如果需要修改 proto 文件）

## 5 分钟快速启动

### 1. 准备数据库

```sql
-- 创建数据库
CREATE DATABASE product_service CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 验证数据库
SHOW DATABASES LIKE 'product_service';
```

### 2. 配置服务

检查 `config/config.yaml` 配置是否正确：

```yaml
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456 # 修改为你的密码
  dbname: product_service
```

### 3. 安装依赖

```bash
cd product-service
go mod download
```

### 4. 生成 Proto 文件（如果需要）

如果 `proto/` 目录下没有 `.pb.go` 文件，需要生成：

```bash
# 安装 protoc 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go-grpc_out=. proto/product.proto
```

### 5. 启动服务

```bash
go run cmd/main.go
```

看到以下输出表示启动成功：

```
Product service running on 8082...
```

## 验证服务

### 方法 1：使用 grpcurl

```bash
# 查看服务列表
grpcurl -plaintext localhost:8082 list

# 测试商品列表接口
grpcurl -plaintext -d '{"page":1,"page_size":10}' localhost:8082 product.ProductService/ListProducts
```

### 方法 2：使用 Go 客户端

创建 `test_client.go`：

```go
package main

import (
    "context"
    "log"
    pb "product-service/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    conn, err := grpc.Dial("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewProductServiceClient(conn)

    // 测试商品列表
    resp, err := client.ListProducts(context.Background(), &pb.ListProductsRequest{
        Page:     1,
        PageSize: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Total products: %d", resp.Total)
}
```

运行：

```bash
go run test_client.go
```

## 完整测试流程

### 1. 启动 user-service（如果需要认证功能）

```bash
cd ../user-service
go run cmd/main.go
```

### 2. 注册并登录获取 Token

```bash
# 注册用户
grpcurl -plaintext -d '{"username":"admin","password":"123456"}' localhost:8081 user.UserService/Register

# 登录获取 Token
grpcurl -plaintext -d '{"username":"admin","password":"123456"}' localhost:8081 user.UserService/Login
```

保存返回的 Token。

### 3. 创建商品

```bash
grpcurl -plaintext \
  -H "authorization: YOUR_TOKEN_HERE" \
  -d '{
    "name": "测试商品",
    "description": "这是一个测试商品",
    "price": 99.99,
    "stock": 100
  }' localhost:8082 product.ProductService/CreateProduct
```

### 4. 查询商品

```bash
# 获取商品列表
grpcurl -plaintext -d '{"page":1,"page_size":10}' localhost:8082 product.ProductService/ListProducts

# 获取商品详情（假设 ID 为 1）
grpcurl -plaintext -d '{"id":1}' localhost:8082 product.ProductService/GetProduct
```

## 常见问题

### Q1: 启动时报错 "Error 1049: Unknown database"

**解决**：确保已创建数据库

```sql
CREATE DATABASE product_service;
```

### Q2: 启动时报错 "dial tcp 127.0.0.1:3306: connect: connection refused"

**解决**：确保 MySQL 服务已启动

```bash
# Windows
net start mysql

# Mac
brew services start mysql

# Linux
sudo systemctl start mysql
```

### Q3: 报错 "protoc: command not found"

**解决**：安装 Protocol Buffers 编译器

- Windows: 从 [GitHub Releases](https://github.com/protocolbuffers/protobuf/releases) 下载
- Mac: `brew install protobuf`
- Linux: `apt-get install protobuf-compiler`

### Q4: 创建商品时报错 "missing token"

**解决**：确保在 metadata 中添加了 authorization 字段

```bash
-H "authorization: YOUR_TOKEN"
```

### Q5: 端口 8082 已被占用

**解决**：修改 `config/config.yaml` 中的端口号

```yaml
server:
  port: 8083 # 改为其他端口
```

同时修改 `cmd/main.go` 中的监听端口：

```go
lis, err := net.Listen("tcp", ":8083")
```

## 数据库表结构

服务启动后会自动创建 `products` 表：

```sql
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `description` text,
  `price` double NOT NULL,
  `stock` int NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

## 下一步

- 阅读 [PRODUCT*SERVICE*详细说明.md](./PRODUCT_SERVICE_详细说明.md) 了解架构设计
- 阅读 [API_TEST_GUIDE.md](./test/API_TEST_GUIDE.md) 学习 API 测试
- 阅读 [SERVICE_COMPARISON.md](./SERVICE_COMPARISON.md) 了解与 user-service 的差异

## 开发模式

### 热重载

使用 `air` 实现热重载：

```bash
# 安装 air
go install github.com/cosmtrek/air@latest

# 创建 .air.toml 配置文件
air init

# 启动热重载
air
```

### 调试模式

使用 Delve 调试器：

```bash
# 安装 delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug cmd/main.go
```

## 生产部署

### Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o product-service cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/product-service .
COPY config/config.yaml config/
EXPOSE 8082
CMD ["./product-service"]
```

构建并运行：

```bash
docker build -t product-service .
docker run -p 8082:8082 product-service
```

### 使用 Docker Compose

参考项目根目录的 `docker-compose.yml` 文件。

## 性能优化建议

1. **数据库连接池**：配置合适的连接池参数
2. **索引优化**：为常用查询字段添加索引
3. **缓存**：使用 Redis 缓存热门商品
4. **分页优化**：使用游标分页替代 offset
5. **并发控制**：使用连接池限制并发数

## 监控和日志

建议添加：

- 结构化日志（zap、logrus）
- 指标收集（Prometheus）
- 链路追踪（Jaeger、Zipkin）
- 健康检查接口

## 技术支持

如有问题，请查看：

- [详细说明文档](./PRODUCT_SERVICE_详细说明.md)
- [API 测试指南](./test/API_TEST_GUIDE.md)
- [服务对比说明](./SERVICE_COMPARISON.md)
