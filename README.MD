# 微服务项目 - Go gRPC 实战

这是一个基于 Go 语言和 gRPC 的微服务项目，包含用户服务、商品服务和订单服务。

## 项目架构

```
microservices/
├── api_gateway/           # API 网关 (8080) - Gin
├── user-service/          # 用户服务 (8081) - gRPC
├── product-service/       # 商品服务 (8082) - gRPC
├── order-service/         # 订单服务 (8083) - gRPC
├── docker-compose.yml     # Docker 编排
└── init.sql              # 数据库初始化
```

## 服务说明

### User Service (端口 8081)

- 用户注册
- 用户登录
- JWT Token 生成

### Product Service (端口 8082)

- 商品增删改查
- 商品列表（分页）
- 库存管理

### Order Service (端口 8083)

- 创建订单
- 查询订单
- 取消订单
- 订单列表

## 技术栈

- **语言**: Go 1.25.5
- **Web 框架**: Gin (API Gateway)
- **RPC 框架**: gRPC
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **认证**: JWT
- **配置**: YAML
- **容器化**: Docker

## 快速开始

### 方式 1: Docker Compose（推荐）

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

### 方式 2: 本地运行

#### 1. 准备数据库

```sql
CREATE DATABASE user_service;
CREATE DATABASE product_service;
CREATE DATABASE order_service;
```

#### 2. 启动服务

```bash
# 启动 API Gateway
cd api_gateway
go run cmd/main.go

# 启动 User Service
cd user-service
go run cmd/main.go

# 启动 Product Service
cd product-service
go run cmd/main.go

# 启动 Order Service
cd order-service
go run cmd/main.go
```

## 测试流程

### 方式 1: 通过 API Gateway (推荐)

```bash
# 1. 注册用户
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 2. 登录获取 Token
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 3. 创建商品
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","description":"最新款","price":5999,"stock":100}'

# 4. 创建订单
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'
```

### 方式 2: 直接调用 gRPC

```bash
# 1. 注册用户
grpcurl -plaintext -d '{"username":"admin","password":"123456"}' localhost:8081 user.UserService/Register

# 2. 登录获取 Token
grpcurl -plaintext -d '{"username":"admin","password":"123456"}' localhost:8081 user.UserService/Login

# 3. 创建商品
grpcurl -plaintext \
  -H "authorization: YOUR_TOKEN" \
  -d '{"name":"iPhone 15","description":"最新款","price":5999,"stock":100}' \
  localhost:8082 product.ProductService/CreateProduct

# 4. 创建订单
grpcurl -plaintext \
  -H "authorization: YOUR_TOKEN" \
  -d '{"product_id":1,"quantity":2}' \
  localhost:8083 order.OrderService/CreateOrder
```

## 项目结构

每个服务都采用相同的分层架构：

```
service/
├── cmd/                    # 程序入口
│   └── main.go
├── config/                 # 配置管理
│   ├── config.go
│   └── config.yaml
├── internal/               # 内部代码
│   ├── handler/           # gRPC 处理器层
│   ├── service/           # 业务逻辑层
│   ├── repository/        # 数据访问层
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件
│   └── pkg/               # 工具包
├── proto/                  # Proto 定义
└── go.mod
```

## 分层职责

### 1. cmd - 程序入口

- 初始化配置
- 初始化数据库
- 注册 gRPC 服务
- 启动服务器

### 2. handler - 处理器层

- 接收 gRPC 请求
- 参数校验
- 调用 service 层

### 3. service - 业务逻辑层

- 实现业务逻辑
- 调用 repository 层
- 跨服务调用

### 4. repository - 数据访问层

- 数据库 CRUD 操作
- SQL 查询
- 使用 GORM

### 5. model - 数据模型

- 定义数据结构
- GORM 标签

### 6. proto - gRPC 定义

- 定义服务接口
- 定义消息结构

## 详细文档

- [项目架构文档](./ARCHITECTURE.md)
- [API Gateway 详细说明](./api_gateway/API_GATEWAY_详细说明.md)
- [User Service 详细说明](./user-service/USER_SERVICE_详细说明.md)
- [Product Service 详细说明](./product-service/PRODUCT_SERVICE_详细说明.md)
- [Order Service 详细说明](./order-service/ORDER_SERVICE_详细说明.md)
- [部署指南](./DEPLOYMENT.md)

## 服务依赖关系

```
Client (HTTP/REST)
    ↓
API Gateway (8080)
    ├─→ User Service (8081)
    ├─→ Product Service (8082)
    └─→ Order Service (8083)
            ├─→ User Service (JWT 认证)
            └─→ Product Service (商品信息)
```

## 端口分配

| 服务            | 端口 |
| --------------- | ---- |
| User Service    | 8081 |
| Product Service | 8082 |
| Order Service   | 8083 |
| API Gateway     | 8080 |
| MySQL           | 3306 |

## 开发指南

### 生成 Proto 文件

```bash
# 安装工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### 添加新服务

1. 复制现有服务目录结构
2. 修改 proto 定义
3. 实现业务逻辑
4. 更新 docker-compose.yml

## 常见问题

### Q: 如何修改数据库密码？

修改各服务的 `config/config.yaml` 文件。

### Q: 如何添加新的 API？

1. 在 proto 文件中定义接口
2. 生成代码
3. 实现 handler 方法
4. 实现 service 逻辑

### Q: 如何调试服务？

使用 grpcurl 或 BloomRPC 测试 gRPC 接口。

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

针对proto: ubuntu
cd /mnt/e/go/微服务/product-service/proto
protoc --go_out=. --go_opt=paths=source_relative \product.proto
protoc --go-grpc_out=. --go-grpc_opt=paths=source_relative \
 product.proto
