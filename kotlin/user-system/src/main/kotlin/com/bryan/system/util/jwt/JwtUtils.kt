package com.bryan.system.util.jwt

import com.bryan.system.config.JwtProperties
import io.jsonwebtoken.Claims
import io.jsonwebtoken.Jwts
import io.jsonwebtoken.SignatureAlgorithm
import io.jsonwebtoken.security.Keys
import jakarta.servlet.http.HttpServletRequest
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Component
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes
import java.nio.charset.StandardCharsets
import java.util.Date
import javax.crypto.SecretKey

@Component
class JwtUtils {
    @Autowired
    fun setJwtProperties(value: JwtProperties) {
        jwtProperties = value
    }

    companion object {
        private lateinit var jwtProperties: JwtProperties

        private fun secretKey(): SecretKey =
            Keys.hmacShaKeyFor(jwtProperties.secretKey.toByteArray(StandardCharsets.UTF_8))

        private fun tokenPrefix(): String = jwtProperties.tokenPrefix

        fun generateToken(userId: String, claims: Map<String, Any> = emptyMap()): String =
            Jwts.builder()
                .setClaims(claims)
                .setSubject(userId)
                .setIssuedAt(Date())
                .setExpiration(Date(System.currentTimeMillis() + jwtProperties.expirationMs))
                .signWith(secretKey(), SignatureAlgorithm.HS256)
                .compact()

        fun getCurrentUserId(): Long = parseCurrentRequestClaims().subject.toLong()

        fun getCurrentUsername(): String = parseCurrentRequestClaims()["username"].toString()

        fun getCurrentUserRoles(): List<String> = getRolesFromClaims(parseCurrentRequestClaims())

        fun getUserIdFromToken(token: String): String = claims(token).subject

        fun getUsernameFromToken(token: String): String = claims(token)["username"] as String

        fun getRolesFromToken(token: String): List<String> = getRolesFromClaims(claims(token))

        fun validateToken(token: String): Boolean = try {
            claims(token)
            true
        } catch (_: Exception) {
            false
        }

        private fun parseCurrentRequestClaims(): Claims {
            val request = currentRequest()
            var token = request.getHeader("Authorization")
            if (token != null && token.startsWith(tokenPrefix())) {
                token = token.substring(tokenPrefix().length)
                return claims(token)
            }
            throw RuntimeException("Missing or invalid Authorization token")
        }

        private fun currentRequest(): HttpServletRequest {
            val attrs = RequestContextHolder.getRequestAttributes() as ServletRequestAttributes?
                ?: throw RuntimeException("No current request")
            return attrs.request
        }

        private fun claims(token: String): Claims =
            Jwts.parser()
                .verifyWith(secretKey())
                .build()
                .parseClaimsJws(token)
                .body

        private fun getRolesFromClaims(claims: Claims): List<String> {
            val roles = claims["roles"]?.toString()
            if (roles.isNullOrBlank()) return emptyList()
            return roles.split(",").map { it.trim() }.filter { it.isNotEmpty() }
                .map { if (it.startsWith("ROLE_")) it else "ROLE_$it" }
        }
    }
}
