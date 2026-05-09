package com.bryan.system.service.auth;

import com.bryan.system.config.properties.SecurityProperties;
import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.entity.UserRole;
import com.bryan.system.domain.enums.user.UserStatusEnum;
import com.bryan.system.domain.request.auth.LoginRequest;
import com.bryan.system.domain.request.auth.RegisterRequest;
import com.bryan.system.exception.BusinessException;
import com.bryan.system.exception.ResourceNotFoundException;
import com.bryan.system.mapper.UserMapper;
import com.bryan.system.mapper.UserRoleMapper;
import com.bryan.system.service.redis.RedisStringService;
import com.bryan.system.util.http.HttpUtils;
import com.bryan.system.util.jwt.JwtUtils;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 用户认证服务类，处理注册、登录、鉴权、当前用户信息等逻辑。
 *
 * @author Bryan Long
 */
@Slf4j
@Service
@RequiredArgsConstructor
public class AuthService implements UserDetailsService {

    private final UserMapper userMapper;
    private final UserRoleMapper userRoleMapper;
    private final PasswordEncoder passwordEncoder;
    private final RedisStringService redisStringService;
    private final SecurityProperties securityProperties;

    /**
     * 用户注册。
     *
     * @param registerRequest 注册请求对象
     * @return 注册成功的用户实体
     * @throws BusinessException 用户名已存在
     * @throws BusinessException 插入数据库失败
     */
    public SysUser register(RegisterRequest registerRequest) {
        // 1. 检查用户名是否已存在
        if(userMapper.selectByUsername(registerRequest.getUsername()) != null) {
            throw new BusinessException("用户名已存在");
        }

        // 2. 查出默认角色
        UserRole defaultRole = userRoleMapper.selectOneByIsDefaultTrue();
        if(defaultRole == null) {
            throw new BusinessException("系统未配置默认角色");
        }

        // 3. 构建用户实体，密码加密
        SysUser sysUser = SysUser.builder()
                .username(registerRequest.getUsername())
                .password(passwordEncoder.encode(registerRequest.getPassword()))
                .phone(registerRequest.getPhone())
                .email(registerRequest.getEmail())
                .roles(defaultRole.getRoleName())
                .status(UserStatusEnum.NORMAL)
                .build();

        this.fillInsert(sysUser);

        // 4. 插入用户数据
        int saved = userMapper.insert(sysUser);
        if (saved == 0) {
            throw new BusinessException("插入数据库失败");
        }

        log.info("用户注册成功: id: {}, username: {} ", sysUser.getId(), sysUser.getUsername());

        // 5. 返回新注册用户
        return sysUser;
    }

