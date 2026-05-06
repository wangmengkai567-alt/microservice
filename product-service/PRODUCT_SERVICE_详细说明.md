# Product Service 项目详细说明文档

## 项目概述

这是一个基于 Go 语言开发的商品服务微服务，使用 gRPC 作为通信协议，提供商品的增删改查功能。项目采用与 user-service 相同的分层架构设计，包含数据层、业务逻辑层、处理层和传输层。

---

## 项目结构

```
product-service/
├── cmd/                    # 应用程序入口
│   └── main.go
├── config/                 # 配置管理
│   ├── config.go
│   └── config.yaml
├── internal/               # 内部代码（不对外暴露）
│   ├── handler/           # gRPC 处理器层
│   ├── middleware/        # 中间件
│   ├── model/             # 数据模型
│   ├── pkg/               # 工具包
│   ├── repository/        # 数据访问层
│   └── service/           # 业务逻辑层
├── proto/                  # Protocol Buffers 定义
└── go.mod                  # Go 模块依赖
```

---

## 核心依赖

```go
// 主要依赖包
- google.golang.org/grpc          // gRPC 框架
- gorm.io/gorm                    // ORM 数据库操作
- gorm.io/driver/mysql            // MySQL 驱动
- github.com/golang-jwt/jwt/v5    // JWT 令牌验证
- gopkg.in/yaml.v3                // YAML 配置解析
```

---

## 详细代码分析

### 1. 程序入口 - `cmd/main.go`

服务启动在 8082 端口，初始化流程与 user-service 类似：

1. 加载配置文件
2. 连接 MySQL 数据库
3. 自动迁移 Product 表结构
4. 依赖注入构建各层实例
5. 注册 gRPC 服务并启动

---

### 2. 配置管理 - `config/`

#### `config.yaml` - 配置文件

```yaml
server:
  port: 8082 # 服务端口

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: 123456
  dbname: product_service # 独立的数据库

jwt:
  secret: "my_secret_key" # 与 user-service 共享密钥
  expire: 7200
```

---

### 3. 数据模型 - `internal/model/product.go`

```go
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"not null"`           // 商品名称
    Description string    `gorm:"type:text"`          // 商品描述
    Price       float64   `gorm:"not null"`           // 商品价格
    Stock       int       `gorm:"not null;default:0"` // 库存数量
    CreatedAt   time.Time                             // 创建时间
    UpdatedAt   time.Time                             // 更新时间
}
```

---

### 4. 数据访问层 - `internal/repository/product_repository.go`

提供以下数据库操作：

- `Create(product)` - 创建商品
- `FindByID(id)` - 根据 ID 查询商品
- `Update(product)` - 更新商品信息
- `Delete(id)` - 删除商品
- `List(page, pageSize)` - 分页查询商品列表

---

### 5. 业务逻辑层 - `internal/service/product_service.go`

#### 创建商品

```go
func (s *ProductService) CreateProduct(name, description string, price float64, stock int) (uint, error)
```

- 验证商品名称不为空
- 验证价格大于 0
- 验证库存不为负数
- 返回创建的商品 ID

#### 获取商品

```go
func (s *ProductService) GetProduct(id uint) (*model.Product, error)
```

- 验证 ID 有效性
- 返回商品详情

#### 更新商品

```go
func (s *ProductService) UpdateProduct(id uint, name, description string, price float64, stock int) error
```

- 先查询商品是否存在
- 更新非空字段
- 保存到数据库

#### 删除商品

```go
func (s *ProductService) DeleteProduct(id uint) error
```

- 验证商品存在
- 执行软删除或硬删除

#### 商品列表

```go
func (s *ProductService) ListProducts(page, pageSize int) ([]model.Product, int64, error)
```

- 默认每页 10 条
- 最大每页 100 条
- 返回商品列表和总数

---

### 6. 处理器层 - `internal/handler/product_handler.go`

实现 5 个 gRPC 接口：

1. `CreateProduct` - 创建商品
2. `GetProduct` - 获取商品详情
3. `UpdateProduct` - 更新商品
4. `DeleteProduct` - 删除商品
5. `ListProducts` - 商品列表（分页）

---

### 7. 中间件 - `internal/middleware/auth.go`

认证策略：

- `GetProduct` 和 `ListProducts` - 公开接口，无需认证
- `CreateProduct`、`UpdateProduct`、`DeleteProduct` - 需要 JWT Token 认证

---

### 8. Protocol Buffers 定义 - `proto/product.proto`

```protobuf
service ProductService {
  rpc CreateProduct (CreateProductRequest) returns (CreateProductResponse);
  rpc GetProduct (GetProductRequest) returns (GetProductResponse);
  rpc UpdateProduct (UpdateProductRequest) returns (UpdateProductResponse);
  rpc DeleteProduct (DeleteProductRequest) returns (DeleteProductResponse);
  rpc ListProducts (ListProductsRequest) returns (ListProductsResponse);
}
```

---

## 数据流程图

### 创建商品流程

```
客户端（携带 Token）
  → gRPC CreateProduct 请求
  → AuthInterceptor 验证 Token
  → Handler.CreateProduct()
  → Service.CreateProduct()
    → 验证输入（名称、价格、库存）
    → Repository.Create()
      → 保存到 MySQL
  ← 返回商品 ID
