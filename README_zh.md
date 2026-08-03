# 用户系统

## 项目简介

本项目为基于 Go（Gin 框架）的用户管理系统，支持用户注册、登录、信息管理、权限控制、数据导出等功能。后端采用 PostgreSQL 作为主数据库，Redis 用于缓存和分布式场景，支持 JWT 无状态认证和基于角色的权限控制。

## 技术栈

### 后端

* Go 1.25
* Gin - Web 框架
* `database/sql` + GORM - ORM
* go-redis/v9 - Redis 客户端
* golang-jwt/v5 - JWT 认证
* Viper - 配置管理
* Zap - 日志
* excelize - Excel 导出
* bcrypt - 密码加密

### 前端

* Vue 3
* TypeScript
* Vite
* Arco Design Vue
* Pinia
* Pinia Plugin Persist
* Vue Router
* Axios

## 项目结构

```
backend/
  main.go              # 应用入口
  auth/                # 认证模块（handler、service、repository、middleware）
  cache/               # Redis 客户端
  config/              # 配置加载
  config.yaml          # 配置文件
  middleware/          # 全局中间件（CORS、错误处理）
  model/               # 数据模型 & DTO
  pkg/                 # 工具包（JWT、Redis、HTTP、响应）
  response/            # 响应封装
  system/              # 系统模块（日志等）
  user/                # 用户模块（handler、service、repository、file）
  go.mod

frontend/
  src/
    api/               # API 请求模块
    assets/            # 静态资源
    components/        # 公共组件
    router/            # 路由定义与守卫
    stores/            # Pinia 状态管理
    styles/            # 全局样式
    types/             # TypeScript 类型定义
    utils/             # 工具函数
    views/             # 页面视图
      login/           # 登录页
      register/        # 注册页
      profile/         # 个人中心
      admin/           # 管理员页面
        users/         # 用户管理
        logs/          # 系统日志
  package.json
  vite.config.ts
  tsconfig.json
```

## 环境要求

* Go 1.25+
* Node.js 18+
* PostgreSQL 17.x / MySQL 8.0.x
* Redis 6.x 或更高

## 配置说明

* 编辑 `backend/config.yaml` 配置数据库、Redis、JWT、服务器端口等参数。
* 数据库建表脚本见 [`sql/create_table.sql`](sql/create_table.sql)。

`config.yaml` 示例：

```yaml
server:
  port: 8080

database:
  driver: postgres
  host: localhost
  port: 5432
  username: postgres
  password: "postgres"
  dbname: user_system
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret_key: "your-secret-key"
  expiration_ms: 86400000
  token_prefix: "Bearer "

cors:
  allowed_origins: "http://localhost:5173"

logging:
  level: info
  file: logs/user-system.log
```

## 启动方式

> 以 PostgreSQL 为例

1. 初始化数据库（PostgreSQL），执行建表脚本：

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```

2. 启动 Redis 服务。
3. 启动后端服务：

   ```sh
   cd backend
   go mod tidy
   go run main.go
   ```

4. 启动前端服务（另开终端）：

   ```sh
   cd frontend
   npm install
   npm run dev
   ```

5. 访问 `http://localhost:5173`（Vite 默认端口）。

### 构建

```sh
# 前端
cd frontend && npm run build

# 后端
cd backend && go build -o bin/server main.go
```

## 常用接口

* 用户注册：`POST /api/auth/register`
* 用户登录：`POST /api/auth/login`
* 验证令牌：`GET /api/auth/validate`
* 获取当前用户：`GET /api/auth/me`
* 修改密码：`PUT /api/auth/password`
* 注销账号：`DELETE /api/auth`
* 退出登录：`GET /api/auth/logout`

* 查询所有用户：`GET /api/users`（管理员权限）
* 根据ID查询用户：`GET /api/users/:userId`（管理员权限）
* 根据用户名查询用户：`GET /api/users/username/:username`（管理员权限）
* 搜索用户：`POST /api/users/search`（管理员权限）
* 创建用户：`POST /api/users`（管理员权限）
* 更新用户：`PUT /api/users/:userId`
* 修改用户角色：`PUT /api/users/roles/:userId`（管理员权限）
* 重置密码：`PUT /api/users/password/:userId`（管理员权限）
* 封禁用户：`PUT /api/users/block/:userId`（管理员权限）
* 解封用户：`PUT /api/users/unblock/:userId`（管理员权限）
* 删除用户：`DELETE /api/users/:userId`（管理员权限，逻辑删除）
* 导出用户数据：`GET /api/users/export`（管理员权限）

* 上传头像：`POST /api/user-profiles/avatar`
* 根据用户ID查询资料：`GET /api/user-profiles/:userId`
* 根据真实姓名查询资料：`GET /api/user-profiles/name/:realName`
* 获取当前用户资料：`GET /api/user-profiles/me`
* 更新用户资料：`PUT /api/user-profiles`

* 获取所有角色：`GET /api/user-roles`（管理员权限）

* 查看最新日志：`GET /api/admin/logs`（管理员权限）
* 查看日志文件列表：`GET /api/admin/logs/files`（管理员权限）

## 前端功能

* Vue 3 + TypeScript SPA，基于 Vite 和 Arco Design Vue。
* 认证流程：登录、注册、登出及路由守卫。
* 个人中心：基本信息、安全设置、登录历史。
* 管理面板：用户管理（搜索、更新、角色变更、封禁/解封、删除）、系统日志查看。
* Axios 拦截器实现 JWT 注入和统一错误处理。

## 其他说明

* JWT 密钥建议在生产环境通过配置文件注入，避免硬编码。
* 逻辑删除字段为 `deleted`：0 表示未删除，1 表示已删除。

## License

本项目采用 MIT 协议。详见 [LICENSE](LICENSE)。
