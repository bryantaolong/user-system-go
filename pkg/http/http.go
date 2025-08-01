package http

import (
	"net/http"
	"strings"
)

// GetClientIP 获取客户端IP地址
func GetClientIP(r *http.Request) string {
	// 尝试从各种可能的Header中获取IP
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" || strings.EqualFold(ip, "unknown") {
		ip = r.Header.Get("Proxy-Client-IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") {
		ip = r.Header.Get("WL-Proxy-Client-IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") {
		ip = r.Header.Get("HTTP_CLIENT_IP")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") {
		ip = r.Header.Get("HTTP_X_FORWARDED_FOR")
	}
	if ip == "" || strings.EqualFold(ip, "unknown") {
		ip = r.RemoteAddr
		// 去除端口号
		if strings.Contains(ip, ":") {
			ip = strings.Split(ip, ":")[0]
		}
	}

	// 对于通过多个代理的情况，第一个IP为客户端真实IP
	if strings.Contains(ip, ",") {
		ip = strings.TrimSpace(strings.Split(ip, ",")[0])
	}

	return ip
}

// GetClientOS 获取客户端操作系统
func GetClientOS(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown"
	}

	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "windows"):
		return "Windows"
	case strings.Contains(ua, "mac"):
		return "Mac OS"
	case strings.Contains(ua, "x11"):
		return "Unix"
	case strings.Contains(ua, "android"):
		return "Android"
	case strings.Contains(ua, "iphone") || strings.Contains(ua, "ipad"):
		return "iOS"
	case strings.Contains(ua, "linux"):
		return "Linux"
	default:
		return "Unknown"
	}
}

// GetClientBrowser 获取客户端浏览器
func GetClientBrowser(r *http.Request) string {
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		return "Unknown"
	}

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