    /**
     * 用户登录，验证用户名和密码，生成 JWT Token。
     *
     * @param loginRequest 登录请求对象
     * @return 登录成功后的 JWT Token
     * @throws BusinessException 用户名不存在或密码错误
     */
    public String login(LoginRequest loginRequest) {
        // 1. 验证用户凭证
        SysUser sysUser = userMapper.selectByUsername(loginRequest.getUsername());

        // 无论用户名是否存在，都进行失败次数记录，防止用户名枚举攻击
        if (sysUser == null) {
            // 用户不存在时，也记录一次失败的登录尝试（可以使用固定的用户ID或IP作为key）
            log.warn("登录失败 - 用户不存在: {}", loginRequest.getUsername());
            // 为防止用户名枚举，仍然抛出相同的错误消息
            throw new BusinessException("用户名或密码错误");
        }

        if(!passwordEncoder.matches(loginRequest.getPassword(), sysUser.getPassword())){
            LocalDateTime now = LocalDateTime.now();
            Integer currentFailCount = sysUser.getLoginFailCount() == null ? 0 : sysUser.getLoginFailCount();
            sysUser.setLoginFailCount(currentFailCount + 1);
            // 手动设置审计字段（更新时）- 登录失败时使用用户自己的ID
            this.fillUpdate(sysUser, sysUser.getId() == null ? "SYSTEM" : sysUser.getId().toString());

            // 如果输入密码错误次数达到限额，则锁定账号
            if(sysUser.getLoginFailCount() >= securityProperties.getLoginFailLimit()) {
                sysUser.setStatus(UserStatusEnum.LOCKED);
                sysUser.setLockedAt(now);
                userMapper.update(sysUser);
                log.warn("用户登录失败次数过多，已锁定: {}", sysUser.getUsername());
                throw new BusinessException("输入密码错误次数过多，账号锁定");
            }
            userMapper.update(sysUser);
            log.warn("用户登录密码错误: {}, 失败次数: {}", sysUser.getUsername(), sysUser.getLoginFailCount());
            throw new BusinessException("用户名或密码错误");
        }

        // 2. 检查现有Token（使用JwtUtils验证有效性）
        String existingToken = redisStringService.get(sysUser.getUsername());
        if (existingToken != null && JwtUtils.validateToken(existingToken)) {
            // 刷新 Redis 中的 Token 过期时间
            redisStringService.setExpire(sysUser.getUsername(), 86400000 / 1000);
            return existingToken;
        }

        // 3. 更新用户登录信息
        LocalDateTime now = LocalDateTime.now();
        sysUser.setLastLoginAt(now);
        sysUser.setLastLoginIp(HttpUtils.getClientIp());
        sysUser.setLastLoginDevice(HttpUtils.getClientOS() + " / " + HttpUtils.getClientBrowser());
        sysUser.setLoginFailCount(0); // 重置密码输入错误次数
        // 登录成功时使用用户自己的ID作为updatedBy
        this.fillUpdate(sysUser, sysUser.getId() == null ? "SYSTEM" : sysUser.getId().toString());
        userMapper.update(sysUser);


        // 4. 生成新的JWT Token
        Map<String, Object> claims = new HashMap<>();
        claims.put("username", sysUser.getUsername());
        claims.put("roles", sysUser.getRoles());

        String token = JwtUtils.generateToken(sysUser.getId().toString(), claims);

        // 5. 存储到Redis（设置与JWT相同的过期时间）
        boolean saved = redisStringService.set(
                sysUser.getUsername(),
                token,
                86400000 / 1000
        );

        if (!saved) {
            throw new BusinessException("Token 存储失败");
        }

        return token;
    }

    /**
     * 获取当前登录用户的 ID。
     *
     * @return 当前用户 ID
     */
    public Long getCurrentUserId() {
        // 1. 从 JWT Token 中提取用户 ID
        return JwtUtils.getCurrentUserId();
    }

    /**
     * 获取当前登录用户的用户名。
     *
     * @return 当前用户名
     */
    public String getCurrentUsername() {
        // 1. 从 JWT Token 中提取用户名
        return JwtUtils.getCurrentUsername();
    }

    /**
     * 获取当前登录用户的完整信息。
     *
     * @return 当前用户实体
     */
    public SysUser getCurrentUser() {
        // 1. 获取当前用户 ID
        Long userId = JwtUtils.getCurrentUserId();

        // 2. 查询数据库返回用户信息
        return userMapper.selectById(userId);
    }

    /**
     * 判断用户是否具有管理员权限。
     *
     * @return 是否为管理员
     */
    public boolean isAdmin() {
        // 1. 调用 JwtUtils 获取用户权限列表
        List<String> roles = JwtUtils.getCurrentUserRoles();

        // 2. 遍历用户权限，判断是否包含 ROLE_ADMIN
        for (String role : roles) {
            if ("ROLE_ADMIN".equals(role)) {
                return true;
            }
        }
        return false;
    }

    /**
     * 判断用户是否具有管理员权限。
     *
     * @param userDetails 当前用户信息
     * @return 是否为管理员
     */
    public boolean isAdmin(UserDetails userDetails) {
        // 1. 遍历用户权限，判断是否包含 ROLE_ADMIN
        return userDetails.getAuthorities().stream()
                .anyMatch(a -> a.getAuthority().equals("ROLE_ADMIN"));
    }

    /**
     * 校验 JWT Token 是否有效。
     *
     * @param token 待校验的 Token
     * @return 是否有效
     */
    public boolean validateToken(String token) {
        // 1. 调用工具类验证 Token 合法性
        return JwtUtils.validateToken(token);
    }

