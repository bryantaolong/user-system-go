# User System

用户管理系统

## 技术栈

- **Gin** - Web 框架
- **GORM** - ORM
- **go-redis/v9** - Redis 客户端
- **golang-jwt/v5** - JWT 认证
- **Viper** - 配置管理
- **Zap** - 日志
- **excelize** - Excel 导出
- **bcrypt** - 密码加密

## 快速开始

```bash
# 安装依赖
go mod tidy

# 启动服务
go run cmd/server/main.go
```

## 配置

编辑 `config.yaml` 文件配置数据库、Redis、JWT 等参数。

## 项目结构

```
├── cmd/server/          # 入口
├── internal/
│   ├── config/          # 配置加载
│   ├── model/           # 数据模型 & DTO
│   ├── repository/      # 数据访问层
│   ├── service/         # 业务逻辑层
│   ├── handler/         # HTTP 处理器
│   ├── middleware/       # 中间件
│   └── pkg/             # 工具包
├── config.yaml          # 配置文件
└── go.mod
```
