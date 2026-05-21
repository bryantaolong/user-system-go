package http

import (
	"net/http"
	"strings"
)

// GetClientIP 获取客户端真实 IP
func GetClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("HTTP_CLIENT_IP")
	}
	if ip == "" || ip == "unknown" {
		ip = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ip == "" || ip == "unknown" {
		ip = r.RemoteAddr
	}
	// 多级代理时取第一个
	if idx := strings.Index(ip, ","); idx != -1 {
		ip = strings.TrimSpace(ip[:idx])
	}
	if ip == "" {
		return "Unknown"
	}
	return ip
}

// GetClientOS 获取客户端操作系统
func GetClientOS(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac"):
		return "macOS"
	case strings.Contains(ua, "x11"):
		return "Unix"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone"):
		return "iOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}

// GetClientAgent 获取客户端浏览器
func GetClientAgent(userAgent string) string {
	ua := strings.ToLower(userAgent)
	switch {
	case strings.Contains(ua, "edg/") || strings.Contains(ua, "edge/"):
		return "Edge"
	case strings.Contains(ua, "opr/") || strings.Contains(ua, "opera"):
		return "Opera"
	case strings.Contains(ua, "chrome") && !strings.Contains(ua, "chromium"):
		return "Chrome"
	case strings.Contains(ua, "firefox") || strings.Contains(ua, "fxios"):
		return "Firefox"
	case strings.Contains(ua, "safari") && !strings.Contains(ua, "chrome"):
		return "Safari"
	case strings.Contains(ua, "msie") || strings.Contains(ua, "trident/7"):
		return "Internet Explorer"
	default:
		return "Unknown"
	}
}