    /**
     * 刷新当前用户的 JWT Token。
     *
     * @return String 是否有效
     */
    public String refreshToken() {
        // 1. 获取当前用户信息
        SysUser sysUser = this.getCurrentUser();

        // 2. 构建 JWT claims，包含用户角色
        Map<String, Object> claims = new HashMap<>();
        claims.put("username", sysUser.getUsername());
        claims.put("roles", sysUser.getRoles());

        // 3. 生成并返回 Token
        return JwtUtils.generateToken(sysUser.getId().toString(), claims);
    }

    /**
     * 根据用户名加载用户信息，用于 Spring Security 登录认证。
     *
     * @param username 用户名
     * @return 用户详情
     * @throws UsernameNotFoundException 用户不存在
     */
    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        // 1. 根据用户名查询用户
        SysUser sysUser = userMapper.selectByUsername(username);

        // 2. 用户不存在则抛出异常
        if (sysUser == null) {
            throw new UsernameNotFoundException("用户不存在: " + username);
        }

        // 3. 返回用户详情（已实现 UserDetails）
        return sysUser;
    }

    /**
     * 修改用户密码。
     *
     * @param oldPassword 旧密码（明文）
     * @param newPassword 新密码（明文）
     * @return 更新后的用户对象
     * @throws ResourceNotFoundException 用户不存在时抛出
     * @throws BusinessException         旧密码验证失败时抛出
     */
    public SysUser changePassword(String oldPassword,
                                  String newPassword) {
        SysUser user = this.getCurrentUser();
        if (!passwordEncoder.matches(oldPassword, user.getPassword())) {
            throw new BusinessException("旧密码不正确");
        }

        LocalDateTime now = LocalDateTime.now();

        // 更新密码
        user.setPassword(passwordEncoder.encode(newPassword));
        user.setPasswordResetAt(now);

        this.fillUpdate(user);
        userMapper.update(user);
        
        // 清除 Redis 中的旧 Token，强制用户重新登录
        boolean deleted = redisStringService.delete(user.getUsername());
        if (!deleted) {
            log.warn("用户ID: {} 密码更新成功，但清除旧 Token 失败", user.getId());
        } else {
            log.info("用户ID: {} 密码更新成功，旧 Token 已清除", user.getId());
        }
        
        return user;
    }

    /**
     * 清除当前用户的 JWT Token 在 Redis 中的缓存。
     *
     * @return boolean 是否成功
     * @throws BusinessException Token 清理失败
     */
    public boolean logout() {
        String username = JwtUtils.getCurrentUsername();
        boolean deleted = redisStringService.delete(username);
        if (!deleted) {
            throw new BusinessException("Token 清除失败");
        }

        return true;
    }

    /**
     * 注销用户。
     *
     * @return 注销结果
     * @throws ResourceNotFoundException 用户不存在时抛出
     */
    public SysUser deleteAccount() {
        SysUser user = this.getCurrentUser();
        if (user == null) {
            throw new ResourceNotFoundException("用户状态异常，无法注销");
        }
        userMapper.deleteById(user.getId(), LocalDateTime.now(), user.getId().toString());
        log.info("用户ID: {} 注销成功", user.getId());
        return user;
    }

    private void fillInsert(SysUser user) {
        LocalDateTime now = LocalDateTime.now();
        String operator = user.getUsername();

        user.setDeleted(0);
        user.setVersion(0);
        user.setPasswordResetAt(now);
        user.setCreatedAt(now);
        user.setUpdatedAt(now);
        user.setUpdatedBy(operator != null ? operator.toString() : "SYSTEM");
        user.setCreatedBy(operator != null ? operator.toString() : "SYSTEM");
    }

    private void fillUpdate(SysUser user) {
        LocalDateTime now = LocalDateTime.now();
        Long operator = JwtUtils.getCurrentUserId();

        Integer currentVersion = user.getVersion();
        user.setVersion(currentVersion == null ? 1 : currentVersion + 1);
        user.setUpdatedAt(now);
        user.setUpdatedBy(operator.toString());
    }

    private void fillUpdate(SysUser user, String operator) {
        LocalDateTime now = LocalDateTime.now();

        Integer currentVersion = user.getVersion();
        user.setVersion(currentVersion == null ? 1 : currentVersion + 1);
        user.setUpdatedAt(now);
        user.setUpdatedBy(operator == null ? "SYSTEM" : operator);
    }
}
