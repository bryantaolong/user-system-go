package com.bryan.system.filter;

import com.bryan.system.domain.entity.SysUser;
import com.bryan.system.domain.enums.HttpStatus;
import com.bryan.system.domain.response.Result;
import com.bryan.system.service.auth.AuthService;
import com.bryan.system.service.redis.RedisStringService;
import com.bryan.system.util.jwt.JwtUtils;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.MediaType;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Collection;
import java.util.List;
import java.util.stream.Collectors;

/**
 * JWT 认证过滤器
 * 负责解析请求头中的 Bearer Token，验证 Redis 白名单，并构建 Spring Security 上下文。
 *
 * @author Bryan Long
 */
@Slf4j
@Component
@RequiredArgsConstructor
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    private final AuthService authService;
    private final ObjectMapper objectMapper;
    private final RedisStringService redisStringService;

    /**
     * 单次请求过滤逻辑
     *
     * @param request     当前请求
     * @param response    当前响应
     * @param filterChain 过滤器链
     * @throws ServletException 过滤器异常
     * @throws IOException      IO 异常
     */
    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) throws ServletException, IOException {
        String token = request.getHeader("Authorization");
        // 如果没有 Token 或者 Token 格式不正确，则直接放行，让后续的 Security 配置处理（例如匿名访问或认证失败）
        if (token == null || !token.startsWith("Bearer ")) {
            filterChain.doFilter(request, response);
            return;
        }

        token = token.substring(7); // 截取掉 "Bearer " 前缀

        try {
            // Redis Token 白名单验证
            String username = JwtUtils.getUsernameFromToken(token);
            String redisToken = redisStringService.get(username);

            if (redisToken == null || !redisToken.equals(token)) {
                this.writeUnauthorized(response, "Token已失效，请重新登录");
                return;
            }

            // 从 Token 的 Claims 中获取角色列表
            List<String> roles = JwtUtils.getRolesFromToken(token);
            // 将角色字符串列表转换为 Spring Security 的 GrantedAuthority 列表
            Collection<? extends GrantedAuthority> authorities = roles.stream()
                    .map(SimpleGrantedAuthority::new)
                    .collect(Collectors.toList());

            // 从数据库加载用户主体，验证用户状态
            // 这里不再需要 authService.getCurrentUser() 来获取角色，
            // 权限信息直接从 Token 中获取，提高性能并符合JWT无状态原则。
            // 但为了安全起见，通常还会从数据库加载用户主体（User对象），以验证用户状态等。
            // 这里我们仍然从数据库加载用户，以确保用户是存在的且状态正常。
            SysUser sysUser = authService.getCurrentUser();
            if (sysUser == null || !sysUser.isEnabled() || !sysUser.isAccountNonLocked()) {
                this.writeUnauthorized(response, "用户状态异常或不存在");
                return;
            }

            // 构建认证对象，使用从 Token 和数据库验证后的权限
            UsernamePasswordAuthenticationToken authentication =
                    new UsernamePasswordAuthenticationToken(sysUser, null, authorities);
            SecurityContextHolder.getContext().setAuthentication(authentication);

        } catch (Exception e) {
            // 处理 Token 无效或过期的情况
            // 不返回具体异常信息，防止信息泄露
            log.warn("Token验证失败: {}", e.getClass().getSimpleName());
            this.writeUnauthorized(response, "Token无效或已过期");
            return;
        }

        filterChain.doFilter(request, response);
    }

    /**
     * 快速写入 401 响应
     *
     * @param response 响应对象
     * @param msg      错误描述
     * @throws IOException 写出异常
     */
    private void writeUnauthorized(HttpServletResponse response, String msg) throws IOException {
        response.setContentType(MediaType.APPLICATION_JSON_VALUE);
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.getWriter().write(
                objectMapper.writeValueAsString(Result.error(HttpStatus.UNAUTHORIZED, msg))
        );
    }
}
