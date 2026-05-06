# API Gateway

基于 Gin 框架的 API 网关，将 gRPC 微服务转换为 RESTful API。

## 功能特性

- RESTful API 接口
- JWT 认证中间件
- CORS 跨域支持
- 请求参数验证
- gRPC 客户端连接池
- 统一错误处理

## 技术栈

- Gin Web Framework
- gRPC Client
- JWT 认证
- YAML 配置

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置服务

修改 `config/config.yaml`：

```yaml
server:
  port: 8080
  mode: debug

services:
  user_service: "localhost:8081"
  product_service: "localhost:8082"
  order_service: "localhost:8083"
```

### 3. 启动服务

```bash
go run cmd/main.go
```

服务将在 `8080` 端口启动。

## API 接口

### 用户接口

#### 注册

```
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

#### 登录

```
POST /api/v1/users/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
```

### 商品接口

#### 商品列表（公开）

```
GET /api/v1/products?page=1&page_size=10
```

#### 获取商品（公开）

```
GET /api/v1/products/1
```

#### 创建商品（需认证）

```
POST /api/v1/products
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "iPhone 15",
  "description": "最新款",
  "price": 5999.00,
  "stock": 100
}
```

#### 更新商品（需认证）

```
PUT /api/v1/products/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "iPhone 15 Pro",
  "price": 7999.00
}
```

#### 删除商品（需认证）

```
DELETE /api/v1/products/1
Authorization: Bearer <token>
```

### 订单接口（所有接口需认证）

#### 创建订单

```
POST /api/v1/orders
Authorization: Bearer <token>
Content-Type: application/json

{
  "product_id": 1,
  "quantity": 2
}
```

#### 获取订单

```
GET /api/v1/orders/1
Authorization: Bearer <token>
```

#### 订单列表

```
GET /api/v1/orders?page=1&page_size=10
Authorization: Bearer <token>
```

#### 取消订单

```
DELETE /api/v1/orders/1
Authorization: Bearer <token>
```

## 项目结构

```
api-gateway/
├── cmd/
│   └── main.go           # 程序入口
├── config/
│   ├── config.go         # 配置加载
│   └── config.yaml       # 配置文件
├── internal/
│   ├── handler/          # HTTP 处理器
│   ├── middleware/       # 中间件
│   ├── router/           # 路由配置
│   └── client/           # gRPC 客户端
└── go.mod
```

## 测试示例

### 使用 curl

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'

# 登录
TOKEN=$(curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}' | jq -r '.data.token')

# 创建商品
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","price":5999,"stock":100}'

# 创建订单
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"product_id":1,"quantity":2}'
```

## 中间件

### 认证中间件

- 验证 JWT Token
- 提取用户信息
- 注入到 Context

### CORS 中间件

- 支持跨域请求
- 配置允许的源、方法、头部

## 配置说明

```yaml
server:
  port: 8080 # 服务端口
  mode: debug # 运行模式: debug/release/test

jwt:
  secret: "my_secret_key" # JWT 密钥

services:
  user_service: "localhost:8081"
  product_service: "localhost:8082"
  order_service: "localhost:8083"

cors:
  allow_origins: ["*"]
  allow_methods: ["GET", "POST", "PUT", "DELETE"]
  allow_headers: ["Content-Type", "Authorization"]
```
