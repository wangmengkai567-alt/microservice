# Product Service vs User Service 对比说明

## 架构相似点

两个服务都采用相同的分层架构：

```
cmd/          → 程序入口
config/       → 配置管理
internal/
  ├── handler/     → gRPC 处理器层
  ├── service/     → 业务逻辑层
  ├── repository/  → 数据访问层
  ├── model/       → 数据模型
  ├── middleware/  → 中间件
  └── pkg/         → 工具包
proto/        → gRPC 定义
```

## 核心差异

### 1. 服务端口

| 服务            | 端口 |
| --------------- | ---- |
| user-service    | 8081 |
| product-service | 8082 |

### 2. 数据库

| 服务            | 数据库名称      |
| --------------- | --------------- |
| user-service    | user_service    |
| product-service | product_service |

### 3. 数据模型

#### User Model

```go
type User struct {
    ID        uint
    Username  string    // 唯一索引
    Password  string    // bcrypt 加密
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### Product Model

```go
type Product struct {
    ID          uint
    Name        string    // 商品名称
    Description string    // 商品描述
    Price       float64   // 商品价格
    Stock       int       // 库存数量
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### 4. API 接口

#### User Service

- `Register` - 用户注册
- `Login` - 用户登录（返回 Token）

#### Product Service

- `CreateProduct` - 创建商品
- `GetProduct` - 获取商品详情
- `UpdateProduct` - 更新商品
- `DeleteProduct` - 删除商品
- `ListProducts` - 商品列表（分页）

### 5. 认证策略

#### User Service

```go
// 白名单：Login 和 Register 不需要认证
if strings.Contains(info.FullMethod, "Login") ||
   strings.Contains(info.FullMethod, "Register") {
    return handler(ctx, req)
}
```

#### Product Service

```go
// 白名单：GetProduct 和 ListProducts 不需要认证（只读操作）
if strings.Contains(info.FullMethod, "GetProduct") ||
   strings.Contains(info.FullMethod, "ListProducts") {
    return handler(ctx, req)
}
```

### 6. 业务逻辑差异

#### User Service 特有功能

- 密码加密（bcrypt）
- JWT Token 生成
- 用户名唯一性验证

#### Product Service 特有功能

- 分页查询
- 价格和库存验证
- 商品 CRUD 完整操作

### 7. 工具包差异

#### User Service pkg/

- `jwt.go` - JWT Token 生成和解析
- `password.go` - 密码加密和验证

#### Product Service pkg/

- `jwt.go` - JWT Token 解析（仅用于验证）
- 无需 `password.go`（不涉及密码操作）

## 服务间集成

### Token 共享机制

1. **统一的 JWT Secret**

两个服务在 `config.yaml` 中使用相同的 JWT Secret：

```yaml
jwt:
  secret: "my_secret_key" # 必须保持一致
  expire: 7200
```

2. **Token 流转**

```
用户 → user-service.Login() → 获取 Token
     ↓
用户携带 Token → product-service.CreateProduct()
     ↓
product-service 验证 Token（使用相同的 Secret）
     ↓
验证通过，执行操作
```

### 调用流程示例

```
1. 用户注册
   POST user-service:8081/Register
   → 创建用户账号

2. 用户登录
   POST user-service:8081/Login
   → 返回 JWT Token

3. 创建商品（需要认证）
   POST product-service:8082/CreateProduct
   Header: authorization = <Token>
   → 验证 Token → 创建商品

4. 查询商品（无需认证）
   GET product-service:8082/ListProducts
   → 直接返回商品列表
```

## 代码复用

### 可以直接复用的代码

1. **配置管理** - `config/config.go` 结构完全相同
2. **JWT 工具** - `internal/pkg/jwt.go` 可直接复用
3. **中间件框架** - `internal/middleware/auth.go` 逻辑相同，只需修改白名单

### 需要定制的代码

1. **数据模型** - 根据业务需求定义不同的结构
2. **Repository** - 根据模型实现不同的数据库操作
3. **Service** - 实现各自的业务逻辑
4. **Handler** - 实现各自的 gRPC 接口
5. **Proto 定义** - 定义各自的服务和消息

## 扩展建议

### 统一的基础设施

可以将公共代码抽取到 `common` 包：

```
common/
├── config/      # 配置管理
├── middleware/  # 通用中间件
├── pkg/
│   ├── jwt/     # JWT 工具
│   ├── logger/  # 日志工具
│   └── errors/  # 错误处理
└── database/    # 数据库连接
```

### 服务发现

引入服务注册中心（如 Consul、Etcd）：

```
user-service    → 注册到 Consul
product-service → 注册到 Consul
api-gateway     → 从 Consul 发现服务
```

### API Gateway

添加统一的 API 网关：

```
客户端
  ↓
API Gateway (8080)
  ├→ user-service (8081)
  └→ product-service (8082)
```

### 配置中心

使用配置中心统一管理配置：

```
Consul/Nacos
  ├→ user-service 配置
  └→ product-service 配置
```

## 总结

| 维度     | User Service        | Product Service            |
| -------- | ------------------- | -------------------------- |
| 核心功能 | 用户认证            | 商品管理                   |
| 端口     | 8081                | 8082                       |
| 数据库   | user_service        | product_service            |
| 认证策略 | 登录/注册公开       | 读操作公开                 |
| 特殊功能 | 密码加密、Token生成 | 分页查询、库存管理         |
| 依赖关系 | 独立服务            | 依赖 user-service 的 Token |

两个服务通过共享 JWT Secret 实现了统一认证，形成了完整的微服务体系。
