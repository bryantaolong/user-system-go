# 用户系统

## 项目简介

本项目为基于 Go（Gin 框架）的用户管理系统，支持用户注册、登录、信息管理、权限控制、数据导出等功能。后端采用 PostgreSQL 作为主数据库，Redis 用于缓存和分布式场景，支持 JWT 无状态认证和基于角色的权限控制。

## 技术栈

* Go 1.22
* Gin - Web 框架
* GORM - ORM
* go-redis/v9 - Redis 客户端
* golang-jwt/v5 - JWT 认证
* Viper - 配置管理
* Zap - 日志
* excelize - Excel 导出
* bcrypt - 密码加密

## 项目结构

```
backend/
  cmd/server/          # 应用入口
  internal/
    config/            # 配置加载
    handler/           # HTTP 处理器
    middleware/        # 中间件（认证、CORS、错误处理）
    model/             # 数据模型 & DTO
    pkg/               # 工具包（JWT、Redis、HTTP、响应）
    repository/        # 数据访问层
    service/           # 业务逻辑层
  config.yaml          # 配置文件
  go.mod
frontend/
  src/
  package.json
```

## 环境要求

* Go 1.22+
* PostgreSQL 17.x / MySQL 8.0.x
* Redis 6.x 或更高

## 配置说明

* 编辑 `backend/config.yaml` 配置数据库、Redis、JWT 等参数。
* 数据库建表脚本见 [`sql/create_table.sql`](sql/create_table.sql)。

## 启动方式

> 以 PostgreSQL 为例

1. 初始化数据库（PostgreSQL），执行建表脚本：

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```
2. 启动 Redis 服务。
3. 安装依赖并运行项目：

   ```sh
   cd backend
   go mod tidy
   go run cmd/server/main.go
   ```

## 常用接口

* 用户注册：`POST /api/auth/register`
* 用户登录：`POST /api/auth/login`
* 获取当前用户：`GET /api/auth/me`
* 查询所有用户：`GET /api/users`（管理员权限）
* 用户搜索：`POST /api/users/search`
* 用户更新、角色变更、密码修改、封禁/解封、逻辑删除等
* 用户数据导出：`GET /api/users/export`（管理员权限）
* 用户档案：头像上传、档案 CRUD

## 其他说明

* JWT 密钥建议在生产环境通过配置文件注入，避免硬编码。
* 逻辑删除字段为 `deleted`：0 表示未删除，1 表示已删除。

## License

本项目采用 MIT 协议。详见 [LICENSE](LICENSE)。
