package com.bryan.system.service.auth

import com.bryan.system.config.SecurityProperties
import com.bryan.system.domain.entity.SysUser
import com.bryan.system.domain.enums.UserStatusEnum
import com.bryan.system.domain.request.auth.LoginRequest
import com.bryan.system.domain.request.auth.RegisterRequest
import com.bryan.system.exception.BusinessException
import com.bryan.system.exception.ResourceNotFoundException
import com.bryan.system.mapper.UserMapper
import com.bryan.system.mapper.UserRoleMapper
import com.bryan.system.service.redis.RedisStringService
import com.bryan.system.util.http.HttpUtils
import com.bryan.system.util.jwt.JwtUtils
import org.slf4j.LoggerFactory
import org.springframework.security.core.userdetails.UserDetails
import org.springframework.security.core.userdetails.UserDetailsService
import org.springframework.security.core.userdetails.UsernameNotFoundException
import org.springframework.security.crypto.password.PasswordEncoder
import org.springframework.stereotype.Service
import java.time.LocalDateTime

@Service
class AuthService(
    private val userMapper: UserMapper,
    private val userRoleMapper: UserRoleMapper,
    private val passwordEncoder: PasswordEncoder,
    private val redisStringService: RedisStringService,
    private val securityProperties: SecurityProperties
) : UserDetailsService {
    private val log = LoggerFactory.getLogger(javaClass)

    fun register(registerRequest: RegisterRequest): SysUser {
        if (userMapper.selectByUsername(registerRequest.username) != null) throw BusinessException("用户名已存在")
        val defaultRole = userRoleMapper.selectOneByIsDefaultTrue() ?: throw BusinessException("系统未配置默认角色")
        val user = SysUser(
            username = registerRequest.username,
            password = passwordEncoder.encode(registerRequest.password),
            phone = registerRequest.phone,
            email = registerRequest.email,
            roles = defaultRole.roleName,
            status = UserStatusEnum.NORMAL
        )
        fillInsert(user)
        if (userMapper.insert(user) == 0) throw BusinessException("插入数据库失败")
        log.info("用户注册成功: id: {}, username: {}", user.id, user.getUsername())
        return user
    }

    fun login(loginRequest: LoginRequest): String {
        val user = userMapper.selectByUsername(loginRequest.username)
            ?: throw BusinessException("用户名或密码错误")

        if (!passwordEncoder.matches(loginRequest.password, user.getPassword())) {
            val now = LocalDateTime.now()
            user.loginFailCount = (user.loginFailCount ?: 0) + 1
            fillUpdate(user, user.id?.toString() ?: "SYSTEM")
            if ((user.loginFailCount ?: 0) >= securityProperties.loginFailLimit) {
                user.status = UserStatusEnum.LOCKED
                user.lockedAt = now
                userMapper.update(user)
                throw BusinessException("输入密码错误次数过多，账号锁定")
            }
            userMapper.update(user)
            throw BusinessException("用户名或密码错误")
        }

        val existingToken = redisStringService.get(user.getUsername())
        if (existingToken != null && JwtUtils.validateToken(existingToken)) {
            redisStringService.setExpire(user.getUsername(), 86400000 / 1000)
            return existingToken
        }

        user.lastLoginAt = LocalDateTime.now()
        user.lastLoginIp = HttpUtils.getClientIp()
        user.lastLoginDevice = "${HttpUtils.getClientOS()} / ${HttpUtils.getClientBrowser()}"
        user.loginFailCount = 0
        fillUpdate(user, user.id?.toString() ?: "SYSTEM")
        userMapper.update(user)

        val token = JwtUtils.generateToken(user.id.toString(), mapOf("username" to user.getUsername(), "roles" to (user.roles ?: "")))
        if (!redisStringService.set(user.getUsername(), token, 86400000 / 1000)) throw BusinessException("Token 存储失败")
        return token
    }

    fun getCurrentUserId(): Long = JwtUtils.getCurrentUserId()
    fun getCurrentUsername(): String = JwtUtils.getCurrentUsername()
    fun getCurrentUser(): SysUser = userMapper.selectById(JwtUtils.getCurrentUserId()) ?: throw ResourceNotFoundException("用户不存在")
    fun isAdmin(): Boolean = JwtUtils.getCurrentUserRoles().any { it == "ROLE_ADMIN" }
    fun isAdmin(userDetails: UserDetails): Boolean = userDetails.authorities.any { it.authority == "ROLE_ADMIN" }
    fun validateToken(token: String): Boolean = JwtUtils.validateToken(token)
    fun refreshToken(): String {
        val user = getCurrentUser()
        return JwtUtils.generateToken(user.id.toString(), mapOf("username" to user.getUsername(), "roles" to (user.roles ?: "")))
    }

    override fun loadUserByUsername(username: String): UserDetails =
        userMapper.selectByUsername(username) ?: throw UsernameNotFoundException("用户不存在: $username")

    fun changePassword(oldPassword: String?, newPassword: String?): SysUser {
        val user = getCurrentUser()
        if (!passwordEncoder.matches(oldPassword, user.getPassword())) throw BusinessException("旧密码不正确")
        user.setPassword(passwordEncoder.encode(newPassword))
        user.passwordResetAt = LocalDateTime.now()
        fillUpdate(user)
        userMapper.update(user)
        redisStringService.delete(user.getUsername())
        return user
    }

    fun logout(): Boolean {
        if (!redisStringService.delete(JwtUtils.getCurrentUsername())) throw BusinessException("Token 清除失败")
        return true
    }

    fun deleteAccount(): SysUser {
        val user = getCurrentUser()
        userMapper.deleteById(user.id, LocalDateTime.now(), user.id.toString())
        return user
    }

    private fun fillInsert(user: SysUser) {
        val now = LocalDateTime.now()
        val operator = user.getUsername().ifBlank { "SYSTEM" }
        user.deleted = 0
        user.version = 0
        user.passwordResetAt = now
        user.createdAt = now
        user.updatedAt = now
        user.createdBy = operator
        user.updatedBy = operator
    }

    private fun fillUpdate(user: SysUser) = fillUpdate(user, JwtUtils.getCurrentUserId().toString())

    private fun fillUpdate(user: SysUser, operator: String?) {
        user.version = (user.version ?: 0) + 1
        user.updatedAt = LocalDateTime.now()
        user.updatedBy = operator ?: "SYSTEM"
    }
}
