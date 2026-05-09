# 用户模块文档

欢迎查阅用户模块文档！本文档集帮助新开发者快速了解和上手用户管理系统。

## 文档目录

### [01. 概览](./01-overview.md)
- 模块简介与核心功能
- 技术栈与模块结构
- 数据库表概览
- API 端点概览
- 快速开始指南

### [02. 数据模型](./02-data-model.md)
- 数据库表结构详解
- 实体类说明
- 枚举类型定义
- 表关系与数据完整性
- 数据库兼容性说明

### [03. API 接口文档](./03-api-reference.md)
- 用户管理 API
- 用户资料 API
- 用户角色 API
- 请求/响应格式
- 错误码说明

### [04. 权限与安全](./04-security.md)
- JWT Token 认证机制
- 角色与权限控制
- 安全特性（密码加密、登录锁定等）
- 安全配置与最佳实践
- 常见安全问题解答

### [05. 业务逻辑说明](./05-business-logic.md)
- 用户生命周期管理
- 角色分配流程
- 用户状态管理
- 密码管理机制
- 数据一致性保证
- 性能优化策略

### [06. 开发指南](./06-development-guide.md)
- 环境准备与配置
- 代码规范与命名约定
- 开发流程详解
- 调试技巧
- 常见问题解答
- 最佳实践

## 快速导航

### 我是新手开发者
1. 从 [概览](./01-overview.md) 开始，了解模块整体架构
2. 阅读 [数据模型](./02-data-model.md)，理解数据结构
3. 参考 [开发指南](./06-development-guide.md) 搭建开发环境

### 我要调用 API
1. 查看 [API 接口文档](./03-api-reference.md) 了解可用接口
2. 阅读 [权限与安全](./04-security.md) 了解认证方式
3. 参考示例代码进行集成

### 我要修改业务逻辑
1. 阅读 [业务逻辑说明](./05-business-logic.md) 理解现有流程
2. 参考 [开发指南](./06-development-guide.md) 了解开发规范
3. 查看 [数据模型](./02-data-model.md) 确认数据结构

### 我遇到了问题
1. 查看 [开发指南 - 常见问题](./06-development-guide.md#常见问题)
2. 查看 [权限与安全 - 常见安全问题](./04-security.md#常见安全问题)
3. 检查日志输出和错误信息

## 核心概念

### 用户（SysUser）
系统的基本用户实体，包含认证信息（用户名、密码）、联系方式（手机、邮箱）、状态信息等。

### 用户资料（UserProfile）
用户的扩展信息，包含真实姓名、性别、生日、头像等个人资料。

### 用户角色（UserRole）
定义用户的权限级别，系统预置管理员（ROLE_ADMIN）和普通用户（ROLE_USER）两种角色。

### 用户状态
- NORMAL（正常）：可正常登录使用
- LOCKED（锁定）：临时锁定，1小时后自动解锁
- BANNED（封禁）：永久封禁，需管理员手动解封

## 技术栈

- Spring Boot 3.x
- Spring Security 6.x
- MyBatis 3.x
- PostgreSQL / MySQL
- JWT
- Lombok

## 项目结构

```
backend/src/main/java/com/bryan/system/
├── controller/user/          # 用户相关控制器
├── service/user/             # 用户相关服务
├── mapper/                   # 数据访问层
├── domain/                   # 领域模型
│   ├── entity/              # 实体类
│   ├── dto/                 # 数据传输对象
│   ├── request/user/        # 请求对象
│   ├── vo/                  # 视图对象
│   └── enums/user/          # 枚举类
└── ...

backend/src/main/resources/
├── mapper/                   # MyBatis XML 映射文件
└── application*.yaml         # 配置文件

frontend/src/
├── api/                      # 前端接口封装
├── views/                    # 页面视图
└── components/               # 组件

sql/
├── create_table.sql          # PostgreSQL 初始化脚本
└── create_table_mysql.sql    # MySQL 初始化脚本
```

## 快速开始

### 1. 初始化数据库

```bash
# PostgreSQL
psql -U postgres -d your_database -f sql/create_table.sql

# MySQL
mysql -u root -p your_database < sql/create_table_mysql.sql
```

### 2. 配置应用

编辑 `backend/src/main/resources/application-dev.yaml`：

```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/your_database
    username: your_username
    password: your_password
```

### 3. 启动应用

```bash
cd backend
./mvnw spring-boot:run
```

### 4. 测试 API

```bash
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 相关资源

- [Spring Boot 文档](https://spring.io/projects/spring-boot)
- [Spring Security 文档](https://spring.io/projects/spring-security)
- [MyBatis 文档](https://mybatis.org/mybatis-3/)
- [PostgreSQL 文档](https://www.postgresql.org/docs/)

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'feat(user): add amazing feature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件至项目维护者
- 参与项目讨论

---

最后更新：2024-01-01
