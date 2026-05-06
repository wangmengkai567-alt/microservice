# 前端Web应用

基于Vue 3 + TypeScript + Pinia + TanStack Query + Tailwind CSS构建的电商微服务系统前端应用。

## 技术栈

- **Vue 3** - 渐进式JavaScript框架
- **TypeScript** - 类型安全
- **Pinia** - 状态管理
- **Vue Router** - 路由管理
- **TanStack Query** - 服务端状态管理和缓存
- **Axios** - HTTP客户端
- **Tailwind CSS** - 实用优先的CSS框架
- **VeeValidate + Zod** - 表单验证
- **Vite** - 构建工具
- **Vitest** - 单元测试
- **fast-check** - 属性测试

## 项目结构

```
frontend-web-app/
├── src/
│   ├── components/      # 可复用UI组件
│   ├── views/           # 页面组件
│   ├── composables/     # 组合式函数
│   ├── stores/          # Pinia状态管理
│   ├── services/        # API服务层
│   ├── types/           # TypeScript类型定义
│   ├── utils/           # 工具函数
│   ├── router/          # 路由配置
│   ├── styles/          # 全局样式
│   ├── App.vue          # 根组件
│   └── main.ts          # 应用入口
├── tests/               # 测试文件
├── public/              # 静态资源
└── package.json
```

## 本地开发

### 安装依赖

```bash
cd frontend-web-app
npm install
```

### 启动开发服务器

```bash
npm run dev
```

应用将在 http://localhost:5173 启动

### 构建生产版本

```bash
npm run build
```

### 预览生产构建

```bash
npm run preview
```

## 测试

### 运行单元测试

```bash
npm run test:unit
```

### 运行属性测试

```bash
npm run test:property
```

### 运行集成测试

```bash
npm run test:integration
```

## 代码质量

### 运行ESLint

```bash
npm run lint
```

### 格式化代码

```bash
npm run format
```

## API配置

前端通过HTTP REST API与后端API Gateway通信：

- 开发环境: http://localhost:8080
- 生产环境: 在.env.production中配置

## 功能特性

- ✅ 用户注册、登录、退出
- ✅ 商品列表、详情、创建、编辑、删除
- ✅ 订单创建、列表、详情、取消
- ✅ JWT Token认证
- ✅ 表单验证和错误提示
- ✅ 响应式布局（移动端、平板、桌面）
- ✅ 可访问性支持（ARIA标签、键盘导航）
- ✅ 加载状态和错误处理
- ✅ 数据缓存和自动重新验证

## 部署

### Docker部署

```bash
docker build -t frontend-web-app .
docker run -p 80:80 frontend-web-app
```

### Docker Compose

```bash
docker-compose up -d
```

## 浏览器支持

- Chrome (最新版)
- Firefox (最新版)
- Safari (最新版)
- Edge (最新版)

## 许可证

MIT
