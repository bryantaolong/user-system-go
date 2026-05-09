package com.bryan.system.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.ResourceHandlerRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

import java.nio.file.Paths;

/**
 * Spring MVC 全局配置类
 * 统一配置跨域、拦截器、格式化器等 Web 层行为。
 *
 * @author Bryan Long
 */
@Configuration
public class WebMvcConfig implements WebMvcConfigurer {

    @Value("${file.upload-dir}")
    private String uploadDir;

    /**
     * 配置静态资源映射
     * 将 /uploads/** 的请求映射到物理上传目录，便于前端访问上传的头像、附件等。
     *
     * @param registry ResourceHandlerRegistry 注册器
     */
    @Override
    public void addResourceHandlers(ResourceHandlerRegistry registry) {
        String absolutePath = Paths.get(uploadDir).toAbsolutePath().toUri().toString();
        registry.addResourceHandler("/uploads/**")
                .addResourceLocations(absolutePath);
    }

    @Value("${cors.allowed-origins:http://localhost:5173}")  // 默认为开发环境前端地址
    private String allowedOrigins;

    /**
     * 配置全局跨域规则
     * 允许前端开发服务器（localhost:5173）跨域访问所有接口，支持常见 HTTP 方法与凭证传递。
     *
     * @param registry CorsRegistry 注册器
     */
    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**")
                .allowedOrigins(allowedOrigins)
                .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
                .allowedHeaders("*")
                .allowCredentials(true)
                .maxAge(3600);
    }
}
