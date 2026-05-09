package com.bryan.system.config.properties;

import lombok.Getter;
import lombok.Setter;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

/**
 * JWT 配置属性类
 * 用于从配置文件中读取 JWT 相关配置，避免硬编码。
 *
 * @author Bryan Long
 */
@Setter
@Getter
@Component
@ConfigurationProperties(prefix = "jwt")
public class JwtProperties {

    /**
     * JWT 签名密钥
     * 建议至少 256 位（32 字符）
     */
    private String secretKey;

    /**
     * JWT Token 过期时间（毫秒）
     * 默认 24 小时
     */
    private Long expirationMs = 86400000L;

    /**
     * JWT Token 前缀
     * 默认 "Bearer "
     */
    private String tokenPrefix = "Bearer ";

}
