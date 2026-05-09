# 用户系统

## 项目简介

本项目为基于 Spring Boot 3 的用户管理系统，支持用户注册、登录、信息管理、权限控制、数据导出等功能。后端采用 PostgreSQL 作为主数据库，Redis 用于缓存和分布式场景，支持 JWT 无状态认证和基于角色的权限控制。

## 技术栈

* Java 17
* Spring Boot 3.5.4
* MyBatis
* PostgreSQL 17.x
* MySQL 8.0.x
* Redis
* Spring Security
* EasyExcel (阿里巴巴 Excel 导出)
* Lombok
* JJWT (JWT 令牌)
* Maven 3.9.x

## 项目结构

```
backend/
  src/
    main/
      java/com/bryan/system/
        config/         # 配置类（安全、Redis、MyBatis 等）
        controller/     # RESTful 控制器
        domain/         # 实体、请求/响应对象、VO
        filter/         # JWT 认证过滤器
        handler/        # 全局异常处理
        mapper/         # MyBatis Mapper 接口
        service/        # 业务服务层
        util/           # 工具类（JWT、HTTP 等）
      resources/
        application.yaml
        application-dev.yaml
        application-mysql.yaml
        mapper/         # MyBatis Mapper XML 文件
    test/
      java/com/bryan/system/
        UserSystemApplicationTests.java
  pom.xml
  mvnw
frontend/
  src/
  package.json
```

## 环境要求

* JDK 17+
* Maven 3.9.9+
* PostgreSQL 17.x/MySQL 8.0.x
* Redis 6.x 或更高

## 配置说明

* 数据库连接、Redis 配置请在 `backend/src/main/resources/application-dev.yaml` 中修改。
* 日志、MyBatis 等通用配置见 `backend/src/main/resources/application.yaml`。
* 数据库建表脚本见 [`sql/create_table.sql`](sql/create_table.sql)。

## 启动方式

> 以 PostgreSQL 为例

1. 初始化数据库（PostgreSQL），执行建表脚本：

   ```sh
   psql -U postgres -d postgres -f sql/create_table.sql
   ```
2. 启动 Redis 服务。
3. 使用 Maven 构建并运行项目：

   ```sh
   cd backend
   ./mvnw spring-boot:run
   ```

   或直接运行打包后的 JAR：

   ```sh
   cd backend
   ./mvnw clean package
   java -jar target/user-system-0.0.1-SNAPSHOT.jar
   ```

## 常用接口

* 用户注册：`POST /api/auth/register`
* 用户登录：`POST /api/auth/login`
* 查询所有用户：`GET /api/user/all`（管理员权限）
* 用户搜索：`POST /api/user/search`
* 用户更新、角色变更、密码修改、封禁/解封、逻辑删除等详见 [`UserController`](backend/src/main/java/com/bryan/system/controller/user/UserController.java)
* 用户数据导出：`GET /api/user/export/all`、`POST /api/user/export/field`（管理员权限）

## 其他说明

* JWT 密钥建议在生产环境通过配置文件注入，避免硬编码。
* 全局异常处理类见 [`GlobalExceptionHandler`](backend/src/main/java/com/bryan/system/handler/GlobalExceptionHandler.java)。
* 逻辑删除字段为 `deleted`：0 表示未删除，1 表示已删除。

## License

本项目采用 MIT 协议。详见 [LICENSE](LICENSE)。
