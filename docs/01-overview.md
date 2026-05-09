# 用户模块概览

## 简介

用户模块是本系统的核心模块之一，负责管理用户的基本信息、认证授权、角色权限以及用户资料等功能。该模块基于 Spring Boot 和 Spring Security 构建，采用标准的三层架构设计。

## 核心功能

- 用户管理：创建、查询、更新、删除用户
- 角色管理：用户角色分配与权限控制
- 用户资料：扩展用户信息（真实姓名、性别、生日、头像等）
- 用户状态：正常、锁定、封禁状态管理
- 安全特性：密码加密、登录失败锁定、乐观锁、逻辑删除

## 技术栈

- Spring Boot 3.x
- Spring Security 6.x
- MyBatis 3.x
- PostgreSQL（主数据库）
- MySQL（可选支持）
- Lombok
- Jakarta Validation

## 模块结构

```
backend/src/main/java/com/bryan/system/
├── controller/user/          # 控制器层
│   ├── UserController.java
│   ├── UserProfileController.java
│   ├── UserRoleController.java
│   └── UserExportController.java
├── service/user/             # 服务层
│   ├── UserService.java
│   ├── UserProfileService.java
│   ├── UserRoleService.java
│   └── UserExportService.java
├── mapper/                   # 数据访问层
│   ├── UserMapper.java
│   ├── UserProfileMapper.java
│   └── UserRoleMapper.java
├── domain/                   # 领域模型
│   ├── entity/
│   │   ├── SysUser.java
│   │   ├── UserProfile.java
│   │   └── UserRole.java
│   ├── dto/
│   │   ├── UserUpdateDTO.java
│   │   └── UserProfileUpdateDTO.java
│   ├── request/user/
│   │   ├── UserCreateRequest.java
│   │   ├── UserUpdateRequest.java
│   │   ├── UserSearchRequest.java
│   │   ├── ChangeRoleRequest.java
│   │   └── ChangePasswordRequest.java
│   ├── vo/
│   │   ├── UserProfileVO.java
│   │   └── UserRoleOptionVO.java
│   └── enums/user/
│       ├── UserStatusEnum.java
│       └── GenderEnum.java
└── ...

backend/src/main/resources/
├── mapper/
│   ├── UserMapper.xml
│   ├── UserProfileMapper.xml
│   └── UserRoleMapper.xml
└── ...

frontend/src/
├── api/
├── views/
└── components/
```

## 数据库表

用户模块涉及三张核心数据表：

1. `sys_user` - 用户基本信息表
2. `user_profile` - 用户资料扩展表
3. `user_role` - 用户角色表

详细表结构请参考 [数据模型文档](./02-data-model.md)。

## API 端点概览

### 用户管理 API (`/api/users`)

- `POST /api/users` - 创建用户（管理员）
- `GET /api/users` - 获取用户列表（分页）
- `GET /api/users/{userId}` - 根据ID查询用户
- `GET /api/users/username/{username}` - 根据用户名查询
- `POST /api/users/search` - 搜索用户（多条件）
- `PUT /api/users/{userId}` - 更新用户信息
- `PUT /api/users/roles/{userId}` - 修改用户角色
- `PUT /api/users/password/{userId}` - 重置密码
- `PUT /api/users/block/{userId}` - 封禁用户
- `PUT /api/users/unblock/{userId}` - 解封用户
- `DELETE /api/users/{userId}` - 删除用户（逻辑删除）

### 用户资料 API (`/api/user-profiles`)

- `GET /api/user-profiles/me` - 获取当前用户资料
- `GET /api/user-profiles/{userId}` - 根据用户ID查询资料
- `GET /api/user-profiles/name/{realName}` - 根据真实姓名查询
- `PUT /api/user-profiles` - 更新当前用户资料
- `POST /api/user-profiles/avatar` - 上传头像

### 用户角色 API (`/api/user-roles`)

- `GET /api/user-roles` - 获取所有角色选项

详细 API 文档请参考 [API 接口文档](./03-api-reference.md)。

## 权限控制

用户模块使用 Spring Security 的 `@PreAuthorize` 注解进行方法级权限控制：

- `ROLE_ADMIN` - 管理员角色，拥有所有权限
- `ROLE_USER` - 普通用户角色，默认角色

权限控制详情请参考 [权限与安全文档](./04-security.md)。

## 快速开始

### 1. 初始化数据库

```bash
# PostgreSQL
psql -U postgres -d your_database -f sql/create_table.sql

# MySQL
mysql -u root -p your_database < sql/create_table_mysql.sql
```

### 2. 配置数据库连接

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
# 创建用户（需要管理员权限）
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "email": "test@example.com"
  }'
```

## 下一步

- [数据模型详解](./02-data-model.md)
- [API 接口文档](./03-api-reference.md)
- [权限与安全](./04-security.md)
- [业务逻辑说明](./05-business-logic.md)
- [开发指南](./06-development-guide.md)
