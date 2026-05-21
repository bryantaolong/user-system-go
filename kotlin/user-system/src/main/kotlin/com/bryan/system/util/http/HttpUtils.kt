package com.bryan.system.util.http

import jakarta.servlet.http.HttpServletRequest
import org.springframework.web.context.request.RequestContextHolder
import org.springframework.web.context.request.ServletRequestAttributes

object HttpUtils {
    private fun request(): HttpServletRequest =
        (RequestContextHolder.getRequestAttributes() as ServletRequestAttributes).request

    @JvmStatic
    fun getClientIp(): String {
        val req = request()
        var ip = listOf(
            "X-Forwarded-For",
            "Proxy-Client-IP",
            "WL-Proxy-Client-IP",
            "HTTP_CLIENT_IP",
            "HTTP_X_FORWARDED_FOR"
        ).firstNotNullOfOrNull { h -> req.getHeader(h)?.takeUnless { it.isBlank() || it.equals("unknown", true) } }
            ?: req.remoteAddr
        if (ip.contains(",")) ip = ip.substringBefore(",").trim()
        return ip.ifBlank { "Unknown" }
    }

    @JvmStatic
    fun getClientOS(): String {
        val ua = request().getHeader("User-Agent")?.lowercase() ?: return "Unknown"
        return when {
            "windows" in ua -> "Windows"
            "mac" in ua -> "macOS"
            "x11" in ua -> "Unix"
            "android" in ua -> "Android"
            "iphone" in ua -> "iOS"
            "linux" in ua -> "Linux"
            else -> "Unknown"
        }
    }

    @JvmStatic
    fun getClientBrowser(): String {
        val ua = request().getHeader("User-Agent")?.lowercase() ?: return "Unknown"
        return when {
            "edg/" in ua || "edge/" in ua -> "Edge"
            "opr/" in ua || "opera" in ua -> "Opera"
            "chrome" in ua && "chromium" !in ua -> "Chrome"
            "firefox" in ua || "fxios" in ua -> "Firefox"
            "safari" in ua && "chrome" !in ua -> "Safari"
            "msie" in ua || "trident/7" in ua -> "Internet Explorer"
            else -> "Unknown"
        }
    }
}