```

### 查询商品流程

```
客户端
  → gRPC GetProduct 请求
  → Handler.GetProduct()（无需认证）
  → Service.GetProduct()
    → Repository.FindByID()
      → 从 MySQL 查询
  ← 返回商品详情
```

### 商品列表流程

```
客户端
  → gRPC ListProducts 请求（page, pageSize）
  → Handler.ListProducts()（无需认证）
  → Service.ListProducts()
    → 验证分页参数
    → Repository.List()
      → 分页查询 MySQL
  ← 返回商品列表和总数
```

---

## 与 User Service 的集成

1. **共享 JWT 密钥**：两个服务使用相同的 JWT Secret，实现统一认证
2. **Token 传递**：客户端从 user-service 获取 Token，在调用 product-service 时携带
3. **认证流程**：
   - 用户在 user-service 登录获取 Token
   - 调用 product-service 的写操作时携带 Token
   - product-service 验证 Token 有效性

---

## 运行说明

1. 创建数据库：

```sql
CREATE DATABASE product_service;
```

2. 生成 proto 文件：

```bash
cd product-service
protoc --go_out=. --go-grpc_out=. proto/product.proto
```

3. 运行服务：

```bash
go run cmd/main.go
```

4. 服务将在 `8082` 端口监听 gRPC 请求

---

## API 使用示例

### 1. 创建商品（需要认证）

```
metadata: authorization = <JWT Token>
CreateProductRequest {
  name: "iPhone 15"
  description: "最新款苹果手机"
  price: 5999.00
  stock: 100
}
```

### 2. 获取商品（无需认证）

```
GetProductRequest {
  id: 1
}
```

### 3. 商品列表（无需认证）

```
ListProductsRequest {
  page: 1
  page_size: 10
}
```

---

## 技术亮点

1. **RESTful 风格的 gRPC API**：完整的 CRUD 操作
2. **分页查询**：支持大数据量的商品列表查询
3. **灵活的认证策略**：读操作公开，写操作需要认证
4. **输入验证**：在 Service 层进行业务规则验证
5. **统一架构**：与 user-service 保持一致的代码结构

---

## 改进建议

1. **库存管理**：添加库存扣减和回滚机制
2. **商品分类**：增加分类字段和分类查询
3. **搜索功能**：支持按名称、价格范围搜索
4. **图片管理**：添加商品图片 URL 字段
5. **缓存层**：使用 Redis 缓存热门商品
6. **日志记录**：记录商品操作日志
7. **软删除**：使用 GORM 的软删除功能
8. **价格历史**：记录商品价格变更历史

---

## 总结

Product Service 是一个功能完整的商品管理微服务，采用与 User Service 相同的架构模式，实现了商品的增删改查和分页查询功能。通过 JWT 认证与用户服务集成，形成了完整的微服务体系。
