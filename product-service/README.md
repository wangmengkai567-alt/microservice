# Product Service

商品管理微服务，提供商品的增删改查功能。

## 功能特性

- 创建商品
- 获取商品详情
- 更新商品信息
- 删除商品
- 商品列表（分页）

## 技术栈

- Go 1.25.5
- gRPC
- GORM
- MySQL
- JWT 认证

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 生成 Proto 文件

需要先安装 protoc 编译器和 Go 插件：

```bash
# 安装 protoc（根据操作系统选择）
# Windows: 下载 https://github.com/protocolbuffers/protobuf/releases
# Mac: brew install protobuf
# Linux: apt-get install protobuf-compiler

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 生成代码
protoc --go_out=. --go-grpc_out=. proto/product.proto
```

### 3. 配置数据库

创建数据库：

```sql
CREATE DATABASE product_service;
```

修改 `config/config.yaml` 中的数据库配置。

### 4. 运行服务

```bash
go run cmd/main.go
```

服务将在 `8082` 端口启动。

## API 接口

### 1. CreateProduct（需要认证）

创建新商品

**请求**：

```protobuf
CreateProductRequest {
  name: "商品名称"
  description: "商品描述"
  price: 99.99
  stock: 100
}
```

**响应**：

```protobuf
CreateProductResponse {
  id: 1
  message: "product created successfully"
}
```

### 2. GetProduct（公开）

获取商品详情

**请求**：

```protobuf
GetProductRequest {
  id: 1
}
```

**响应**：

```protobuf
GetProductResponse {
  id: 1
  name: "商品名称"
  description: "商品描述"
  price: 99.99
  stock: 100
}
```

### 3. UpdateProduct（需要认证）

更新商品信息

**请求**：

```protobuf
UpdateProductRequest {
  id: 1
  name: "新商品名称"
  description: "新描述"
  price: 89.99
  stock: 50
}
```

### 4. DeleteProduct（需要认证）

删除商品

**请求**：

```protobuf
DeleteProductRequest {
  id: 1
}
```

### 5. ListProducts（公开）

获取商品列表

**请求**：

```protobuf
ListProductsRequest {
  page: 1
  page_size: 10
}
```

**响应**：

```protobuf
ListProductsResponse {
  products: [...]
  total: 100
}
```

## 认证说明

- 读操作（GetProduct、ListProducts）：无需认证
- 写操作（CreateProduct、UpdateProduct、DeleteProduct）：需要在 metadata 中携带 JWT Token

Token 获取方式：从 user-service 登录接口获取

**gRPC Metadata 示例**：

```
authorization: <JWT Token>
```

## 项目结构

```
product-service/
├── cmd/                    # 程序入口
├── config/                 # 配置文件
├── internal/
│   ├── handler/           # gRPC 处理器
│   ├── service/           # 业务逻辑
│   ├── repository/        # 数据访问
│   ├── model/             # 数据模型
│   ├── middleware/        # 中间件
│   └── pkg/               # 工具包
└── proto/                  # Proto 定义
```

## 配置说明

`config/config.yaml`：

```yaml
server:
  port: 8082 # 服务端口

database:
  host: 127.0.0.1 # 数据库地址
  port: 3306 # 数据库端口
  user: root # 数据库用户
  password: 123456 # 数据库密码
  dbname: product_service # 数据库名称

jwt:
  secret: "my_secret_key" # JWT 密钥（需与 user-service 一致）
  expire: 7200 # Token 过期时间（秒）
```

## 开发说明

详细的架构说明和代码分析请参考：[PRODUCT*SERVICE*详细说明.md](./PRODUCT_SERVICE_详细说明.md)
