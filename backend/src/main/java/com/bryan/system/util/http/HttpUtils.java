package com.bryan.system.util.http;

import jakarta.servlet.http.HttpServletRequest;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

import java.util.Objects;

/**
 * HTTP 工具类
 * 提供获取客户端 IP、操作系统、浏览器等常用能力。
 *
 * @author Bryan Long
 */
public class HttpUtils {

    /**
     * 获取客户端真实 IP
     * 依次解析 X-Forwarded-For、Proxy-Client-IP、WL-Proxy-Client-IP、
     * HTTP_CLIENT_IP、HTTP_X_FORWARDED_FOR，最后回退到 RemoteAddr。
     *
     * @return 客户端 IP；未知返回 "Unknown"
     */
    public static String getClientIp() {
        ServletRequestAttributes attributes = (ServletRequestAttributes) Objects.requireNonNull(
                RequestContextHolder.getRequestAttributes());
        HttpServletRequest request = attributes.getRequest();

        String ip = request.getHeader("X-Forwarded-For");
        if (ip == null || ip.isEmpty() || "unknown".equalsIgnoreCase(ip)) {
            ip = request.getHeader("Proxy-Client-IP");
        }
        if (ip == null || ip.isEmpty() || "unknown".equalsIgnoreCase(ip)) {
            ip = request.getHeader("WL-Proxy-Client-IP");
        }
        if (ip == null || ip.isEmpty() || "unknown".equalsIgnoreCase(ip)) {
            ip = request.getHeader("HTTP_CLIENT_IP");
        }
        if (ip == null || ip.isEmpty() || "unknown".equalsIgnoreCase(ip)) {
            ip = request.getHeader("HTTP_X_FORWARDED_FOR");
        }
        if (ip == null || ip.isEmpty() || "unknown".equalsIgnoreCase(ip)) {
            ip = request.getRemoteAddr();
        }

        // 多级代理时，取第一个真实 IP
        if (ip != null && ip.contains(",")) {
            ip = ip.substring(0, ip.indexOf(",")).trim();
        }
        return ip == null ? "Unknown" : ip;
    }

    /**
     * 获取客户端操作系统
     *
     * @return 操作系统名称；未知返回 "Unknown"
     */
    public static String getClientOS() {
        ServletRequestAttributes attributes = (ServletRequestAttributes) Objects.requireNonNull(
                RequestContextHolder.getRequestAttributes());
        HttpServletRequest request = attributes.getRequest();

        String userAgent = request.getHeader("User-Agent");
        if (userAgent == null) {
            return "Unknown";
        }

        String ua = userAgent.toLowerCase();

        if (ua.contains("windows")) {
            return "Windows";
        } else if (ua.contains("mac")) {
            return "macOS";
        } else if (ua.contains("x11")) {
            return "Unix";
        } else if (ua.contains("android")) {
            return "Android";
        } else if (ua.contains("iphone")) {
            return "iOS";
        } else if (ua.contains("linux")) {
            return "Linux";
        } else {
            return "Unknown";
        }
    }

    /**
     * 获取客户端浏览器
     *
     * @return 浏览器名称；未知返回 "Unknown"
     */
    public static String getClientBrowser() {
        ServletRequestAttributes attributes = (ServletRequestAttributes) Objects.requireNonNull(
                RequestContextHolder.getRequestAttributes());
        HttpServletRequest request = attributes.getRequest();

        String userAgent = request.getHeader("User-Agent");
        if (userAgent == null) {
            return "Unknown";
        }

        String ua = userAgent.toLowerCase();

        if (ua.contains("edg/") || ua.contains("edge/")) {
            return "Edge";
        } else if (ua.contains("opr/") || ua.contains("opera")) {
            return "Opera";
        } else if (ua.contains("chrome") && !ua.contains("chromium")) {
            return "Chrome";
        } else if (ua.contains("firefox") || ua.contains("fxios")) {
            return "Firefox";
        } else if (ua.contains("safari") && !ua.contains("chrome")) {
            return "Safari";
        } else if (ua.contains("msie") || ua.contains("trident/7")) {
            return "Internet Explorer";
        } else {
            return "Unknown";
        }
    }
}