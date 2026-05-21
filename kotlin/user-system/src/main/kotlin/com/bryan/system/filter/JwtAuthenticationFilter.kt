package com.bryan.system.filter

import com.bryan.system.domain.enums.HttpStatus
import com.bryan.system.domain.response.Result
import com.bryan.system.service.auth.AuthService
import com.bryan.system.service.redis.RedisStringService
import com.bryan.system.util.jwt.JwtUtils
import com.fasterxml.jackson.databind.ObjectMapper
import jakarta.servlet.FilterChain
import jakarta.servlet.http.HttpServletRequest
import jakarta.servlet.http.HttpServletResponse
import org.slf4j.LoggerFactory
import org.springframework.http.MediaType
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken
import org.springframework.security.core.authority.SimpleGrantedAuthority
import org.springframework.security.core.context.SecurityContextHolder
import org.springframework.stereotype.Component
import org.springframework.web.filter.OncePerRequestFilter

@Component
class JwtAuthenticationFilter(
    private val authService: AuthService,
    private val objectMapper: ObjectMapper,
    private val redisStringService: RedisStringService
) : OncePerRequestFilter() {
    private val log = LoggerFactory.getLogger(javaClass)

    override fun doFilterInternal(request: HttpServletRequest, response: HttpServletResponse, filterChain: FilterChain) {
        var token = request.getHeader("Authorization")
        if (token == null || !token.startsWith("Bearer ")) {
            filterChain.doFilter(request, response)
            return
        }
        token = token.substring(7)
        try {
            val username = JwtUtils.getUsernameFromToken(token)
            if (redisStringService.get(username) != token) {
                writeUnauthorized(response, "Token已失效，请重新登录")
                return
            }
            val authorities = JwtUtils.getRolesFromToken(token).map(::SimpleGrantedAuthority)
            val user = authService.getCurrentUser()
            if (!user.isEnabled || !user.isAccountNonLocked) {
                writeUnauthorized(response, "用户状态异常或不存在")
                return
            }
            SecurityContextHolder.getContext().authentication =
                UsernamePasswordAuthenticationToken(user, null, authorities)
        } catch (e: Exception) {
            log.warn("Token验证失败: {}", e.javaClass.simpleName)
            writeUnauthorized(response, "Token无效或已过期")
            return
        }
        filterChain.doFilter(request, response)
    }

    private fun writeUnauthorized(response: HttpServletResponse, msg: String) {
        response.contentType = MediaType.APPLICATION_JSON_VALUE
        response.status = HttpServletResponse.SC_UNAUTHORIZED
        response.writer.write(objectMapper.writeValueAsString(Result.error<String>(HttpStatus.UNAUTHORIZED, msg)))
    }
}
