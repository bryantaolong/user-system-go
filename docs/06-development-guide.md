# 开发指南

## 环境准备

### 1. 开发工具

- JDK 17+
- Maven 3.9+
- PostgreSQL 17+ 或 MySQL 8+
- IDE（推荐 IntelliJ IDEA）
- Git

### 2. 数据库初始化

#### PostgreSQL

```bash
# 创建数据库
createdb user_system

# 执行初始化脚本
psql -U postgres -d user_system -f sql/create_table.sql
```

#### MySQL

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE user_system CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 执行初始化脚本
mysql -u root -p user_system < sql/create_table_mysql.sql
```

### 3. 配置文件

编辑 `backend/src/main/resources/application-dev.yaml`：

```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/user_system
    username: postgres
    password: your_password
    driver-class-name: org.postgresql.Driver

  jpa:
    show-sql: true
    
mybatis:
  configuration:
    log-impl: org.apache.ibatis.logging.stdout.StdOutImpl

jwt:
  secret: your-secret-key-change-in-production
  expiration: 86400000
```

### 4. 启动应用

```bash
# 开发模式启动
cd backend
./mvnw spring-boot:run

# 或者使用 IDE 运行 UserSystemApplication.java
```

## 代码规范

### 1. 命名规范

#### Java 类命名

```java
// 实体类：名词，UpperCamelCase
public class SysUser { }
public class UserProfile { }

// 控制器：资源名 + Controller
public class UserController { }

// 服务类：资源名 + Service
public class UserService { }

// Mapper：资源名 + Mapper
public interface UserMapper { }

// DTO：用途 + DTO
public class UserUpdateDTO { }

// Request：用途 + Request
public class UserCreateRequest { }

// Response/VO：用途 + VO
public class UserProfileVO { }
```

#### 方法命名

```java
// 查询：get/find/query/list
public SysUser getUserById(Long id) { }
public List<SysUser> listUsers() { }

// 创建：create/add
public SysUser createUser(UserCreateRequest req) { }

// 更新：update/modify
public SysUser updateUser(Long id, UserUpdateDTO dto) { }

// 删除：delete/remove
public Long deleteUser(Long id) { }

// 判断：is/has/can
public boolean existsById(Long id) { }
```

### 2. 包结构规范

```
com.bryan.system
├── controller          # 控制器层
│   ├── user           # 用户相关控制器
│   ├── auth           # 认证相关控制器
│   └── admin          # 管理员相关控制器
├── service            # 服务层
│   ├── user           # 用户相关服务
│   ├── auth           # 认证相关服务
│   └── redis          # Redis 相关服务
├── mapper             # 数据访问层
├── domain             # 领域模型
│   ├── entity         # 实体类
│   ├── dto            # 数据传输对象
│   ├── request        # 请求对象
│   ├── response       # 响应对象
│   ├── vo             # 视图对象
│   └── enums          # 枚举类
├── config             # 配置类
├── filter             # 过滤器
├── handler            # 处理器
├── interceptor        # 拦截器
├── exception          # 异常类
└── util               # 工具类
```

### 3. 代码风格

#### 缩进与格式

```java
// 使用 4 个空格缩进
public class UserService {
    
    private final UserMapper userMapper;
    
    public SysUser getUserById(Long userId) {
        return Optional.ofNullable(userMapper.selectById(userId))
                .orElseThrow(() -> new ResourceNotFoundException("用户不存在"));
    }
}
```

#### 注释规范

```java
/**
 * 用户服务实现类，处理用户注册、登录、信息管理等业务逻辑。
 *
 * @author Bryan Long
 */
@Service
public class UserService {
    
