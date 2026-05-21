package com.bryan.system.config

import org.springframework.boot.context.properties.ConfigurationProperties

@ConfigurationProperties(prefix = "jwt")
data class JwtProperties(
    var secretKey: String = "",
    var expirationMs: Long = 86400000,
    var tokenPrefix: String = "Bearer "
)

@ConfigurationProperties(prefix = "security")
data class SecurityProperties(
    var loginFailLimit: Int = 5,
    var loginFailResetMinutes: Long = 30,
    var accountLockDurationMinutes: Long = 30
)
