package com.bryan.system.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.client.RestTemplate;

/**
 * 应用级通用配置类
 * 用于注册非业务型、全局使用的第三方 Bean，如 RestTemplate、线程池等。
 *
 * @author Bryan Long
 */
@Configuration
public class AppConfig {

    /**
     * 注册 RestTemplate 单例
     * 供整个应用远程调用 HTTP 接口使用，统一连接池与消息转换器配置。
     *
     * @return RestTemplate 实例
     */
    @Bean
    public RestTemplate restTemplate() {
        return new RestTemplate();
    }
}
