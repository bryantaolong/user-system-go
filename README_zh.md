# User System (Go Version)

[English README here (英文版说明)](./README.md)

## 项目简介

本项目为基于 Go 语言的用户管理系统，支持用户注册、登录、信息管理、权限控制、数据导出等功能。后端采用 PostgreSQL 作为主数据库，Redis 用于缓存，支持 JWT 无状态认证和基于角色的权限控制。

## 技术栈

- Go 1.20+
- Gin Web 框架
- GORM (ORM 框架)
- PostgreSQL 17.x
- Redis 6.x 或更高
- JWT (JSON Web Token)

## 目录结构

```
cmd/
  main.go                # 程序入口
internal/
  config/                # 配置相关
  handler/               # 路由处理器（用户、认证等）
  middleware/            # 中间件（认证、权限等）
  model/
    entity/              # 实体定义
    request/             # 请求结构体
    response/            # 响应结构体
  router/                # 路由注册
  service/               # 业务逻辑
    redis/               # Redis 相关服务
pkg/
  db/                    # 数据库初始化
  http/                  # HTTP 工具
  jwt/                   # JWT 工具
go.mod, go.sum           # Go 依赖管理
README.md, README_zh.md  # 项目说明
LICENSE                  # 许可证
```

## 主要功能

- 用户注册、登录、登出
- 用户信息查询与修改
- 密码修改
- 角色管理与权限控制
- 用户禁用/解禁、逻辑删除
- 用户数据分页与搜索
- 用户数据导出（如 CSV/Excel）
- JWT 认证与中间件
- Redis 缓存支持

## 环境要求

- Go 1.20 及以上
- PostgreSQL 17.x
- Redis 6.x 或更高

## 配置说明

- 数据库、Redis 等配置请在 `internal/config/config.go` 或相关环境变量中设置。
- 其他通用配置可参考 `internal/config` 目录下内容。

## 启动方式

1. 初始化数据库（PostgreSQL），建表 SQL 可根据 `model/entity` 结构体自动生成或手动编写。
2. 启动 Redis 服务。
3. 安装依赖并运行：

   ```sh
   go mod tidy
   go run cmd/main.go
   ```

   或编译后运行：

   ```sh
   go build -o user-system-go cmd/main.go
   ./user-system-go
   ```

## 常用接口

- 用户注册：`POST /api/auth/register`
- 用户登录：`POST /api/auth/login`
- 查询所有用户：`GET /api/user/all`（管理员权限）
- 用户搜索：`POST /api/user/search`
- 用户信息更新、角色变更、密码修改、封禁/解封、逻辑删除等接口详见 `internal/handler/user_handler.go`
- 用户数据导出：如 `GET /api/user/export/all`、`POST /api/user/export/field`（管理员权限）

## 其他说明

- JWT 密钥建议通过配置文件或环境变量注入，避免硬编码。
- 全局异常处理与统一响应格式可在 `internal/handler` 或中间件实现。
- 逻辑删除字段建议为 `deleted`，0 表示未删除，1 表示已删除。

## License

本项目采用 MIT 协议。
详见 [LICENSE](LICENSE) 。
