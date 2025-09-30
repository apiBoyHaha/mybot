# MyBot - JWT认证API服务

基于Go语言开发的JWT Token认证API服务，提供用户注册、登录和受保护API访问功能。具体认证流程如下：
    用户登录 → 验证凭证 → 生成Token → 返回Token → 后续请求携带Token → 验证Token → 访问资源
## 项目结构

```
mybot/
├── cmd/app1/main.go              # 应用入口点
├── internal/
│   ├── auth/                     # 认证相关模块
│   │   ├── middleware.go         # JWT认证中间件
│   │   └── handler.go            # 认证处理器
│   ├── user/                     # 用户管理模块
│   │   ├── service.go            # 用户业务逻辑
│   │   ├── repository.go         # 用户数据访问层
│   │   └── model.go              # 用户数据模型
│   └── database/                 # 数据库相关
│       └── connection.go         # 数据库连接配置
├── pkg/
│   ├── config/                   # 配置管理
│   │   └── cfg.go
│   └── utils/                    # 工具函数
│       └── jwt.go                # JWT工具函数
├── api/routes/                   # 路由定义
│   ├── auth_routes.go            # 认证路由
│   └── api_routes.go             # 主路由设置
├── migrations/                   # 数据库迁移文件
│   └── tables.sql                # 用户表结构
├── .env.example                  # 环境配置示例
├── go.mod                        # Go模块定义
└── readme.md                     # 项目说明
```

## 功能特性

- ✅ JWT Token认证
- ✅ 用户注册和登录
- ✅ 密码加密存储（bcrypt）
- ✅ 角色权限管理
- ✅ 受保护API路由
- ✅ 健康检查接口

## 快速开始

### 1. 环境准备

确保已安装：
- Go 1.19+
- MySQL 5.7+

### 2. 数据库设置

```sql
-- 创建数据库
CREATE DATABASE mybot_db;

-- 执行迁移文件
-- 运行 migrations/tables.sql 中的SQL语句
```

### 3. 配置环境

复制环境配置文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件，配置数据库连接等信息：
```env
SERVER_PORT=8080
DATABASE_URL=username:password@tcp(localhost:3306)/mybot_db
JWT_SECRET=your-super-secret-jwt-key-here
ENVIRONMENT=development
```

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 启动服务

```bash
go run cmd/app1/main.go
```

服务将启动在 http://localhost:8080

## API文档

### 认证接口

#### 用户注册
- **URL**: `POST /api/register`
- **Body**:
```json
{
  "username": "testuser",
  "password": "password123",
  "email": "test@example.com",
  "role": "user"
}
```

#### 用户登录
- **URL**: `POST /api/login`
- **Body**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```
- **Response**:
```json
{
  "token": "jwt-token-here",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "role": "user",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### 受保护接口

#### 获取受保护数据
- **URL**: `GET /api/protected/data`
- **Headers**: `Authorization: Bearer <jwt-token>`
- **Response**:
```json
{
  "message": "这是受保护的数据",
  "status": "success"
}
```

### 健康检查

- **URL**: `GET /health`
- **Response**:
```json
{
  "status": "healthy"
}
```

## 开发说明

### 添加新的受保护路由

在 `api/routes/api_routes.go` 的 `protectedMux` 中添加新的路由：

```go
protectedMux.HandleFunc("/api/protected/new-endpoint", newProtectedHandler)
```

### 自定义认证逻辑

修改 `internal/auth/middleware.go` 中的 `JWTAuthMiddleware` 函数。

### 数据库操作

用户相关的数据库操作在 `internal/user/repository.go` 中实现。

## 安全注意事项

1. 生产环境务必修改默认的JWT密钥
2. 使用HTTPS协议传输敏感数据
3. 定期更新依赖包版本
4. 设置合适的Token过期时间
5. 实施密码强度策略

## 许可证

MIT License