    /**
     * 根据用户ID获取用户信息。
     *
     * @param userId 用户 ID
     * @return 用户实体对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public SysUser getUserById(Long userId) {
        // 实现代码
    }
}
```

## 开发流程

### 1. 添加新功能

#### Step 1: 定义实体类

```java
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class SysUser {
    private Long id;
    private String username;
    // ...
}
```

#### Step 2: 创建 Mapper 接口

```java
@Mapper
public interface UserMapper {
    int insert(SysUser user);
    SysUser selectById(Long id);
    int update(SysUser user);
    int deleteById(Long id);
}
```

#### Step 3: 编写 MyBatis XML

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN"
        "http://mybatis.org/dtd/mybatis-3-mapper.dtd">
<mapper namespace="com.bryan.system.mapper.UserMapper">
    
    <resultMap id="BaseResultMap" type="com.bryan.system.domain.entity.SysUser">
        <id column="id" property="id"/>
        <result column="username" property="username"/>
    </resultMap>
    
    <insert id="insert" useGeneratedKeys="true" keyProperty="id">
        INSERT INTO sys_user (username, password) 
        VALUES (#{username}, #{password})
    </insert>
    
</mapper>
```

#### Step 4: 实现 Service 层

```java
@Service
@RequiredArgsConstructor
public class UserService {
    
    private final UserMapper userMapper;
    
    public SysUser createUser(UserCreateRequest req) {
        SysUser user = SysUser.builder()
                .username(req.getUsername())
                .password(passwordEncoder.encode(req.getPassword()))
                .build();
        
        userMapper.insert(user);
        return user;
    }
}
```

#### Step 5: 实现 Controller 层

```java
@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
public class UserController {
    
    private final UserService userService;
    
    @PostMapping
    @PreAuthorize("hasRole('ADMIN')")
    public Result<SysUser> createUser(@RequestBody @Valid UserCreateRequest req) {
        return Result.success(userService.createUser(req));
    }
}
```

#### Step 6: 编写测试

```java
@SpringBootTest
class UserServiceTest {
    
    @Autowired
    private UserService userService;
    
    @Test
    void testCreateUser() {
        UserCreateRequest req = new UserCreateRequest();
        req.setUsername("testuser");
        req.setPassword("password123");
        
        SysUser user = userService.createUser(req);
        
        assertNotNull(user.getId());
        assertEquals("testuser", user.getUsername());
    }
}
```

### 2. 修改现有功能

#### Step 1: 理解现有代码

- 阅读相关实体类、Service、Controller
- 查看 MyBatis XML 映射文件
- 了解业务逻辑和数据流

#### Step 2: 修改代码

- 遵循现有代码风格
- 保持层次分离（Controller -> Service -> Mapper）
- 添加必要的注释

#### Step 3: 测试验证

```bash
# 运行单元测试
cd backend
./mvnw test

# 运行特定测试类
cd backend
./mvnw test -Dtest=UserServiceTest

# 运行特定测试方法
cd backend
./mvnw test -Dtest=UserServiceTest#testCreateUser
```

## 调试技巧

### 1. 日志调试

```java
@Slf4j
@Service
public class UserService {
    
    public SysUser createUser(UserCreateRequest req) {
        log.debug("创建用户请求: {}", req);
        
        SysUser user = userMapper.insert(user);
        
        log.info("用户创建成功: id={}, username={}", user.getId(), user.getUsername());
        
        return user;
    }
}
```

### 2. SQL 日志

在 `application-dev.yaml` 中启用 SQL 日志：

```yaml
mybatis:
  configuration:
    log-impl: org.apache.ibatis.logging.stdout.StdOutImpl

logging:
  level:
    com.bryan.system.mapper: DEBUG
```

### 3. 断点调试

在 IDE 中设置断点：

- Controller 方法入口
- Service 业务逻辑关键点
- Mapper 方法调用前后

### 4. 接口测试

使用 curl 或 Postman 测试：

```bash
# 创建用户
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'

# 查询用户
curl -X GET http://localhost:8080/api/users/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 常见问题

### Q1: MyBatis 映射找不到？

**问题**: `Invalid bound statement (not found)`

**解决**:
1. 检查 Mapper 接口的 `@Mapper` 注解
2. 确认 XML 文件的 namespace 与接口全限定名一致
3. 检查 XML 文件位置是否在 `backend/src/main/resources/mapper/` 下
4. 清理并重新构建项目：先进入 `backend` 目录，再执行 `./mvnw clean install`

### Q2: 数据库连接失败？

**问题**: `Connection refused` 或 `Access denied`

**解决**:
1. 确认数据库服务已启动
2. 检查 `application-dev.yaml` 中的连接配置
3. 验证用户名和密码是否正确
4. 检查数据库是否已创建

### Q3: 权限验证失败？

**问题**: `Access Denied` 或 `403 Forbidden`

**解决**:
1. 检查 JWT Token 是否有效
2. 确认用户角色是否正确
3. 查看 `@PreAuthorize` 注解的权限表达式
4. 检查 Spring Security 配置

### Q4: 事务不生效？

**问题**: 数据未回滚或部分提交

**解决**:
1. 确保方法上有 `@Transactional` 注解
2. 检查异常是否被捕获（捕获后事务不会回滚）
3. 确认方法是通过 Spring 代理调用的（不是内部调用）
4. 检查事务传播行为配置

## 最佳实践

### 1. 使用 DTO 传输数据

```java
// 不要直接使用实体类作为请求参数
// ❌ 错误
@PostMapping
public Result<SysUser> createUser(@RequestBody SysUser user) { }

// ✅ 正确
@PostMapping
public Result<SysUser> createUser(@RequestBody UserCreateRequest req) { }
```

### 2. 参数校验

```java
public class UserCreateRequest {
    
    @NotBlank(message = "用户名不能为空")
    @Size(min = 4, max = 20, message = "用户名长度必须在4-20之间")
    private String username;
    
    @NotBlank(message = "密码不能为空")
    @Size(min = 6, message = "密码长度不能少于6位")
    private String password;
    
    @Email(message = "邮箱格式不正确")
    private String email;
}
```

### 3. 异常处理

```java
// Service 层抛出业务异常
public SysUser getUserById(Long userId) {
    return Optional.ofNullable(userMapper.selectById(userId))
            .orElseThrow(() -> new ResourceNotFoundException("用户不存在"));
}

// 全局异常处理器统一处理
@RestControllerAdvice
public class GlobalExceptionHandler {
    @ExceptionHandler(ResourceNotFoundException.class)
    public Result<?> handleNotFoundException(ResourceNotFoundException e) {
        return Result.error(HttpStatus.NOT_FOUND, e.getMessage());
    }
}
```

### 4. 使用构建器模式

```java
// ✅ 推荐使用 Builder
SysUser user = SysUser.builder()
        .username("testuser")
        .password(encodedPassword)
        .status(UserStatusEnum.NORMAL)
        .build();

// ❌ 避免使用 setter 链
SysUser user = new SysUser();
user.setUsername("testuser");
user.setPassword(encodedPassword);
user.setStatus(UserStatusEnum.NORMAL);
```

### 5. 日志记录

```java
// 关键操作记录日志
log.info("用户创建成功: id={}, username={}", user.getId(), user.getUsername());

// 异常记录详细信息
log.error("用户创建失败: username={}", req.getUsername(), e);

// 调试信息使用 debug 级别
log.debug("查询用户: userId={}", userId);
```

## 提交代码

### 1. Commit 规范

遵循 Conventional Commits：

```bash
# 新功能
git commit -m "feat(user): add user search functionality"

# 修复 bug
git commit -m "fix(user): handle null pointer in getUserById"

# 重构
git commit -m "refactor(user): optimize user query performance"

# 文档
git commit -m "docs(user): update API documentation"
```

### 2. 提交前检查

```bash
# 运行测试
cd backend
./mvnw test

# 检查代码格式
cd backend
./mvnw spotless:check

# 构建项目
cd backend
./mvnw clean package
```

### 3. Pull Request

PR 应包含：
- 清晰的标题和描述
- 关联的 Issue 编号
- 测试验证说明
- 必要的截图（UI 变更）
